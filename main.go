package main

import (
	"crypto/tls"
	"os"
	"fmt"

	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/auth"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/client"
	pb "github.com/stargate/stargate-grpc-go-client/stargate/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var stargateClient *client.StargateClient

func main() {
    //Authenticated Connection
	/*grpcEndpoint := "localhost:8090"
	authEndpoint := "localhost:8081"
	username := "cassandra"
	passwd := "cassandra"

	conn, err := grpc.Dial(grpcEndpoint, grpc.WithInsecure(), grpc.WithBlock(),
		grpc.WithPerRPCCredentials(
			auth.NewTableBasedTokenProviderUnsafe(
				fmt.Sprintf("http://%s/v1/auth", authEndpoint), username, passwd,
			),
		),
	)
	*/
	// Astra DB configuration
	const astra_uri = "$ASTRA_CLUSTER_ID-$ASTRA_REGION.apps.astra.datastax.com:443"
	const bearer_token = "AstraCS:xxxxx"

	config := &tls.Config{
		InsecureSkipVerify: false,
	}

	fmt.Printf("before grpc conn\n")

	conn, err := grpc.Dial(astra_uri, grpc.WithTransportCredentials(credentials.NewTLS(config)),
		grpc.WithBlock(),
		grpc.WithPerRPCCredentials(
			auth.NewStaticTokenProvider(bearer_token),
		),
	)

	fmt.Print("Post connection : ", err, "\n")

	stargateClient, err = client.NewStargateClientWithConn(conn)

	if err != nil {
		fmt.Printf("error creating client %v", err)
		os.Exit(1)
	}

	selectQuery := &pb.Query{
		Cql: "SELECT * FROM user_details.user;",
	}

	response, err := stargateClient.ExecuteQuery(selectQuery)
	if err != nil {
		fmt.Printf("error executing query %v", err)
		return
	}
	fmt.Printf("select executed\n")

	result := response.GetResultSet()
	fmt.Printf("result: %v \n", result)

}
