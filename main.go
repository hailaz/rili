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
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

// main description
//
// createTime: 2024-03-28 17:16:59
func main() {
	now := time.Now()
	for i := 0; i < 2; i++ {
		Rili(now.AddDate(0, 0, i))
	}
}

// MyImage description
type MyImage struct {
	image.Image
	Font *truetype.Font
}

// NewMyImage description
//
// createTime: 2024-03-28 17:54:05
func NewMyImage(srcImage string, fontPath string) *MyImage {
	srcFile, err := os.Open(srcImage)
	if err != nil {
		log.Println(err)
		return nil
	}
	img, err := png.Decode(srcFile)
	if err != nil {
		log.Println(err)
		return nil
	}
	// 读取字体
	fontBytes, err := os.ReadFile(fontPath)
	if err != nil {
		return nil
	}
	myFont, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil
	}
	return &MyImage{
		Image: img,
		Font:  myFont,
	}
}

// Rili 日历
//
// createTime: 2024-03-28 17:54:05
func Rili(t time.Time) {
	img := NewMyImage("src/null.png", "src/simsun.ttc")
	// 给图片添加文字，这里要指定字体文件的路径
	// 这里我使用的是simsun.ttf是宋体
	// 这里换成你自己电脑上的字体文件
	err := img.AddLabel("今年可真难熬,是吧？", 110, 55, color.RGBA{0, 0, 0, 255}, 40)

	if err != nil {
		log.Println(err)
		return
	}

	err = img.AddLabel(t.Format("船长,今天才01月02号"), 40, 156, color.RGBA{0, 0, 0, 255}, 40)

	if err != nil {
		log.Println(err)
		return
	}

	img.SaveFile(t.Format("output/2006-01-02.png"))

}

// SaveFile description
//
// createTime: 2024-03-28 17:58:35
func (img *MyImage) SaveFile(filePath string) {
	// 判断目录是否存在
	_, err := os.Stat("output")
	if err != nil {
		if os.IsNotExist(err) {
			// 创建目录
			err = os.MkdirAll("output", os.ModePerm)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}

	f, err := os.Create(filePath)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	err = png.Encode(f, img.Image)
	if err != nil {
		log.Println(err)
		return
	}
}

func (img *MyImage) AddLabel(label string, x, y int, fontColor color.Color, size float64) error {
	bound := img.Bounds()
	// 创建一个新的图片
	rgba := image.NewRGBA(image.Rect(0, 0, bound.Dx(), bound.Dy()))

	draw.Draw(rgba, rgba.Bounds(), img, bound.Min, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(img.Font)
	c.SetFontSize(size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	uni := image.NewUniform(fontColor)
	c.SetSrc(uni)
	c.SetHinting(font.HintingNone)

	// 在指定的位置显示
	pt := freetype.Pt(x, y+int(c.PointToFixed(size)>>6))
	if _, err := c.DrawString(label, pt); err != nil {
		return err
	}

	img.Image = rgba
	return nil
}
