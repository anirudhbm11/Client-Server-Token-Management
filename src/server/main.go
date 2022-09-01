package main

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/anirudhbm11/Client-Server-Token-Manager/Protos"
	"gopkg.in/yaml.v3"
)

type server struct {
	pb.UnimplementedTokenMgmtServer
}

type Token struct {
	id     string
	name   string
	Domain struct {
		low  uint64
		mid  uint64
		high uint64
	}
	State struct {
		partialvalue uint64
		finalvalue   uint64
		stamp        uint64
	}
}

type ImposeToken struct {
	Id         string
	Name       string
	Low        uint64
	Mid        uint64
	High       uint64
	stamp      uint64
	finalvalue uint64
}

type Tokenlist struct {
	Token   int64
	Writer  string
	Readers []string
}

var token_map = make(map[string]Token)

// Creating Read-Write Mutex
var rwLock sync.RWMutex

func printKeys() {
	// Printing token IDs currently in the server
	log.Println("List of IDs in server: ")
	for key, _ := range token_map {
		fmt.Printf("%s ", key)
		log.Printf(key)
	}
	fmt.Println("")
}

func createservers(reqtoken Tokenlist) {
	readservers := reqtoken.Readers

	for _, ips := range readservers {
		// Wg.Add(1)
		readport := strings.Split(ips, ":")[1]
		host := strings.Split(ips, ":")[0]

		conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, readport), time.Second)
		if err != nil {
			cmnd := exec.Command("go", "run", "../server/main.go", "-port", readport)
			cmnd.Start()
		}

		if conn != nil {
			defer conn.Close()
			fmt.Println("Server at port " + readport + " is already up and running")
		}
	}
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

