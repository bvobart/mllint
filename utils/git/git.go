package git

import "os/exec"

// Detect detects whether this directory is inside a Git repository
func Detect(dir string) bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = dir
	_, err := cmd.Output()
	return err == nil
}
