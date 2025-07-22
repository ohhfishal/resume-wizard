package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ohhfishal/resume-wizard/cmd"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	if err := cmd.Run(ctx, os.Stdout, os.Args[1:]); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
