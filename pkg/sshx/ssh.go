package sshx

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/qf0129/gox/pkg/logx"
	"golang.org/x/crypto/ssh"
)

type RunCommandOption struct {
	Host       string
	Port       int
	User       string
	Password   string
	PrivateKey string
	Command    string
	Timeout    time.Duration
	Stdout     io.Writer
	Stderr     io.Writer
}

// 通过ssh连接目标机器执行命令
func RunCommand(opt *RunCommandOption) error {
	authMethod := ssh.Password(opt.Password)
	if opt.PrivateKey != "" {
		signer, err := ssh.ParsePrivateKey([]byte(opt.PrivateKey))
		if err != nil {
			logx.Error("sshParsePrivateKeyErr", "err", err)
			return err
		}
		authMethod = ssh.PublicKeys(signer)
	}

	config := &ssh.ClientConfig{
		User:            opt.User,
		Auth:            []ssh.AuthMethod{authMethod},
		Timeout:         opt.Timeout,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", opt.Host, opt.Port), config)
	if err != nil {
		logx.Error("sshDialErr", "err", err)
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		logx.Error("sshNewSessionErr", "err", err)
		return err
	}
	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	if opt.Stdout != nil {
		session.Stdout = opt.Stdout
	}
	if opt.Stderr != nil {
		session.Stderr = opt.Stderr
	}
	logx.Info("sshRun", "cmd", opt.Command)
	err = session.Run(opt.Command)
	if err != nil {
		logx.Error("sshRunErr", "err", err)
		return err
	}
	return nil
}
