package cli

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func EnsureBinary(name string) error {
	if _, err := exec.LookPath(name); err != nil {
		return fmt.Errorf("%s not found in PATH", name)
	}
	return nil
}

func RunCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func RunCommandOutput(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var buffer bytes.Buffer
	cmd.Stdout = &buffer
	cmd.Stderr = &buffer
	if err := cmd.Run(); err != nil {
		return buffer.String(), err
	}
	return buffer.String(), nil
}
