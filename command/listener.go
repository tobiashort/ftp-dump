package command

import (
	"io"
	"strings"

	"github.com/tobiashort/ftp-dump/logger"
	"github.com/tobiashort/ftp-dump/response"
	"github.com/tobiashort/ftp-dump/session"
)

var cmdReg = NewCommandRegistry()

type CommandListener interface {
	Listen()
}

type commandListener struct {
	session *session.Session
}

func (listener commandListener) Listen() {
	buf := make([]byte, 1024)
	for {
		bytesRead, err := listener.session.Read(buf)
		if err != nil {
			if err != io.EOF {
				logger.Error(err)
			}
			listener.session.Close()
			return
		}
		rawCmd := string(buf[:bytesRead])
		rawCmd = strings.TrimSpace(rawCmd)
		if strings.HasPrefix(rawCmd, "PASS ") {
			logger.Info("PASS ****")
		} else {
			logger.Info(rawCmd)
		}
		split := strings.Split(rawCmd, " ")
		code := split[0]
		parameter := strings.Join(split[1:], " ")
		if cmd, ok := cmdReg.Resolve(CmdCode(code)); ok {
			cmd.Exec(listener.session, parameter)
		} else {
			resp := response.NewCommandNotImplementedResponse()
			resp.Send(listener.session)
		}
	}
}

func NewCommandListener(session *session.Session) CommandListener {
	return commandListener{session}
}
