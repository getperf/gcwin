package agent

import (
  "github.com/getperf/gcagent/exporter"
  _ "github.com/getperf/gcagent/exporter/all"
)

var sampleProjectConfig = `
# gcagent.toml for inventory collector on Windows, UNIX Server
#
# Format
# Write in the form of a parameter name = value ("string").
# At the beginning of a sentence '#' is a comment line. 
# Parameters in comment lines show default values
# --------- Log usage monitoring setting -------------------------------
# LOG save time [h]
# save_hour = 72

# Log retransmission time [h]
# recovery_hour = 3

# --------- Log output setting -----------------------------------------
# Log level. fatal, err, warn, info, dbg
# log_level = "warn"

# Log rotation generation
# log_rotation = 5

# --------- HTTP proxy, timeout setting --------------------------------
# The presence or absence of proxy settings
# proxy_enable = false

# Proxy server address (in the case of the blank is applied HTTP_PROXY environment variable)
# proxy_url = ""

# HTTP connection time-out [s]
# http_timeout = 300

# HTTP retry count
# http_retry = 3

# --------- Web service connection settings --------------------------------
# Admin management web service url.
# service_url = "https://192.168.10.10:59443/"

# Site key.
# site_key = "site1"

#---------- Exporter job schedule, Each platform collector setting.
[jobs]

  # [jobs.sample]
  # # Collecting enable (true or false)
  # enable = true
  # 
  # #Interval sec (> 300)
  # interval = 300
  # 
  # #Timeout sec
  # timeout = 340
`

func (p *Project) SampleConfig() string {
  config := sampleProjectConfig
  for _, exporter := range exporter.Exporters {
    config += exporter().SampleScheduleConfig()
  }
  return config
}
