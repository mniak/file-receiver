package main

import (
	"fmt"
	"os"
	"path/filepath"

	"receivefiles"

	"github.com/mdp/qrterminal/v3"
	"github.com/spf13/cobra"
)

func main() {
	var p receivefiles.AppParams
	cmd := cobra.Command{
		Use: "receivefiles",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			p.ReceivedFilesDir, err = filepath.Abs(p.ReceivedFilesDir)
			if err != nil {
				return err
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			app := receivefiles.NewApp(p)

			url := "https://flores7.mniak.dev"
			fmt.Printf("Access %s to send files\nor scan the QR Code below\n", url)
			showQRCode(url)
			fmt.Printf("The files will be stored on %s\n", p.ReceivedFilesDir)

			cobra.CheckErr(app.Start())
			cobra.CheckErr(app.Wait())
		},
	}
	cmd.Flags().IntVar(&p.Port, "port", 10777, "HTTP port of the server")
	cmd.Flags().StringVar(&p.ReceivedFilesDir, "save-to", "./uploads", "Where to save the received files")
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
