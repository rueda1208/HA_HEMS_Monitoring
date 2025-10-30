#!/usr/bin/env bash
set -euo pipefail

# Path for the dynamically generated full config file
CONFIG_FILE="/etc/telegraf/telegraf.conf"

echo "[INFO] Starting Telegraf configuration generation..."

# --- 1. SETUP ---
mkdir -p /etc/telegraf

# --- 2. GENERATE CORE CONFIG (Agent & Outputs) ---
echo "[INFO] Generating [agent] and [[outputs]] sections..."

# Reset file (overwrite)
cat > "${CONFIG_FILE}" <<EOF
# ===============================
#  Telegraf Configuration File
#  (auto-generated)
# ===============================

EOF

# Agent config (optional, from ENV)
if [[ -n "${AGENT_CONFIG:-}" ]]; then
  echo "[INFO] Adding Agent config..."
  cat >> "${CONFIG_FILE}" <<EOF
# Agent Configuration
${AGENT_CONFIG}

EOF
fi

# Outputs (e.g., InfluxDB, PostgreSQL)
if [[ -n "${OUTPUT_PLUGINS:-}" ]]; then
  echo "[INFO] Adding Output Plugins..."
  cat >> "${CONFIG_FILE}" <<EOF
# Output Plugins
${OUTPUT_PLUGINS}

EOF
fi

# --- 3. APPEND INPUT / PROCESSOR PLUGINS ---
if [[ -n "${INPUT_PLUGINS:-}" ]]; then
  echo "[INFO] Adding Input Plugins..."
  cat >> "${CONFIG_FILE}" <<EOF
# Input Plugins
${INPUT_PLUGINS}

EOF
fi

if [[ -n "${PROCESSOR_PLUGINS:-}" ]]; then
  echo "[INFO] Adding Processor Plugins..."
  cat >> "${CONFIG_FILE}" <<EOF
# Processor Plugins
${PROCESSOR_PLUGINS}

EOF
fi

# --- 4. OVERWRITE CUSTOM PLUGIN CONFIG FILE IN VOLUME ---
PLUGIN_CONFIG_FILE="${PLUGIN_CONFIG_FILE:-custom_telegraf.conf}"
PLUGIN_CONFIG_PATH="/share/${PLUGIN_CONFIG_FILE}"

echo "[INFO] Preparing to overwrite plugin file at: ${PLUGIN_CONFIG_PATH}"

if [[ -f "${PLUGIN_CONFIG_PATH}" ]]; then
  echo "[INFO] Existing custom plugin file found. It will be overwritten."
else
  echo "[INFO] No existing custom plugin file found. Creating new one."
fi

cp "${CONFIG_FILE}" "${PLUGIN_CONFIG_PATH}"

echo "[INFO] Custom plugin file successfully written to ${PLUGIN_CONFIG_PATH}"

# --- 5. START TELEGRAF ---
echo "[INFO] Starting Telegraf with config: ${CONFIG_FILE}"
exec telegraf --config "${CONFIG_FILE}"