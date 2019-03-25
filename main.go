package main
import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/zaker/oauth2local/oauth2"

	"github.com/zaker/oauth2local/config"

	"github.com/zaker/oauth2local/register"
)

var redirectCallback = flag.String("r", "", "Handles redirect from azure ad")

func main() {
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Println("Couldn't load config", err)
		fmt.Print("Press 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		log.Fatal(err)
	}

	cli, err := oauth2.NewClient(cfg)
	if err != nil {
		log.Println("Couldn't start client", err)
		fmt.Print("Press 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		log.Fatal(err)
	}

	if len(*redirectCallback) > 0 {
		log.Println("Handle redirect", *redirectCallback)
		code, err := cli.CodeFromCallback(*redirectCallback)
		if err != nil {
			log.Println("Couldn't retreive code from url", err)
			fmt.Print("Press 'Enter' to continue...")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			return
		}
		accessToken, err := cli.GetToken(code)
		if err != nil {
			log.Println("Error parsing url", err)
			fmt.Print("Press 'Enter' to continue...")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			return
		}
		fmt.Println("Access Token", accessToken)
		time.Sleep(time.Second * 10)
		fmt.Print("Press 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')

		return
	}

	if err != nil && !os.IsNotExist(err) {
		log.Fatal("Error loading .env file", err)
	}

	register.RegMe(cfg.HandleScheme, os.Args[0])
	cli.OpenLoginProvider()
}
