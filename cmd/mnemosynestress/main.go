package main

import (
	"fmt"
	"os"
	"strconv"

	"time"

	"github.com/piotrkowalczuk/mnemosyne/internal/discovery"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var config configuration

func main() {
	config.init()
	config.parse()

	pool, err := connect(config)
	if err != nil {
		fmt.Printf("mnemosyned connection failure: %v\n", err)
		os.Exit(1)
	}

	if len(pool) == 0 {
		fmt.Println("empty connection pool")
		os.Exit(1)
	}

	max := 131
	var sessions []*mnemosynerpc.Session
	for i, conn := range pool {
		for j := 0; j < max; j++ {
			res, err := conn.Start(context.Background(), &mnemosynerpc.StartRequest{
				Session: &mnemosynerpc.Session{
					SubjectId: strconv.Itoa(j),
				},
			})
			if err != nil {
				fmt.Printf("session creation error: %s\n", err.Error())
				os.Exit(1)
			}
			sessions = append(sessions, res.Session)

			if config.verbose {
				fmt.Printf("conn %d: session successfully created: %s\n", i+1, res.Session.AccessToken)
			}
		}
	}
	for i, ses := range sessions {
		for j, conn := range pool {
			ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
			got, err := conn.Get(ctx, &mnemosynerpc.GetRequest{
				AccessToken: ses.AccessToken,
			})
			if err != nil {
				fmt.Printf("session retrieval error: %s\n", err.Error())
				os.Exit(1)
			}
			if got.Session.AccessToken != ses.AccessToken {
				fmt.Printf("access token do not match, expected %s but got %s\n", ses.AccessToken, got.Session.AccessToken)
				os.Exit(1)
			}

			if config.verbose {
				fmt.Printf("conn %d: session %d/%d successfully retrieved: %s\n", j+1, len(sessions), i+1, ses.AccessToken)
			}
		}
	}
}

type service struct {
	ServiceAddress string
	ServicePort    int
}

func connect(c configuration) ([]mnemosynerpc.SessionManagerClient, error) {
	var (
		err       error
		addresses []string
	)

	switch {
	case c.cluster.static.enabled:
		addresses = c.cluster.static.members
	case c.cluster.discovery.enabled:
		addresses, err = discovery.DiscoverHTTP(c.cluster.discovery.http)
		if err != nil {
			return nil, err
		}
	}

	opts := []grpc.DialOption{
		grpc.WithUserAgent("mnemosynestress"),
	}
	if config.tls.enabled {
		creds, err := credentials.NewClientTLSFromFile(config.tls.cert, config.tls.key)
		if err != nil {
			return nil, err
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	clients := make([]mnemosynerpc.SessionManagerClient, 0, len(addresses))
	for _, addr := range addresses {
		if c.verbose {
			fmt.Printf("attempt to connect to: %s\n", addr)
		}
		conn, err := grpc.Dial(addr, opts...)
		if err != nil {
			return nil, err
		}

		clients = append(clients, mnemosynerpc.NewSessionManagerClient(conn))

		fmt.Printf("connection established: %s\n", addr)
	}

	return clients, nil
}
