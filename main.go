package main
import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/equinor/oauth2local/storage"

	"github.com/equinor/oauth2local/config"
	"github.com/equinor/oauth2local/ipc"
	"github.com/equinor/oauth2local/oauth2"
	"github.com/equinor/oauth2local/register"
)

var redirectCallback = flag.String("r", "", "Handles redirect from azure ad")

func runClient(cfg *config.Config) error {
	ipcClient, err := ipc.NewClient()
	if err != nil {
		return err
	}
	defer ipcClient.Close()
	log.Println("Running client", cfg.ClientType())
	switch cfg.ClientType() {
	case config.Redirect:
		err := ipcClient.SendCallback(cfg.RedirectURL)
		if err != nil {

			return err
		}
	case config.AccessToken:
		a, err := ipcClient.GetAccessToken()
		if err != nil {

			fmt.Print(a)
		}
	}
	return nil
}

func runServer(cfg *config.Config) error {

	err := cfg.LoadEnv()
	if err != nil {
		return fmt.Errorf("Couldn't load config: %v", err)

	}

	cli, err := oauth2.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("Error with oauth client: %v", err)
	}

	fmt.Println("starting browser...")
	cli.OpenLoginProvider()
	s := ipc.NewServer(*cli, &storage.Memory{})

	return s.Serve()
}

func main() {

	isClient := ipc.HasSovereign()

	flag.Parse()

	cfg := config.Init(isClient)
	cfg.RedirectURL = *redirectCallback
	if cfg.AsClient() {
		err := runClient(cfg)
		if err != nil {

			fmt.Println("Error...", err)
			bufio.NewReader(os.Stdin).ReadBytes('\n')
		}
		fmt.Println("Success...", err)
		bufio.NewReader(os.Stdin).ReadBytes('\n')

	} else {
		register.RegMe(cfg.HandleScheme, os.Args[0])
		log.Fatalf("Cannot serve: %v", runServer(cfg))
	}
}
