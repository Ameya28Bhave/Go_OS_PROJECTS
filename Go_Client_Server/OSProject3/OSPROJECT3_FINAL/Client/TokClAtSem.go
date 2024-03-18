package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	pb "github.com/AmeyaBhave28UMBC/grpc-go/TokServAtSem"

	"google.golang.org/grpc"
)

func readToken(serverAddr string, id string) {
	// Set up a connection to the server.
	// create a new insecure credentials object
	// creds := credentials.NewClientTLSFromCert(nil, "")

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

	//Convert Token ID into type int
	id64, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		// Handle the error
		fmt.Println("Error:", err)
		return
	}
	tok_id := int32(id64)

	// Create a new Read Request message

	readreq := &pb.ReadRequest{
		TokenId: tok_id,
		Ack:     false,
	}

	//fmt.Println("Call Created and function executing")
	// Call the ReadToken RPC method on the server
	resp, err := client.Read(ctx, readreq)
	if err != nil {
		log.Fatalf("could not read token: %v", err)
	}

	// Handle the response from the server
	fmt.Printf("Token read successfully. Response: %s\n", resp.GetValue())
}

func writeToken(serverAddr string, value string, id string) {
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

	id64, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		// Handle the error
		fmt.Println("Error:", err)
		return
	}
	tok_id := int32(id64)

	writereq := &pb.WriteRequest{
		TokenId:     tok_id,
		Value:       value,
		WriteServer: serverAddr,
		Ack:         false,
	}

	//fmt.Println("Call Created and function executing")
	// Call the ReadToken RPC method on the server
	resp, err := client.Write(ctx, writereq)
	if err != nil {
		log.Fatalf("could not read token: %v", err)
	}
	// Handle the response from the server
	fmt.Printf("Token written successfully. Response: %s\n", resp.GetValue())
}

func main() {
	// define command line flags

	//write := flag.Bool("write", false, "Write data in the token")
	read := flag.Bool("read", false, "Read data in the token")
	write := flag.Bool("write", false, "Write data in the token")
	id := flag.String("id", "", "the ID of the token to create")
	host := flag.String("host", "localhost", "the hostname of the server")
	port := flag.String("port", "50051", "the port number of the server")
	value := flag.String("value", "value", "the value to be written")
	//name := flag.String("name", "", "the Name of the token")
	//low := flag.Uint64("low", 0, "Lower bound of value")
	//mid := flag.Uint64("mid", 0, "Mid bound of value")
	//high := flag.Uint64("high", 0, "Higher bound of value")

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

	// build server address 127.0.0.1:50051
	serverAddr := fmt.Sprintf("%s:%s", *host, *port)

	// call the gRPC server with the provided parameters
	// Create a condition check for the request from the flags whether it asks for reading a token
	if *read {
		readToken(serverAddr, *id)
	}

	if *write {
		if *value == " " {
			log.Fatalf("error: the -value flag is required")
		}
		writeToken(serverAddr, *value, *id)
	}
}
