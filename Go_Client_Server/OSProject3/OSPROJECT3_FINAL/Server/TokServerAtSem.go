package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"sync"
	"time"

	pb "github.com/AmeyaBhave28UMBC/grpc-go/TokServAtSem"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"

	"github.com/golang/protobuf/ptypes/timestamp"
)

type server struct {
	pb.UnimplementedTokenServiceServer
}

type Token struct {
	id int32 //`json :"id"`
	//name   string  //`json :"name"`
	//domain *Domain //`json :"domain"`
	//state  *State  //`json :"state"`
	/* New Features added to the struct */
	value     string
	timestamp int64
	writer    string
	reader    []string
}

type tokenManager struct {
	Tokens []*Token
	mutex  *sync.Mutex
}

var tm = &tokenManager{
	Tokens: []*Token{},
	mutex:  &sync.Mutex{},
}

// ******************************************* New Additional Struct Types *******************************************

type TokenInfo struct {
	TokID   string   `yaml:"token"`
	Writer  string   `yaml:"writer"`
	Readers []string `yaml:"readers"`
}

type Config struct {
	Tokens []*TokenInfo `yaml:"tokens"`
}

// ******************************************* Create Tokens from the YAML File *******************************************

func mapYAMLtoTok(config Config) {
	//fmt.Println(*config.Tokens[1])

	for i := 0; i < len(config.Tokens); i++ {

		ts := &timestamp.Timestamp{
			Seconds: time.Now().Unix(),
			Nanos:   0,
		}

		// Convert the timestamp to an int64 value
		timestampInt := time.Unix(ts.Seconds, int64(ts.Nanos)).UnixNano()

		id64, err := strconv.ParseInt(config.Tokens[i].TokID, 10, 32)
		if err != nil {
			// Handle the error
			fmt.Println("Error:", err)
			return
		}
		tok_id := int32(id64)

		newToken := &Token{
			id:        tok_id,
			value:     fmt.Sprintf("%d", rand.Intn(20)),
			timestamp: timestampInt,
			writer:    config.Tokens[i].Writer,
			reader:    config.Tokens[i].Readers,
		}

		//fmt.Println(*newToken)
		tm.Tokens = append(tm.Tokens, newToken)
	}
	//fmt.Println(tm.Tokens[0].reader)
	//fmt.Println(tm.Tokens[1].id)
	//fmt.Println(tm.Tokens[2].writer)
}

// ******************************************* Call all Servers to Read and Return the token value *******************************************

func callReadServ(id int32, Readers string) (string, int64, error) {
	// Read all the read servers and call the Servers from the readers list as if you are the client
	var serverAddr string = Readers

	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure(), grpc.WithBlock()) //grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create a new gRPC client
	client := pb.NewTokenServiceClient(conn)

	// Set up a context for the RPC call inorder to maintain the Fail Silent model
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Create a new Read Request message
	readreq := &pb.ReadRequest{
		TokenId: id,
		Ack:     true,
	}

	//fmt.Println("Call Created and function executing")
	// Call the ReadToken RPC method on the server
	resp, err := client.Read(ctx, readreq)
	if err != nil {
		log.Fatalf("could not read token: %v", err)
	}

	// Return the values and timestamps to the contacted server
	return resp.GetValue(), resp.GetTimestamp(), nil
}

// ******************************************* Return the index of most recent Timestamp *******************************************
func Eval(timestamps []string) (int, error) {
	var max int64
	var index int
	max, err := strconv.ParseInt(timestamps[0], 10, 64)
	if err != nil {
		// Handle the error
		fmt.Println("Error:", err)
		return 0, err
	}
	for i := 1; i < len(timestamps); i++ {
		currentVal, err := strconv.ParseInt(timestamps[i], 10, 64)
		if err != nil {
			// Handle the error
			fmt.Println("Error:", err)
			return 0, err
		}
		if currentVal > max {
			max = currentVal
			index = i
		}
	}
	return index, nil
}

// ******************************************* Impose the Read value *******************************************

