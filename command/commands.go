package command

import (
	"fmt"
	"io"
	"net"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/tobiashort/ftp-dump/config"
	"github.com/tobiashort/ftp-dump/logger"
	"github.com/tobiashort/ftp-dump/response"
	"github.com/tobiashort/ftp-dump/session"
)

type CmdCode string

const (
	ABOR CmdCode = "ABOR"
	ACCT CmdCode = "ACCT"
	ADAT CmdCode = "ADAT"
	ALLO CmdCode = "ALLO"
	APPE CmdCode = "APPE"
	AUTH CmdCode = "AUTH"
	AVBL CmdCode = "AVBL"
	CCC  CmdCode = "CCC"
	CDUP CmdCode = "CDUP"
	CONF CmdCode = "CONF"
	CSID CmdCode = "CSID"
	CWD  CmdCode = "CWD"
	DELE CmdCode = "DELE"
	DSIZ CmdCode = "DSIZ"
	ENC  CmdCode = "ENC"
	EPRT CmdCode = "EPRT"
	EPSV CmdCode = "EPSV"
	FEAT CmdCode = "FEAT"
	HELP CmdCode = "HELP"
	HOST CmdCode = "HOST"
	LANG CmdCode = "LANG"
	LIST CmdCode = "LIST"
	LPRT CmdCode = "LPRT"
	LPSV CmdCode = "LPSV"
	MDTM CmdCode = "MDTM"
	MFCT CmdCode = "MFCT"
	MFF  CmdCode = "MFF"
	MFMT CmdCode = "MFMT"
	MIC  CmdCode = "MIC"
	MKD  CmdCode = "MKD"
	MLSD CmdCode = "MLSD"
	MLST CmdCode = "MLST"
	MODE CmdCode = "MODE"
	NLST CmdCode = "NLST"
	NOOP CmdCode = "NOOP"
	OPTS CmdCode = "OPTS"
	PASS CmdCode = "PASS"
	PASV CmdCode = "PASV"
	PBSZ CmdCode = "PBSZ"
	PORT CmdCode = "PORT"
	PROT CmdCode = "PROT"
	PWD  CmdCode = "PWD"
	QUIT CmdCode = "QUIT"
	REIN CmdCode = "REIN"
	REST CmdCode = "REST"
	RETR CmdCode = "RETR"
	RMD  CmdCode = "RMD"
	RMDA CmdCode = "RMDA"
	RNFR CmdCode = "RNFR"
	RNTO CmdCode = "RNTO"
	SITE CmdCode = "SITE"
	SIZE CmdCode = "SIZE"
	SMNT CmdCode = "SMNT"
	SPSV CmdCode = "SPSV"
	STAT CmdCode = "STAT"
	STOR CmdCode = "STOR"
	STOU CmdCode = "STOU"
	STRU CmdCode = "STRU"
	SYST CmdCode = "SYST"
	THMB CmdCode = "THMB"
	TYPE CmdCode = "TYPE"
	USER CmdCode = "USER"
	XCUP CmdCode = "XCUP"
	XMKD CmdCode = "XMKD"
	XPWD CmdCode = "XPWD"
	XRCP CmdCode = "XRCP"
	XRMD CmdCode = "XRMD"
	XRSQ CmdCode = "XRSQ"
	XSEM CmdCode = "XSEM"
	XSEN CmdCode = "XSEN"
)

type Command interface {
	Exec(*session.Session, string)
}

type cwdCommand struct{}

func (cmd cwdCommand) Exec(session *session.Session, parameter string) {
	response.NewOkResponse("CWD Ignored.").Send(session)
}

type eprtCommand struct{}

func (cmd eprtCommand) Exec(session *session.Session, parameter string) {
	dataConnParts := strings.Split(parameter, "|")
	session.SetClientHost(dataConnParts[2])
	session.SetClientPort(dataConnParts[3])
	response.NewOkResponse("EPRT ok.").Send(session)
}

type passCommand struct{}

func (cmd passCommand) Exec(session *session.Session, parameter string) {
	session.SetPassword(parameter)
	if session.IsAuthenticated() {
		response.NewUserLoggedInResponse().Send(session)
		return
	}
	response.NewInvalidUsernameOrPasswordResponse().Send(session)
}

type portCommand struct{}

func (cmd portCommand) Exec(session *session.Session, parameter string) {
	dataConnParts := strings.Split(parameter, ",")
	dataConnHost := fmt.Sprintf("%s.%s.%s.%s", dataConnParts[0], dataConnParts[1], dataConnParts[2], dataConnParts[3])
	dataConnPortPartOne, err := strconv.ParseInt(dataConnParts[4], 10, 32)
	dataConnPortPartTwo, err := strconv.ParseInt(dataConnParts[5], 10, 32)
	if err != nil {
		logger.Error(err)
		response.NewSyntaxErrorResponse(err.Error()).Send(session)
		return
	}
	dataConnPort := fmt.Sprintf("%d", dataConnPortPartOne*256+dataConnPortPartTwo)
	session.SetClientHost(dataConnHost)
	session.SetClientPort(dataConnPort)
	response.NewOkResponse("PORT ok.").Send(session)
}

type quitCommand struct{}

func (cmd quitCommand) Exec(session *session.Session, parameter string) {
	response.NewServiceClosesControlConnectionResponse().Send(session)
}

type storCommand struct{}

func (cmd storCommand) Exec(session *session.Session, fileName string) {
	dataConn, err := net.Dial("tcp", "["+session.GetClientHost()+"]:"+session.GetClientPort())
	if err != nil {
		logger.Error(err)
		response.NewCannotOpenDataConnectionResponse().Send(session)
		return
	}
	defer dataConn.Close()
	tmpFilePath := path.Join(config.Dumpdir, fileName+".part")
	tmpFile, err := os.Create(tmpFilePath)
	if err != nil {
		logger.Error(err)
		response.NewSyntaxErrorResponse(err.Error()).Send(session)
	}
	defer tmpFile.Close()
	response.NewDataConnectionAlreadyOpenResponse().Send(session)
	_, err = io.Copy(tmpFile, dataConn)
	if err != nil {
		logger.Error(err)
		response.NewActionAbortedErrorInProcessingResponse().Send(session)
		return
	}
	filePath := path.Join(config.Dumpdir, fileName)
	err = os.Rename(tmpFilePath, filePath)
	if err != nil {
		logger.Error(err)
		response.NewActionAbortedErrorInProcessingResponse().Send(session)
		return
	}
	response.NewClosingDataConnectionFileActionSuccessfulResponse().Send(session)
}

type systCommand struct{}

func (cmd systCommand) Exec(session *session.Session, parameter string) {
	response.NewSystemTypeResponse("UNIX Type: L8").Send(session)
}

type typeCommand struct{}

func (cmd typeCommand) Exec(session *session.Session, parameter string) {
	if parameter == "I" {
		// binary
		response.NewOkResponse("TYPE ok.").Send(session)
		return
	}
	response.NewSyntaxErrorResponse(fmt.Sprintf("Unknown TYPE: %s", parameter)).Send(session)
}

type userCommand struct{}

func (cmd userCommand) Exec(session *session.Session, parameter string) {
	session.SetUsername(parameter)
	response.NewUserNameOkayNeedPasswordResponse().Send(session)
}
