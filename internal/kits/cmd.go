package kits

import (
	"os"
	"path/filepath"

	"go-bootstraper/internal/utils"
)

type cmdBase struct {
	executeDir string
}

func newCmdBase(dir string) *cmdBase {
	return &cmdBase{executeDir: dir}
}

func (cb *cmdBase) generateProto() error {
	if x, err := filepath.Glob("proto/*.proto"); err != nil {
		return err
	} else {
		if len(x) == 0 {
			// not doing
			return nil
		}
	}
	return utils.
		NewxCmd("make", "generate").
		RunAtDir(cb.executeDir)
}
func (cb *cmdBase) initGoogleApis() error {
	return utils.
		NewxCmd("git", "clone", "git@github.com:gogo/googleapis.git", "include/googleapis").
		HookBefore(func() error {
			_ = os.MkdirAll(filepath.Join(cb.executeDir, "include/googleapis"), os.ModePerm)

			return nil
		}).
		RunAtDir(cb.executeDir)
}
func (cb *cmdBase) init() error {
	return utils.
		NewxCmd("make", "init").
		RunAtDir(cb.executeDir)
}
func (cb *cmdBase) update() error {
	return utils.
		NewxCmd("make", "update").
		RunAtDir(cb.executeDir)
}
func (cb *cmdBase) fmt() error {
	return utils.
		NewxCmd("make", "fmt").
		RunAtDir(cb.executeDir)
}
func (cb *cmdBase) setupDependencies() error {
	return utils.
		NewxCmd("make", "install-go-tools").
		RunAtDir(cb.executeDir)
}
func (cb *cmdBase) formatGoFiles() error {
	return utils.
		NewxCmd("gofumpt", "-l", "-w", "cmd", "pkg", "configs", "internal").
		HookBefore(func() error {
			os.MkdirAll(filepath.Join(cb.executeDir, "cmd"), os.ModePerm)
			os.MkdirAll(filepath.Join(cb.executeDir, "pkg"), os.ModePerm)
			os.MkdirAll(filepath.Join(cb.executeDir, "internal"), os.ModePerm)
			os.MkdirAll(filepath.Join(cb.executeDir, "configs"), os.ModePerm)
			return nil
		}).
		RunAtDir(cb.executeDir)
}

func (cb *cmdBase) fixImports() error {
	return utils.
		NewxCmd("goimports", "-v", "-w", "cmd", "pkg", "configs", "internal").
		HookBefore(func() error {
			os.MkdirAll(filepath.Join(cb.executeDir, "cmd"), os.ModePerm)
			os.MkdirAll(filepath.Join(cb.executeDir, "pkg"), os.ModePerm)
			os.MkdirAll(filepath.Join(cb.executeDir, "internal"), os.ModePerm)
			os.MkdirAll(filepath.Join(cb.executeDir, "configs"), os.ModePerm)
			return nil
		}).
		RunAtDir(cb.executeDir)
}
