package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kardianos/service"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func main() {
	program := NewProgram()

	cmd := cobra.Command{
		Use: "file-receiver",
		Run: func(cmd *cobra.Command, args []string) {
			if service.Interactive() {
				cobra.CheckErr(RunSimple(program))
			} else {
				cobra.CheckErr(RunAsService(program))
			}
		},
	}

	cmd.Flags().IntVarP(&program.HTTP.Port, "port", "p", 10777, "The server of the http server")
	cmd.Execute()
}

func RunAsService(program *Program) error {
	svcConfig := service.Config{}
	svc, err := service.New(program, &svcConfig)
	if err != nil {
		return err
	}
	return svc.Run()
}

func RunSimple(program *Program) error {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGTERM)

	svc := NotAService{}
	err := program.Start(&svc)
	if err != nil {
		return errors.WithMessage(err, "failed to start")
	}

	sig := <-sigchan
	fmt.Println("Stopping server due to signal", sig)
	err = program.Stop(&svc)
	return errors.WithMessage(err, "failed to stop")
}
