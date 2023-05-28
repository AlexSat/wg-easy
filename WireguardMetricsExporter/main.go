package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type clientMetrics struct {
	PubKey, PresharedKey, ClientIp, ClientName string
	LatestHandshake                            time.Time
	ReceivedBytes, SendBytes                   int64
}

const (
	SUBSYSTEM = "wireguard"
)

var ReceivedBytes = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name:      "received_bytes_gauge",
	Help:      "Amount of bytes received by client",
	Subsystem: SUBSYSTEM,
}, []string{"client_name"})

var SendBytes = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name:      "send_bytes_gauge",
	Help:      "Amount of bytes send by client",
	Subsystem: SUBSYSTEM,
}, []string{"client_name"})

var LastHandshakeTimestamp = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name:      "last_handshake_unixtimestamp",
	Help:      "Date time in unixtimestamp seconds when last handshake activity registered",
	Subsystem: SUBSYSTEM,
}, []string{"client_name"})

func GetAndServeTestMetric() {
	for {
		time.Sleep(time.Second * 5)

		wgOutput, err := exec.Command("wg", "show", "all", "dump").Output()
		if err != nil {
			fmt.Println(err)
			continue
		}
		wgData := bytes.NewReader(wgOutput)
		scanner := bufio.NewScanner(wgData)
		rawMetrics := make(map[string]clientMetrics)
		for scanner.Scan() {
			line := scanner.Text()
			lines := strings.Split(line, "\t")
			dataLines := make([]string, 0, 9)
			for _, l := range lines {
				if l != "" {
					dataLines = append(dataLines, l)
				}
			}
			if len(dataLines) == 9 {
				latest_handshake_uint, err := strconv.ParseInt(dataLines[5], 0, 64)
				if err != nil {
					continue
				}
				latest_handshake := time.Unix(latest_handshake_uint, 0)
				received, err := strconv.ParseInt(dataLines[6], 0, 64)
				if err != nil {
					continue
				}
				send, err := strconv.ParseInt(dataLines[7], 0, 64)
				if err != nil {
					continue
				}
				cm := clientMetrics{
					PresharedKey:    dataLines[2],
					PubKey:          dataLines[1],
					ClientIp:        dataLines[3],
					LatestHandshake: latest_handshake,
					ReceivedBytes:   received,
					SendBytes:       send,
				}
				rawMetrics[cm.PubKey] = cm
			}
		}

		var data interface{}
		plan, err := os.ReadFile("/etc/wireguard/wg0.json")
		if err != nil {
			switch err.(type) {
			case *fs.PathError:
				break
			default:
				fmt.Println(err)
			}
			continue
		}
		err = json.Unmarshal(plan, &data)
		if err != nil {
			fmt.Println(err)
			continue
		}

		clients, ok := data.(map[string]interface{})["clients"]
		if !ok {
			fmt.Println(err)
			continue
		}
		for _, cdata := range clients.(map[string]interface{}) {
			cdatatyped := cdata.(map[string]interface{})
			if publicKeyRaw, ok := cdatatyped["publicKey"]; ok {
				publicKey := publicKeyRaw.(string)
				if nameRaw, ok := cdatatyped["name"]; ok {
					name := nameRaw.(string)
					// найти по publicKey клиента и проставить ему name
					if rawClient, ok := rawMetrics[publicKey]; ok {
						rawClient.ClientName = name
						rawMetrics[publicKey] = rawClient
					}
				}
			}
		}

		for _, client := range rawMetrics {
			ReceivedBytes.WithLabelValues(client.ClientName).Set(float64(client.ReceivedBytes))
			SendBytes.WithLabelValues(client.ClientName).Set(float64(client.SendBytes))
			LastHandshakeTimestamp.WithLabelValues(client.ClientName).Set(float64(client.LatestHandshake.Unix()))
		}
	}
}

func main() {
	fmt.Println("Starting metrics exporter...")
	r := prometheus.NewRegistry()
	r.MustRegister(ReceivedBytes, SendBytes, LastHandshakeTimestamp)
	handler := promhttp.HandlerFor(r, promhttp.HandlerOpts{})
	go GetAndServeTestMetric()
	http.Handle("/metrics", handler)
	err := http.ListenAndServe(":9001", nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Stopping metrics exporter...")
}
