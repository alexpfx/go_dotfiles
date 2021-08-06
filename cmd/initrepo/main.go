package main

import (
	"flag"
	"fmt"
	"github.com/alexpfx/go_dotfiles/internal/dotfile"
	"github.com/alexpfx/go_dotfiles/internal/util"
	"os"
	"path/filepath"
)

const git = "/usr/bin/git"
const defaultAlias = "cfg"

func main() {

	var repo string
	var gitDir string
	var workTree string
	var alias string
	var force bool

	h, err := os.UserHomeDir()
	util.Check(err, "")

	flag.StringVar(&repo, "r", "https://github.com/alexpfx/sway_dotfiles.git", "repository")
	flag.StringVar(&gitDir, "d", filepath.Join(h, ".cfg"), "gitDir")
	flag.StringVar(&workTree, "t", h, "workTree")
	flag.StringVar(&alias, "a", defaultAlias, "command alias")
	flag.BoolVar(&force, "f", false, "remove ditDir if it exists")
	flag.Parse()

	conf := dotfile.Config{
		WorkTree: workTree,
		GitDir:   gitDir,
	}

	if force && util.DirExists(gitDir) {
		err := os.RemoveAll(gitDir)
		util.Check(err, "cannot remove gitDir")
	}

	_, serr, err := util.Call(git, []string{"clone", "--bare", repo, gitDir})
	util.Check(err, serr)

	aliasArgs := []string{
		"--git-dir=" + conf.GitDir + "/",
		"--work-tree=" + conf.WorkTree,
	}

	_, serr, err = util.Call(git, append(aliasArgs, "config", "--local", "status.showUntrackedFiles", "no"))

	util.WriteConfig(alias, &conf)

	checkout(alias, aliasArgs, workTree, &conf)

}

func checkout(alias string, aliasArgs []string, workTree string, conf *dotfile.Config) {
	var existUntracked []string
	_, serr, err := util.Call(git, append(aliasArgs, "checkout"))

	if err != nil {
		existUntracked = util.GetExistUntracked(workTree, serr)
		if len(existUntracked) == 0 {
			util.Check(err, err.Error())
		}

		util.BackupFiles(fmt.Sprintf(".%s%s_bkp/", workTree, alias), existUntracked)

		for _, untracked := range existUntracked {
			os.RemoveAll(untracked)
		}

		_, serr, err = util.Call(git, append(aliasArgs, "checkout"))
		util.Check(err, serr)
	}

}
