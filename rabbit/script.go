package rabbit

import (
	"encoding/base64"
	"os/exec"
	"strings"
)

type ScriptCaller interface {
	ExecWith([]byte) ([]byte, error)
}

type Script struct {
	Config *Config
}

func (s Script) ExecWith(msg []byte) ([]byte, error) {
	var (
		script = s.Config.Consumer.Worker.Script
		args   []string
	)

	if subs := strings.Split(script, " "); len(subs) > 1 {
		script, args = subs[0], subs[1:]
	}

	args = append(args, base64.StdEncoding.EncodeToString(msg))
	return exec.Command(script, args...).CombinedOutput()
}
