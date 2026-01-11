package uciadapter

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

type Engine struct {
	cmd *exec.Cmd
	in  io.Writer
	out *bufio.Scanner
}

func Start(path string, args ...string) (*Engine, error) {
	var cmd = exec.Command(path, args...)
	in, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	return &Engine{
		cmd: cmd,
		in:  in,
		out: bufio.NewScanner(out),
	}, nil
}

func (e *Engine) Wait() error {
	return e.cmd.Wait()
}

func (e *Engine) Uci() EngineInfo {
	fmt.Fprintln(e.in, "uci")
	for e.out.Scan() {
		var msg = e.out.Text()
		if msg == "uciok" {
			return EngineInfo{}
		}
	}
	// TODO err
	return EngineInfo{}
}

func (e *Engine) SetOption(option Option) error {
	fmt.Fprintf(e.in, "setoption name %v value %v\n", option.Name, option.Value)
	return nil
}

func (e *Engine) UciNewgame() {
	fmt.Fprintln(e.in, "ucinewgame")
}

func (e *Engine) IsReady() {
	fmt.Fprintln(e.in, "isready")
	for e.out.Scan() {
		var msg = e.out.Text()
		if msg == "readyok" {
			return
		}
	}
}

func (e *Engine) Stop() {
	fmt.Fprintln(e.in, "stop")
}

func (e *Engine) Quit() {
	fmt.Fprintln(e.in, "quit")
}
