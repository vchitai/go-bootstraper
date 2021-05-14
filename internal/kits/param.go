package kits

import (
	"go-bootstraper/internal/models"
	"go-bootstraper/internal/utils"
)

type paramMaker struct {
	drawing *models.Drawing

	protoMap   map[utils.Name]models.Proto
	serviceMap map[utils.Name]models.Server
}

func newParamMaker(drawing *models.Drawing) *paramMaker {
	protoMap := make(map[utils.Name]models.Proto, len(drawing.ProtoServices))
	serviceMap := make(map[utils.Name]models.Server, len(drawing.Server))

	for _, proto := range drawing.ProtoServices {
		protoMap[proto.Name] = proto
	}
	for _, service := range drawing.Server {
		serviceMap[service.Name] = service
	}
	return &paramMaker{
		drawing:    drawing,
		protoMap:   protoMap,
		serviceMap: serviceMap,
	}
}

func (p paramMaker) createBaseParams() BaseParams {
	srvs := make([]utils.Name, 0, len(p.drawing.Server))
	for _, server := range p.drawing.Server {
		srvs = append(srvs, server.Name)
	}
	return BaseParams{
		Service: p.drawing.Egg.Name.LowerSnakeCase(),
		Domain:  p.drawing.Egg.Domain.LowerSnakeCase(),
		Servers: srvs,
	}
}

func (p paramMaker) createAllProtoParams() []ProtoParams {
	base := p.createBaseParams()

	protos := make([]ProtoParams, 0, len(p.drawing.ProtoServices))

	for _, proto := range p.drawing.ProtoServices {
		protos = append(protos,
			ProtoParams{
				BaseParams: base,
				Name:       proto.Name.LowerSnakeCase(),
				Default:    proto.Default,
			},
		)
	}
	return protos
}

func (p paramMaker) createProtoParams(protoNames ...utils.Name) []ProtoParams {
	base := p.createBaseParams()

	protos := make([]ProtoParams, 0, len(p.drawing.ProtoServices))

	for _, protoName := range protoNames {
		proto, ok := p.protoMap[protoName]
		if !ok {
			continue
		}
		protos = append(protos,
			ProtoParams{
				BaseParams: base,
				Name:       proto.Name.LowerSnakeCase(),
				Default:    proto.Default,
			},
		)
	}
	return protos
}

func (p paramMaker) createSingleServerParams(serverName utils.Name) ServerParams {
	server, ok := p.serviceMap[serverName]
	if !ok {
		return ServerParams{}
	}
	base := p.createBaseParams()
	return ServerParams{
		BaseParams: base,
		Name:       server.Name.LowerSnakeCase(),
		Protos:     p.createProtoParams(server.Protos...),
		Default:    server.Default,
	}
}
func (p paramMaker) createCRUDServerParams() []CRUDServerParams {
	return []CRUDServerParams{}
}
func (p paramMaker) createAllServerParams() []ServerParams {
	res := make([]ServerParams, 0, len(p.drawing.Server))
	for _, server := range p.drawing.Server {
		res = append(res, p.createSingleServerParams(server.Name))
	}
	return res
}
func (p paramMaker) createCmdParams() []CmdParams {
	res := make([]CmdParams, 0, len(p.drawing.Server))
	for _, server := range p.drawing.Server {
		if !server.WithInit {
			continue
		}
		res = append(res, CmdParams{
			BaseParams: p.createBaseParams(),
			Name:       server.Name.LowerSnakeCase(),
			Default:    server.Default,
		})
	}
	return res
}

type BaseParams struct {
	Service utils.Name
	Domain  utils.Name
	Servers []utils.Name
}

