package rust

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	mpath "github.com/somatech1/mikros-cli/pkg/path"
	"github.com/somatech1/mikros-cli/pkg/process"
)

func cargoInit(destinationPath, name string) error {
	// Switch to the destination path to create the crate
	cwd, err := mpath.ChangeDir(filepath.Dir(destinationPath))
	if err != nil {
		return err
	}

	defer func() {
		_ = os.Chdir(cwd)
	}()

	_, err = process.Exec("cargo", "init", name)
	return err
}

func cargoAdd(destinationPath, name, version, git, path string, features []string) error {
	// Switch to the destination path to add the dependency
	cwd, err := mpath.ChangeDir(destinationPath)
	if err != nil {
		return err
	}

	defer func() {
		_ = os.Chdir(cwd)
	}()

	args := []string{
		"cargo",
		"add",
	}
	if version != "" {
		args = append(args, fmt.Sprintf("%v@%v", name, version))
	}
	if git != "" {
		args = append(args, "--git", git)
	}
	if path != "" {
		args = append(args, "--path", path)
	}
	if len(features) > 0 {
		args = append(args, "--features", strings.Join(features, ","))
	}

	out, err := process.Exec(args...)
	println(string(out))
	return err
}