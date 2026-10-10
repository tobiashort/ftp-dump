package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/tobiashort/th-utils/lib/must"
)

var (
	port     string
	username string
	password string
	dumpdir  string
)

func prompt(msg string, def string) string {
	fmt.Printf("%s [%s]: ", msg, def)
	r := bufio.NewReader(os.Stdin)
	a := must.Do2(r.ReadString('\n'))
	a = strings.TrimSpace(a)
	if a != "" {
		return a
	}
	return def
}

func main() {
	os1 := runtime.GOOS
	arch := runtime.GOARCH
	port = "21"
	username = "anonymous"
	password = "anonymous"
	dumpdir = must.Do2(os.Getwd())

	os1 = prompt("OS", os1)
	arch = prompt("Arch", arch)
	port = prompt("Port", port)
	username = prompt("Username", username)
	password = prompt("Password", password)
	dumpdir = prompt("Dumpdir", dumpdir)

	ldflags :=
		fmt.Sprintf(
			""+
				"-X 'github.com/tobiashort/ftp-dump/config.Port=%s' "+
				"-X 'github.com/tobiashort/ftp-dump/config.AllowedUsername=%s' "+
				"-X 'github.com/tobiashort/ftp-dump/config.AllowedPassword=%s'"+
				"-X 'github.com/tobiashort/ftp-dump/config.Dumpdir=%s'",
			port, username, password, dumpdir)

	cmd := exec.Command("go", "build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "CGO_ENABLED=0")
	cmd.Env = append(cmd.Env, "GOOS="+os1)
	cmd.Env = append(cmd.Env, "GOARCH="+arch)
	cmd.Args = append(cmd.Args, "-o", "ftp-dump")
	cmd.Args = append(cmd.Args, "-ldflags", ldflags)
	cmd.Args = append(cmd.Args, "./bin")
	must.Do(cmd.Run())
}
