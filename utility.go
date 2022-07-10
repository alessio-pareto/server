package server

import (
	"io"
	"os/exec"
	"strings"
)

type Process struct {
	cmd exec.Cmd
	exicC chan struct{}
	stdin io.WriteCloser
	exitErr error
}

func NewProcess(name string, args ...string) *Process {
	return &Process {
		cmd: *exec.Command(name, args...),
		exicC: make(chan struct{}),
	}
}

func (p *Process) Start() (err error) {
	p.stdin, err = p.cmd.StdinPipe()
	if err != nil {
		return
	}

	err = p.cmd.Start()
	if err != nil {
		return
	}

	go func() {
		p.exitErr = p.cmd.Wait()
		p.exicC <- struct{}{}
	}()

	return
}

func (p *Process) Stop() (int, error) {
	_, err := p.stdin.Write([]byte("\x03"))
	if err != nil {
		return -1, err
	}

	<- p.exicC
	return p.cmd.ProcessState.ExitCode(), p.exitErr
}

func (p *Process) Kill() (int, error) {
	p.cmd.Process.Kill()
	<- p.exicC
	return p.cmd.ProcessState.ExitCode(), p.exitErr
}

func ParseCommandArgs(args ...string) []string {
	a := make([]string, 0)
	for _, s := range args {
		for i, s1 := range strings.Split(s, "'") {
			if i%2 == 1 {
				a = append(a, s1)
				continue
			}

			for j, s2 := range strings.Split(s1, "\"") {
				if j%2 == 1 {
					a = append(a, s2)
					continue
				}

				for _, s3 := range strings.Split(s2, " ") {
					if s3 != "" {
						a = append(a, s3)
					}
				}
			}
		}
	}

	return a
}