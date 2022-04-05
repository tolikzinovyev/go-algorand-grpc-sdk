package main

import (
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/tolikzinovyev/go-algorand-grpc-sdk/client"
	"github.com/tolikzinovyev/go-algorand-grpc-sdk/types"
)

func main() {
	conn, err := grpc.Dial(
		"127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := client.MakeClient(conn)

	address, _ := types.DecodeAddress(
		"7777777777777777777777777777777777777777777777777774MSJUVU")

	resp, err := client.AccountInformation(address)
	if err != nil {
		log.Fatalf("rpc failed: %v", err)
	}

	fmt.Printf("%+v\n", resp)
}
