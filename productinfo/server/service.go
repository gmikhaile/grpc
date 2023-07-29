package main

import (
	"context"
	"main/productinfo/server/ecommerce"

	"github.com/gofrs/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	ecommerce.UnimplementedProductInfoServer
	productMap map[string]*ecommerce.Product
}

func (s *Server) AddProduct(ctx context.Context, in *ecommerce.Product) (*ecommerce.ProductID, error) {
	out, err := uuid.NewV4()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate product id: %v", err)
	}

	in.Id = out.String()
	if s.productMap == nil {
		s.productMap = make(map[string]*ecommerce.Product)
	}

	s.productMap[in.Id] = in

	return &ecommerce.ProductID{Value: in.Id}, status.New(codes.OK, "").Err()
}

func (s *Server) GetProduct(ctx context.Context, in *ecommerce.ProductID) (*ecommerce.Product, error) {
	if value, ok := s.productMap[in.Value]; ok {
		return value, status.New(codes.OK, "").Err()
	}

	return nil, status.Errorf(codes.NotFound, "product doesn't exist %s", in.Value)
}
