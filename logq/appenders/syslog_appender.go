package appenders

import (
	"bufio"
	"bytes"
	"net"
	"github.com/blazecrystal/beyondts-go/utils"
)

const (
	NET_UDP = "udp"
)

type SyslogAppender struct {
	name, layout, serverIp, serverPort string
	log                                chan string
}

func NewSyslogAppender(attrs map[string]string) *SyslogAppender {
	return &SyslogAppender{name: attrs["name"], layout: attrs["layout"], serverIp: attrs["serverIp"], serverPort: attrs["serverPort"], log: make(chan string)}
}

func (a *SyslogAppender) GetType() string {
	return Type_Syslog
}

func (a *SyslogAppender) WriteLog(fills map[string]string) {
	a.log <- fillLayout(a.layout, fills)
}

func (a *SyslogAppender) Stop() {
	a.log <- string(CMD_END)
}

func (a *SyslogAppender) Run(flag chan string) {
	servAddr, err := net.ResolveUDPAddr(NET_UDP, utils.Concat(a.serverIp, ":", a.serverPort))
	if err != nil {
		debugMsg("can't resolve udp address, ", a.serverIp, ":", a.serverPort, ", this appender will be skipped\n", err)
		return
	}
	con, err := net.DialUDP(NET_UDP, nil, servAddr)
	defer con.Close()
	if err != nil {
		debugMsg("can't create udp connection to ", a.serverIp, ":", a.serverPort, ", this appender will be skipped\n", err)
		return
	}
	writer := bufio.NewWriter(con)
	flag <- utils.Concat("+", a.name)
	for {
		tmp := <-a.log
		if bytes.Equal([]byte(tmp), CMD_END) {
			break
		}
		writer.WriteString(tmp)
		writer.Flush()
	}
	flag <- utils.Concat("-", a.name)

}

func (a *SyslogAppender) Update(attrs map[string]string) {
	a.name, a.layout, a.serverIp, a.serverPort = attrs["name"], attrs["layout"], attrs["serverIp"], attrs["serverPort"]
}
