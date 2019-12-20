package main

import (
	"image/png"
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/spakin/netpbm"
)

func execRootCmd(cmd *cobra.Command, args []string) error {
	srcPath := args[0]
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	img, err := png.Decode(src)
	if err != nil {
		return err
	}
	fmt.Printf("(W, H) = (%d, %d)\n", img.Bounds().Max.X, img.Bounds().Max.Y)

	dstPath := args[1]
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	err = netpbm.Encode(dst, img, &netpbm.EncodeOptions{
		Format: netpbm.PPM,
		Plain: true,
	})
	if err != nil {
		return err
	}

	return nil
}

func newCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "png2ppm",
		Short: "Converts png to ppm",
		Args: cobra.ExactArgs(2),
		RunE: execRootCmd,
	}
	return rootCmd
}

func main() {
	rootCmd := newCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
