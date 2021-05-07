package options

import (
	"flag"

	"github.com/gobike/envflag"
)

type CmdOption struct {
	AppID  string
	UserID string
	Token  string
}

func Parse() *CmdOption {
	var o CmdOption

	flag.StringVar(&o.AppID, "APP_ID", "", "the AppID of RTM instance")
	flag.StringVar(&o.UserID, "USER_ID", "", "the UserID of RTM instance")
	envflag.Parse()

	return &o
}
