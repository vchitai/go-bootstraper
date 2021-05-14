package kits

import (
	"fmt"
	"os"

	"go-bootstraper/configs"
	store2 "go-bootstraper/internal/store"
	"go-bootstraper/internal/utils"
)

const filePerm = os.ModePerm

type Worker struct {
	resource *configs.App
	store    store2.SketchStore
}

func NewWorker(resource *configs.App) *Worker {
	w := &Worker{resource: resource}
	return w
}

func (w *Worker) Close() error {
	return nil
}

func (w *Worker) Boostrap(_name, _domain string) error {
	name, domain := utils.Name(_name), utils.Name(_domain)
	projectFolder := name.LowerDashNotation().String()

	if err := os.MkdirAll(projectFolder, filePerm); err != nil {
		return fmt.Errorf("cannot create project folder %w", err)
	}
	if err := os.Chdir(projectFolder); err != nil {
		return fmt.Errorf("cannot access project folder %w", err)
	}

	drawing := seedDrawing(name, domain)
	if err := store2.NewSketch("").Store(drawing); err != nil {
		return fmt.Errorf("cannot create config file %w", err)
	}
	b, err := NewBootstrapTeam(w.resource, drawing)
	if err != nil {
		return err
	}

	if err := b.Build(); err != nil {
		return err
	}

	b, err = NewBuildTeam(w.resource, drawing)
	if err != nil {
		return err
	}
	return b.Build()
}
