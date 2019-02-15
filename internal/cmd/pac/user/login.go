package user

import (
	"bufio"
	"os"
	"strings"
	"syscall"
	"wio/pkg/log"
	"wio/pkg/npm/login"
	"wio/pkg/npm/registry"

	"golang.org/x/crypto/ssh/terminal"
)

type loginArgs struct {
	name  string
	pass  string
	email string
}

func (c Login) getArgs() (*loginArgs, error) {
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
	tokens, err := login.GetToken(args.name, args.pass, args.email, registry.WioPackageRegistry)
	if err != nil {
		log.WriteFailure()
		return err
	}
	log.WriteSuccess()
	log.Info(log.Cyan, "Saving login token ... ")
	if err := tokens.Save(); err != nil {
		log.WriteFailure()
		return err
	}
	log.WriteSuccess()
	log.Infoln(log.Yellow, "User logged in")

	return nil
}