func (p BaseParams) toMap() map[string]interface{} {
	serverMap := make([]map[string]string, 0, len(p.Servers))
	for _, srv := range p.Servers {
		serverMap = append(serverMap, map[string]string{
			"ServerNameTitle":   srv.UpperCamelCase().String(),
			"serverNameCamel":   srv.LowerCamelCase().String(),
			"server_name_lower": srv.LowerSnakeCase().String(),
			"server_name_dash":  srv.LowerDashNotation().String(),
		})
	}
	return map[string]interface{}{
		"ServiceTitle":  p.Service.UpperCamelCase().String(),
		"serviceCamel":  p.Service.LowerCamelCase().String(),
		"service_lower": p.Service.LowerSnakeCase().String(),
		"service_dash":  p.Service.LowerDashNotation().String(),
		"DomainTitle":   p.Domain.UpperCamelCase().String(),
		"domainCamel":   p.Domain.LowerCamelCase().String(),
		"domain_lower":  p.Domain.LowerSnakeCase().String(),
		"domain_dash":   p.Domain.LowerDashNotation().String(),
		"Servers":       serverMap,
	}
}

type ModelsParams struct {
	BaseParams
	Name utils.Name
}

func (p ModelsParams) toMap() map[string]interface{} {
	mm := p.BaseParams.toMap()
	mm["ModelNameTitle"] = p.Name.UpperCamelCase().String()
	mm["modelNameCamel"] = p.Name.LowerCamelCase().String()
	mm["model_name_lower"] = p.Name.LowerSnakeCase().String()
	mm["model_name_dash"] = p.Name.LowerDashNotation().String()
	return mm
}

type DBParams struct {
	BaseParams
	Name   utils.Name
	Driver utils.Name
}

func (p DBParams) toMap() map[string]interface{} {
	mm := p.BaseParams.toMap()
	mm["DBNameTitle"] = p.Name.UpperCamelCase().String()
	mm["dbNameCamel"] = p.Name.LowerCamelCase().String()
	mm["db_name_lower"] = p.Name.LowerSnakeCase().String()
	mm["db_name_dash"] = p.Name.LowerDashNotation().String()
	mm["Driver"] = p.Driver.LowerSnakeCase()
	return mm
}

type ConfigParams struct {
	BaseParams
	DBs []DBParams
}

func (p ConfigParams) toMap() map[string]interface{} {
	mm := p.BaseParams.toMap()
	dbs := make([]map[string]interface{}, 0, len(p.DBs))
	for _, db := range p.DBs {
		dbs = append(dbs, db.toMap())
	}
	mm["DBs"] = dbs
	return mm
}

type StoresParams struct {
	BaseParams
	Name          utils.Name
	DB            DBParams
	Models        []ModelsParams
	CRUDInterface bool
	CRUDImpl      bool
}

func (p StoresParams) toMap() map[string]interface{} {
	mm := p.BaseParams.toMap()
	mm["StoreNameTitle"] = p.Name.UpperCamelCase().String()
	mm["storeNameCamel"] = p.Name.LowerCamelCase().String()
	mm["store_name_lower"] = p.Name.LowerSnakeCase().String()
	mm["store_name_dash"] = p.Name.LowerDashNotation().String()
	mm["CRUDInterface"] = p.CRUDInterface
	mm["CRUDImpl"] = p.CRUDImpl
	mm["DB"] = p.DB.toMap()
	for a, b := range p.DB.toMap() {
		mm[a] = b
	}
	mdls := make([]map[string]interface{}, 0, len(p.Models))
	for _, model := range p.Models {
		m := model.toMap()
		m["StoreNameTitle"] = mm["StoreNameTitle"]
		m["storeNameCamel"] = mm["storeNameCamel"]
		m["store_name_lower"] = mm["store_name_lower"]
		m["store_name_dash"] = mm["store_name_dash"]
		mdls = append(mdls, m)
	}
	mm["Models"] = mdls
	return mm
}

type ProtoParams struct {
	BaseParams
	Name          utils.Name
	Models        []ModelsParams
	CRUDInterface bool
	Default       bool
}

func (p ProtoParams) toMap() map[string]interface{} {
	mm := p.BaseParams.toMap()
	mm["ProtoNameTitle"] = p.Name.UpperCamelCase().String()
	mm["protoNameCamel"] = p.Name.LowerCamelCase().String()
	mm["proto_name_lower"] = p.Name.LowerSnakeCase().String()
	mm["proto_name_dash"] = p.Name.LowerDashNotation().String()
	mm["CRUDInterface"] = p.CRUDInterface
	if p.Default {
		mm["Default"] = "client"
	}

	mdls := make([]map[string]interface{}, 0, len(p.Models))
	for _, model := range p.Models {
		mdls = append(mdls, model.toMap())
	}
	mm["Models"] = mdls
	return mm
}

