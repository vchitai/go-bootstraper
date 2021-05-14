package services

import (
	"context"

	"domain_dash/service_dash/configs"
	"domain_dash/service_dash/pb"
)

var _ serverNameCamelInterface = &serverNameCamelImpl{}

type serverNameCamelImpl struct {
	cfg *configs.Config
}

func NewServerNameTitle(
	cfg *configs.Config,
) ServiceServer {
	return newServerNameTitleServer(
		&serverNameCamelImpl{
			cfg: cfg,
		},
	)
}

/* range .Protos */
func (s serverNameCamelImpl) HelloProtoNameTitle(ctx context.Context, req *pb.HelloProtoNameTitleRequest) (*pb.HelloProtoNameTitleResponse, error) {
	return &pb.HelloProtoNameTitleResponse{}, nil
}

/* end */
