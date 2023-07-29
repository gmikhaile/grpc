package ecommerce

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Server struct {
	OrderManagementServer
}

func (s *Server) ProcessOrders(stream OrderManagement_ProcessOrdersServer) error {
	for {
		orderID, err := stream.Recv()
		if err == io.EOF {
			stream.Send(
				&CombinedShipment{
					Id:     "shipment done",
					Status: "status done",
				},
			)
			return nil
		}

		if err != nil {
			return err
		}

		stream.Send(
			&CombinedShipment{
				Id:     fmt.Sprintf("shipment 1: %s", orderID),
				Status: "status 1",
			},
		)

		stream.Send(
			&CombinedShipment{
				Id:     fmt.Sprintf("shipment 2: %s", orderID),
				Status: "status 2",
			},
		)
	}
}

func (s *Server) UpdateOrders(stream OrderManagement_UpdateOrdersServer) error {
	for {
		order, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&wrapperspb.StringValue{Value: "done"})
		}

		fmt.Printf("recv: %v\n", order)
	}
}

func (s *Server) SearchOrders(query *wrapperspb.StringValue, stream OrderManagement_SearchOrdersServer) error {
	stream.Send(&Order{
		Id: "1",
		Items: []string{
			"laptop",
			"iphone",
		},
		Desc:        "two",
		Price:       1.12,
		Destination: "Bali",
	})

	stream.Send(&Order{
		Id: "2",
	})

	return nil
}

func (s *Server) GetOrder(context.Context, *wrapperspb.StringValue) (*Order, error) {
	return &Order{
		Id: "one",
		Items: []string{
			"laptopx2",
			"phonex2",
		},
		Desc:        "two",
		Price:       1.12,
		Destination: "Bali",
	}, nil
}