func ImposeReadsAll(id int32, mostRecentVal string, serverAddr string) (string, error) {
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

	writereq := &pb.WriteRequest{
		TokenId:     id,
		Value:       mostRecentVal,
		WriteServer: serverAddr,
		Ack:         true,
	}
	fmt.Println("Just before making the Write request")
	// Call the ReadToken RPC method on the server
	resp, err := client.Write(ctx, writereq)
	if err != nil {
		log.Fatalf("could not read token: %v", err)
	}
	fmt.Println("Token value being imposed after the Read quorum is " + resp.GetValue() + " on Server Address " + serverAddr)
	return resp.GetValue(), nil
}

// ******************************************* Impose the Write value *******************************************

func ImposeWriteAll(id int32, mostRecentVal string, serverAddr string) (string, error) {

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

	writereq := &pb.WriteRequest{
		TokenId:     id,
		Value:       mostRecentVal,
		WriteServer: serverAddr,
		Ack:         true,
	}

	//fmt.Println("Call Created and function executing")

	// Call the Write RPC method on the server
	resp, err := client.Write(ctx, writereq)
	if err != nil {
		log.Fatalf("could not read token: %v", err)
	}
	return resp.GetValue(), nil
}

// ******************************************* Functions for the Tokens *******************************************

func (s *server) Read(ctx context.Context, tok *pb.ReadRequest) (*pb.Token, error) {
	//fmt.Println(tm.Tokens[tok.GetTokenId()].value)

	// Get the Token ID from client for the Value to be returned
	var id = tok.GetTokenId()
	// Handle the condition where the Read function is returning the Value from each server to the Server Contacted by the Client
	var b bool = tok.GetAck()

	//Below condition will be used when the Server is trying to get all the read values from all the Readers of that token
	if b {
		return &pb.Token{
			Value:     tm.Tokens[id].value,
			Timestamp: tm.Tokens[id].timestamp,
		}, nil
	}

	// Initialise list for all the values and timestamps returned by each read server
	var timestamps []string
	var tmstmp int64
	var val string
	var values []string
	for i := 0; i < len(tm.Tokens[id].reader); i++ {

		// Call function to all the servers in  the readers' list
		val, tmstmp, _ = callReadServ(id, tm.Tokens[id].reader[i])

		//fmt.Println("Print after the callReadServ func")
		timestamps = append(timestamps, fmt.Sprint(tmstmp))
		values = append(values, val)
	}

	// Check if the Read Quorum condition is being fullfilled
	if len(timestamps) > ((len(tm.Tokens[id].reader) / 2) + 1) {

		// Call function to get the latest timestamp on the read value among all the servers
		mostRecentTmStpIndex, err := Eval(timestamps)

		//fmt.Println("Print after the Eval func")
		if err != nil {
			// Handle the error
			fmt.Println("Error:", err)
			return &pb.Token{
				Id: tm.Tokens[0].id}, errors.New("THERE WAS AN ERROR IN RETURNING THE TOKEN")
		}
		// Get the value associated with the latest timestamp
		mostRecentVal := values[mostRecentTmStpIndex]

		// Update the local value
		tm.Tokens[id].value = mostRecentVal
		tm.Tokens[id].timestamp, err = strconv.ParseInt(timestamps[mostRecentTmStpIndex], 10, 64)
		if err != nil {
			// Handle the error
			fmt.Println("Error:", err)
			return nil, errors.New("THERE WAS AN ERROR IN RETURNING THE TOKEN")
		}

		// Impose the latest Read value on all other servers that hold the token value
		var count int
		for i := 0; i < len(tm.Tokens[id].reader); i++ {
			//fmt.Println("Print before the ImposeReadsAll func")
			Acks, _ := ImposeReadsAll(tm.Tokens[id].id, mostRecentVal, tm.Tokens[id].reader[i])
			//fmt.Println("Print after the ImposeReadsAll func")
			if Acks == mostRecentVal {
				count++
			}
		}

		// CHECK IF THE IMPOSED VALUE ON ALL THE READ SERVERS SATISFY THE WRITE QUORUM
		if count >= (len(tm.Tokens[id].reader) / 2) {
			return &pb.Token{
				Value: tm.Tokens[id].value,
			}, nil
		}
		return nil, errors.New("THERE WAS AN ERROR WITH THE WRITE QUORUMS")
	} else {
		return nil, errors.New("THERE WAS AN ERROR WITH THE READ QUORUMS")
	}
}

