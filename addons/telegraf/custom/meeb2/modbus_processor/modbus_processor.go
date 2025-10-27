package modbus_processor

import (
	"fmt"
	"log"
	"strconv"
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

type ModbusProcessor struct{}

// SampleConfig provides the default configuration options
func (p *ModbusProcessor) SampleConfig() string {
	return `
    ## Processor to modify Modbus data structure
    `
}

// Apply is where we modify the metric structure
func (p *ModbusProcessor) Apply(metrics ...telegraf.Metric) []telegraf.Metric {
	var tags = map[string]string{
		"_type": "measure",
	}

	energy_meter_Metric := metric.New("energy_meter", tags, nil, time.Now())
	water_heater_Metric := metric.New("water_heater", tags, nil, time.Now())
	space_heating_Metric := metric.New("space_heating", tags, nil, time.Now())
	electric_vehicle_Metric := metric.New("electric_vehicle_v1g", tags, nil, time.Now())
	electric_battery_Metric := metric.New("electric_battery", tags, nil, time.Now())

	// var total_power float64
	// var water_heater_power float64
	// var electric_vehicle_power float64

	// var space_heating_power = map[string]float64{
	// 	"Chauffage SDB 2":          0.0,
	// 	"Chauffage Garage":         0.0,
	// 	"Chauffage Sous-sol":       0.0,
	// 	"Chauffage chambre 2":      0.0,
	// 	"Chauffage chambre 3":      0.0,
	// 	"Chauffage hall-cuisine":   0.0,
	// 	"Chauffage Salle a manger": 0.0,
	// }

	// var electric_battery_XW503 = map[string]float64{
	// 	"AC1 Power":     0.0,
	// 	"AC Load Power": 0.0,
	// }
	// var electric_battery_GW503 = map[string]float64{
	// 	"MPPT PV Power":          0.0,
	// 	"Battery Bank1 SOC":      0.0,
	// 	"Battery Bank1 Power":    0.0,
	// 	"Battery Bank 1 Voltage": 0.0,
	// }
	// var electric_battery_MPPT503 = map[string]float64{
	// 	"PV Power":          0.0,
	// 	"PV Voltage":        0.0,
	// 	"DC Output Power":   0.0,
	// 	"DC Output Voltage": 0.0,
	// }

	// var total_power_exists bool
	// var water_heater_exists bool
	// var space_heating_exists bool
	// var electric_vehicle_exists bool
	// var electric_battery_XW503_exists bool
	// var electric_battery_GW503_exists bool
	// var electric_battery_MPPT503_exists bool

	// Modify data structure
	for _, metric_input := range metrics {
		// Get modbus data
		if metric_input.Name() == "modbus" {
			if tag_value, ok := metric_input.GetTag("name"); ok {
				if tag_value == "eGauge_PP" {
					// Total power
					if value, ok := metric_input.GetField("P_Total de la maison(MEEB1/PSR)"); ok {
						if _total_power, err := toFloat64(value); err == nil {
							energy_meter_Metric.AddField("total_power", _total_power)
						} else {
							log.Printf("P_Total de la maison(MEEB1/PSR) - Error: %v\n", err)
						}
					} else {
						log.Printf("P_Total de la maison(MEEB1/PSR) missing")
					}

					// Water heater
					if value, ok := metric_input.GetField("Chauffe-eau"); ok {
						if _water_heater_power, err := toFloat64(value); err == nil {
							water_heater_Metric.AddField("power", _water_heater_power)
							energy_meter_Metric.AddField("water_heater_power", _water_heater_power)
						} else {
							log.Printf("Chauffe-eau - Error: %v\n", err)
						}
					} else {
						log.Printf("Chauffe-eau missing")
					}

					// Space heating
					space_heating_mapping := map[string]string{
						"Chauffage SDB 2":                "bathroom",
						"Chauffage Garage":               "garage",
						"Chauffage chambre 2":            "bedroom_2",
						"Chauffage chambre 3":            "bedroom_3",
						"Chauffage Sous-sol sud":         "basement_1",
						"Chauffage sous-sol nord":        "basement_2",
						"Chauffage Salle a manger":       "dining_room",
						"Chauffage hall-SDB rdc-cuisine": "kitchen",
					}

					i := 0
					space_heating_total_power := 0.0
					for variable, var_db_name := range space_heating_mapping {
						if value, ok := metric_input.GetField(variable); ok {
							if _space_heating_power, err := toFloat64(value); err == nil {
								space_heating_name_db := "power_" + var_db_name
								space_heating_total_power += _space_heating_power
								space_heating_Metric.AddField(space_heating_name_db, _space_heating_power)
								i++
							} else {
								log.Printf("%s - Error: %v\n", variable, err)
							}
						} else {
							log.Printf("%s missing", variable)
						}
					}

					if i == len(space_heating_mapping) {
						energy_meter_Metric.AddField("space_heating_pp_power", space_heating_total_power)
					}

					// Electric vehicle
					if value, ok := metric_input.GetField("Borne véhicule"); ok {
						if _electric_vehicle_power, err := toFloat64(value); err == nil {
							electric_vehicle_Metric.AddField("power", _electric_vehicle_power)
							energy_meter_Metric.AddField("electric_vehicle_power", _electric_vehicle_power)
						} else {
							log.Printf("Borne véhicule - Error: %v\n", err)
						}
					} else {
						log.Printf("Borne véhicule missing")
					}

					// Return the new combined metric
					return []telegraf.Metric{
						energy_meter_Metric,
						water_heater_Metric,
						space_heating_Metric,
						electric_vehicle_Metric,
						electric_battery_Metric,
					}
				} else if tag_value == "eGauge_PU" {
					// Space heating
					space_heating_mapping := map[string]string{
						"Chauffage salon 1":   "living_room",
						"Chauffage chambre 1": "bedroom_1",
					}

					i := 0
					space_heating_total_power := 0.0
					for variable, var_db_name := range space_heating_mapping {
						if value, ok := metric_input.GetField(variable); ok {
							if _space_heating_power, err := toFloat64(value); err == nil {
								space_heating_name_db := "power_" + var_db_name
								space_heating_total_power += _space_heating_power
								space_heating_Metric.AddField(space_heating_name_db, _space_heating_power)
								i++
							} else {
								log.Printf("%s - Error: %v\n", variable, err)
							}
						} else {
							log.Printf("%s missing", variable)
						}
					}

					if i == len(space_heating_mapping) {
						energy_meter_Metric.AddField("space_heating_pu_power", space_heating_total_power)
					}
					// Return the new combined metric
					return []telegraf.Metric{
						energy_meter_Metric,
						space_heating_Metric,
					}
				} else if tag_value == "InsightHome_XW503" {
					electric_battery_XW503 := map[string]string{
						"AC1 Power":     "ac_grid_power",
						"AC Load Power": "ac_load_power",
					}

					i := 0
					for variable, var_db_name := range electric_battery_XW503 {
						if value, ok := metric_input.GetField(variable); ok {
							if _electric_battery, err := toFloat64(value); err == nil {
								electric_battery_Metric.AddField(var_db_name, _electric_battery)
								i++
							} else {
								log.Printf("%s - Error: %v\n", variable, err)
							}
						} else {
							log.Printf("%s missing", variable)
						}
					}
					// Return the new combined metric
					return []telegraf.Metric{
						electric_battery_Metric,
					}
				} else if tag_value == "InsightHome_GW503" {
					electric_battery_GW503 := map[string]string{
						"MPPT PV Power":          "mppt_pv_power",
						"Battery Bank1 SOC":      "state_of_charge",
						"Battery Bank1 Power":    "battery_power",
						"Battery Bank 1 Voltage": "battery_voltage",
					}

					i := 0
					for variable, var_db_name := range electric_battery_GW503 {
						if value, ok := metric_input.GetField(variable); ok {
							if _electric_battery, err := toFloat64(value); err == nil {
								electric_battery_Metric.AddField(var_db_name, _electric_battery)
								i++
							} else {
								log.Printf("%s - Error: %v\n", variable, err)
							}
						} else {
							log.Printf("%s missing", variable)
						}
					}
					// Return the new combined metric
					return []telegraf.Metric{
						electric_battery_Metric,
					}
				} else if tag_value == "InsightHome_MPPT503" {
					electric_battery_MPPT503 := map[string]string{
						"PV Power":          "pv_power",
						"PV Voltage":        "pv_voltage",
						"DC Output Power":   "dc_output_power",
						"DC Output Voltage": "dc_output_voltage",
					}

					i := 0
					for variable, var_db_name := range electric_battery_MPPT503 {
						if value, ok := metric_input.GetField(variable); ok {
							if _electric_battery, err := toFloat64(value); err == nil {
								electric_battery_Metric.AddField(var_db_name, _electric_battery)
								i++
							} else {
								log.Printf("%s - Error: %v\n", variable, err)
							}
						} else {
							log.Printf("%s missing", variable)
						}
					}
					// Return the new combined metric
					return []telegraf.Metric{
						electric_battery_Metric,
					}
				} else {
					return nil
				}
			} else {
				return nil
			}

		} else if metric_input.Name() == "water_heater" || metric_input.Name() == "space_heating" {
			return []telegraf.Metric{
				metric_input,
			}
		} else {
			return nil
		}
	}
	return nil
	// // Return the new combined metric
	// return []telegraf.Metric{
	// 	energy_meter_Metric,
	// 	water_heater_Metric,
	// 	space_heating_Metric,
	// 	electric_vehicle_Metric,
	// 	electric_battery_Metric,
	// }
}

func init() {
	processors.Add("modbus_processor", func() telegraf.Processor {
		return &ModbusProcessor{}
	})
}
