package appenders

import (
	"bufio"
	"bytes"
	"os"
	"github.com/blazecrystal/beyondts-go/utils"
)

type StdoutAppender struct {
	name, layout string
	log          chan string
}

func NewStdoutAppender(attrs map[string]string) *StdoutAppender {
	return &StdoutAppender{name: attrs["name"], layout: attrs["layout"], log: make(chan string)}
}

func (a *StdoutAppender) GetType() string {
	return Type_Stdout
}

//%t% %p% [%ln%] %sl% %c%
func (a *StdoutAppender) WriteLog(fills map[string]string) {
	a.log <- fillLayout(a.layout, fills)
}

func (a *StdoutAppender) Stop() {
	a.log <- string(CMD_END)
}

func (a *StdoutAppender) Run(flag chan string) {
	writer := bufio.NewWriter(os.Stdout)
	flag <- utils.Concat("+", a.name)
	for {
		tmp := <-a.log
		if bytes.Equal([]byte(tmp), CMD_END) {
			break
		}
		writer.WriteString(tmp)
		writer.WriteRune('\n')
		writer.Flush()
	}
	flag <- utils.Concat("-", a.name)
}

func (a *StdoutAppender) Update(attrs map[string]string) {
	a.name, a.layout = attrs["name"], attrs["layout"]
}
