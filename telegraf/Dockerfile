# -------- Stage 1: Build Telegraf with custom plugins --------
FROM golang:1.23-alpine AS build
ENV CGO_ENABLED=0 

# Install dependencies and tools
RUN apk add --no-cache git make

# Copy the Telegraf source code into the Docker container
WORKDIR /go/src
RUN git clone --branch v1.31.2 --depth=1 https://github.com/influxdata/telegraf.git

# Optional: copy your custom processors
# ARG BUILDING
# COPY buildings/${BUILDING}/all.go telegraf/plugins/processors/all/
# COPY buildings/${BUILDING}/http_processor telegraf/plugins/processors/http_processor/
# COPY buildings/${BUILDING}/modbus_processor telegraf/plugins/processors/modbus_processor/

# Build Telegraf with your custom plugins
WORKDIR /go/src/telegraf
RUN go build -o telegraf ./cmd/telegraf && strip telegraf


# -------- Stage 2: Runtime image --------
ARG BUILD_FROM=ghcr.io/home-assistant/amd64-base:3.19
FROM ${BUILD_FROM}

# Install runtime dependencies
RUN apk add --no-cache jq gettext

# Prepare directories
RUN mkdir -p /etc/telegraf /templates

# Copy compiled Telegraf binary
COPY --from=build /go/src/telegraf/telegraf /usr/bin/telegraf

# Copy add-on scripts and templates
COPY run.sh /run.sh
COPY templates/telegraf.conf.tpl /templates/telegraf.conf.tpl

# Set permissions
RUN chmod a+x /run.sh

# Set the default command (run script)
CMD ["/run.sh"]