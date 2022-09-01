package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/anirudhbm11/Client-Server-Token-Manager/Protos"

	"gopkg.in/yaml.v3"
)

type Tokenlist struct {
	Token   int64
	Writer  string
	Readers []string
}

func createservers(reqtoken Tokenlist) {

	writeserver := reqtoken.Writer

	writeport := strings.Split(writeserver, ":")[1]
	host := strings.Split(writeserver, ":")[0]
	fmt.Println(writeport)

	conn2, err := net.DialTimeout("tcp", net.JoinHostPort(host, writeport), time.Second)
	if err != nil {
		cmnd := exec.Command("go", "run", "../server/main.go", "-port", writeport)
		cmnd.Start()
	}

	if conn2 != nil {
		defer conn2.Close()
		fmt.Println("Server at port " + writeport + " is already up and running")
	}

}

func get_token_details(id string, tokendata []Tokenlist) Tokenlist {
	fmt.Println("Token to create: " + id)

	intid, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		log.Fatalf("Error in type conversion")
	}

	for _, token := range tokendata {
		if token.Token == intid {
			return token
		}
	}

	log.Fatalf("Wrong token requested")

	return Tokenlist{}

}

func readingyaml() []Tokenlist {
	yfile, err := ioutil.ReadFile("../YAMLFiles/commands.yml")

	if err != nil {
		log.Fatalf("Error reading Yaml file")
	}

	data := make([]Tokenlist, 10)

	err2 := yaml.Unmarshal(yfile, &data)

	if err2 != nil {
		log.Fatalf("Error in unmarshalling yaml file data")
	}

	return data
}

func selecting_server(reqtoken Tokenlist) string {
	read_servers := reqtoken.Readers

	randindex := rand.Intn(len(read_servers))
	fmt.Println(read_servers[randindex])

	return read_servers[randindex]
}

func main() {
	// Created a flag set for the use of RPC call input by the user.
	rpc_call := flag.NewFlagSet("rpc_call", flag.ContinueOnError)
	id := rpc_call.String("id", "", "ID of the token")
	name := rpc_call.String("name", "", "Name of the token")
	low := rpc_call.Int64("low", 0, "Low of the domain field in token")
	mid := rpc_call.Int64("mid", 0, "Mid of the domain field in token")
	high := rpc_call.Int64("high", 0, "High of the domain field in token")
	host := rpc_call.String("host", "localhost", "Host of the server")
	port := rpc_call.String("port", "50000", "Port of server")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	// Parsing flags after the RPC call
	err := rpc_call.Parse(os.Args[2:])

	if err != nil {
		log.Fatalf("Error in parsing")
	}

	tokendata := readingyaml()

	reqtoken := get_token_details(*id, tokendata)

	switch os.Args[1] {
	case "-create":
		createservers(reqtoken)

		time.Sleep(time.Second)

		createserver := reqtoken.Writer

		host := strings.Split(createserver, ":")[0]
		port := strings.Split(createserver, ":")[1]

		// Connecting to the server
		connection, err := grpc.Dial(host+":"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer connection.Close()
		c := pb.NewTokenMgmtClient(connection)

		// Populating the values in proto buf token for create rpc call
		token := &pb.CreateToken{
			Id:    *id,
			Stamp: uint64(time.Now().Unix()),
		}
		// Creating and Sending Create request
		r, err := c.CreateRequest(ctx, token)
		if err != nil {
			fmt.Println(err)
			log.Fatalf("Could not create token with ID: " + *id)
		}
		log.Printf("Creating Token: %s", r.GetRes())

	case "-write":
		writeserver := reqtoken.Writer

		host := strings.Split(writeserver, ":")[0]
		port := strings.Split(writeserver, ":")[1]
		// Connecting to the server
		connection, err := grpc.Dial(host+":"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Did not connectionect: %v", err)
		}
		defer connection.Close()
		c := pb.NewTokenMgmtClient(connection)

		// Populating the values in proto buf token for write rpc call
		token := &pb.WriteToken{
			Id:    *id,
			Name:  *name,
			Low:   uint64(*low),
			Mid:   uint64(*mid),
			High:  uint64(*high),
			Stamp: uint64(time.Now().Unix()),
		}

		// Creating and Sending Write request
		r, err := c.WriteRequest(ctx, token)
		if err != nil {
			log.Fatalf("Could not write the token with ID: " + *id)
		}
		log.Printf("Write success: %d", r.GetRes())
	case "-read":
		readserver := selecting_server(reqtoken)

		host := strings.Split(readserver, ":")[0]
		port := strings.Split(readserver, ":")[1]

		// Connecting to the server
		connection, err := grpc.Dial(host+":"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Did not connectionect: %v", err)
		}
		defer connection.Close()
		c := pb.NewTokenMgmtClient(connection)

		// Populating the values in proto buf token for read rpc call
		token := &pb.ReadToken{
			Id:    *id,
			Stamp: uint64(time.Now().Unix()),
		}

		// Creating and Sending Read request
		r, err := c.ReadRequest(ctx, token)
		if err != nil {
			log.Fatalf("Could not read the token with ID: " + *id)
		}
		log.Printf("Read Token: %d", r.GetRes())
	case "-drop":
		// Connecting to the server
		connection, err := grpc.Dial(*host+":"+*port, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Did not connectionect: %v", err)
		}
		defer connection.Close()
		c := pb.NewTokenMgmtClient(connection)

		// Populating the values in proto buf token for drop rpc call
		token := &pb.DropToken{
			Id:    *id,
			Stamp: uint64(time.Now().Unix()),
		}

		// Creating and Sending Drop request
		r, err := c.DropRequest(ctx, token)
		if err != nil {
			log.Fatalf("could not drop token with ID: " + *id)
		}
		log.Printf("Dropped token: %s", r.GetRes())
	}

}
