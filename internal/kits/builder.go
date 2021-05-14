package kits

import (
	"os"

	"go-bootstraper/configs"
	"go-bootstraper/internal/models"
	store2 "go-bootstraper/internal/store"
	"go-bootstraper/internal/utils"
)

type buildingStep func() error

type builder []buildingStep

func (b builder) add(step ...buildingStep) builder {
	return append(b, step...)
}
func (b builder) Build() error {
	for _, w := range b {
		if err := w(); err != nil {
			return err
		}
	}
	return nil
}

func NewBootstrapTeam(material *configs.App, drawing *models.Drawing) (builder, error) {
	return NewTeam(true, material, drawing)
}

func NewBuildTeam(material *configs.App, drawing *models.Drawing) (builder, error) {
	if len(drawing.Server) == 0 {
		allProto := make([]utils.Name, 0, len(drawing.ProtoServices))
		for _, s := range drawing.ProtoServices {
			allProto = append(allProto, s.Name)
		}
		drawing.Server = []models.Server{
			{
				Name:     drawing.Egg.Name,
				Protos:   allProto,
				WithInit: true,
				Default:  true,
			},
		}
	}
	return NewTeam(false, material, drawing)
}

func NewTeam(bootstrap bool, material *configs.App, drawing *models.Drawing) (builder, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	p := newParser(
		material,
		newParamMaker(drawing),
		bootstrap,
	)
	j := newCmdBase(dir)
	var b = builder{
		p.generateProjectStructure, // generate base structure
	}
	if bootstrap {
		b = b.add(
			j.initGoogleApis,
			j.init,   // init project
			j.fmt,    // reformat files
			j.update, // update dependencies
		)
	}
	b = b.add(
		p.generateServiceInterfaces, // generate proto & client
		j.generateProto,             // use tools to generate service interface
		p.generateServiceImpls,      // generate service implement
		j.fmt,                       // reformat files
		j.update,                    // update dependencies
		func() error { return store2.NewSketch("file://.bt.sample.yml").Store(seedConfigSample()) },
	)

	return b, nil
}
