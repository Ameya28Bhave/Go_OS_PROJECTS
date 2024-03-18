package main

//importing packages
import (
	"context"
	"flag"
	"fmt"

	"token_manager/token"
	"token_manager/utils"

	"google.golang.org/grpc"
)

// writing main function for client
func main() {
	//using port as given in pdf
	port_ptr := flag.String("port", "50051",
		"The number of the port on which the server is currently running")
	host_ptr := flag.String("host", "localhost",
		"The host where server is running at")

	create_ptr := flag.Bool("create", false, "set for creating token")
	drop_ptr := flag.Bool("drop", false, "set to drop the token")
	write_ptr := flag.Bool("write", false, "set to write the token")
	read_ptr := flag.Bool("read", false, "set to read the token")

	id_ptr := flag.String("id", "undefined", "id of the token")
	name_ptr := flag.String("name", "undefined", "name of the token")
	low_ptr := flag.Uint64("low", 1, "low value of the domain of token")
	mid_ptr := flag.Uint64("mid", 1, "mid value of the domain of token")
	high_ptr := flag.Uint64("high", 1, "high value of the domain of token")
	flag.Parse()

	var con *grpc.ClientConn
	con, error := grpc.Dial(fmt.Sprintf("%s:%s", *host_ptr, *port_ptr),
		grpc.WithInsecure())
	utils.IsSuccess(error)

	defer con.Close()

	cl := token.NewTokenServiceClient(con)

	req := token.Request{}
	req.Domain = &token.Request_Domain{}
	req.TokenState = &token.Request_State{}

	resp := &token.Response{}
	if *create_ptr {
		req.Id = *id_ptr
		resp, error = cl.Create(context.Background(), &req)
	} else if *drop_ptr {
		req.Id = *id_ptr
		resp, error = cl.Drop(context.Background(), &req)
	} else if *write_ptr {
		req.Id = *id_ptr
		req.Name = *name_ptr
		req.Domain.Low = *low_ptr
		req.Domain.Mid = *mid_ptr
		req.Domain.High = *high_ptr
		resp, error = cl.Write(context.Background(), &req)
	} else if *read_ptr {
		req.Id = *id_ptr
		resp, error = cl.Read(context.Background(), &req)
	}
	fmt.Println("Server Response Time:", (resp.GetBody()))
	utils.IsSuccess(error)

}
