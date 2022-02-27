package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"golang.org/x/term"
	"io/fs"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/user"
	"path"
	"sync"
	"time"
)

var wg sync.WaitGroup
var mutex sync.RWMutex

var width int
var height int
var length = 7
var speed = 3
var delay = 12

type Config struct {
	Speed  int `json:"speed"`
	Delay  int `json:"delay"`
	Length int `json:"length"`
}

func init() {
	getTerminalSize()
	readConfig()
}

func main() {
	flag.IntVar(&speed, "p", speed, "speed of the line")
	flag.IntVar(&delay, "d", delay, "when to hide the line")
	flag.IntVar(&length, "l", length, "the line length")
	flag.Parse()

	fmt.Print("\x1b[?25l")
	for y := 0; y <= height; y += 2 * length {
		// 向右
		goRight(y)
		// 左旋转
		leftRotate(y+length, width)
		if y+length <= height {
			// 向左
			goLeft(y + length)
			// 右旋转
			rightRotate(y+2*length, 0)
		}
	}
	fmt.Printf("\x1bc")
	fmt.Print("\x1b[?25h")

}

func goRight(y int) {
	for x := 0; x <= width; x++ {
		goStraight(y, x)
	}
	wg.Wait()
}

func goLeft(y int) {
	for x := width; x >= 0; x-- {
		goStraight(y, x)
	}
	wg.Wait()
}

func goStraight(y int, x int) {
	time.Sleep(time.Millisecond * time.Duration(speed))
	wg.Add(1)
	go func(y int, x int) {
		for k := 0; k < length; k++ {
			prints(y+k, x, "#")
		}
		wg.Done()
	}(y, x)

	wg.Add(1)
	go cleanStraight(y, x)
}

func cleanStraight(y int, x int) {
	time.Sleep(time.Millisecond * time.Duration(delay))

	for k := 0; k < length; k++ {
		for n := 0; n < 3; n++ {
			prints(y+k, x, " ")
		}
	}
	wg.Done()
}

func leftRotate(y0 int, x0 int) {
	for angle := 180; angle <= 360; angle += 20 {
		rotate(y0, x0, angle)
	}
	wg.Wait()
}

func rightRotate(y0 int, x0 int) {
	for angle := 180; angle >= 0; angle -= 20 {
		rotate(y0, x0, angle)
	}
	wg.Wait()
}

func rotate(y0 int, x0 int, angle int) {
	time.Sleep(time.Millisecond * time.Duration(speed))
	wg.Add(1)
	go func(y0 int, x0 int, angle int) {
		for k := 0; k < length; k++ {
			x, y := calcuCircle(y0, x0, k, angle)
			prints(int(math.Floor(y+0.5)), int(math.Floor(x+0.5)), "#")
		}
		wg.Done()
	}(y0, x0, angle)

	wg.Add(1)
	go cleanRotate(y0, x0, angle)
}

func cleanRotate(y0 int, x0 int, angle int) {
	time.Sleep(time.Millisecond * time.Duration(delay))

	for k := 0; k < length; k++ {
		for n := 0; n < 3; n++ {
			x, y := calcuCircle(y0, x0, k, angle)
			prints(int(math.Floor(y+0.5)), int(math.Floor(x+0.5)), " ")
		}
	}
	wg.Done()
}

func calcuCircle(y0 int, x0 int, k int, angle int) (float64, float64) {
	y := float64(y0) + float64(k)*math.Cos(math.Pi*float64(angle)/180.0)
	x := float64(x0) + float64(k)*math.Sin(math.Pi*float64(angle)/180.0)
	return x, y
}

func prints(y int, x int, str string) {
	mutex.Lock()
	fmt.Printf("\x1b[%d;%dH", y, x)
	fmt.Printf(str)
	mutex.Unlock()
}

//getTerminalSize get terminal size
func getTerminalSize() {
	var err error
	width, height, err = term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		panic(err)
	}
}

func readConfig() {
	// 获取用户目录
	currentUser, err := user.Current()
	if err != nil {
		return
	}
	// 拼接配置文件
	p := path.Join(currentUser.HomeDir, "/.config", "/wclear.json")

	if !fileExist(p) {
		err = os.MkdirAll(path.Join(currentUser.HomeDir, "/.config"), os.ModePerm)
		if err != nil {
			log.Panicln(err)
		}
		// 创建配置文件
		f, err := os.OpenFile(p, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				log.Panicln(err)
			}
		}(f)
		if err != nil {
			log.Panicln("create file failL:", err)
		}
		// 写入配置初始化内容
		config, err := json.MarshalIndent(Config{speed, delay, length}, "", "	")
		if err != nil {
			log.Panicln(err)
		}
		_, err = f.Write(config)
		if err != nil {
			return
		}
		return
	}
	buf, err := ioutil.ReadFile(p)
	if err != nil {
		log.Panicln("load config conf failed: ", err)
	}
	config := &Config{}
	err = json.Unmarshal(buf, config)
	if err != nil {
		config, err := json.MarshalIndent(Config{speed, delay, length}, "", "	")
		if err != nil {
			log.Panicln(err)
		}
		err = ioutil.WriteFile(p, config, fs.ModePerm)
		if err != nil {
			log.Panicln(err)
			return
		}
	}

	speed = config.Speed
	delay = config.Delay
	length = config.Length
}

func fileExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}
