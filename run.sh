#!/usr/bin/env bash
set -euo pipefail

OPTIONS_FILE=/data/options.json
TPL=/templates/telegraf.conf.tpl
CONF=/etc/telegraf/telegraf.conf

# Create destination folder
mkdir -p /etc/telegraf

# Extract each block from options.json
export AGENT=$(jq -r '.agent // ""' "$OPTIONS_FILE")
export OUTPUT_PLUGINS=$(jq -r '.output_plugins // ""' "$OPTIONS_FILE")
export INPUT_PLUGINS=$(jq -r '.input_plugins // ""' "$OPTIONS_FILE")
export PROCESSOR_PLUGINS=$(jq -r '.processor_plugins // ""' "$OPTIONS_FILE")

# Generate telegraf.conf from template
envsubst < "$TPL" > "$CONF"

echo "====================================="
echo " Generated final telegraf.conf:"
echo "====================================="
cat "$CONF"
echo "====================================="

# Start Telegraf
exec /usr/bin/telegraf --config "$CONF"
