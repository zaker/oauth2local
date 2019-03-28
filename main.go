package main
import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/equinor/oauth2local/config"
	"github.com/equinor/oauth2local/ipc"
	"github.com/equinor/oauth2local/oauth2"
	"github.com/equinor/oauth2local/register"
)

var redirectCallback = flag.String("r", "", "Handles redirect from azure ad")

func client(cfg *config.Config) error {

	switch cfg.ClientType() {
	case config.Redirect:
		err := ipc.SendCode(cfg.RedirectURL)
		if err != nil {

			return err
		}
		fmt.Print("Press 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
	return nil
}

func main() {

	isClient, err := ipc.HasSovereign()

	if err != nil {
		log.Println("Sovereign check failed", err)

	}
	flag.Parse()

	cfg := config.Init(isClient)
	cfg.RedirectURL = *redirectCallback
	if isClient {
		err = client(cfg)
		if err != nil {

			fmt.Println("Error...", err)
			bufio.NewReader(os.Stdin).ReadBytes('\n')
		}
		return
	}

	err = cfg.LoadEnv()
	if err != nil {
		log.Println("Couldn't load config", err)
		fmt.Print("Press 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		log.Fatal(err)
	}
	// cli, _ := oauth2.NewClient(cfg)
	// cli.OpenLoginProvider()

	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	cli, err := oauth2.NewClient(cfg)
	if err != nil {
		log.Println("Couldn't start client", err)
		fmt.Print("Press 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		log.Fatal(err)
	}

	// if isRedirect {
	// 	// log.Println("Handle redirect", *redirectCallback)
	// 	code, err := oauth2.CodeFromURL(*redirectCallback, cfg.HandleScheme)
	// 	if err != nil {
	// 		log.Println("Couldn't retreive code from url", err)
	// 		fmt.Print("Press 'Enter' to continue...")
	// 		bufio.NewReader(os.Stdin).ReadBytes('\n')
	// 		return
	// 	}
	// 	accessToken, err := cli.GetToken(code)
	// 	if err != nil {
	// 		log.Println("Error parsing url", err)
	// 		fmt.Print("Press 'Enter' to continue...")
	// 		bufio.NewReader(os.Stdin).ReadBytes('\n')
	// 		return
	// 	}
	// 	fmt.Println("Access Token", accessToken)
	// 	fmt.Print("Press 'Enter' to continue...")
	// 	bufio.NewReader(os.Stdin).ReadBytes('\n')

	// 	return
	// }

	register.RegMe(cfg.HandleScheme, os.Args[0])
	fmt.Println("starting browser...")
	s := ipc.NewServer(*cli)
	s.Run()
	return
}
