package ipc
import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/equinor/oauth2local/ipc/localauth"
	"github.com/equinor/oauth2local/oauth2"
	"github.com/equinor/oauth2local/storage"
	"google.golang.org/grpc"
)

type Server struct {
	oauthCli oauth2.Client
	store    storage.Storage
}

type client struct{}

func NewServer(cli oauth2.Client) (s *Server) {
	s = new(Server)
	s.oauthCli = cli
	return
}

func (s *Server) GetAccessToken(ctx context.Context, _ *pb.Empty) (*pb.ATResponse, error) {
	r := new(pb.ATResponse)
	c, err := s.store.GetCode()
	a, err := s.oauthCli.GetToken(c)
	r.AccessToken = a
	return r, err
}

func (s *Server) UpdateCode(ctx context.Context, cr *pb.UCRequest) (*pb.Empty, error) {
	r := new(pb.Empty)
	fmt.Println("Received:", cr.Code)
	return r, nil

}

func (s *Server) Ping(ctx context.Context, _ *pb.Empty) (*pb.PingResponse, error) {
	r := new(pb.PingResponse)
	r.Message = "pong"
	return r, nil
}

func HasSovereign() (bool, error) {
	conn, err := grpc.Dial("pipe", grpc.WithInsecure(), grpc.WithContextDialer(localPipeDial))
	if err != nil {
		return false, err
	}
	defer conn.Close()
	c := pb.NewLocalAuthClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Ping(ctx, new(pb.Empty))
	if err != nil {
		return false, err
	}
	return r.Message == "pong", nil
}

func SendCode(code string) error {
	conn, err := grpc.Dial("pipe", grpc.WithInsecure(), grpc.WithContextDialer(localPipeDial))
	if err != nil {
		return err
	}
	defer conn.Close()
	c := pb.NewLocalAuthClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	_, err = c.UpdateCode(ctx, &pb.UCRequest{Code: code})
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Run() {
	lis, err := listener()
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	gs := grpc.NewServer()
	pb.RegisterLocalAuthServer(gs, s)
	if err := gs.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
