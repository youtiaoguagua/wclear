package main

import (
	"fmt"
	"github.com/muesli/termenv"
	"golang.org/x/term"
	"math"
	"os"
	"sync"
	"time"
)

var width int
var height int
var wg sync.WaitGroup

var lenght = 7

func init() {
	getTerminalSize()
}

//getTerminalSize get terminal size
func getTerminalSize() {

	var err error

	width, height, err = term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	termenv.HideCursor()

	for i := 0; i <= height; i += 2 * lenght {
		// 向右
		goRight(i)
		// 左旋转
		//leftRotate(i+lenght, width)
		// 向左
		goLeft(i + lenght)
		// 右旋转
		//rightRotate(i+2*lenght, 0)
	}
	fmt.Printf("\x1bc")
	termenv.ShowCursor()
}

func leftRotate(x0 int, y0 int) {
	for i := 180; i <= 360; i += 5 {
		rotate(x0, y0, i)
	}
	wg.Wait()
}

func rightRotate(x0 int, y0 int) {
	for i := 180; i >= 0; i -= 5 {
		rotate(x0, y0, i)
	}
	wg.Wait()
}

func goRight(i int) {
	for j := 0; j <= width; j++ {
		goStraight(i, j)
		wg.Wait()
	}
}
func goLeft(i int) {
	for j := width; j >= 0; j-- {
		goStraight(i, j)
		wg.Wait()
	}
}

func calcuCircle(x0 int, y0 int, k int, i int) (float64, float64) {
	x := float64(x0) + float64(k)*math.Cos(math.Pi*float64(i)/180.0)
	y := float64(y0) + float64(k)*math.Sin(math.Pi*float64(i)/180.0)
	return x, y
}

func rotate(x0 int, y0 int, i int) {
	time.Sleep(time.Millisecond * 2)
	wg.Add(1)
	go func(x0 int, y0 int) {
		for k := 0; k < lenght; k++ {
			x, y := calcuCircle(x0, y0, k, i)
			prints(int(math.Floor(x+0.5)), int(math.Floor(y+0.5)), "#")
		}
		wg.Add(1)
		go func() {
			time.Sleep(time.Millisecond * 3)
			for k := 0; k < lenght; k++ {
				x, y := calcuCircle(x0, y0, k, i)
				prints(int(math.Floor(x+0.5)), int(math.Floor(y+0.5)), " ")
			}
			for k := 0; k < lenght; k++ {
				x, y := calcuCircle(x0, y0, k, i)
				prints(int(math.Floor(x+0.5)), int(math.Floor(y+0.5)), " ")
			}
			wg.Done()
		}()
		wg.Done()
	}(x0, y0)
}

func goStraight(i int, j int) {
	time.Sleep(time.Millisecond * 3)
	wg.Add(1)
	go func(x int, y int) {
		for k := 0; k < lenght; k++ {
			prints(x+k, y, "#")
		}
		wg.Add(1)
		go func() {
			time.Sleep(time.Millisecond * 3)
			for k := 0; k < lenght; k++ {
				prints(x+k+1, y, " ")
			}
			for k := 0; k < lenght; k++ {
				prints(x+k-1, y, " ")
			}
			wg.Done()
		}()
		wg.Done()
	}(i, j)
}
func prints(x int, y int, str string) {
	//termenv.MoveCursor(x, y)
	fmt.Printf("\x1b[%d;%dH", x, y)
	//p := termenv.ColorProfile()
	//s := termenv.String(str)
	//s = s.Foreground(p.FromColor(color.RGBA{255, 128, 0, 255}))
	fmt.Printf(str)
}
