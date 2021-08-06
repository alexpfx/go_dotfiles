package dotfile

type Config struct {
	GitDir string `json:"git_dir"`
	WorkTree string `json:"work_tree"`
}
