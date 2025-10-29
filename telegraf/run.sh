#!/usr/bin/env sh
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

# # Start Telegraf
# exec /usr/bin/telegraf --config "$CONF"

echo "-----------------------------------------------------------"
echo "Starting Telegraf add-on for Home Assistant..."
echo "-----------------------------------------------------------"

# Ensure Telegraf binary exists and is executable
if [ ! -x "/usr/bin/telegraf" ]; then
  echo "Setting execution permissions for /usr/bin/telegraf..."
  chmod +x /usr/bin/telegraf
fi


# Print version info for logging
echo "Telegraf version:"
telegraf --version || echo "⚠️  Could not run telegraf binary"

# Start Telegraf with live config reload
echo "Running: telegraf --config \"$CONF\" --watch-config poll"
exec telegraf --config "$CONF" --watch-config poll
