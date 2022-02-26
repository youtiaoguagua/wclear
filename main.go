package main

import (
	"fmt"
	"golang.org/x/sys/windows"
	"math"
	"os"
	"sync"
	"time"
)

var width int
var height int
var wg sync.WaitGroup

var length = 7

var mutex sync.RWMutex

func init() {
	getTerminalSize()
}

func main() {
	fmt.Print("\x1b[?25l")
	for i := 0; i <= height; i += 2 * length {
		// 向右
		goRight(i)
		// 左旋转
		leftRotate(i+length, width)
		// 向左
		goLeft(i + length)
		// 右旋转
		rightRotate(i+2*length, 0)
	}
	fmt.Printf("\x1bc")
	fmt.Print("\x1b[?25h")

}

func leftRotate(x0 int, y0 int) {
	for i := 180; i <= 360; i += 20 {
		rotate(x0, y0, i)
	}
	wg.Wait()
}

func rightRotate(x0 int, y0 int) {
	for i := 180; i >= 0; i -= 20 {
		rotate(x0, y0, i)
	}
	wg.Wait()
}

func goRight(i int) {
	for j := 0; j <= width; j++ {
		goStraight(i, j)
	}
	wg.Wait()
}
func goLeft(i int) {
	for j := width; j >= 0; j-- {
		goStraight(i, j)
	}
	wg.Wait()
}

func calcuCircle(x0 int, y0 int, k int, i int) (float64, float64) {
	x := float64(x0) + float64(k)*math.Cos(math.Pi*float64(i)/180.0)
	y := float64(y0) + float64(k)*math.Sin(math.Pi*float64(i)/180.0)
	return x, y
}

func rotate(x0 int, y0 int, i int) {
	time.Sleep(time.Millisecond * 6)
	wg.Add(1)
	go func(x0 int, y0 int, i int) {
		for k := 0; k < length; k++ {
			x, y := calcuCircle(x0, y0, k, i)
			prints(int(math.Floor(x+0.5)), int(math.Floor(y+0.5)), "#")
		}
		wg.Add(1)
		go cleanRotate(x0, y0, i)
		wg.Done()
	}(x0, y0, i)
}

func goStraight(i int, j int) {
	time.Sleep(time.Millisecond * 4)
	wg.Add(1)
	go func(x int, y int) {
		for k := 0; k < length; k++ {
			prints(x+k, y, "#")
		}
		wg.Add(1)
		go cleanStraight(x, y)
		wg.Done()
	}(i, j)
}

func cleanRotate(x0 int, y0 int, i int) {
	time.Sleep(time.Millisecond * 7)

	for k := 0; k < length; k++ {
		for n := 0; n < 3; n++ {
			x, y := calcuCircle(x0, y0, k, i)
			prints(int(math.Floor(x+0.5)), int(math.Floor(y+0.5)), " ")
		}
	}
	wg.Done()
}

func cleanStraight(x int, y int) {
	time.Sleep(time.Millisecond * 10)

	for k := 0; k < length; k++ {
		for n := 0; n < 3; n++ {
			prints(x+k, y, " ")
		}
	}
	wg.Done()
}

func prints(x int, y int, str string) {
	mutex.Lock()
	fmt.Printf("\x1b[%d;%dH", x, y)
	fmt.Printf(str)
	mutex.Unlock()

}

//getTerminalSize get terminal size
func getTerminalSize() {
	var info windows.ConsoleScreenBufferInfo
	if err := windows.GetConsoleScreenBufferInfo(windows.Handle(int(os.Stdout.Fd())), &info); err != nil {
		panic(err)
	}
	width = int(info.Window.Right - info.Window.Left + 1)
	height = int(info.Window.Bottom - info.Window.Top + 1)
}
