package main

import (
	"context"
	"log"
	"main/productinfo/client/ecommerce"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const address = "localhost:12345"

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := ecommerce.NewProductInfoClient(conn)

	name := "apple iphone 14"
	desc := "description"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.AddProduct(ctx, &ecommerce.Product{
		Name:        name,
		Description: desc,
	})
	if err != nil {
		log.Fatalf("failed to add product: %v", err)
	}

	product, err := client.GetProduct(ctx, &ecommerce.ProductID{
		Value: r.Value,
	})
	if err != nil {
		log.Fatalf("failed to get product: %v", err)
	}

	log.Printf("product: %s", product.String())
}