type CRUDServerParams struct {
	BaseParams
	Model ModelsParams
	Store StoresParams
}

func (p CRUDServerParams) toMap() map[string]interface{} {
	mm := p.BaseParams.toMap()
	mm["Store"] = p.Store.toMap()
	mm["Model"] = p.Model.toMap()

	for a, b := range p.Store.toMap() {
		mm[a] = b
	}
	for a, b := range p.Model.toMap() {
		mm[a] = b
	}
	return mm
}

type ServerParams struct {
	BaseParams
	Name          utils.Name
	Protos        []ProtoParams
	Stores        []StoresParams
	CRUDInterface bool
	CRUDImpl      bool
	CRUDs         []CRUDServerParams
	Default       bool
}

func (p ServerParams) toMap() map[string]interface{} {
	mm := p.BaseParams.toMap()
	mm["ServerNameTitle"] = p.Name.UpperCamelCase().String()
	mm["serverNameCamel"] = p.Name.LowerCamelCase().String()
	mm["server_name_lower"] = p.Name.LowerSnakeCase().String()
	mm["server_name_dash"] = p.Name.LowerDashNotation().String()
	stores := make([]map[string]interface{}, 0, len(p.Stores))
	for _, store := range p.Stores {
		stores = append(stores, store.toMap())
	}
	mm["Stores"] = stores

	protos := make([]map[string]interface{}, 0, len(p.Protos))
	for _, proto := range p.Protos {
		m := proto.toMap()
		m["ServerNameTitle"] = mm["ServerNameTitle"]
		m["serverNameCamel"] = mm["serverNameCamel"]
		m["server_name_lower"] = mm["server_name_lower"]
		m["server_name_dash"] = mm["server_name_dash"]
		protos = append(protos, m)
	}
	mm["Protos"] = protos
	mm["CRUDInterface"] = p.CRUDInterface
	mm["CRUDImpl"] = p.CRUDImpl

	mdlStore := make(map[utils.Name]StoresParams, 0)
	for _, store := range p.Stores {
		for _, mdl := range store.Models {
			mdlStore[mdl.Name] = store
		}
	}
	cruds := make([]map[string]interface{}, 0, len(p.CRUDs))
	for _, itf := range p.Protos {
		for _, mdl := range itf.Models {
			m := CRUDServerParams{
				Model: mdl,
				Store: mdlStore[mdl.Name],
			}.toMap()
			m["ServerNameTitle"] = mm["ServerNameTitle"]
			m["serverNameCamel"] = mm["serverNameCamel"]
			m["server_name_lower"] = mm["server_name_lower"]
			m["server_name_dash"] = mm["server_name_dash"]
			cruds = append(cruds, m)
		}
	}
	mm["CRUDs"] = cruds

	if p.Default {
		mm["Default"] = "service"
	}
	return mm
}

type CmdParams struct {
	BaseParams
	Name    utils.Name
	Stores  []StoresParams
	DBs     []DBParams
	Default bool
}

func (p CmdParams) toMap() map[string]interface{} {
	mm := p.BaseParams.toMap()
	mm["ServerNameTitle"] = p.Name.UpperCamelCase().String()
	mm["serverNameCamel"] = p.Name.LowerCamelCase().String()
	mm["server_name_lower"] = p.Name.LowerSnakeCase().String()
	mm["server_name_dash"] = p.Name.LowerDashNotation().String()
	storeMap := make([]map[string]interface{}, 0)
	for _, store := range p.Stores {
		storeMap = append(storeMap, store.toMap())
	}
	mm["Stores"] = storeMap
	dbMap := make([]map[string]interface{}, 0)
	for _, db := range p.DBs {
		dbMap = append(dbMap, db.toMap())
	}
	mm["DBs"] = dbMap
	if p.Default {
		mm["Default"] = "server"
	}
	return mm
}
