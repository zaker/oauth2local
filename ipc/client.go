package ipc

import (
	"context"
	"fmt"
	"time"

	jww "github.com/spf13/jwalterweatherman"
	pb "github.com/zaker/oauth2local/ipc/localauth"
	"google.golang.org/grpc"
)

type Client struct {
	grpcConn *grpc.ClientConn
	inner    pb.LocalAuthClient
}

func NewClient() (c *Client, err error) {
	c = &Client{}
	c.grpcConn, err = grpc.Dial("pipe", grpc.WithInsecure(), grpc.WithContextDialer(localPipeDial))
	if err != nil {
		return nil, err
	}
	c.inner = pb.NewLocalAuthClient(c.grpcConn)
	return
}

func (c *Client) Close() {
	c.grpcConn.Close()
}

func HasSovereign() bool {
	c, err := NewClient()
	if err != nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.inner.Ping(ctx, new(pb.Empty))
	if err != nil {
		return false
	}
	return r.Message == "pong"
}

func (c *Client) SendCallback(callbackURL string) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	r, err := c.inner.Callback(ctx, &pb.CBRequest{Url: callbackURL})
	if err != nil {
		return err
	}
	if r == nil {
		return fmt.Errorf("No message received")
	}
	jww.DEBUG.Println("Sent callback:", callbackURL)
	return nil
}

func (c *Client) GetAccessToken() (string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	r, err := c.inner.GetAccessToken(ctx, &pb.Empty{})
	if err != nil {
		return "", err
	}
	jww.DEBUG.Println("Got access token")
	return r.AccessToken, nil
}
