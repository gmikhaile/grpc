package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"main/unary/client/ecommerce"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const port = ":12345"

func main() {
	conn, err := grpc.Dial(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()
	client := ecommerce.NewOrderManagementClient(conn)

	fmt.Println("get order")
	order, err := client.GetOrder(context.TODO(), &wrapperspb.StringValue{Value: "1"})
	if err != nil {
		log.Fatalf("failed to get order: %v", err)
	}
	fmt.Printf("%v\n", order)

	fmt.Println("search orders")
	orderStream, err := client.SearchOrders(context.TODO(), &wrapperspb.StringValue{Value: "1"})
	if err != nil {
		log.Fatalf("failed to search order: %v", err)
	}

	order, err = orderStream.Recv()
	if err != nil {
		log.Fatalf("failed to recv order 1: %v", err)
	}
	fmt.Printf("%v\n", order)

	order, err = orderStream.Recv()
	if err != nil {
		log.Fatalf("failed to recv order 2: %v", err)
	}
	fmt.Printf("%v\n", order)

	fmt.Println("update orders")
	updateStream, err := client.UpdateOrders(context.TODO())
	if err != nil {
		log.Fatalf("failed to update order: %v", err)
	}
	updateStream.Send(&ecommerce.Order{
		Id: "25",
	})
	updateStream.Send(&ecommerce.Order{
		Id: "95",
	})

	res, _ := updateStream.CloseAndRecv()
	fmt.Printf("%v\n", res)

	fmt.Println("process orders")
	processStream, err := client.ProcessOrders(context.TODO())
	if err != nil {
		log.Fatalf("failed to process order: %v", err)
	}
	if err := processStream.Send(&wrapperspb.StringValue{
		Value: "first order",
	}); err != nil {
		log.Fatalf("failed to send 1 order to process: %v", err)
	}

	if err := processStream.Send(&wrapperspb.StringValue{
		Value: "second order",
	}); err != nil {
		log.Fatalf("failed to send 2 order to process: %v", err)
	}

	ch := make(chan struct{})

	go asyncBiDirRPC(processStream, ch)
	time.Sleep(time.Millisecond * 1000)

	if err := processStream.Send(&wrapperspb.StringValue{
		Value: "third order",
	}); err != nil {
		log.Fatalf("failed to send 2 order to process: %v", err)
	}
	if err := processStream.CloseSend(); err != nil {
		log.Fatal(err)
	}

	ch <- struct{}{}
}

func asyncBiDirRPC(stream ecommerce.OrderManagement_ProcessOrdersClient, ch chan struct{}) {
	for {
		combinedShipment, err := stream.Recv()

		if err == io.EOF {
			break
		}

		fmt.Printf("combined shipment: %v\n", combinedShipment)
	}

	<-ch
}
