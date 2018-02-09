package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"

	colorful "github.com/lucasb-eyer/go-colorful"
	qrand "github.com/quillaja/goutil/rand"
)

// Make 0-10 random walkers go for a walk. Can produce a single image or a
// sequence of images which can be used for video generation (eg ffmpeg).

func main() {
	const speed = 1
	var w, h, steps, nWalkers int
	var noDiag, video bool

	flag.IntVar(&w, "w", 800, "Width of area.")
	flag.IntVar(&h, "h", 600, "Height of area")
	flag.IntVar(&steps, "steps", 1000, "Number of steps.")
	flag.IntVar(&nWalkers, "n", 1, "Number of walkers. [1,10]")
	flag.BoolVar(&noDiag, "nd", false, "Set to restrict walker to only vertical and horizontal movement (no diagonal).")
	flag.BoolVar(&video, "video", false, "Set to produce sequential images for video production.")
	flag.Parse()

	// ensure nWalkers is appropriate
	switch {
	case nWalkers <= 0:
		nWalkers = 1
	case nWalkers > 10:
		nWalkers = 10
	}

	// create "video" directory if it doesnt exist
	if video {
		if err := os.Mkdir("video", 0755); err != nil && os.IsNotExist(err) {
			panic(err)
		}
	}

	// init RND and thread control stuff
	rand.Seed(time.Now().UnixNano())
	sem := make(chan struct{}, runtime.NumCPU())
	var wg sync.WaitGroup

	walkers := make([]image.Point, nWalkers)
	for i := range walkers {
		walkers[i].X, walkers[i].Y = qrand.IntNM(0, w+1), qrand.IntNM(0, h+1)
	}

	img := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(img, img.Bounds(), image.NewUniform(color.White), image.ZP, draw.Src)

	for i := 0; i < steps; i++ {
		// make a colorful footprint
		for _, pos := range walkers {
			img.Set(pos.X, pos.Y, colorful.Hsv(360*(float64(i)/float64(steps)), 1, 1))
		}

		if video {
			// fucking bullshit to require this crap just to copy a struct
			imgCopy := image.RGBA{
				Pix:    make([]uint8, len(img.Pix)),
				Stride: img.Stride,
				Rect:   img.Rect}
			copy(imgCopy.Pix, img.Pix)

			// start gofunc to write image
			sem <- struct{}{}
			go func(name string, img image.RGBA, w []image.Point) {
				wg.Add(1)
				defer wg.Done()
				// add "cross hairs" to image
				for _, p := range w {
					img.Set(p.X+1, p.Y, color.Black)
					img.Set(p.X-1, p.Y, color.Black)
					img.Set(p.X, p.Y+1, color.Black)
					img.Set(p.X, p.Y-1, color.Black)
				}
				// write image
				writeJpg(name, &img)
				<-sem
			}(fmt.Sprintf("video/%010d.jpg", i), imgCopy, walkers)
		}

		for i := range walkers {
			pos := &walkers[i] // wtf? fucking stupid

			// take a step
			if noDiag {
				if rand.Intn(2) == 0 {
					pos.X += qrand.IntNM(0-speed, 1+speed)
				} else {
					pos.Y += qrand.IntNM(0-speed, 1+speed)
				}
			} else {
				pos.X += qrand.IntNM(0-speed, 1+speed)
				pos.Y += qrand.IntNM(0-speed, 1+speed)
			}

			// clamp to image
			switch {
			case pos.X > w:
				pos.X = w
			case pos.X < 0:
				pos.X = 0
			}
			switch {
			case pos.Y > h:
				pos.Y = h
			case pos.Y < 0:
				pos.Y = 0
			}
		}

	}

	if video {
		wg.Wait()
	} else {
		writeJpg("output.jpg", img)
	}

}

func writeJpg(name string, img *image.RGBA) {
	file, _ := os.Create(name)
	jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
	file.Close()
}
