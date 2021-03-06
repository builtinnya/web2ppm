package main

import (
	"bytes"
	"context"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/spakin/netpbm"
	"github.com/spf13/cobra"
	"image/png"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strings"
)

func fullScreenshot(urlstr string, quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}
			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).Do(ctx)
			if err != nil {
				return err
			}

			*res, err = page.CaptureScreenshot().
				WithQuality(quality).
				WithClip(&page.Viewport{
					X:      contentSize.X,
					Y:      contentSize.Y,
					Width:  contentSize.Width,
					Height: contentSize.Height,
					Scale:  1,
				}).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	}
}

type cmdParams struct {
	Comments string
}

func execRootCmd(cmd *cobra.Command, args []string, params cmdParams) error {
	urlstr := args[0]
	var comments []string
	if params.Comments != "" {
		content, err := ioutil.ReadFile(params.Comments)
		if err != nil {
			return err
		}
		comments = strings.Split(strings.TrimSpace(string(content)), "\n")
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.NoSandbox,
	)

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancelAlloc()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	var pngBuf []byte
	if err := chromedp.Run(ctx, fullScreenshot(urlstr, 90, &pngBuf)); err != nil {
		return err
	}

	img, err := png.Decode(bytes.NewReader(pngBuf))
	if err != nil {
		return err
	}

	err = netpbm.Encode(os.Stdout, img, &netpbm.EncodeOptions{
		Format:   netpbm.PPM,
		Comments: comments,
		Plain:    true,
	})
	if err != nil {
		return err
	}

	return nil
}

func newCmd() *cobra.Command {
	var params cmdParams
	rootCmd := &cobra.Command{
		Use:   "web2ppm [flags] <url>",
		Short: "Takes a screenshot of an entire web page and outputs to stdout as Plain Portable Pixel Map",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return execRootCmd(cmd, args, params)
		},
	}
	rootCmd.Flags().StringVarP(&params.Comments, "comments", "c", "", "A file path of comments to be inserted")
	return rootCmd
}

func main() {
	rootCmd := newCmd()
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
