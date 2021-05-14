package models

import "go-bootstraper/internal/utils"

type Egg struct {
	Name   utils.Name `yaml:"name"`   // Server Name
	Domain utils.Name `yaml:"domain"` // Server Domain
}

type Proto struct {
	Name    utils.Name   `yaml:",omitempty"`
	Models  []utils.Name `yaml:",omitempty"`
	Default bool         `yaml:",omitempty"`
}

type Server struct {
	Name     utils.Name   `yaml:",omitempty"`
	Stores   []utils.Name `yaml:",omitempty"`
	Protos   []utils.Name `yaml:",omitempty"`
	WithInit bool         `yaml:",omitempty"`
	Default  bool         `yaml:",omitempty"`
}

type Drawing struct {
	Egg           Egg      `yaml:",omitempty"`
	Server        []Server `yaml:",omitempty"`
	ProtoServices []Proto  `yaml:",omitempty"`
	Previous      *Drawing `yaml:",omitempty"` // for regenerating only differences later
}
