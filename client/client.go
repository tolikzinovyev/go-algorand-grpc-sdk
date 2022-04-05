package client

import (
	"context"

	"google.golang.org/grpc"

	"github.com/tolikzinovyev/go-algorand-grpc-sdk/client/internal/proto"
	"github.com/tolikzinovyev/go-algorand-grpc-sdk/types"
)

type Client struct {
	c proto.DefaultClient
}

func MakeClient(conn grpc.ClientConnInterface) Client {
	return Client{c: proto.NewDefaultClient(conn)}
}

func (c Client) AccountInformation(address types.Address) (AccountResponse, error) {
	response, err := c.c.AccountInformation(
		context.Background(), &proto.AccountRequest{Address: address[:]})
	if err != nil {
		return AccountResponse{}, err
	}

	return unconvertAccountResponse(response), nil
}
