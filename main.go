package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/zkep/login-model/models"
	"github.com/urfave/cli"
)

func main() {
	cmd := models.NewCommond()
	app := cli.NewApp()
	app.Name = "login"
	app.Usage = "a tool to login model"
	app.Version = "v0.0.1"
	app.Writer = os.Stdout
	app.ErrWriter = os.Stderr
	app.Commands = []cli.Command{
		cmd.QQZoneCmd(),
	}
	app.CommandNotFound = func(c *cli.Context, command string) {
		fmt.Fprintf(c.App.ErrWriter, "command %q can not be found.\n", command)
		cli.ShowAppHelp(c)
	}
	app.Run(os.Args)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM)
	<-ch
	cmd.Destroy()
}
