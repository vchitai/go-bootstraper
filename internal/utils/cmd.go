package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

type xCmd struct {
	*exec.Cmd
	before func() error
	after  func(error) error
}

func (xCmd *xCmd) RunAtDir(dir string) error {
	if err := xCmd.before(); err != nil {
		return err
	}
	xCmd.Dir = dir
	xCmd.Stdout = os.Stdout
	xCmd.Stderr = os.Stderr
	err := xCmd.Run()
	if err != nil {
		err = fmt.Errorf("while running cmd %s, %w", xCmd.String(), err)
	}
	return xCmd.after(err)
}

func (xCmd *xCmd) HookBefore(before func() error) *xCmd {
	xCmd.before = before
	return xCmd
}
func (xCmd *xCmd) HookAfter(after func(error) error) *xCmd {
	xCmd.after = after
	return xCmd
}

func NewxCmd(name string, args ...string) *xCmd {
	cmd := xCmd{Cmd: exec.Command(name, args...)}
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr
	log.Printf("running %s", cmd.String())
	cmd.before = func() error {
		return nil
	}
	cmd.after = func(err error) error {
		return err
	}
	return &cmd
}
