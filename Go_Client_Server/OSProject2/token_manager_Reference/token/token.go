package token

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"

	"token_manager/utils"
)

type Server struct {
	UnimplementedTokenServiceServer
}

type TokenType struct {
	request *Request
	mutex   *sync.Mutex
}

var tokenStore = make(map[string](TokenType))
var queryNumber uint64 = 0

func GetTokenState(t *Request) string {
	tStr := fmt.Sprintf("{id: %s, name: %s, "+
		"domain: {low: %d, mid: %d, high: %d}, "+
		"state: {partial_val: %d, final_val: %d} }",
		t.GetId(), t.GetName(), t.GetDomain().GetLow(),
		t.GetDomain().GetMid(), t.GetDomain().GetHigh(),
		t.GetTokenState().GetPartialval(), t.GetTokenState().GetFinalval())
	return tStr
}

func (s *Server) Create(ctx context.Context, req *Request) (*Response, error) {
	queryNumber += 1
	currQueryNumber := queryNumber
	reqQuery := fmt.Sprintf("{Action: create, Id: %s}", (req.GetId()))
	fmt.Println("Request Received --> Request Number: ",
		currQueryNumber, ", Request: ", reqQuery)
	_, keyExists := tokenStore[req.GetId()]
	if keyExists {
		fmt.Println("Processed request number ", currQueryNumber, ", error occured")
		return &Response{Body: "Token is already present"},
			errors.New("token creation failed")
	}
	tokenStore[req.GetId()] = TokenType{request: req, mutex: &sync.Mutex{}}
	fmt.Println("Processed request number ", currQueryNumber,
		", Token State Now: ", GetTokenState(tokenStore[req.GetId()].request))
	fmt.Println("Tokenstore contains: ", reflect.ValueOf(tokenStore).MapKeys())
	return &Response{Body: fmt.Sprintf("Token created with id: %s",
		(req.GetId()))}, nil
}

func (s *Server) Drop(ctx context.Context, req *Request) (*Response, error) {
	queryNumber += 1
	currQueryNumber := queryNumber
	reqQuery := fmt.Sprintf("{Action: drop, Id: %s}", (req.GetId()))
	fmt.Println("Request Received --> Request Number: ",
		currQueryNumber, ", Request: ", reqQuery)
	_, keyExists := tokenStore[req.GetId()]
	if !keyExists {
		fmt.Println("Processed request number ", currQueryNumber, ", error occured")
		return &Response{Body: "Token is absent, nothing to delete"},
			errors.New("token drop failed")
	}
	tokenStore[req.GetId()].mutex.Lock()
	defer tokenStore[req.GetId()].mutex.Unlock()
	delete(tokenStore, req.GetId())
	fmt.Println("Processed request number ", currQueryNumber,
		", Token State Now: ", GetTokenState(tokenStore[req.GetId()].request))
	fmt.Println("Tokenstore contains: ", reflect.ValueOf(tokenStore).MapKeys())
	return &Response{Body: fmt.Sprintf("Token dropped with id: %s",
		(req.GetId()))}, nil
}

func (s *Server) Write(ctx context.Context, req *Request) (*Response, error) {
	queryNumber += 1
	currQueryNumber := queryNumber
	reqQuery := fmt.Sprintf(
		"{Action: write, Id: %s, Name: %s, Low: %d, Mid: %d, High: %d}",
		req.GetId(), req.GetName(), req.GetDomain().GetLow(),
		req.GetDomain().GetMid(), req.GetDomain().GetHigh())
	fmt.Println("Request Received --> Request Number: ",
		currQueryNumber, ", Request: ", reqQuery)
	val, keyExists := tokenStore[req.GetId()]
	if !keyExists {
		fmt.Println("Processed request number ", currQueryNumber,
			", error occured")
		return &Response{Body: "Token is not available"},
			errors.New("token write failed")
	}
	tokenStore[req.GetId()].mutex.Lock()
	defer tokenStore[req.GetId()].mutex.Unlock()
	val.request.Name = req.GetName()
	val.request.Domain = &Request_Domain{
		Low:  req.GetDomain().GetLow(),
		Mid:  req.GetDomain().GetMid(),
		High: req.GetDomain().GetHigh(),
	}
	val.request.TokenState = &Request_State{
		Partialval: utils.FindArgminxHash(
			val.request.GetName(), val.request.GetDomain().GetLow(),
			val.request.GetDomain().GetMid()),
		Finalval: 0,
	}
	tokenStore[req.GetId()] = val
	fmt.Println("Processed request number ", currQueryNumber,
		", Token State Now: ", GetTokenState(tokenStore[req.GetId()].request))
	fmt.Println("Tokenstore contains: ", reflect.ValueOf(tokenStore).MapKeys())
	return &Response{Body: fmt.Sprintf(
		"Token updated with partial value: %d",
		(tokenStore[req.GetId()].request.GetTokenState().GetPartialval()))}, nil
}

func (s *Server) Read(ctx context.Context, req *Request) (*Response, error) {
	queryNumber += 1
	currQueryNumber := queryNumber
	reqQuery := fmt.Sprintf("{Action: read, Id: %s}", (req.GetId()))
	fmt.Println(
		"Request Received --> Request Number: ",
		currQueryNumber, ", Request: ", reqQuery)
	val, keyExists := tokenStore[req.GetId()]
	if !keyExists {
		fmt.Println("Processed request number ", currQueryNumber,
			", error occured")
		return &Response{Body: "Token is not available"},
			errors.New("token read failed")
	}
	tokenStore[req.GetId()].mutex.Lock()
	defer tokenStore[req.GetId()].mutex.Unlock()
	minMidHigh := utils.FindArgminxHash(
		val.request.GetName(), val.request.GetDomain().GetMid(),
		val.request.GetDomain().GetHigh())
	if minMidHigh < val.request.GetTokenState().GetPartialval() {
		val.request.TokenState.Finalval = minMidHigh
	} else {
		val.request.TokenState.Finalval = val.request.GetTokenState().GetPartialval()
	}
	tokenStore[req.GetId()] = val
	fmt.Println("Processed request number ", currQueryNumber,
		", Token State Now: ", GetTokenState(tokenStore[req.GetId()].request))
	fmt.Println("Tokenstore contains: ", reflect.ValueOf(tokenStore).MapKeys())
	return &Response{Body: fmt.Sprintf(
		"Token updated with final value: %d",
		(tokenStore[req.GetId()].request.GetTokenState().GetFinalval()))}, nil
}
