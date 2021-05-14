package kits

import (
	"go-bootstraper/configs"
	"go-bootstraper/internal/utils"
)

type parser struct {
	walker     *utils.JohnyWalker
	paramMaker *paramMaker
}

func newParser(resource *configs.App, paramMaker *paramMaker, bootstrap bool) *parser {
	return &parser{
		paramMaker: paramMaker,
		walker: utils.NewWalker().
			Override(true).
			Force(bootstrap).
			WithFS(resource.AssetFS).
			From(resource.AssetPath).
			To("."),
	}
}

func (p *parser) basket(b basket) *utils.JohnyWalker {
	return p.walker.From(string(b))
}
func (p *parser) parseBasket(b basket, params map[string]interface{}) error {
	return p.basket(b).WithParams(params).Go()
}
func (p *parser) generateBaseStructure(params BaseParams) error {
	return p.parseBasket(base, params.toMap())
}

func (p *parser) updateConfigs(params ConfigParams) error {
	return p.parseBasket(config, params.toMap())
}

func (p *parser) generateModel(params ModelsParams) error {
	return p.
		basket(model).
		WithParams(params.toMap()).
		Go()
}

func (p *parser) generateStore(params StoresParams) error {
	return p.
		basket(store).
		WithParams(params.toMap()).
		Go()
}
func (p *parser) generateServiceInterface(params ProtoParams) error {
	return p.
		basket(proto).
		WithParams(params.toMap()).
		Go()
}

func (p *parser) generateServiceImpl(params ServerParams) error {
	return p.
		basket(service).
		WithParams(params.toMap()).
		Go()
}

func (p *parser) generateServiceInit(params CmdParams) error {
	return p.
		basket(cmd).
		WithParams(params.toMap()).
		Go()
}

func (p *parser) generateProjectStructure() error {
	return p.generateBaseStructure(p.paramMaker.createBaseParams())
}
func (p *parser) generateServiceInterfaces() error {
	for _, proto := range p.paramMaker.createAllProtoParams() {
		if err := p.generateServiceInterface(proto); err != nil {
			return err
		}
	}
	return nil
}
func (p *parser) generateServiceImpls() error {
	for _, server := range p.paramMaker.createAllServerParams() {
		if err := p.generateServiceImpl(server); err != nil {
			return err
		}

	}
	return nil
}

type basket string

const (
	base    basket = "base"
	cmd     basket = "cmd"
	config  basket = "configs"
	mapping basket = "mapping"
	model   basket = "models"
	proto   basket = "proto"
	service basket = "services"
	store   basket = "stores"
)
