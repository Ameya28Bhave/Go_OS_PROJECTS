package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/go-yaml/yaml"
	"github.com/golang/protobuf/ptypes/timestamp"
)

/*
type AccessPoint struct {
	ipAddressPort string
	//port      string
}

type Domain struct {
	low  uint64
	mid  uint64
	high uint64
}

type State struct {
	partial_value uint64
	final_value   uint64
}
*/

type Token struct {
	id    string //`json :"id"`
	value string //`json :"name"`
	//domain *Domain //`json :"domain"`
	//state  *State  //`json :"state"`
	timestamp int64    //*timestamp.Timestamp
	writer    string   //*AccessPoint
	reader    []string //*AccessPoint
}

type tokenManager struct {
	Tokens []*Token
}

type TokenInfo struct {
	TokID   string   `yaml:"token"`
	Writer  string   `yaml:"writer"`
	Readers []string `yaml:"readers"`
}

type Config struct {
	Tokens []*TokenInfo `yaml:"tokens"`
}

var tm = &tokenManager{
	Tokens: []*Token{},
}

/*
func createtoken() {
	var acpt = &AccessPoint{
		ipAddress: "",
		port:      "",
	}
	var acptlist []*AccessPoint
	ipPrefix := "1.1.1."
	port := "8080"
	//var acpt *AccessPoint
	for i := 1; i < 5; i++ {
		ip := string(ipPrefix + fmt.Sprintf("%d", i))
		acpt = &AccessPoint{
			ipAddress: ip,
			port:      port,
		}
		acptlist = append(acptlist, acpt)
	}
	//fmt.Println(acptlist)

		newToken1 := &Token{
			id:     "ABC123",
			name:   "Ameya",
			domain: &Domain{low: 0, mid: 0, high: 0},
			state:  &State{partial_value: 0, final_value: 0},
			writer: &AccessPoint{ipAddress: "0.0.0.0", port: "50051"},
			reader: acptlist,
		}
		fmt.Println(*newToken1.reader[0])

}
*/

func mapYAMLtoTok(config Config) {
	fmt.Println(*config.Tokens[1])
	/*
		var acpt = &AccessPoint{
			ipAddressPort: "",
		}
		var acptlist []*AccessPoint
	*/
	for i := 0; i < len(config.Tokens); i++ {
		//acptlist = nil
		/*
			for j := 0; j < len(config.Tokens[i].Readers); j++ {
				acpt = &AccessPoint{
					ipAddressPort: config.Tokens[i].Readers[j],
				}
				acptlist = append(acptlist, acpt)
			}
		*/

		// Create a new timestamp with the current time
		ts := &timestamp.Timestamp{
			Seconds: time.Now().Unix(),
			Nanos:   0,
		}

		// Convert the timestamp to an int64 value
		timestampInt := time.Unix(ts.Seconds, int64(ts.Nanos)).UnixNano()

		// Print the int64 value
		fmt.Println(timestampInt)

		newToken := &Token{
			id:        config.Tokens[i].TokID,
			value:     fmt.Sprintf("%d", rand.Intn(20)),
			timestamp: timestampInt,
			writer:    config.Tokens[i].Writer,
			reader:    config.Tokens[i].Readers,
		}

		fmt.Println(newToken)
		tm.Tokens = append(tm.Tokens, newToken)
	}
	fmt.Println(tm.Tokens[0].reader)
	fmt.Println(tm.Tokens[1].id)
	fmt.Println(tm.Tokens[2].timestamp)
	fmt.Println(tm.Tokens[2].writer)
}

func main() {
	//var yamlFile string
	// read the YAML file
	yamlFile, err := os.ReadFile("Tok_IP_Port.yaml")
	if err != nil {
		panic(err)
	}

	// parse the YAML file into a Config struct
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	//fmt.Println(config.Tokens[0].Readers)
	//fmt.Println(len(config.Tokens))
	mapYAMLtoTok(config)
	//createtok()
	//droptok()
	//readtok()
	//writetok()
	//fmt.Println(tm.Tokens)
}