func get_token_details(id string, tokendata []Tokenlist) Tokenlist {
	log.Println("Token to create/read: " + id)

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

func get_replica_servers(reqtoken Tokenlist) []string {
	read_servers := reqtoken.Readers

	return read_servers
}

func replicate_create(in *pb.CreateToken, replica_servers []string) {
	for _, server := range replica_servers {
		server_split := strings.Split(server, ":")
		port := server_split[1]
		host := server_split[0]
		replicate_helper(in, port, host)
	}
}

func replicate_helper(in *pb.CreateToken, port string, host string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	connection, err := grpc.Dial(host+":"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := pb.NewTokenMgmtClient(connection)

	// Populating the values in proto buf token for create rpc call
	token := &pb.CreateToken{
		Id:    in.GetId(),
		Stamp: in.GetStamp(),
	}

	// Creating and Sending Create request
	r, err := c.CreateReplicate(ctx, token)
	if err != nil {
		log.Fatalf("Could not replicate token with ID: " + in.GetId())
	}
	log.Printf("Replicating Token: %s", r.GetRes())
	connection.Close()
}

func (s *server) CreateReplicate(ctx context.Context, in *pb.CreateToken) (*pb.CreateReplicateResponse, error) {
	f, err := os.OpenFile("log_file.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	rwLock.Lock()
	defer rwLock.Unlock()
	fmt.Println("Token received for replication is: " + in.GetId())

	create_token := Token{
		id: in.GetId(),
	}

	create_token.State.stamp = in.GetStamp()

	token_map[in.GetId()] = create_token

	fmt.Println("Token details of replicating created token: ")
	log.Println(token_map)
	fmt.Println(create_token)

	printKeys()

	log.Println("Done Replicating token with ID: ", in.GetId())

	return &pb.CreateReplicateResponse{Res: "Replication Success" + " in port " + *port}, nil

}

func (s *server) CreateRequest(ctx context.Context, in *pb.CreateToken) (*pb.CreateResponse, error) {
	/* Function for accepting Create Request from the client. Here new tokens are created and are stored in token map*/

	// Write lock
	rwLock.Lock()
	defer rwLock.Unlock()

	f, err := os.OpenFile("log_file.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println("Token received for creating is: " + in.GetId())

	create_token := Token{
		id: in.GetId(),
	}
	create_token.State.stamp = in.GetStamp()

	token_map[in.GetId()] = create_token

	log.Println("Token details of created token: ")
	log.Println(token_map)

	printKeys()

	log.Println("Done creating token with ID: ", in.GetId())

	tokendata := readingyaml()

	reqtoken := get_token_details(in.GetId(), tokendata)

	createservers(reqtoken)

	time.Sleep(time.Second * 3)

	replica_servers := get_replica_servers(reqtoken)

	replicate_create(in, replica_servers)

	return &pb.CreateResponse{Res: "Success"}, nil
}

func get_timestamps(read_servers []string, write_server string, in *pb.ReadToken) ImposeToken {
	nservers := len(read_servers) + 1

	inchan := make(chan ImposeToken, 1)
	var maxtimestamp uint64
	maxtimestamp = 0
	var final_token ImposeToken
	count := 1

	go fetch_token(write_server, inchan, in)

	for i := 0; i < nservers-1; i++ {
		if count == (nservers/2)+1 {
			break
		}
		go fetch_token(read_servers[i], inchan, in)
		impose := <-inchan
		log.Printf("Final value of token is: ")
		log.Println(impose.finalvalue)
		log.Println(impose.stamp)
		log.Println(" in Port: " + read_servers[i])
		if impose.stamp > maxtimestamp {
			maxtimestamp = impose.stamp
			final_token = impose
		}
		count += 1
	}

	log.Println(final_token)
	return final_token
}

func fetch_token(server string, inchan chan ImposeToken, in *pb.ReadToken) {
	server_split := strings.Split(server, ":")
	port := server_split[1]
	host := server_split[0]

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	connection, err := grpc.Dial(host+":"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := pb.NewTokenMgmtClient(connection)

	// Populating the values in proto buf token for create rpc call
	token := &pb.ReadWriteToken{
		Id: in.GetId(),
	}

	// Creating and Sending Create request
	r, err := c.ReadWrite(ctx, token)
	if err != nil {
		log.Fatalf("Could not replicate token with ID: " + in.GetId())
	}
	connection.Close()

	impose := ImposeToken{
		Id:         r.GetId(),
		Name:       r.GetName(),
		Low:        r.GetLow(),
		Mid:        r.GetMid(),
		High:       r.GetHigh(),
		stamp:      r.GetStamp(),
		finalvalue: r.GetFinalvalue(),
	}

	inchan <- impose

}

func (s *server) ReadWrite(ctx context.Context, in *pb.ReadWriteToken) (*pb.ReadWriteResponse, error) {
	f, err := os.OpenFile("log_file.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	_, ok := token_map[in.GetId()]
	if ok {
		token := token_map[in.GetId()]
		stamp := token.State.stamp
		finalvalue := token.State.finalvalue

		impose := &pb.ReadWriteResponse{
			Id:         token.id,
			Name:       token.name,
			Low:        token.Domain.low,
			Mid:        token.Domain.mid,
			High:       token.Domain.high,
			Stamp:      stamp,
			Finalvalue: finalvalue,
		}
		return impose, nil
	} else {
		return &pb.ReadWriteResponse{}, errors.New("Token doesn't exist")
	}
}

func writeall(read_servers []string, write_server string, final_token ImposeToken) {
	servers := append(read_servers, write_server)
	for _, server := range servers {
		host := strings.Split(server, ":")[0]
		port := strings.Split(server, ":")[1]
		// Connecting to the server

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		connection, err := grpc.Dial(host+":"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Did not connectionect: %v", err)
		}
		defer connection.Close()
		c := pb.NewTokenMgmtClient(connection)

		// Populating the values in proto buf token for write rpc call
		token := &pb.WriteToken{
			Id:    final_token.Id,
			Name:  final_token.Name,
			Low:   final_token.Low,
			Mid:   final_token.Mid,
			High:  final_token.High,
			Stamp: final_token.stamp,
		}

		// Creating and Sending Write request
		r, err := c.WriteReplicate(ctx, token)
		if err != nil {
			log.Fatalf("Could not write the token with ID: " + final_token.Id)
		}
		log.Printf("Write back success: %d", r.GetRes())
	}
}

func (s *server) ReadRequest(ctx context.Context, in *pb.ReadToken) (*pb.ReadResponse, error) {
	/* Function for accepting Read Request from the client. Here requested tokens are read and final value is calculated
	and sent back to the client */

	// Read lock
	f, err := os.OpenFile("log_file.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	rwLock.RLock()
	defer rwLock.RUnlock()
	_, ok := token_map[in.GetId()]
	if ok {
		curr_time := uint64(time.Now().Unix())
		token := token_map[in.GetId()]
		min_final_val := token.State.finalvalue

		if *port == "8085" && in.GetId() == "1234" && curr_time-token.State.stamp > 5 && curr_time-token.State.stamp < 20 {
			time.Sleep(time.Second * 10)
			return &pb.ReadResponse{}, errors.New("Token reading failed")
		}
		fmt.Println("Reading the token with ID: ", in.GetId())

		tokendata := readingyaml()

		reqtoken := get_token_details(in.GetId(), tokendata)

		read_servers := get_replica_servers(reqtoken)
		write_server := reqtoken.Writer

		final_token := get_timestamps(read_servers, write_server, in)

		printKeys()
		fmt.Println("Reading done for token ID", in.GetId())

		if min_final_val != final_token.finalvalue {
			writeall(read_servers, write_server, final_token)
		}

		log.Println("\n")
		log.Println("\n")

		return &pb.ReadResponse{Res: final_token.finalvalue}, nil
	} else {
		return &pb.ReadResponse{}, errors.New("Token doesn't exist")
	}
}

func (s *server) DropRequest(ctx context.Context, in *pb.DropToken) (*pb.DropResponse, error) {
	/* Function for accepting Drop Request from the client. Here requested tokens are dropped from the map and
	the result is sent back to the client */

	f, err := os.OpenFile("log_file.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	// Write lock
	rwLock.Lock()
	defer rwLock.Unlock()
	_, ok := token_map[in.GetId()]

	if ok {
		fmt.Println("Deleting token: ", in.GetId())
		to_delete := token_map[in.GetId()]

		delete(token_map, in.GetId())

		fmt.Println("Token details of dropped token: ")

		fmt.Println(to_delete)

		fmt.Println("Deleted token is: " + in.GetId())

		printKeys()

		return &pb.DropResponse{Res: "Success"}, nil
	} else {
		return &pb.DropResponse{}, errors.New("Token doesn't exist")
	}

}

func replicating_writes(in *pb.WriteToken) {
	tokendata := readingyaml()

	reqtoken := get_token_details(in.GetId(), tokendata)

	replica_servers := get_replica_servers(reqtoken)

	replicate_writes_helper(replica_servers, in)

}

func replicate_writes_helper(replica_servers []string, in *pb.WriteToken) {
	for _, server := range replica_servers {
		host := strings.Split(server, ":")[0]
		port := strings.Split(server, ":")[1]
		// Connecting to the server

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		connection, err := grpc.Dial(host+":"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Did not connectionect: %v", err)
		}
		defer connection.Close()
		c := pb.NewTokenMgmtClient(connection)

		// Populating the values in proto buf token for write rpc call
		token := &pb.WriteToken{
			Id:    in.GetId(),
			Name:  in.GetName(),
			Low:   uint64(in.GetLow()),
			Mid:   uint64(in.GetMid()),
			High:  uint64(in.GetHigh()),
			Stamp: uint64(in.GetStamp()),
		}

		// Creating and Sending Write request
		r, err := c.WriteReplicate(ctx, token)
		if err != nil {
			log.Println("Could not write the token with ID: " + in.GetId())
		} else {
			log.Printf("Replication of write success: %d", r.GetRes())
		}
	}
}

func (s *server) WriteReplicate(ctx context.Context, in *pb.WriteToken) (*pb.WriteReplicateResponse, error) {
	/* Function for accepting Write Request from the client. Here requested tokens are retreived if already created and
	values are written to the token which is sent by the client*/

	// Write lock

	f, err := os.OpenFile("log_file.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	rwLock.Lock()
	defer rwLock.Unlock()
	token, ok := token_map[in.GetId()]
	if ok {
		curr_time := uint64(time.Now().Unix())
		if *port == "8088" && in.GetId() == "100" && curr_time-token.State.stamp > 10 && curr_time-token.State.stamp < 40 {
			time.Sleep(time.Second * 5)
			return &pb.WriteReplicateResponse{}, errors.New("Write Failed")
		}
		token.name = in.GetName()
		log.Println("Writing with name: " + token.name + " in port: " + *port)
		token.Domain.low = in.GetLow()
		token.Domain.high = in.GetHigh()
		token.Domain.mid = in.GetMid()
		token.State.stamp = in.GetStamp()

		var min_partial_val uint64
		var min_final_val uint64

		min_partial_val = math.MaxUint64
		min_final_val = math.MaxUint64

		for x := in.GetLow(); x < in.GetMid(); x++ {
			partial_val := Hash(in.GetName(), x)

			if partial_val < min_partial_val {
				min_partial_val = partial_val
			}
		}

		token.State.partialvalue = min_partial_val

		for x := in.GetMid(); x < in.GetHigh(); x++ {
			final_val := Hash(token.name, x)

			if final_val < min_final_val {
				min_final_val = final_val
			}
		}

		token.State.finalvalue = min_final_val

		token_map[in.GetId()] = token

		log.Println("Token details after write are: ")
		log.Println(token)

		printKeys()

		log.Println("Writing with name: ", token.name, " done.")

		return &pb.WriteReplicateResponse{Res: min_partial_val}, nil
	} else {
		return &pb.WriteReplicateResponse{}, errors.New("Token doesn't exist")
	}
}

func (s *server) WriteRequest(ctx context.Context, in *pb.WriteToken) (*pb.WriteResponse, error) {
	/* Function for accepting Write Request from the client. Here requested tokens are retreived if already created and
	values are written to the token which is sent by the client*/

	// Write lock

	f, err := os.OpenFile("log_file.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	rwLock.Lock()
	defer rwLock.Unlock()
	token, ok := token_map[in.GetId()]

	if ok {
		token.name = in.GetName()
		token.Domain.low = in.GetLow()
		token.Domain.high = in.GetHigh()
		token.Domain.mid = in.GetMid()
		token.State.stamp = in.GetStamp()

		var min_partial_val uint64
		var min_final_val uint64

		min_partial_val = math.MaxUint64
		min_final_val = math.MaxUint64

		for x := in.GetLow(); x < in.GetMid(); x++ {
			partial_val := Hash(in.GetName(), x)

			if partial_val < min_partial_val {
				min_partial_val = partial_val
			}
		}

		for x := in.GetMid(); x < in.GetHigh(); x++ {
			final_val := Hash(token.name, x)

			if final_val < min_final_val {
				min_final_val = final_val
			}
		}

		token.State.finalvalue = min_final_val
		token.State.partialvalue = min_partial_val

		token_map[in.GetId()] = token

		log.Println("Token details after write are: ")
		log.Println(token)

		printKeys()

		log.Println("Writing with name: ", token.name, " done ", " in port ", *port)
		replicating_writes(in)
		log.Println("\n")
		log.Println("\n")

		return &pb.WriteResponse{Res: min_partial_val}, nil
	} else {
		return &pb.WriteResponse{}, errors.New("Token doesn't exist")
	}
}

func Hash(name string, nonce uint64) uint64 {
	// Function for calculating the hash values of the name
	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("%s %d", name, nonce)))
	return binary.BigEndian.Uint64(hasher.Sum(nil))
}

var port *string

func main() {
	port = flag.String("port", "50000", "Port of server")

	// Port of the server is parsed
	flag.Parse()

	lis, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTokenMgmtServer(s, &server{})
	log.Printf("Server is listening at %v", lis.Addr())
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
