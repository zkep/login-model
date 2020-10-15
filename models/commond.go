package models

import (
	"fmt"
	"github.com/urfave/cli"
)

type Commond struct {
	s *Selenium
}

func NewCommond() *Commond {
	s := NewSelenium(Port(PickUnusedPort()), Size(50, 50))
	c := &Commond{s: s}
	return c
}

func (self *Commond) Destroy() {
	self.s.Destroy()
}

// QQZone
func (self *Commond) QQZoneCmd() cli.Command {
	return cli.Command{
		Name:      "qqzone",
		Usage:     "qqzone login ",
		ArgsUsage: "FILE1 [FILE2] [FILE3]...",
		Action: func(c *cli.Context) {
			if c.NArg() < 2 {
				fmt.Fprintln(c.App.ErrWriter, "keys requires at least 1 argument")
				cli.ShowCommandHelp(c, "login")
				return
			}
			self.s.Start()
			qqzone := NewQQZone(self.s)
			qqzone.Login(c.Args().Get(0), c.Args().Get(1))
			qqzone.GET()
		},
	}
}


