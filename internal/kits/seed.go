package kits

import (
	"go-bootstraper/internal/models"
	"go-bootstraper/internal/utils"
)

func seedDrawing(name, domain utils.Name) *models.Drawing {
	return &models.Drawing{
		Egg: models.Egg{
			Name:   name,
			Domain: domain,
		},
		ProtoServices: []models.Proto{
			{
				Name:    name.LowerDashNotation(),
				Default: true,
			},
		},
	}
}

func seedConfigSample() *models.Drawing {
	return &models.Drawing{
		Egg: models.Egg{
			Name:   "service-name",
			Domain: "domain-name",
		},
		ProtoServices: []models.Proto{
			{
				Name: "pb-service-name",
			},
			{
				Name: "pb-service-name-2",
			},
		},
	}
}
