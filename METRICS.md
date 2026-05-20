# Starlink Exporter - Exported Metrics Reference

## Metric Overview

All metrics are prefixed with `starlink_dish_` and use Prometheus gauge type unless otherwise specified.

## Connection & Status Metrics

### starlink_dish_up
- **Type**: Gauge (0 or 1)
- **Description**: Indicates successful communication with Starlink dish
- **Labels**: None
- **Example**: `starlink_dish_up 1`

### starlink_dish_state
- **Type**: Gauge
- **Description**: Current state of the dish
- **Values**: 
  - 0: Unknown
  - 1: Booting
  - 2: Searching
  - 3: Connected
- **Example**: `starlink_dish_state 3`

### starlink_dish_scrape_duration_seconds
- **Type**: Gauge
- **Description**: Time taken to scrape metrics from the dish
- **Unit**: Seconds
- **Example**: `starlink_dish_scrape_duration_seconds 0.125`

## Network Performance Metrics

### starlink_dish_pop_ping_latency_seconds
- **Type**: Gauge
- **Description**: Latency to Point of Presence (PoP)
- **Unit**: Seconds
- **Range**: 0.02 - 0.12 seconds (20-120 ms typical)
- **Example**: `starlink_dish_pop_ping_latency_seconds 0.045`

### starlink_dish_pop_ping_drop_ratio
- **Type**: Gauge
- **Description**: Ratio of dropped ping packets
- **Range**: 0.0 - 1.0 (0-100%)
- **Range**: 0 - 0.05 (0-5% typical)
- **Example**: `starlink_dish_pop_ping_drop_ratio 0.002`

### starlink_dish_snr
- **Type**: Gauge
- **Description**: Signal-to-Noise Ratio
- **Unit**: dB (Decibels)
- **Range**: 5 - 15 dB typical
- **Example**: `starlink_dish_snr 9.5`

## Throughput Metrics

### starlink_dish_downlink_throughput_bytes
- **Type**: Gauge
- **Description**: Download throughput
- **Unit**: Bytes per second
- **Typical Range**: 6.25M - 31.25M bytes/sec (50-250 Mbps)
- **Example**: `starlink_dish_downlink_throughput_bytes 15625000` (125 Mbps)

### starlink_dish_uplink_throughput_bytes
- **Type**: Gauge
- **Description**: Upload throughput
- **Unit**: Bytes per second
- **Typical Range**: 625K - 5M bytes/sec (5-40 Mbps)
- **Example**: `starlink_dish_uplink_throughput_bytes 1250000` (10 Mbps)

## Obstruction Metrics

### starlink_dish_currently_obstructed
- **Type**: Gauge (0 or 1)
- **Description**: Whether the dish is currently obstructed
- **Values**: 
  - 0: Clear view
  - 1: Currently obstructed
- **Example**: `starlink_dish_currently_obstructed 0`

### starlink_dish_fraction_obstruction_ratio
- **Type**: Gauge
- **Description**: Overall obstruction percentage
- **Range**: 0.0 - 1.0 (0-100%)
- **Typical Range**: 0 - 0.2 (0-20%)
- **Example**: `starlink_dish_fraction_obstruction_ratio 0.05` (5% obstructed)

### starlink_dish_last_24h_obstructed_seconds
- **Type**: Gauge
- **Description**: Total seconds obstructed in the last 24 hours
- **Unit**: Seconds
- **Example**: `starlink_dish_last_24h_obstructed_seconds 120`

### starlink_dish_wedge_fraction_obstruction_ratio
- **Type**: Gauge with labels
- **Description**: Obstruction ratio for each 30-degree wedge section
- **Labels**: 
  - `wedge`: Wedge number (0-11)
  - `wedge_name`: Compass direction range (e.g., "0_to_30")
- **Example**: 
  ```
  starlink_dish_wedge_fraction_obstruction_ratio{wedge="0",wedge_name="0_to_30"} 0.1
  starlink_dish_wedge_fraction_obstruction_ratio{wedge="1",wedge_name="30_to_60"} 0.05
  ```

### starlink_dish_wedge_abs_fraction_obstruction_ratio
- **Type**: Gauge with labels
- **Description**: Absolute obstruction fraction for each wedge
- **Labels**: Same as `wedge_fraction_obstruction_ratio`

### starlink_dish_valid_seconds
- **Type**: Gauge
- **Description**: Seconds since last obstruction data update
- **Unit**: Seconds

## Device Information

### starlink_dish_info
- **Type**: Gauge (always 1)
- **Description**: Device information and metadata
- **Labels**:
  - `device_id`: Unique dish identifier
  - `hardware_version`: Hardware revision (e.g., "rev2_proto_v2")
  - `software_version`: Current firmware version (e.g., "2024.07.01")
  - `country_code`: Device country code (e.g., "US")
  - `utc_offset`: UTC offset in seconds (e.g., "-28800")
