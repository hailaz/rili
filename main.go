package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"time"

	"github.com/golang/freetype"
	"golang.org/x/image/font"
)

// main description
//
// createTime: 2024-03-28 17:16:59
func main() {
	now := time.Now()
	for i := 0; i < 100; i++ {
		creat(now.AddDate(0, 0, i))
	}
}

func creat(t time.Time) {
	src, err := os.Open("src/null.png")
	if err != nil {
		log.Println(err)
		return
	}
	img, err := png.Decode(src)
	if err != nil {
		log.Println(err)
		return
	}
	// 给图片添加文字，这里要指定字体文件的路径
	// 这里我使用的是simsun.ttf是宋体
	// 这里换成你自己电脑上的字体文件
	outimage, err := addLabel(img, "今年可真难熬,是吧？", 110, 55, color.RGBA{0, 0, 0, 255}, 40, "src/simsun.ttc")

	if err != nil {
		log.Println(err)
		return
	}

	outimage, err = addLabel(outimage, t.Format("船长,今天才01月02号"), 40, 156, color.RGBA{0, 0, 0, 255}, 40, "simsun.ttc")

	if err != nil {
		log.Println(err)
		return
	}

	// 创建目录
	err = os.MkdirAll("output", os.ModePerm)
	if err != nil {
		log.Println(err)
		return
	}

	f, err := os.Create(t.Format("output/2006-01-02.png"))
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	err = png.Encode(f, outimage)
	if err != nil {
		log.Println(err)
		return
	}
}

func addLabel(img image.Image, label string, x, y int, fontColor color.Color, size float64, fontPath string) (image.Image, error) {
	bound := img.Bounds()
	// 创建一个新的图片
	rgba := image.NewRGBA(image.Rect(0, 0, bound.Dx(), bound.Dy()))
	// 读取字体
	fontBytes, err := os.ReadFile(fontPath)
	if err != nil {
		return rgba, err
	}
	myFont, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return rgba, err
	}

	draw.Draw(rgba, rgba.Bounds(), img, bound.Min, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(myFont)
	c.SetFontSize(size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	uni := image.NewUniform(fontColor)
	c.SetSrc(uni)
	c.SetHinting(font.HintingNone)

	// 在指定的位置显示
	pt := freetype.Pt(x, y+int(c.PointToFixed(size)>>6))
	if _, err := c.DrawString(label, pt); err != nil {
		return rgba, err
	}

	return rgba, nil
}
