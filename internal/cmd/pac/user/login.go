package user

import (
	"bufio"
	"os"
	"strings"
	"syscall"
	"wio/internal/cmd"
	"wio/pkg/log"
	"wio/pkg/npm/login"

	"golang.org/x/crypto/ssh/terminal"
)

type loginArgs struct {
	dir   string
	name  string
	pass  string
	email string
}

func (c Login) getArgs() (*loginArgs, error) {
	dir, err := cmd.GetDirectory(c)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(os.Stdin)
	log.Info(log.Cyan, "Username: ")
	username, _ := reader.ReadString('\n')

	log.Info(log.Cyan, "Email: ")
	email, _ := reader.ReadString('\n')

	log.Info(log.Cyan, "Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, err
	}
	log.Infoln()

	return &loginArgs{
		dir:   dir,
		name:  strings.Trim(username, "\n"),
		pass:  string(bytePassword),
		email: strings.Trim(email, "\n"),
	}, nil
}

func (c Login) Execute() error {
	args, err := c.getArgs()
	if err != nil {
		return err
	}
	log.Info(log.Cyan, "Sending login info ... ")
	token, err := login.GetToken(args.name, args.pass, args.email)
	if err != nil {
		log.WriteFailure()
		return err
	}
	log.WriteSuccess()
	log.Info(log.Cyan, "Saving login token ... ")
	if err := token.Save(args.dir); err != nil {
		log.WriteFailure()
		return err
	}
	log.WriteSuccess()
	log.Infoln(log.Yellow, "User logged in")

	return nil
}
