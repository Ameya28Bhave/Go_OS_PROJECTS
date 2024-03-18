1. Check the go.mod file as we initially need all the dependencies in the go.mod file. We can use : go mod tidy or : go get <packagename> 
    google.golang.org/grpc v1.54.0
    google.golang.org/protobuf v1.30.0
    github.com/AmeyaBhave28UMBC/grpc-go v0.0.0-20230510013047-71b7febc330b   
    github.com/go-yaml/yaml v2.1.0+incompatible
    github.com/golang/protobuf v1.5.2
    golang.org/x/net v0.8.0 // indirect
    golang.org/x/sys v0.6.0 // indirect
    golang.org/x/text v0.8.0 // indirect
    google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f // indirect
    gopkg.in/yaml.v2 v2.4.0

2. Execute the command to navigate to the Server Directory in your terminal : cd /OSPROJECT3_FINAL/Server 
3. Execute the command : go build
4. Clone the terminal for 2 more terminals as there are in total 3 servers being spawned
5. Execute the command : ./Server -port 50051 , this is for the first server
6. Execute the command : ./Server -port 50052 , in the first cloned terminal, this is for the second server
7. Execute the command : ./Server -port 50053 , in the second cloned terminal, this is for the third server
8. On a new terminal execute the command to navigate to the Client Directory in your terminal: cd /OSPROJECT3_FINAL/Client
9. For the client execute the command : go build
10. In order to execute the read command : ./Client -read -id 0 -host 127.0.0.1 -port 50051
11. In order to execute the write command : ./Client -write -id 0 -value 5 -host 127.0.0.1 -port 50051

NOTE : The Server Directory also contains the YAML file which is read by the server code to populate the list of Tokens in each Server
spawned.