- **Example**: 
  ```
  starlink_dish_info{country_code="US",device_id="00aaabbbcccddd",
  hardware_version="rev2_proto_v2",software_version="2024.07.01",
  utc_offset="-28800"} 1
  ```

### starlink_dish_uptime_seconds
- **Type**: Gauge
- **Description**: Device uptime since last boot
- **Unit**: Seconds
- **Example**: `starlink_dish_uptime_seconds 604800` (7 days)

## Positioning Metrics

### starlink_dish_bore_sight_azimuth_deg
- **Type**: Gauge
- **Description**: Dish azimuth (horizontal) pointing angle
- **Unit**: Degrees (0-360)
- **Example**: `starlink_dish_bore_sight_azimuth_deg 180.5`

### starlink_dish_bore_sight_elevation_deg
- **Type**: Gauge
- **Description**: Dish elevation (vertical) pointing angle
- **Unit**: Degrees (-90 to 90)
- **Example**: `starlink_dish_bore_sight_elevation_deg 65.3`

## Network Cell/Location Metrics

### starlink_dish_cell_id
- **Type**: Gauge
- **Description**: Cell ID the dish is connected to
- **Example**: `starlink_dish_cell_id 42501`

### starlink_dish_pop_rack_id
- **Type**: Gauge
- **Description**: Point of Presence rack identifier
- **Example**: `starlink_dish_pop_rack_id 15`

### starlink_dish_initial_satellite_id
- **Type**: Gauge
- **Description**: Initial satellite identifier for this connection session
- **Example**: `starlink_dish_initial_satellite_id 1001`

### starlink_dish_initial_gateway_id
- **Type**: Gauge
- **Description**: Initial gateway identifier for this connection session
- **Example**: `starlink_dish_initial_gateway_id 5`

### starlink_dish_backup_beam
- **Type**: Gauge (0 or 1)
- **Description**: Whether currently using backup beam
- **Example**: `starlink_dish_backup_beam 0`

## Slot/Timing Metrics

### starlink_dish_time_to_slot_end_seconds
- **Type**: Gauge
- **Description**: Seconds remaining in current time slot
- **Unit**: Seconds
- **Example**: `starlink_dish_time_to_slot_end_seconds 45`

### starlink_dish_first_nonempty_slot_seconds
- **Type**: Gauge
- **Description**: Seconds until next non-empty slot
- **Unit**: Seconds
- **Example**: `starlink_dish_first_nonempty_slot_seconds 2`

## Alert Metrics

All alert metrics are gauges with values 0 (inactive) or 1 (active).

### starlink_dish_alert_motors_stuck
- **Description**: Dish motor movement is stuck

### starlink_dish_alert_thermal_throttle
- **Description**: Dish is thermal throttling

### starlink_dish_alert_thermal_shutdown
- **Description**: Dish experienced thermal shutdown

### starlink_dish_alert_mast_not_near_vertical
- **Description**: Mounting mast is not near vertical

### starlink_dish_alert_unexpected_location
- **Description**: Dish location is unexpected

### starlink_dish_alert_slow_eth_speeds
- **Description**: Ethernet connection is slow

## Obstruction Duration Metrics

### starlink_dish_prolonged_obstruction_duration_seconds
- **Type**: Gauge
- **Description**: Average duration of prolonged obstruction events
- **Unit**: Seconds
- **Example**: `starlink_dish_prolonged_obstruction_duration_seconds 45`

### starlink_dish_prolonged_obstruction_interval_seconds
- **Type**: Gauge
- **Description**: Average interval between prolonged obstruction events
- **Unit**: Seconds
- **Example**: `starlink_dish_prolonged_obstruction_interval_seconds 3600`

## Metric Query Examples

### PromQL Queries

Get current download speed in Mbps:
```promql
rate(starlink_dish_downlink_throughput_bytes[1m]) / 125000
```

Get current upload speed in Mbps:
```promql
rate(starlink_dish_uplink_throughput_bytes[1m]) / 125000
```

Get average latency over last hour:
```promql
avg_over_time(starlink_dish_pop_ping_latency_seconds[1h])
```

Get maximum packet loss in last hour:
```promql
max_over_time(starlink_dish_pop_ping_drop_ratio[1h])
```

Check if dish is currently connected:
```promql
starlink_dish_up
```

Get current obstruction percentage:
```promql
starlink_dish_fraction_obstruction_ratio * 100
```

## Prometheus Scrape Configuration

```yaml
scrape_configs:
  - job_name: 'starlink'
    scrape_interval: 10s
    scrape_timeout: 5s
    static_configs:
      - targets: ['localhost:9817']
    metrics_path: '/metrics'
```

## Notes

- All timing-based metrics use seconds as the unit
- All throughput metrics use bytes/second
- Obstruction metrics are ratios (0.0-1.0) unless specified otherwise
- Alert metrics are binary (0 or 1)
- Metrics are updated on each scrape cycle
- Historical metrics are retained per Prometheus configuration
