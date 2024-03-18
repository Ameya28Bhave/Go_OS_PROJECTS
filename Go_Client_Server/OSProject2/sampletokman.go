package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

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
}

var tm = &tokenManager{
	Tokens: []*Token{},
}

func Hash(name string, nonce uint64) uint64 {
	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("%s %d", name, nonce)))
	return binary.BigEndian.Uint64(hasher.Sum(nil))
}

func createtok() {
	newToken1 := &Token{
		id:     "ABC123",
		name:   "Ameya",
		domain: &Domain{low: 0, mid: 0, high: 0},
		state:  &State{partial_value: 0, final_value: 0}}
	newToken2 := &Token{
		id:     "ABCD1234",
		name:   "Rachit",
		domain: &Domain{low: 0, mid: 0, high: 0},
		state:  &State{partial_value: 0, final_value: 0}}
	tm.Tokens = append(tm.Tokens, newToken1)
	tm.Tokens = append(tm.Tokens, newToken2)
	for i := 0; i < 2; i++ {
		fmt.Println(tm.Tokens[i].id)
	}
}

func droptok() {
	result := []*Token{}
	//tm := &tokenManager{
	//	Tokens: []*Token{},
	//}
	for _, token := range tm.Tokens {
		fmt.Println(token.id)
	}

	for _, token := range tm.Tokens {
		if token.id != "ABC123" {
			result = append(result, token)
		}
	}
	for i := 0; i < 1; i++ {
		fmt.Println(result[i].id)
	}

	tm.Tokens = result
	for _, token := range tm.Tokens {
		if token.id == "ABC127" {
			fmt.Println("error")
		}

	}
	//fmt.Println(tm.Tokens[0])
}

func readtok() {
	var min uint64
	var temp uint64
	var t *Token

	for _, token := range tm.Tokens {
		if token.id == "ABC123" {
			t = token
		}
	}
	fmt.Println(t.id)

	if t == nil {
		fmt.Println("Error")
		//return nil, errors.New("error in reading the token")
	}

	temp = Hash(t.name, t.domain.mid)
	min = temp
	for i := t.domain.mid + 1; i < t.domain.high; i++ {
		temp = Hash(t.name, i)
		//fmt.Println(temp)
		if min > temp {
			min = temp
		}
	}
	if t.state.partial_value > min {
		t.state.final_value = min
	} else {
		t.state.final_value = t.state.partial_value
	}
	fmt.Println(uint64(t.state.final_value))
	//return &pb.ReadTokenResponse{RdMessage: strconv.FormatUint(uint64(t.state.final_value), 10)}, nil
}

func writetok() {
	var t *Token
	var temp uint64
	for _, token := range tm.Tokens {
		if token.id == "ABC123" {
			token.name = "Ameya"
			token.domain.low = 0
			token.domain.mid = 10
			token.domain.high = 100
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
	t.state.final_value = 0
	//fmt.Println(uint64(t.state.partial_value))
	fmt.Println(t.name)
}

func main() {
	createtok()
	//droptok()
	readtok()
	writetok()
	//fmt.Println(tm.Tokens)
}
