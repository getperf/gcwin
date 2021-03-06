package agent

import (
  "github.com/getperf/gcagent/exporter"
  _ "github.com/getperf/gcagent/exporter/all"
)

var sampleProjectConfig = `
# Param.ini for Performace monitoring system on Windows, UNIX Server
#
# Format
# Write in the form of a parameter name = value (string).
# At the beginning of a sentence; there is line is a comment line. (In the middle of a line there; it is recognized as a character.)
# --------- Log usage monitoring setting -------------------------------
# Threshold of the disk usage. If the specified value [%] following a disk full error.
disk_capacity = 0

# LOG save time [H]
save_hour = 24

# Log retransmission time [H]
recovery_hour = 3

# Number of output line of the error log
max_error_log = 5

# --------- Log output setting -----------------------------------------
# Log level. None 0, FATAL 1, CRIT 2, ERR 3, WARN 4, NOTICE 5, INFO 6, DBG 7
log_level = 7

# The standard output of the log (for debugging)
debug_console = false

# Log size (bytes)
# log_size = 100000

# Log rotation generation
# log_rotation = 5

# --------- HA state monitoring setting --------------------------------
# The presence or absence of state detection of node
# hanode_enable = false

# Setting of node script for the detection of state ('{HOME}/ptune/script' placed under)
# hanode_cmd = hastatus.sh

# --------- Post-processing setting ------------------------------------
# The presence or absence of post-processing
# post_enable = false

# Post-processing command
# post_cmd = scp _zip_ hogehoge@test.getperf.com: ~/work/tmp
# post_cmd = "C:\Program Files (x86)\WinSCP\winscp.exe" /script="C:\ptune\script\upload.script" / parameter "_zip_"

# --------- HTTP proxy, timeout setting --------------------------------
# The presence or absence of proxy settings
# proxy_enable = false

# Proxy server address (in the case of the blank is applied HTTP_PROXY environment variable)
# proxy_host = ""

# Proxy server port
# proxy_port = 8080

# HTTP connection time-out
# soap_timeout = 300

# --------- Web service connection settings --------------------------------
# WEB Service enable, If it is false, the agent doesn't send performance data.
# remhost_enable = false

# Admin management web service url.
# url_cm = "https://192.168.10.10:57443/axis2/services/GetperfService"

# Data management web service url.
# url_pm = "https://192.168.10.10:58443/axis2/services/GetperfService"

# Site key.
# site_key = "site1"

#---------- Account template setting.
[accounts]
#  [accounts.admin]
#  user = "admin"
#  password = "password"

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
  for _, exp := range exporter.Exporters {
    // config += exporter().SampleScheduleConfig()
    config += exp().Config(exporter.SCHEDULE)
  }
  return config
}
