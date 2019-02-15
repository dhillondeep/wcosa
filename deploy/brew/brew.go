package main

import (
	"bufio"
	"gopkg.in/src-d/go-git.v4"
	"log"
	"os"
	"path/filepath"
	"strings"
	"wio/internal/config/meta"
	"wio/pkg/util/sys"
	"wio/pkg/util/template"
)

const (
	execName  = "wio_darwin_x86_64"
	extension = "tar.gz"
)

func main() {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	currPath := filepath.Dir(ex)

	cloneOptions := &git.CloneOptions{
		URL:               "https://github.com/wio/homebrew-wio",
		Progress:          os.Stdout,
		Depth:             1,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	}

	repoPath := sys.Path(currPath, "homebrew-wio")

	if err := os.RemoveAll(repoPath); err != nil {
		log.Fatal(err)
	}

	if _, err := git.PlainClone(repoPath, false, cloneOptions); err != nil {
		log.Fatal(err)
	}

	checkSumsPath := sys.Path(currPath, "../", "bin", "checksums.txt")

	checkSumsFile, err := os.Open(checkSumsPath)
	if err != nil {
		log.Fatal(err)
	}
	defer checkSumsFile.Close()

	scanner := bufio.NewScanner(checkSumsFile)
	var checkSum string
	for scanner.Scan() {
		tokens := strings.Split(scanner.Text(), " ")

		if strings.Trim(tokens[2], " ") == execName+"."+extension {
			checkSum = tokens[0]
		}
	}

	rbFile := sys.Path(repoPath, "Formula", "wio.rb")

	if err := sys.NormalIO.CopyFile(sys.Path(currPath, "brew-template.tpl"), rbFile, true); err != nil {
		log.Fatal(err)
	}

	if err := template.IOReplace(rbFile, map[string]string{
		"version":   meta.Version,
		"execName":  execName,
		"extension": extension,
		"checksum":  checkSum,
	}); err != nil {
		log.Fatal(err)
	}
}
