package session

import (
	"net"

	"github.com/tobiashort/ftp-dump/config"
	"github.com/tobiashort/ftp-dump/logger"
)

var sessionId int = 0

type Session struct {
	ID         int
	conn       net.Conn
	username   string
	password   string
	clientHost string
	clientPort string
}

func New(conn net.Conn) *Session {
	sessionId = sessionId + 1
	logger.Info("New session", sessionId)
	return &Session{
		ID:   sessionId,
		conn: conn,
	}
}

func (s *Session) Close() {
	err := s.conn.Close()
	if err != nil {
		logger.Error(err)
	}
	logger.Info("Session", s.ID, "closed.")
}

func (s *Session) Read(buffer []byte) (int, error) {
	return s.conn.Read(buffer)
}

func (s *Session) Write(bytes []byte) error {
	_, err := s.conn.Write(bytes)
	return err
}

func (s *Session) IsAuthenticated() bool {
	if s.username == config.AllowedUsername && s.password == config.AllowedPassword {
		return true
	}
	return false
}

func (s *Session) GetUsername() string {
	return s.username
}

func (s *Session) SetUsername(username string) {
	s.username = username
}

func (s *Session) GetPassword() string {
	return s.password
}

func (s *Session) SetPassword(password string) {
	s.password = password
}

func (s *Session) GetClientHost() string {
	return s.clientHost
}

func (s *Session) SetClientHost(host string) {
	s.clientHost = host
}

func (s *Session) GetClientPort() string {
	return s.clientPort
}

func (s *Session) SetClientPort(port string) {
	s.clientPort = port
}
