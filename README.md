# Repository is archived

Hard to support and fetch upstream updates -> moved to new generation 'all-in-one' https://github.com/AlexSat/wg-easy-behind-wstunnel

# WireGuard Easy

You have found the easiest way to install & manage WireGuard on any Linux host!

<p align="center">
  <img src="./assets/screenshot.png" width="802" />
  <img src="./assets/grafana.png" width="802" />
</p>


## Features

* All-in-one: WireGuard + Web UI + Moonitoring (Prometheus + Grafana).
* Easy installation, simple to use.
* List, create, edit, delete, enable & disable clients.
* Show a client's QR code.
* Download a client's configuration file.
* Statistics for which clients are connected (Web UI admin + Grafana).
* Tx/Rx charts for each connected client.
* Gravatar support.

## Requirements

* A host with a kernel that supports WireGuard (all modern kernels).
* A host with Docker installed.

## Installation

### 1. Install Docker

If you haven't installed Docker yet, install it by running:

```bash
$ curl -sSL https://get.docker.com | sh
$ sudo usermod -aG docker $(whoami)
$ exit
```

And log in again.

### 2. Run WireGuard Easy with docker compose

Replace environment variables in `docker-compose.yml` and start it with `docker compose up -d`

The Web UI admin for Wireguard will now be available on `http://0.0.0.0:51821`.
The Web monitoring (Grafana) for Wireguard will now be available on `http://0.0.0.0:3000` (`login`:`pass` is `admin`:`admin`).

> 💡 Your configuration files will be saved in `~/wireguard`, `~/prometheus` and `./grafana`

### 3. Sponsor

Are you enjoying this project? [Buy him a beer!](https://github.com/sponsors/WeeJeWel) 🍻

## Options

These options can be configured by setting environment variables using `-e KEY="VALUE"` in the `docker run` command. (BE CAREFULLY AND CHECK ENVIRONMENT VARIABLES IN `docker-compose.yml`)

| Env | Default | Example | Description |
| - | - | - | - |
| `PASSWORD` | `p123123` | `foobar123` | When set, requires a password when logging in to the Web UI. |
| `WG_HOST` | - | `vpn.myserver.com` | The public hostname of your VPN server. |
| `WG_DEVICE` | `eth0` | `ens6f0` | Ethernet device the wireguard traffic should be forwarded through. |
| `WG_PORT` | `51820` | `12345` | The public UDP port of your VPN server. WireGuard will always listen on `51820` inside the Docker container. |
| `WG_MTU` | `null` | `1420` | The MTU the clients will use. Server uses default WG MTU. |
| `WG_PERSISTENT_KEEPALIVE` | `30` | `0` | Value in seconds to keep the "connection" open. If this value is 0, then connections won't be kept alive. |
| `WG_DEFAULT_ADDRESS` | `10.8.0.x` | `10.6.0.x` | Clients IP address range. |
| `WG_DEFAULT_DNS` | `1.1.1.1` | `8.8.8.8, 8.8.4.4` | DNS server clients will use. |
| `WG_ALLOWED_IPS` | `0.0.0.0/0, ::/0` | `192.168.15.0/24, 10.0.1.0/24` | Allowed IPs clients will use. |
| `WG_PRE_UP` | `...` | - | See [config.js](https://github.com/WeeJeWel/wg-easy/blob/master/src/config.js#L19) for the default value. |
| `WG_POST_UP` | `...` | `iptables ...` | See [config.js](https://github.com/WeeJeWel/wg-easy/blob/master/src/config.js#L20) for the default value. |
| `WG_PRE_DOWN` | `...` | - | See [config.js](https://github.com/WeeJeWel/wg-easy/blob/master/src/config.js#L27) for the default value. |
| `WG_POST_DOWN` | `...` | `iptables ...` | See [config.js](https://github.com/WeeJeWel/wg-easy/blob/master/src/config.js#L28) for the default value. |

> If you change `WG_PORT`, make sure to also change the exposed port.

## Updating

To update to the latest version, simply run:

```bash
docker compose down
docker composse pull
docker compose up -d
```

## Common Use Cases

* [Using WireGuard-Easy with Pi-Hole](https://github.com/WeeJeWel/wg-easy/wiki/Using-WireGuard-Easy-with-Pi-Hole)
* [Using WireGuard-Easy with nginx/SSL](https://github.com/WeeJeWel/wg-easy/wiki/Using-WireGuard-Easy-with-nginx-SSL)
