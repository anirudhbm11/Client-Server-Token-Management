# Client-Server-Token-Management

The Project folder contains src, bin and folders.
*  The bin folder consists the compiled and installed server and client go files.
* The src folder has 4 folders, namely client, server, YAMLFiles and Protos.

## Breakdown of folders
* Client folder: Consists of main.go which contains the code for spawning the client.
* Server folder: Consists of main.go which contains the code for spawning the server.
* Protos folder: Consists of messages.proto which is proto file and is used to generate 2 more proto files.
* YAMLFiles folder: This folder consists of the details of tokens which will be replicated and written. Each token has its own write server where to only that server it is written. Also, tokens have read server where the tokens are read from only those servers.

## To build and install client code:
* Go to /src/client folder
* go build
* go install

## To build and install server code:
* Go to /src/server folder
* go build
* go install

## To install dependencies:
Go mod init github.com/anirudhbm11/Client-Server-Token-Manager
There is 1 bash files in src/project3 folder,
* client_script.sh: Used to execute sequence of client commands.
Go to project3/client and execute “bash ../client_script.sh”
All the outputs of the replicated servers are stored in the log_file.txt.

## Algorithm:
* Create: During creating of the token, the token is first created in the write server and the token is replicated only to the read servers in that quorum. The servers which are not in the list of read and write will not receive the creation of token request.
* Read: During read request, one of the read server is randomly selected and is requested for the token. Before servicing the token, the read-impose-write majority algorithm is implemented where the server requests for the token details from the readers using go routines of that token.
Once the majority of the tokens in the quorum are received, the token with latest timestamp is selected and if the timestamp of the requested token is different from the latest timestamp, the token with new timestamp is broadcasted and replicated in all the servers of that token quorum.
* Write: All the token details are written to the write server and then the write server then broadcasts the token details to all read servers.
The servers are created using exec.Command() where the server process will be created in the background. After execution, kill the process with the name “main.go” to stop the servers.

## Fail silent:
It is simulated as
if *port == "8088" && in.GetId() == "100" && curr_time-token.State.stamp > 5 && curr_time-token.State.stamp < 20 \
So for port 8088 and token 100, if the read is done on 100 after 10 seconds and below 40 seconds after token creation, the server will not respond, and the client will timeout resulting in fail silent failure. Other tokens can be accessed during that time since we have the token id in the if condition. \
Also, as shown in the output log file, after creating 2 tokens, 1234 and 41, token 41 can be read after its creation. However, the token 1234 cannot be read since, I have simulated the fail-silent for that token not to be read until 20 seconds after creation. It can be read after sometime where we can see the output of the command in the logs and bash terminal. \
Similarly, write fail silent is also implemented in the code which results in the stale data in one of the replica servers and hence read-impose write all can be simulated if the read server which has a stale data is read.
