#!/usr/bin/env bashio
# set -euo pipefail # bashio handles error checking, but keeping good practice

# Path for the dynamically generated full config file
CONFIG_FILE="/etc/telegraf/telegraf.conf"

bashio::log.info "Starting Telegraf configuration generation from HA UI options..."

# --- 1. SETUP ---
# Create destination folder (if not already handled by Dockerfile/Add-on base image)
mkdir -p /etc/telegraf

# --- 2. GENERATE CORE CONFIG (Agent & Outputs) ---
bashio::log.info "Generating [agent] and [[outputs]] sections..."

# Start writing the base configuration file
if bashio::config.has_value agent; then
    AGENT_CONFIG=$(bashio::config agent)
    bashio::log.info "Appending custom Agent Config..."
    cat >> "${CONFIG_FILE}" << EOF

# ===============================
#  Telegraf Configuration File
#  (auto-generated from Add-on)
# ===============================

# Agent Configuration (From HA UI)
${AGENT_CONFIG}
EOF
fi

# Outputs (Example: InfluxDB or PostgreSQL)
if bashio::config.has_value 'output_plugins'; then
    OUTPUT_PLUGINS=$(bashio::config 'output_plugins')
    bashio::log.info "Appending custom Output Plugins..."
    cat >> "${CONFIG_FILE}" << EOF

# Output Plugins (From HA UI)
${OUTPUT_PLUGINS}
EOF
fi

# --- 3. APPEND COMPLEX PLUGIN CONFIGS ---
# This is where we handle the raw input for plugins, assuming the user pasted valid TOML/YAML.

# Input Plugins
if bashio::config.has_value 'input_plugins'; then
    INPUT_PLUGINS=$(bashio::config 'input_plugins')
    bashio::log.info "Appending custom Input Plugins..."
    cat >> "${CONFIG_FILE}" << EOF

# Input Plugins (From HA UI)
${INPUT_PLUGINS}
EOF
fi

# Processor Plugins
if bashio::config.has_value 'processor_plugins'; then
    PROCESSOR_PLUGINS=$(bashio::config 'processor_plugins')
    bashio::log.info "Appending custom Processor Plugins..."
    cat >> "${CONFIG_FILE}" << EOF

# Processor Plugins (From HA UI)
${PROCESSOR_PLUGINS}
EOF
fi


# # --- 4. APPEND CUSTOM PLUGIN CONFIG FROM VOLUME ---

# PLUGIN_CONFIG_FILE=$(bashio::config 'plugin_config_file')
# PLUGIN_CONFIG_PATH="/share/${PLUGIN_CONFIG_FILE}"

# bashio::log.info "Checking for custom plugins file at: ${PLUGIN_CONFIG_PATH}"

# if bashio::fs.file_exists "${PLUGIN_CONFIG_PATH}"; then
#     bashio::log.info "Custom plugin configuration found. Appending to main config."
#     # Add a separator and append the contents of the user-provided file
#     cat >> "${CONFIG_FILE}" << EOF

# # ----------------------------------
# # Custom Plugins (Inputs/Processors)
# # Loaded from /share/${PLUGIN_CONFIG_FILE}
# # ----------------------------------

# EOF
#     cat "${PLUGIN_CONFIG_PATH}" >> "${CONFIG_FILE}"
# else
#     bashio::log.warn "No custom plugin file found at ${PLUGIN_CONFIG_PATH}. Running with only UI configuration."
# fi

# --- 4. START TELEGRAF ---
bashio::log.info "Generated Telegraf config saved to ${CONFIG_FILE}. Starting Telegraf..."
exec telegraf --config "${CONFIG_FILE}"
