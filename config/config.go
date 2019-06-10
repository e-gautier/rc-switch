package config

import "fmt"

const SysLogTag string = "RCSWITCH"
var SysLogTagSniffer = fmt.Sprintf("[%s] SNIFFER", SysLogTag)
var SysLogTagSender = fmt.Sprintf("[%s] SENDER", SysLogTag)