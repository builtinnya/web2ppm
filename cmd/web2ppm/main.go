package main

import (
	"context"
	"io/ioutil"
	"log"
	"math"
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func fullScreenshot(urlstr string, quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// get layout metrics
			_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}

			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			// force viewport emulation
			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}

			// capture screenshot
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

func execRootCmd(cmd *cobra.Command, args []string) error {
	url := args[0]

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.NoSandbox,
	)

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancelAlloc()
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	var buf []byte
	if err := chromedp.Run(ctx, fullScreenshot(url, 90, &buf)); err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile("screenshot.png", buf, 0644); err != nil {
		log.Fatal(err)
	}

	return nil
}

func newCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "web2ppm",
		Short: "Takes a screenshot of a web page as Portable Pixel Map",
		Args: cobra.ExactArgs(1),
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
