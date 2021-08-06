package util

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/alexpfx/go_dotfiles/internal/dotfile"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const appConfigDir = "go_dotfile"
const appConfigTemplate = "config_%s"

func Call(cmdStr string, args []string) (string, string, error) {
	cmd := exec.Command(cmdStr, args...)

	var sout bytes.Buffer
	var serr bytes.Buffer

	cmd.Stdout = &sout
	cmd.Stderr = &serr
	err := cmd.Run()

	return string(sout.Bytes()), string(serr.Bytes()), err
}

func Check(err error, msg string) {
	if err != nil {
		log.Print(msg)
		log.Fatal(err)
	}
}

func WriteConfig(aliasName string, conf *dotfile.Config) {
	cfgDir := resolveConfigDir()
	err := os.MkdirAll(cfgDir, 0700)
	Check(err, "")

	cfgPath := resolveConfigPath(cfgDir, aliasName)

	f, err := os.OpenFile(cfgPath, os.O_CREATE|os.O_RDWR, 0660)
	Check(err, "")
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "    ")

	err = enc.Encode(&conf)
	Check(err, "")
}

func QuoteArgs(args []string) []string {
	for i, a := range args {
		if strings.ContainsRune(a, ' ') {
			args[i] = strconv.Quote(a)
		}
	}
	return args
}

func BackupFiles(backupDir string, paths []string) {
	if len(paths) == 0 {
		return
	}

	backupErr := "cannot backup. stopping..."

	err := os.MkdirAll(backupDir, 0700)
	Check(err, backupErr)
	for _, path := range paths {
		if len(path) == 0 {
			continue
		}
		source, err := os.Open(path)
		Check(err, backupErr)

		backupFilePath := filepath.Join(backupDir, path)
		backupFileDir := filepath.Dir(backupFilePath)

		err = os.MkdirAll(backupFileDir, 0700)
		Check(err, backupErr+" "+backupFileDir)

		target, err := os.OpenFile(backupFilePath, os.O_CREATE|os.O_RDWR, 0660)
		Check(err, backupErr+" "+backupFileDir)
		_, err = io.Copy(target, source)
		Check(err, backupErr)

		err = target.Close()
		Check(err, backupErr+" "+backupFileDir)
		err = source.Close()
		Check(err, backupErr+" "+backupFileDir)
	}
}

func GetExistUntracked(workTree string, gitMessage string) []string {
	scanner := bufio.NewScanner(strings.NewReader(gitMessage))
	paths := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "\t") {
			paths = append(paths, filepath.Join(workTree, strings.TrimPrefix(line, "\t")))
		}
	}
	return paths
}

func LoadConfig(aliasName string) *dotfile.Config {
	cfgPath := resolveConfigPath(resolveConfigDir(), aliasName)

	f, err := os.Open(cfgPath)
	Check(err, "")
	defer f.Close()

	dec := json.NewDecoder(f)

	conf := dotfile.Config{}
	err = dec.Decode(&conf)
	Check(err, "")

	return &conf
}

func resolveConfigDir() string {
	userCfg, err := os.UserConfigDir()
	Check(err, "")
	return filepath.Join(userCfg, appConfigDir)
}
func resolveConfigPath(configDir, aliasName string) string {
	return filepath.Join(configDir, fmt.Sprintf(appConfigTemplate, aliasName))
}

func DirExists(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		Check(err, "")
	}
	return stat.IsDir()
}

func FileExists(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		Check(err, "")
	}
	return !stat.IsDir()

}
