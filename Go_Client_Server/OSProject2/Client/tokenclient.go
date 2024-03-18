package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	pb "github.com/AmeyaBhave28UMBC/grpc-go/Token_Service"

	"google.golang.org/grpc"
)

// Function to create a gRPC call to Create a Token
func createToken(serverAddr string, id string) {
	// create a new insecure credentials object
	//creds := credentials.NewClientTLSFromCert(nil, "")
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure(), grpc.WithBlock()) //grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create a new gRPC client
	client := pb.NewTokenServiceClient(conn)

	// Set up a context for the RPC call
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Create a new Token message
	token := &pb.Token{
		Id: id,
	}

	// Call the CreateToken RPC method on the server
	resp, err := client.CreateToken(ctx, token)
	if err != nil {
		log.Fatalf("could not create token: %v", err)
	}

	// Handle the response from the server
	fmt.Printf("Token created successfully. Response: %s\n", resp.CrMessage)
}

// Function to create a gRPC call to Drop a Token
func dropToken(serverAddr string, id string) {
	// create a new insecure credentials object
	// creds := credentials.NewClientTLSFromCert(nil, "")
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create a new gRPC client
	client := pb.NewTokenServiceClient(conn)

	// Set up a context for the RPC call
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Create a new Token message
	token := &pb.Token{
		Id: id,
	}

	// Call the DropToken RPC method on the server
	resp, err := client.DropToken(ctx, token)
	if err != nil {
		log.Fatalf("could not drop token: %v", err)
	}

	// Handle the response from the server
	fmt.Printf("Token dropped successfully. Response: %s\n", resp.DrMessage)
}

// Function to create a gRPC call to Read in a Token

func readToken(serverAddr string, id string) {
	// Set up a connection to the server.
	// create a new insecure credentials object
	// creds := credentials.NewClientTLSFromCert(nil, "")
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create a new gRPC client
	client := pb.NewTokenServiceClient(conn)

	// Set up a context for the RPC call
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Create a new Token message
	token := &pb.Token{
		Id: id,
	}

	// Call the ReadToken RPC method on the server
	resp, err := client.ReadToken(ctx, token)
	if err != nil {
		log.Fatalf("could not read token: %v", err)
	}

	// Handle the response from the server
	fmt.Printf("Token read successfully. Response: %s\n", resp.RdMessage)
}

// Function to create a gRPC call to Write in a Token
func writeToken(serverAddr string, id string, name string, low uint64, mid uint64, high uint64) {
	// create a new insecure credentials object
	// creds := credentials.NewClientTLSFromCert(nil, "")
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create a new gRPC client
	client := pb.NewTokenServiceClient(conn)

	// Set up a context for the RPC call
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Create a new Token message
	token := &pb.Token{
		Id:   id,
		Name: name,
		Low:  low,
		Mid:  mid,
		High: high,
	}
	// Call the WriteToken RPC method on the server
	resp, err := client.WriteToken(ctx, token)
	if err != nil {
		log.Fatalf("could not write token: %v", err)
	}

	// Handle the response from the server
	fmt.Printf("Token written successfully. Response: %s\n", resp.WrMessage)

}

func main() {
	// define command line flags
	create := flag.Bool("create", false, "create a new token")
	write := flag.Bool("write", false, "Write data in the token")
	read := flag.Bool("read", false, "Read data in the token")
	drop := flag.Bool("drop", false, "Drop the token")
	id := flag.String("id", "", "the ID of the token to create")
	host := flag.String("host", "localhost", "the hostname of the server")
	port := flag.String("port", "50051", "the port number of the server")
	name := flag.String("name", "", "the Name of the token")
	low := flag.Uint64("low", 0, "Lower bound of value")
	mid := flag.Uint64("mid", 0, "Mid bound of value")
	high := flag.Uint64("high", 0, "Higher bound of value")

	// parse command line flags
	flag.Parse()

	// check if the id flag is set
	if *id == "" {
		log.Fatalf("error: the -id flag is required")
	}

	// check if the port flag is empty
	if *port == "" {
		log.Fatalf("error: the -port flag is required")
	}

	// check if the host flag is empty
	if *host == "" {
		log.Fatalf("error: the -host flag is required")
	}

	// build server address
	serverAddr := fmt.Sprintf("%s:%s", *host, *port)

	// call the gRPC server with the provided parameters
	// Create a condition check for the request from the flags whether it asks for creating a token
	if *create {
		createToken(serverAddr, *id)
	}

	// call the gRPC server with the provided parameters
	// Create a condition check for the request from the flags whether it asks for dropping a token
	if *drop {
		dropToken(serverAddr, *id)
	}

	// call the gRPC server with the provided parameters
	// Create a condition check for the request from the flags whether it asks for reading a token
	if *read {
		readToken(serverAddr, *id)
	}

	// call the gRPC server with the provided parameters
	// Create a condition check for the request from the flags whether it asks for writing a token
	if *write {
		writeToken(serverAddr, *id, *name, *low, *mid, *high)
	}
}
