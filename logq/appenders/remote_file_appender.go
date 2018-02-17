package appenders

import (
	"bytes"
	"github.com/blazecrystal/beyondts-go/utils"
)

type RemoteFileAppender struct {
	name, layout, serverIp, serverPort, user, pwd, file, dailyRolling, threshold, zip, keep string
	log                                                                                     chan string
}

func NewRemoteFileAppender(attrs map[string]string) *RemoteFileAppender {
	return &RemoteFileAppender{name: attrs["name"], layout: attrs["layout"], file: attrs["file"], dailyRolling: attrs["dailyRolling"], threshold: attrs["threshold"], zip: attrs["zip"], keep: attrs["keep"], serverIp: attrs["serverIp"], serverPort: attrs["serverPort"], user: attrs["user"], pwd: attrs["pwd"], log: make(chan string)}
}

func (a *RemoteFileAppender) GetType() string {
	return Type_RemoteFile
}

func (a *RemoteFileAppender) WriteLog(fills map[string]string) {
	a.log <- fillLayout(a.layout, fills)
}

func (a *RemoteFileAppender) Stop() {
	a.log <- string(CMD_END)
}

func (a *RemoteFileAppender) Run(flag chan string) {
	// TODO open remote file & build writer

	flag <- utils.Concat("+", a.name)

	for {
		tmp := <-a.log
		if bytes.Equal([]byte(tmp), CMD_END) {
			break
		}
		// TODO write tmp to remote file
	}
	// TODO close file & connection
	flag <- utils.Concat("-", a.name)
}

func (a *RemoteFileAppender) Update(attrs map[string]string) {
	a.name, a.layout, a.file, a.dailyRolling, a.threshold, a.zip, a.keep, a.serverIp, a.serverPort, a.user, a.pwd = attrs["name"], attrs["layout"], attrs["file"], attrs["dailyRolling"], attrs["threshold"], attrs["zip"], attrs["keep"], attrs["serverIp"], attrs["serverPort"], attrs["user"], attrs["pwd"]
}
