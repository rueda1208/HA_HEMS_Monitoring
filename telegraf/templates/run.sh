#!/bin/sh
set -e

CONFIG_FILE="/etc/telegraf/telegraf.conf"

echo "[INFO] Starting Telegraf configuration generation..."

# # --- 1. SETUP ---
# mkdir -p /etc/telegraf

# # --- 2. GENERATE CORE CONFIG (Agent & Outputs) ---
# echo "[INFO] Generating [agent] and [[outputs]] sections..."

# # Reset file (overwrite)
# cat > "$CONFIG_FILE" <<EOF
# # ===============================
# #  Telegraf Configuration File
# #  (auto-generated)
# # ===============================

# EOF

# # Agent config (optional)
# if [ -n "${AGENT_CONFIG:-}" ]; then
#   echo "[INFO] Adding Agent config..."
#   printf "%s\n\n" "# Agent Configuration" "$AGENT_CONFIG" >> "$CONFIG_FILE"
# fi

# # Outputs
# if [ -n "${OUTPUT_PLUGINS:-}" ]; then
#   echo "[INFO] Adding Output Plugins..."
#   printf "%s\n\n" "# Output Plugins" "$OUTPUT_PLUGINS" >> "$CONFIG_FILE"
# fi

# # --- 3. APPEND INPUT / PROCESSOR PLUGINS ---
# if [ -n "${INPUT_PLUGINS:-}" ]; then
#   echo "[INFO] Adding Input Plugins..."
#   printf "%s\n\n" "# Input Plugins" "$INPUT_PLUGINS" >> "$CONFIG_FILE"
# fi

# if [ -n "${PROCESSOR_PLUGINS:-}" ]; then
#   echo "[INFO] Adding Processor Plugins..."
#   printf "%s\n\n" "# Processor Plugins" "$PROCESSOR_PLUGINS" >> "$CONFIG_FILE"
# fi

# # --- 4. OVERWRITE CUSTOM PLUGIN CONFIG FILE IN VOLUME ---
# PLUGIN_CONFIG_FILE="${PLUGIN_CONFIG_FILE:-custom_telegraf.conf}"
# PLUGIN_CONFIG_PATH="/share/${PLUGIN_CONFIG_FILE}"

# echo "[INFO] Preparing to overwrite plugin file at: ${PLUGIN_CONFIG_PATH}"

# if [ -f "$PLUGIN_CONFIG_PATH" ]; then
#   echo "[INFO] Existing custom plugin file found. It will be overwritten."
# else
#   echo "[INFO] No existing custom plugin file found. Creating new one."
# fi

# cp "$CONFIG_FILE" "$PLUGIN_CONFIG_PATH"

# echo "[INFO] Custom plugin file successfully written to ${PLUGIN_CONFIG_PATH}"

# --- 5. START TELEGRAF ---
echo "[INFO] Starting Telegraf with config: ${CONFIG_FILE}"
exec telegraf --config "$CONFIG_FILE"
