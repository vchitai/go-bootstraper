package client

import (
	"context"

	"domain_dash/service_dash/pb"
	"google.golang.org/grpc"
)

type ProtoNameTitleClient struct {
	pb.ProtoNameTitleClient
}

func NewProtoNameTitleClient(address string) (*ProtoNameTitleClient, error) {
	conn, err := grpc.DialContext(context.Background(), address,
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
	)

	if err != nil {
		return nil, err
	}

	return &ProtoNameTitleClient{
		ProtoNameTitleClient: pb.NewProtoNameTitleClient(conn),
	}, nil
}
