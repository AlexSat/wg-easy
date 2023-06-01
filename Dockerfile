# There's an issue with node:16-alpine.
# On Raspberry Pi, the following crash happens:

# #FailureMessage Object: 0x7e87753c
# #
# # Fatal error in , line 0
# # unreachable code
# #
# #
# #


FROM docker.io/library/node:14-alpine@sha256:dc92f36e7cd917816fa2df041d4e9081453366381a00f40398d99e9392e78664 AS build_node_modules

# Copy Web UI
COPY src/ /app/
WORKDIR /app
RUN npm ci --production

FROM golang:1.20-alpine AS build_metrics_exporter
WORKDIR /app
COPY WireguardMetricsExporter/go.mod WireguardMetricsExporter/go.sum ./
RUN go mod download

COPY WireguardMetricsExporter/*.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /wg_metrics_exporter

# Copy build result to a new image.
# This saves a lot of disk space.
FROM docker.io/library/node:14-alpine@sha256:dc92f36e7cd917816fa2df041d4e9081453366381a00f40398d99e9392e78664
COPY --from=build_node_modules /app /app
COPY --from=build_metrics_exporter /wg_metrics_exporter /app

# Move node_modules one directory up, so during development
# we don't have to mount it in a volume.
# This results in much faster reloading!
#
# Also, some node_modules might be native, and
# the architecture & OS of your development machine might differ
# than what runs inside of docker.
RUN mv /app/node_modules /node_modules

# Enable this to run `npm run serve`
RUN npm i -g nodemon

# Install Linux packages
RUN apk add -U --no-cache \
  wireguard-tools

ENV S6_OVERLAY_VERSION=3.1.5.0

# Install supervisor
ADD https://github.com/just-containers/s6-overlay/releases/download/v${S6_OVERLAY_VERSION}/s6-overlay-noarch.tar.xz /tmp
RUN tar -C / -Jxpf /tmp/s6-overlay-noarch.tar.xz
ADD https://github.com/just-containers/s6-overlay/releases/download/v${S6_OVERLAY_VERSION}/s6-overlay-x86_64.tar.xz /tmp
RUN tar -C / -Jxpf /tmp/s6-overlay-x86_64.tar.xz

#Create s6-overlay configuration to run wireguard with ui and metrics exporter
COPY s6-rc.d /etc/s6-overlay/s6-rc.d
RUN mkdir -p /var/log/node_service
RUN touch /var/log/node_service/current
RUN chown -R nobody:nogroup /var/log
VOLUME /var/log/

RUN \
	chmod +x /etc/s6-overlay/s6-rc.d/node_service/run && \
	chmod +x /etc/s6-overlay/s6-rc.d/metrics_exporter/run && \
	chmod +x /etc/s6-overlay/s6-rc.d/node_service_logger/run

# Expose Ports
EXPOSE 51820/udp
EXPOSE 51821/tcp

# Set Environment
ENV DEBUG=Server,WireGuard

# Run Web UI
WORKDIR /app

ENTRYPOINT ["/init"]
