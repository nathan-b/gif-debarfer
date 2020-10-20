package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"os"
)

func pixel_is_transparent(img *image.Paletted, x int, y int) bool {
	_, _, _, a := img.At(x, y).RGBA()
	return a == 0
}

func pixel_is_light(img *image.Paletted, x int, y int) bool {
	r, g, b, _ := img.At(x, y).RGBA()

	sum := r + g + b

	return (sum > 0x18000) && (sum < 0x2e000)
}

func should_turn_pixel_transparent(img *image.Paletted, x int, y int) bool {
	xmin := img.Rect.Min.X
	ymin := img.Rect.Min.Y
	xmax := img.Rect.Max.X - 1
	ymax := img.Rect.Max.Y - 1

	i_am_light := pixel_is_light(img, x, y)

	// A light pixel between a transparent pixel and a dark pixel -> transparent
	if (y > ymin) && (y < ymax) && i_am_light {
		if pixel_is_transparent(img, x, y-1) && !pixel_is_light(img, x, y+1) {
			return true
		}
		if pixel_is_transparent(img, x, y+1) && !pixel_is_light(img, x, y-1) {
			return true
		}
	}
	if (x > xmin) && (x < xmax) && i_am_light {
		if pixel_is_transparent(img, x-1, y) && !pixel_is_light(img, x+1, y) {
			return true
		}
		if pixel_is_transparent(img, x+1, y) && !pixel_is_light(img, x-1, y) {
			return true
		}
	}

	// A light pixel between a dark pixel and the edge -> transparent
	if i_am_light {
		if y == ymin && !pixel_is_light(img, x, y+1) {
			return true
		}
		if y == ymax && !pixel_is_light(img, x, y-1) {
			return true
		}
		if x == xmin && !pixel_is_light(img, x+1, y) {
			return true
		}
		if x == xmax && !pixel_is_light(img, x-1, y) {
			return true
		}
	}

	return false
}

func debarf_frame(frame *image.Paletted) error {
	xmin := frame.Rect.Min.X
	ymin := frame.Rect.Min.Y
	xmax := frame.Rect.Max.X
	ymax := frame.Rect.Max.Y

	fmt.Printf("(%d, %d), (%d, %d)\n", xmin, xmax, ymin, ymax)
	for y := ymin; y < ymax; y++ {
		for x := xmin; x < xmax; x++ {
			c := frame.At(x, y)
			r, g, b, a := c.RGBA()
			if pixel_is_transparent(frame, x, y) {
				fmt.Printf("...")
			} else if should_turn_pixel_transparent(frame, x, y) {
				fmt.Printf("xxx")
				frame.Set(x, y, color.RGBA{0, 0, 0, 0})
			} else {
				r /= 0x1000
				g /= 0x1000
				b /= 0x1000
				if a != 0xffff {
					fmt.Printf("barf")
				}
				fmt.Printf("%x%x%x", r, g, b)
			}
		}
		fmt.Printf("\n")
	}

	return nil
}

func debarf_image(img *gif.GIF) error {
	fmt.Printf("\n")
	for i := 0; i < len(img.Image); i++ {
		err := debarf_frame(img.Image[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func process_file(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	g, err := gif.DecodeAll(reader)
	if err != nil {
		return err
	}
	err = debarf_image(g)

	of, err := os.Create("out.gif")
	if err != nil {
		return err
	}
	defer of.Close()
	writer := bufio.NewWriter(of)
	err = gif.EncodeAll(writer, g)
	if err != nil {
		return err
	}
	return err
}

func usage() {
	fmt.Printf("Usage: debarfer <inputfile0> [<inputfile1> ...]\n")
}

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	for i := 1; i < len(os.Args); i++ {
		fmt.Printf("Processing file %s...", os.Args[i])
		err := process_file(os.Args[i])
		if err == nil {
			fmt.Printf("Success!\n")
		} else {
			fmt.Printf("%s\n", err.Error())
		}
	}
}

