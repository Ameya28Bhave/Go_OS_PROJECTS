package main

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"strconv"
	"sync"

	pb "github.com/AmeyaBhave28UMBC/grpc-go/Token_Service"
	"google.golang.org/grpc"
)

// ******************************************* Server ADT *******************************************

type server struct {
	pb.UnimplementedTokenServiceServer
}

// ******************************************* Token ADT *******************************************

type Domain struct {
	low  uint64
	mid  uint64
	high uint64
}

type State struct {
	partial_value uint64
	final_value   uint64
}

type Token struct {
	id     string  //`json :"id"`
	name   string  //`json :"name"`
	domain *Domain //`json :"domain"`
	state  *State  //`json :"state"`
}

type tokenManager struct {
	Tokens []*Token
	mutex  *sync.Mutex
}

var tm = &tokenManager{
	Tokens: []*Token{},
	mutex:  &sync.Mutex{},
}

// ******************************************* Hash Function *******************************************

func Hash(name string, nonce uint64) uint64 {
	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("%s %d", name, nonce)))
	return binary.BigEndian.Uint64(hasher.Sum(nil))
}

// ******************************************* Functions for the Tokens *******************************************

func (s *server) CreateToken(ctx context.Context, tok *pb.Token) (*pb.CreateTokenResponse, error) {
	// Lock the mutex to ensure exclusive access to the token list
	//tm.mutex.Lock()
	//defer tm.mutex.Unlock()
	result := "Successfull created the token"
	// Create a new token
	token := &Token{
		id:     tok.GetId(),
		name:   "",
		domain: &Domain{low: 0, mid: 0, high: 0},
		state:  &State{partial_value: 0, final_value: 0},
	}

	//var tm *tokenManager
	// check if the Tokens slice is nil
	// adding the newly created token to our array []*Token
	tm.Tokens = append(tm.Tokens, token)
	//tokenStore[tok.GetId()] = tokenManager{Tokens: token}
	// Send a success response message back to the Client
	return &pb.CreateTokenResponse{CrMessage: result}, nil
}

func (s *server) DropToken(ctx context.Context, tok *pb.Token) (*pb.DropTokenResponse, error) {
	// Lock the mutex to ensure exclusive access to the token list
	tm.mutex.Lock()
	defer tm.mutex.Unlock()
	res := "Successfully dropped the token"
	// Drop the token by creating a new []*Token object where we store all the tokens
	// from our TokenManager excpet the one for the id we want to drop and then assign
	// the replace the old array in the TokenManager struct with the newly created array
	result := []*Token{}
	//var dtm *tokenManager
	//dtm := &tokenManager{
	//	Tokens: []*Token{},
	//}
	//tm := &tokenManager{}

	for _, token := range tm.Tokens {
		if token.id != tok.GetId() {
			result = append(result, token)
		}
	}
	tm.Tokens = result
	for _, token := range tm.Tokens {
		if token.id == tok.GetId() {
			return nil, errors.New("error in removing the tokens")
		}
	}
	return &pb.DropTokenResponse{DrMessage: res}, nil
}

func (s *server) ReadToken(ctx context.Context, tok *pb.Token) (*pb.ReadTokenResponse, error) {
	// Lock the mutex to ensure exclusive access to the token list
	tm.mutex.Lock()
	defer tm.mutex.Unlock()
	// Read the token
	//var rtm *tokenManager
	var min uint64
	var temp uint64
	var t *Token

	for _, token := range tm.Tokens {
		if token.id == tok.GetId() {
			t = token
		}
	}

	//fmt.Println(t.id)

	if t == nil {
		return nil, errors.New("error in reading the token")
	}

	temp = Hash(t.name, t.domain.mid)
	min = temp
	for i := t.domain.mid + 1; i < t.domain.high; i++ {
		temp = Hash(t.name, i)
		if min > temp {
			min = temp
		}
	}
	if t.state.partial_value > min {
		t.state.final_value = min
	} else {
		t.state.final_value = t.state.partial_value
	}

	return &pb.ReadTokenResponse{RdMessage: strconv.FormatUint(uint64(t.state.final_value), 10)}, nil
}

func (s *server) WriteToken(ctx context.Context, tok *pb.Token) (*pb.WriteTokenResponse, error) {
	//var wrtm *tokenManager
	// Lock the mutex to ensure exclusive access to the token list
	tm.mutex.Lock()
	defer tm.mutex.Unlock()
	var t *Token
	var temp uint64
	for _, token := range tm.Tokens {
		if token.id == tok.GetId() {
			token.name = tok.GetName()
			token.domain.low = tok.GetLow()
			token.domain.mid = tok.GetMid()
			token.domain.high = tok.GetHigh()
			t = token
		}
	}
	t.state.partial_value = Hash(t.name, t.domain.low)
	for i := t.domain.low + 1; i < t.domain.mid; i++ {
		temp = Hash(t.name, i)
		if t.state.partial_value > temp {
			t.state.partial_value = temp
		}
	}
	t.state.final_value = 0 // resetting final value

	return &pb.WriteTokenResponse{WrMessage: strconv.FormatUint(uint64(t.state.partial_value), 10)}, nil
}

// ******************************************* Main Function Call *******************************************

func main() {
	// var port string
	port := flag.String("port", "50051", "the port number of the server")

	// parse command line flags
	flag.Parse()

	tokenManager := new(tokenManager)
	rpc.Register(tokenManager)
	listener, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	x := &server{}
	pb.RegisterTokenServiceServer(s, x)
	fmt.Printf("Server listening on port " + *port)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
