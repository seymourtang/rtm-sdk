package main

import (
	"k8s.io/klog/v2"

	"agora.io/rtm-sdk/cmd/options"
	"agora.io/rtm-sdk/internal/rtm"
	"agora.io/rtm-sdk/internal/utils/signal"
)

func main() {
	klog.InitFlags(nil)
	cmdOption := options.Parse()
	options := &rtm.OperatorOptions{
		AppID:  cmdOption.AppID,
		UserID: cmdOption.UserID,
		Token:  cmdOption.Token,
	}
	operator := rtm.New(options)
	errGroup, stop := signal.SetupStopSignalContext()
	errGroup.Go(func() error {
		return operator.Run(stop)
	})
	if err := errGroup.Wait(); err != nil {
		klog.Fatal(err)
	}
}
