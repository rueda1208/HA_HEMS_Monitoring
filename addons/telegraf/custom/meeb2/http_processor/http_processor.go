package http_processor

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/metric"
	"github.com/influxdata/telegraf/plugins/processors"
)

func toFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f, nil
		}
		return 0, fmt.Errorf("could not convert the string '%s' to float64", v)
	default:
		return 0, fmt.Errorf("unknown type: %T", v)
	}
}

type HTTProcessor struct{}

// SampleConfig provides the default configuration options
func (p *HTTProcessor) SampleConfig() string {
	return `
    ## Processor to modify HTTP data structure
    `
}

// Apply is where we modify the metric structure
func (p *HTTProcessor) Apply(metrics ...telegraf.Metric) []telegraf.Metric {
	var tags = map[string]string{
		"_type": "measure",
	}

	var mapping = map[string]string{
		"switch.water_heater": "water_heater",
		"climate.dining_room": "space_heating",
		"climate.living_room": "space_heating",
		"climate.powder_room": "space_heating",
		"climate.basement_1":  "space_heating",
		"climate.basement_2":  "space_heating",
		"climate.bathroom":    "space_heating",
		"climate.bedroom_1":   "space_heating",
		"climate.bedroom_2":   "space_heating",
		"climate.bedroom_3":   "space_heating",
		"climate.kitchen":     "space_heating",
		"climate.garage":      "space_heating",
	}

	var mapping_attributes = map[string]string{
		"attributes_local_temperature":         "temperature",
		"attributes_occupied_heating_setpoint": "setpoint",
	}

	// Modify data structure
	for _, metric_input := range metrics {
		// Get http data
		if metric_input.Name() == "home_assistant" {
			if tag_value, ok := metric_input.GetTag("entity_id"); ok {
				// Water heater
				if tag_value == "switch.water_heater" {
					new_netric := metric.New(mapping[tag_value], tags, nil, time.Now())

					if value, ok := metric_input.GetField("attributes_temperature"); ok {
						if _field_value, err := toFloat64(value); err == nil {
							new_netric.AddField("temperature", _field_value)
						} else {
							log.Printf("Water heater - Error: %v\n", err)
						}
					} else {
						log.Printf("Water heater missing")
					}
					return []telegraf.Metric{
						new_netric,
					}
				} else if strings.HasPrefix(tag_value, "climate.") {
					new_netric := metric.New(mapping[tag_value], tags, nil, time.Now())
					thermostat := strings.Split(tag_value, ".")[1]

					filds_list := []string{"attributes_local_temperature", "attributes_occupied_heating_setpoint"}
					for _, field_name := range filds_list {
						if value, ok := metric_input.GetField(field_name); ok {
							if _field_value, err := toFloat64(value); err == nil {
								new_netric.AddField(mapping_attributes[field_name]+"_"+thermostat, _field_value)
							} else {
								log.Printf("%s - Error: %v\n", field_name, err)
							}
						} else {
							log.Printf("%s missing", field_name)
						}
					}
					return []telegraf.Metric{
						new_netric,
					}
				} else {
					return nil
				}
			} else {
				return nil
			}
		} else if metric_input.Name() == "modbus" {
			return []telegraf.Metric{
				metric_input,
			}
		}
	}
	return nil
}

func init() {
	processors.Add("http_processor", func() telegraf.Processor {
		return &HTTProcessor{}
	})
}
