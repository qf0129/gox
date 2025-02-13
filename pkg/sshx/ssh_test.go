package sshx

import (
	"testing"
	"time"
)

func TestRunCommand(t *testing.T) {
	RunCommand(&RunCommandOption{
		Host:     "192.168.11.104",
		Port:     22,
		User:     "root",
		Password: "root",
		Command:  "ls /",
		Timeout:  10 * time.Second,
	})
}