func (s *server) Write(ctx context.Context, tok *pb.WriteRequest) (*pb.Token, error) {
	// Get the token id
	var id = tok.GetTokenId()

	// Get the token value and update the current value of the token
	tm.Tokens[id].value = tok.GetValue()

	// Create a new timestamp for
	ts := &timestamp.Timestamp{
		Seconds: time.Now().Unix(),
		Nanos:   0,
	}
	fmt.Println("The value which is being written is "+tm.Tokens[id].value+" on the Token ID ", tm.Tokens[id].id)

	// Convert the timestamp to an int64 value
	timestampInt := time.Unix(ts.Seconds, int64(ts.Nanos)).UnixNano()
	// Publish a new timestamp on the token as the value is updated
	tm.Tokens[id].timestamp = timestampInt

	// Handle the condition where the write function is called from the write server to all the read servers of that token
	// to update the value and timestamp of that token
	if tok.GetAck() {
		return &pb.Token{
			Id:          tm.Tokens[id].id,
			Value:       tm.Tokens[id].value,
			Timestamp:   tm.Tokens[id].timestamp,
			ReadServers: []string{},
			WriteServer: "",
		}, nil
	}

	var val string
	var values []string

	// Get the readers of the token
	// Invoke write requests to the readers in order for the readers to update their local token values and send ack

	for i := 0; i < len(tm.Tokens[id].reader); i++ {
		val, _ = ImposeWriteAll(tm.Tokens[id].id, tok.GetValue(), tm.Tokens[id].reader[i])
		values = append(values, val)
	}

	// Check the number of acknowledgements as part of Write Quorums
	var count int
	for j := 0; j < len(values); j++ {
		if values[j] == tok.GetValue() {
			count++
		}
	}
	if count >= (len(values) / 2) {
		// Return the token message backt to the client
		return &pb.Token{
			Id:        tm.Tokens[id].id,
			Value:     tm.Tokens[id].value,
			Timestamp: tm.Tokens[id].timestamp,
		}, nil
	} else {
		return nil, errors.New("THERE WAS AN ERROR WITH THE WRITE QUORUMS")
	}
}

// ******************************************* MAIN *******************************************

func main() {
	// var port string
	port := flag.String("port", "5005_", "the port number of the server")
	//ipaddr := flag.String("ipaddr", "127.0.0.1 , 127.0.0.2 , 127.0.0.3 , 127.0.0.4", "the port number of the server")

	// parse command line flags
	flag.Parse()

	yamlFile, err := os.ReadFile("Tok_IP_Port.yaml")
	if err != nil {
		panic(err)
	}
	fmt.Println("File read successfully")
	// parse the YAML file into a Config struct
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	// Call the function for implementing the replication scheme of the tokens from the YAMl file into the TokenManager
	mapYAMLtoTok(config)

	tokenManager := new(tokenManager)
	rpc.Register(tokenManager)
	// Creation of 1st Server
	if *port == "50051" {
		listener1, err := net.Listen("tcp", "127.0.0.1:"+*port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s1 := grpc.NewServer()
		x1 := &server{}
		pb.RegisterTokenServiceServer(s1, x1)
		fmt.Println("Server listening on ipaddr:port " + "127.0.0.1:" + *port)
		if err := s1.Serve(listener1); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}

	// Creation of 2nd Server
	if *port == "50052" {
		listener2, err := net.Listen("tcp", "127.0.0.1:"+*port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s2 := grpc.NewServer()
		x2 := &server{}
		pb.RegisterTokenServiceServer(s2, x2)
		fmt.Println("Server listening on ipaddr:port " + "127.0.0.1::" + *port)
		if err := s2.Serve(listener2); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}

	// Creation of 3rd Server
	if *port == "50053" {
		listener3, err := net.Listen("tcp", "127.0.0.1:"+*port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s3 := grpc.NewServer()
		x3 := &server{}
		pb.RegisterTokenServiceServer(s3, x3)
		fmt.Println("Server listening on ipaddr:port " + "127.0.0.1:" + *port)
		if err := s3.Serve(listener3); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}

}
