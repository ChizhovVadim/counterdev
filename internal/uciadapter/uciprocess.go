package uciadapter

import (
	"fmt"
	"io"
	"os/exec"
)

type UciProcess struct {
	path string
	args []string
	cmd  *exec.Cmd
	in   io.WriteCloser
	out  io.ReadCloser
}

func NewUciProcess(path string, args []string) *UciProcess {
	return &UciProcess{path: path, args: args}
}

func (uci *UciProcess) Start() error {
	var cmd = exec.Command(uci.path, uci.args...)
	uci.cmd = cmd
	in, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	uci.in = in
	out, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	uci.out = out
	err = cmd.Start()
	if err != nil {
		return err
	}
	return nil
}

func (uci *UciProcess) Close() error {
	fmt.Fprintln(uci.in, "quit")
	return uci.cmd.Wait()
}
