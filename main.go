package main

import (
	"fmt"
	"os"

	"github.com/mdp/qrterminal/v3"
	"github.com/spf13/cobra"
)

type Flags struct {
	Port             int
	ReceivedFilesDir string
}

func main() {

	var flags Flags
	cmd := cobra.Command{
		Use: "receivefiles",
		Run: func(cmd *cobra.Command, args []string) {
			os.MkdirAll(flags.ReceivedFilesDir, os.ModePerm)
			url := "https://flores7.mniak.dev"
			fmt.Printf("Acesse %s para enviar arquivos \n ou leia o QR Code abaixo\n", url)
			showQRCode(url)

			var server Server
			server.run()
		},
	}
	cmd.Flags().IntVar(&flags.Port, "port", 10777, "HTTP port of the server")
	cmd.Flags().StringVar(&flags.ReceivedFilesDir, "save-to", "./uploads", "Where to save the received files")
	cobra.CheckErr(cmd.Execute())
}

func showQRCode(text string) {
	config := qrterminal.Config{
		HalfBlocks: true,
		Level:      qrterminal.M,
		Writer:     os.Stdout,
	}
	qrterminal.GenerateWithConfig(text, config)
}
