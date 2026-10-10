package response

import (
	"fmt"

	"github.com/tobiashort/ftp-dump/session"
)

type RetCode int

type Response interface {
	Send(*session.Session) error
}

type response struct {
	session session.Session
	code    RetCode
	message string
}

func (r response) Send(session *session.Session) error {
	var msg string
	if r.message != "" {
		msg = fmt.Sprintf("%d %s\n", r.code, r.message)
	} else {
		msg = fmt.Sprintf("%d\n", r.code)
	}
	return session.Write([]byte(msg))
}

func NewDataConnectionAlreadyOpenResponse() Response {
	return response{
		code:    125,
		message: "Data connection already open; transfer starting.",
	}
}

func NewOkResponse(message string) Response {
	return response{
		code:    200,
		message: message,
	}
}

func NewSystemTypeResponse(systype string) Response {
	return response{
		code:    215,
		message: systype,
	}
}

func NewServiceReadyResponse() Response {
	return response{
		code:    220,
		message: "Service ready for new user.",
	}
}

func NewServiceClosesControlConnectionResponse() Response {
	return response{
		code:    221,
		message: "Service closing control connection.",
	}
}

func NewClosingDataConnectionFileActionSuccessfulResponse() Response {
	return response{
		code:    226,
		message: "Closing data connection. Requested file action successful.",
	}
}

func NewUserLoggedInResponse() Response {
	return response{
		code:    230,
		message: "User logged in, proceed. Logged out if appropriate.",
	}
}

func NewUserNameOkayNeedPasswordResponse() Response {
	return response{
		code:    331,
		message: "User name okay, need password.",
	}
}

func NewCannotOpenDataConnectionResponse() Response {
	return response{
		code:    425,
		message: "Can't open data connection.",
	}
}

func NewInvalidUsernameOrPasswordResponse() Response {
	return response{
		code:    430,
		message: "Invalid username or password.",
	}
}

func NewActionAbortedErrorInProcessingResponse() Response {
	return response{
		code:    451,
		message: "Requested action aborted. Local error in processing.",
	}
}

func NewSyntaxErrorResponse(message string) Response {
	return response{
		code:    500,
		message: message,
	}
}

func NewCommandNotImplementedResponse() Response {
	return response{
		code:    502,
		message: "Command not implemented.",
	}
}
