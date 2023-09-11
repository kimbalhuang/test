package main

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/emulation"
	"io/ioutil"
	"log"
	"math"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func main() {

	// 禁用chrome headless
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// navigate to a page, wait for an element, click

	// capture screenshot of an element
	var buf []byte
	// capture entire browser viewport, returning png with quality=90
	if err := chromedp.Run(ctx, fullScreenshot1(`https://jqlang.github.io/jq/manual/`, 90, &buf)); err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile("./Screenshot.jpeg", buf, 0644); err != nil {
		log.Fatal(err)
	}
	log.Println("图片写入完成")

}

// fullScreenshot takes a screenshot of the entire browser viewport.
// Liberally copied from puppeteer's source.
// Note: this will override the viewport emulation settings.
func fullScreenshot1(urlstr string, quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.EmulateViewport(1200, 45855),
		chromedp.ActionFunc(func(ctx context.Context) error {
			//// get layout metrics
			size1, size2, size3, size4, cssVisualViewportSize, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}

			fmt.Println(size1, size2, size3, size4, cssVisualViewportSize, contentSize)
			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			fmt.Println(width, height)
			// force viewport emulation
			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				//WithScreenOrientation(&emulation.ScreenOrientation{
				//	Type:  emulation.OrientationTypePortraitPrimary,
				//	Angle: 0,
				//}).
				WithScreenHeight(height).Do(ctx)
			if err != nil {
				return err
			}

			// capture screenshot
			*res, _ = page.CaptureScreenshot().
				WithQuality(quality).
				WithClip(&page.Viewport{
					X:      contentSize.X,
					Y:      contentSize.Y,
					Width:  contentSize.Width,
					Height: contentSize.Height,
					Scale:  1,
				}).
				Do(ctx)
			//if err != nil {
			//	return err
			//}
			return nil
		}),
	}
}
