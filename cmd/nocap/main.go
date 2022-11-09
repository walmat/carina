package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/oliamb/cutter"
	"image"
	"image/jpeg"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	// create chrome instance
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))...)
	defer cancel()
	ctx, cancel = chromedp.NewContext(
		ctx,
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// create a timeout
	// ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	// defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate("https://patrickhlauke.github.io/recaptcha/"),
	)
	if err != nil {
		log.Fatal(err)
	}

	var iframes []*cdp.Node
	if err = chromedp.Run(ctx, chromedp.Nodes("iframe", &iframes, chromedp.ByQuery)); err != nil {
		log.Fatal(err)
	}

	if err = chromedp.Run(ctx,
		chromedp.WaitVisible(".recaptcha-checkbox-border", chromedp.FromNode(iframes[0])),
	); err != nil {
		log.Fatal(err)
	}
	err = chromedp.Run(ctx,
		chromedp.Click(".recaptcha-checkbox-border", chromedp.FromNode(iframes[0])),
		// TODO: find a better way to do this, the proper way would be what is commented below but it doesnt work
		// for whatever reason
		chromedp.ActionFunc(func(ctx context.Context) error {
			// wait for challenge, yes this is not proper but idc anymore i havent slept
			time.Sleep(time.Second * 1)
			return nil
		}),
		// chromedp.WaitVisible(".rc-imageselect", chromedp.FromNode(iframes[0])),
	)
	if err != nil {
		log.Fatal(err)
	}

	clickedAny := false
	clickedNodes := make(map[int64]bool)
	cachedImages := make(map[string]image.Image)
	for {
		var challengeNameNodes []*cdp.Node
		var verifyNodes []*cdp.Node
		err = chromedp.Run(ctx,
			chromedp.Nodes(".rc-imageselect-desc-no-canonical > strong", &challengeNameNodes, chromedp.FromNode(iframes[0])),
			chromedp.Nodes("#recaptcha-verify-button", &verifyNodes, chromedp.FromNode(iframes[0])),
		)
		if err != nil {
			log.Fatal(err)
		}

		challengeName := strings.ReplaceAll(challengeNameNodes[0].Children[0].NodeValue, " ", "_")

		var imgNodes []*cdp.Node
		err = chromedp.Run(ctx,
			chromedp.Nodes(".rc-image-tile-wrapper > img", &imgNodes, chromedp.FromNode(iframes[0])),
		)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("images")

		type data struct {
			Challenge string `json:"challenge"`
			Image     string `json:"image"`
		}

		imgRowSize := math.Sqrt(float64(len(imgNodes)))
		for idx, imgNode := range imgNodes {
			x := idx % int(imgRowSize)
			y := int(math.Floor(float64(idx) / imgRowSize))
			_, previouslyClicked := clickedNodes[imgNode.NodeID.Int64()]
			fmt.Println(imgNode.NodeID.Int64(), previouslyClicked)

			imgUrl := imgNode.AttributeValue("src")
			img := image.Image(nil)
			if cachedImage, ok := cachedImages[imgUrl]; ok {
				img = cachedImage
			} else {
				resp, err := http.Get(imgUrl)
				if err != nil {
					panic(err)
				}
				decodedImg, err := jpeg.Decode(resp.Body)
				resp.Body.Close()
				if err != nil {
					panic(err)
				}

				img = decodedImg
				cachedImages[imgUrl] = decodedImg
			}

			pixelsPerSprite := img.Bounds().Max.X
			if !previouslyClicked {
				pixelsPerSprite /= int(imgRowSize)
			} else {
				x = 0
				y = 0
			}
			fmt.Println(x, y, pixelsPerSprite)
			newImg, err := cutter.Crop(img, cutter.Config{
				Width:  pixelsPerSprite,
				Height: pixelsPerSprite,
				Anchor: image.Point{X: pixelsPerSprite * x, Y: pixelsPerSprite * y},
				Mode:   cutter.TopLeft,
			})
			if err != nil {
				panic(err)
			}

			var jpegBuf bytes.Buffer
			if err = jpeg.Encode(&jpegBuf, newImg, nil); err != nil {
				panic(err)
			}

			jsonBytes, _ := json.Marshal(data{Challenge: challengeName, Image: hex.EncodeToString(jpegBuf.Bytes())})
			for {
				resp, err := http.Post("http://localhost:8080/nocap", "application/json", bytes.NewReader(jsonBytes))
				if err != nil {
					continue
				}
				respBytes, err := io.ReadAll(resp.Body)
				if err != nil {
					continue
				}
				if string(respBytes) == "true" {
					err = chromedp.Run(ctx,
						chromedp.MouseClickNode(imgNode),
					)
					if err == nil {
						clickedNodes[imgNode.NodeID.Int64()] = true
						clickedAny = true
						break
					}
				} else {
					break
				}
			}
		}
		if clickedAny {
			var dynamicNodes []*cdp.Node
			timeoutCtx, cancelFunc := context.WithTimeout(ctx, time.Millisecond*100)
			err = chromedp.Run(timeoutCtx,
				chromedp.Nodes(".rc-imageselect-dynamic-selected", &dynamicNodes, chromedp.FromNode(iframes[0])),
			)
			cancelFunc()
			if len(dynamicNodes) == 0 {
				err = chromedp.Run(ctx,
					chromedp.MouseClickNode(verifyNodes[0]),
				)
				if err != nil {
					panic(err)
				}
				fmt.Println("clicked submit")
				time.Sleep(time.Millisecond * 500)
			} else {
				fmt.Println("found dynamic nodes, waiting for them to resolve...")
				err = chromedp.Run(ctx,
					chromedp.WaitNotPresent(".rc-imageselect-dynamic-selected", chromedp.FromNode(iframes[0])),
				)
				if err != nil {
					panic(err)
				}
				time.Sleep(time.Millisecond * 250)
			}
		}
	}

	/*
		fmt.Println(challengeNameNodes[0].Children[0].NodeValue)
		fmt.Println(len(imgNodes))
		imgUrl, _ := imgNodes[0].Attribute("src")

		resp, err := http.Get(imgUrl)
		if err != nil {
			log.Fatal(err)
		}
		img, err := jpeg.Decode(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}

		type data struct {
			Challenge string `json:"challenge"`
			Image     string `json:"image"`
		}
		if len(imgNodes) == 9 { // 3x3
			pixelsPerSprite := img.Bounds().Max.X / 3
			for x := 0; x < 3; x++ {
				for y := 0; y < 3; y++ {
					newImg, err := cutter.Crop(img, cutter.Config{
						Width:  pixelsPerSprite,
						Height: pixelsPerSprite,
						Anchor: image.Point{X: pixelsPerSprite * x, Y: pixelsPerSprite * y},
						Mode:   cutter.TopLeft,
					})
					if err != nil {
						log.Fatal(err)
					}

					var imgBuffer bytes.Buffer
					if err = jpeg.Encode(&imgBuffer, newImg, &jpeg.Options{Quality: 100}); err != nil {
						log.Fatal(err)
					}
					jsonBytes, _ := json.Marshal(data{Challenge: strings.ReplaceAll(challengeNameNodes[0].Children[0].NodeValue, " ", "_"), Image: hex.EncodeToString(imgBuffer.Bytes())})
					resp, _ := http.Post("http://localhost:8080/nocap", "application/json", bytes.NewReader(jsonBytes))
					respBytes, _ := io.ReadAll(resp.Body)
					resp.Body.Close()
					if string(respBytes) == "true" {
						err = chromedp.Run(ctx,
							chromedp.MouseClickNode(imgNodes[(x * 3) + y].Parent),
						)
						if err != nil {
							log.Fatal(err)
						}
						fmt.Println("clicked", (x * 3) + y)
					}
				}
			}
		} else if len(imgNodes) == 16 { // 4x4
			pixelsPerSprite := img.Bounds().Max.X / 4
			for x := 0; x < 4; x++ {
				for y := 0; y < 4; y++ {
					newImg, err := cutter.Crop(img, cutter.Config{
						Width:  pixelsPerSprite,
						Height: pixelsPerSprite,
						Anchor: image.Point{X: pixelsPerSprite * x, Y: pixelsPerSprite * y},
						Mode:   cutter.TopLeft,
					})
					if err != nil {
						log.Fatal(err)
					}

					var imgBuffer bytes.Buffer
					if err = jpeg.Encode(&imgBuffer, newImg, &jpeg.Options{Quality: 100}); err != nil {
						log.Fatal(err)
					}
					jsonBytes, _ := json.Marshal(data{Challenge: strings.ReplaceAll(challengeNameNodes[0].Children[0].NodeValue, " ", "_"), Image: hex.EncodeToString(imgBuffer.Bytes())})
					resp, _ := http.Post("http://localhost:8080/nocap", "application/json", bytes.NewReader(jsonBytes))
					respBytes, _ := io.ReadAll(resp.Body)
					resp.Body.Close()
					fmt.Println(string(respBytes))
				}
			}
		}
	*/

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
