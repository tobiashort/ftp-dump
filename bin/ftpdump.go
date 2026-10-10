package main

import (
	"net"

	"github.com/tobiashort/ftp-dump/command"
	"github.com/tobiashort/ftp-dump/response"
	"github.com/tobiashort/ftp-dump/session"
	"github.com/tobiashort/th-utils/lib/assert"

	"github.com/tobiashort/ftp-dump/config"
	"github.com/tobiashort/ftp-dump/logger"
)

func main() {
	assert.True(config.Port != "", "port not set")
	assert.True(config.AllowedUsername != "", "username not set")
	assert.True(config.AllowedPassword != "", "password not set")
	assert.True(config.Dumpdir != "", "dumpdir not set")

	addr := "0.0.0.0:" + config.Port
	cmdbus, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatal(err)
	}
	defer cmdbus.Close()
	logger.Info("FTP Dump server listening on " + addr)
	for {
		conn, err := cmdbus.Accept()
		if err != nil {
			logger.Error(err)
		} else {
			go handleNewConnection(conn)
		}
	}
}

func handleNewConnection(conn net.Conn) {
	session := session.New(conn)
	err := response.NewServiceReadyResponse().Send(session)
	if err != nil {
		logger.Error(err)
		session.Close()
		return
	}
	listener := command.NewCommandListener(session)
	listener.Listen()
}
