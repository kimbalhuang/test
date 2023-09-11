// Command emulate is a chromedp example demonstrating how to emulate a
// specific device such as an iPhone.
package main

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/page"
	"log"
	"os"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
)

func main() {
	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// run
	var b1, b2 []byte
	//// get layout metrics
	size1, size2, size3, size4, cssVisualViewportSize, contentSize, err := page.GetLayoutMetrics().Do(ctx)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(size1, size2, size3, size4, cssVisualViewportSize, contentSize)
	//width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))
	//fmt.Println(width, height)

	if err := chromedp.Run(ctx,
		// emulate iPhone 7 landscape
		//chromedp.Emulate(device.IPhone7landscape),
		chromedp.Navigate(`https://jqlang.github.io/jq/manual/`),
		//chromedp.FullScreenshot(&b1, 90),

		// reset
		chromedp.Emulate(device.Reset),

		// set really large viewport
		chromedp.EmulateViewport(1200, 45855),
		chromedp.Navigate(`https://jqlang.github.io/jq/manual/`),
		chromedp.CaptureScreenshot(&b2),
	); err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile("screenshot1.png", b1, 0o644); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile("screenshot2.png", b2, 0o644); err != nil {
		log.Fatal(err)
	}
	log.Printf("wrote screenshot1.png and screenshot2.png")
}
