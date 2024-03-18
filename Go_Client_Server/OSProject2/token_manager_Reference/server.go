package main

//import packages
import (
	"flag"
	"fmt"
	"net"

	"token_manager/token"
	"token_manager/utils"

	"google.golang.org/grpc"
)

func main() {
	port_ptr := flag.String("port", "50051", "port number which we will use")
	flag.Parse()

	fmt.Println("\nServer got started on the port", (*port_ptr), "\n")
	len, error := net.Listen("tcp", fmt.Sprintf(":%s", (*port_ptr)))
	utils.IsSuccess(error)

	s := token.Server{}
	server := grpc.NewServer()

	token.RegisterTokenServiceServer(server, &s)

	error = server.Serve(len)
	utils.IsSuccess(error)
}
