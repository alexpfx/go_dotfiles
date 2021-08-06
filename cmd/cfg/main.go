package main

import (
	"flag"
	"fmt"
	"github.com/alexpfx/go_dotfiles/internal/dotfile"
	"github.com/alexpfx/go_dotfiles/internal/util"

	"log"
)

const git = "/usr/bin/git"

const defaultAlias = "cfg"

func main() {
	var gitDir string
	var workTree string
	var alias string
	var updateConfig bool
	var help bool

	flag.StringVar(&gitDir, "d", "", "gitDir")
	flag.StringVar(&workTree, "t", "", "workTree")
	flag.StringVar(&alias, "a", defaultAlias, "command alias")
	flag.BoolVar(&updateConfig, "u", false, "write new config file and exit")
	flag.BoolVar(&help, "h", false, "print usage and exit")

	flag.Parse()

	if help {
		flag.PrintDefaults()
		return
	}

	if updateConfig {
		checkArgs(gitDir, workTree, alias)
		conf := dotfile.Config{
			WorkTree: workTree,
			GitDir:   gitDir,
		}
		util.WriteConfig(alias, &conf)
		return
	}

	conf := util.LoadConfig(alias)

	tail := flag.Args()
	aliasArgs := []string{
		"--git-dir=" + conf.GitDir + "/",
		"--work-tree=" + conf.WorkTree,
	}
	if len(tail) == 0 {
		return
	}
	out, stderr, err := util.Call(git, append(aliasArgs, tail...))
	util.Check(err, stderr)
	fmt.Println(out)
}

func checkArgs(args ...string) {
	for _, s := range args {
		if s == "" {
			log.Fatal("all parameters must be provided")
		}
	}

}
