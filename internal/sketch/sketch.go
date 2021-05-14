package sketch

import (
	"fmt"

	"go-bootstraper/internal/models"
	"go-bootstraper/internal/utils"
)

// for validating the implementation
type ISketcher interface {
	AddProtoServices(...models.Proto) error
	RemoveProtoServices(...string) error
	AddServices(...models.Server) error
	RemoveServices(...string) error
}

// wraps strategy to deal with project sketch map
type sketcher struct {
	drawing *models.Drawing
}

var _ ISketcher = &sketcher{}

func New(drawing *models.Drawing) *sketcher {
	if drawing == nil {
		// no drawing, no sketcher
		return nil
	}
	return &sketcher{drawing: drawing}
}

func (w *sketcher) AddProtoServices(proto ...models.Proto) error {
	protoSet := utils.NewSetStr(0)
	for _, proto := range proto {
		if protoSet.Contains(proto.Name.LowerCamelCase().String()) {
			return fmt.Errorf("dupplicated proto entry")
		}
		protoSet.Put(proto.Name.LowerCamelCase().String())
	}

	for _, proto := range w.drawing.ProtoServices {
		if protoSet.Contains(proto.Name.LowerCamelCase().String()) {
			return fmt.Errorf("duplicated proto entry")
		}
	}
	w.drawing.ProtoServices = append(w.drawing.ProtoServices, proto...)

	return nil
}

func (w *sketcher) RemoveProtoServices(protos ...string) error {
	protoSet := utils.FromList(protos)

	remainProtoServices := make([]models.Proto, 0, len(w.drawing.ProtoServices))
	for _, proto := range w.drawing.ProtoServices {
		if !protoSet.Contains(proto.Name.LowerCamelCase().String()) {
			remainProtoServices = append(remainProtoServices, proto)
		}
	}

	w.drawing.ProtoServices = remainProtoServices
	return nil
}

func (w *sketcher) AddServices(services ...models.Server) error {
	newStoresSet := utils.NewSetStr(0)
	for _, service := range services {
		if newStoresSet.Contains(service.Name.LowerCamelCase().String()) {
			return fmt.Errorf("dupplicated service entry")
		}
		newStoresSet.Put(service.Name.LowerCamelCase().String())
	}

	for _, service := range w.drawing.Server {
		if newStoresSet.Contains(service.Name.LowerCamelCase().String()) {
			return fmt.Errorf("duplicated service entry")
		}
	}
	w.drawing.Server = append(w.drawing.Server, services...)
	return nil
}

func (w *sketcher) RemoveServices(services ...string) error {
	serviceSet := utils.FromList(services)

	remainServices := make([]models.Server, 0, len(w.drawing.Server))
	for _, service := range w.drawing.Server {
		if !serviceSet.Contains(service.Name.LowerCamelCase().String()) {
			remainServices = append(remainServices, service)
		}
	}

	w.drawing.Server = remainServices
	return nil
}

func (w *sketcher) Init(services ...models.Server) error {
	newStoresSet := utils.NewSetStr(0)
	for _, service := range services {
		if newStoresSet.Contains(service.Name.LowerCamelCase().String()) {
			return fmt.Errorf("dupplicated service entry")
		}
		newStoresSet.Put(service.Name.LowerCamelCase().String())
	}

	for idx, model := range w.drawing.Server {
		if newStoresSet.Contains(model.Name.LowerCamelCase().String()) {
			model.WithInit = true
			w.drawing.Server[idx] = model
		}
	}

	return nil
}

func (w *sketcher) retrieve() *models.Drawing {
	return w.drawing
}
