package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/djgoulart/fc2-grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect to gRPC server: %v", err)
	}

	defer connection.Close()

	client := pb.NewUserServiceClient(connection)
	//AddUser(client)
	//AddUserVerbose(client)
	//AddUsers(client)
	AddUsersStreamBoth(client)
}

// AddUser ...
func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Jose",
		Email: "jose@email.com",
	}

	res, err := client.AddUser(context.Background(), req)

	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	fmt.Println(res)
}

// AddUserVerbose ...
func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Jose",
		Email: "jose@email.com",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)

	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	for {
		stream, err := responseStream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Could not receive gRPC msg: %v", err)
		}

		fmt.Println("Status: ", stream.Status)
	}
}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		&pb.User{
			Id:    "d1",
			Name:  "Diego 1",
			Email: "email1@diego.com",
		},
		&pb.User{
			Id:    "d2",
			Name:  "Diego 2",
			Email: "email2@diego.com",
		},
		&pb.User{
			Id:    "d3",
			Name:  "Diego 3",
			Email: "email3@diego.com",
		},
		&pb.User{
			Id:    "d4",
			Name:  "Diego 4",
			Email: "email4@diego.com",
		},
		&pb.User{
			Id:    "d5",
			Name:  "Diego 5",
			Email: "email5@diego.com",
		},
	}

	stream, err := client.AddUsers(context.Background())

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	fmt.Println(res)
}

func AddUsersStreamBoth(client pb.UserServiceClient) {
	stream, err := client.AddUsersStreamBoth(context.Background())
	if err != nil {
		log.Fatalf("Errpr creating request %v", err)
	}

	reqs := []*pb.User{
		&pb.User{
			Id:    "d1",
			Name:  "Diego 1",
			Email: "email1@diego.com",
		},
		&pb.User{
			Id:    "d2",
			Name:  "Diego 2",
			Email: "email2@diego.com",
		},
		&pb.User{
			Id:    "d3",
			Name:  "Diego 3",
			Email: "email3@diego.com",
		},
		&pb.User{
			Id:    "d4",
			Name:  "Diego 4",
			Email: "email4@diego.com",
		},
		&pb.User{
			Id:    "d5",
			Name:  "Diego 5",
			Email: "email5@diego.com",
		},
	}

	wait := make(chan int)

	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user: ", req.Name)
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Error receiving data %v", err)
				break
			}
			fmt.Printf("Receiving user %v with status: %v\n", res.GetUser().GetName(), res.GetStatus())

		}
		close(wait)
	}()

	<-wait
}
