package watcher

import (
	"fmt"
	"github.com/mitchellh/cli"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type Command struct{}

func CommandFactory() (cli.Command, error) {
	return new(Command), nil
}

func (c *Command) Run(args []string) int {
	return NewService(args...).Run(args)
}

func (c *Command) Help() string {
	return help
}

func (c *Command) Synopsis() string {
	return synopsis
}

const synopsis = "watcher service"
const help = `
Usage: watcher -t=https://google.com -d=600 
`

type Service struct {
	larkUrl  string
	link     string
	duration time.Duration
}

func NewService(args ...string) *Service {
	if len(args) >= 2 {
		duration, _ := strconv.ParseInt(args[1], 10, 64)
		service := &Service{
			link:     args[0],
			duration: time.Duration(duration * int64(time.Second)),
		}
		service.initConfig()
		return service
	} else {
		return nil
	}
}

func (s *Service) Run(args []string) int {

	waitToStop := func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
		<-signalChan
	}

	stopAllService := func() {
		fmt.Println("all service has stopped, exit.")
	}

	go s.startWatcher()

	waitToStop()
	stopAllService()

	return 1
}

func (s *Service) initConfig() {
	var config struct {
		lark struct {
			robot string `json:"robot"`
		} `json:"lark"`
	}
	if content, err := ioutil.ReadFile("config.yaml"); err != nil {
		panic(err)
	} else {
		if err = yaml.Unmarshal(content, &config); err != nil {
			panic(err)
		}
		s.larkUrl = config.lark.robot
	}
}
