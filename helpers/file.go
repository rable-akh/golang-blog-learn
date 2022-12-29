package helpers

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"math"
	"mime/multipart"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ImageProcessing(c *gin.Context, direname string, fileHandel *multipart.FileHeader) (string, error) {
	filename := strings.Replace(uuid.New().String(), "", "", -1)

	folderErr := CreateFolder(direname)
	if folderErr != nil {
		log.Fatal(folderErr)
		return "", folderErr
	}

	// return filename, folderErr
	saveFileErr := c.SaveUploadedFile(fileHandel, "public/blogs/"+fileHandel.Filename)

	if saveFileErr != nil {
		fmt.Println(saveFileErr)
	}

	fileExt := strings.Split(fileHandel.Filename, ".")[1]

	filename = filename + "." + fileExt

	f, err := os.Open("public/blogs/" + fileHandel.Filename)
	fmt.Println(f)
	if err != nil {
		log.Fatal(err)
	}

	//encoding message is discarded, because OP wanted only jpg, else use encoding in resize function
	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	//this is the resized image
	resImg := Resize(img, 100, 100)

	//this is the resized image []bytes
	var imgBytes []byte
	if fileExt == "jpeg" || fileExt == "jpg" {
		imgBytes = JpegBytes(resImg)
	} else if fileExt == "png" {
		imgBytes = PNGBytes(resImg)
	} else if fileExt == "gif" {
		imgBytes = GifBytes(resImg)
	} else {
		return fileHandel.Filename, err
	}

	//optional written to file
	err = ioutil.WriteFile("public/blogs/thumbs/"+filename, imgBytes, 0777)
	if err != nil {
		log.Fatal(err)
	}
	return filename, err
}

func CreateFolder(dirname string) error {
	_, err := os.Stat("public")
	if err == nil {
		errDir := os.MkdirAll("public/"+dirname+"/thumbs", 0755)
		if errDir != nil {
			return errDir
		}
	}

	return nil
}

func Resize(img image.Image, length int, width int) image.Image {
	//truncate pixel size
	minX := img.Bounds().Min.X
	minY := img.Bounds().Min.Y
	maxX := img.Bounds().Max.X
	maxY := img.Bounds().Max.Y
	for (maxX-minX)%length != 0 {
		maxX--
	}
	for (maxY-minY)%width != 0 {
		maxY--
	}
	scaleX := (maxX - minX) / length
	scaleY := (maxY - minY) / width

	imgRect := image.Rect(0, 0, length, width)
	resImg := image.NewRGBA(imgRect)
	draw.Draw(resImg, resImg.Bounds(), &image.Uniform{C: color.White}, image.ZP, draw.Src)
	for y := 0; y < width; y += 1 {
		for x := 0; x < length; x += 1 {
			averageColor := GetAverageColor(img, minX+x*scaleX, minX+(x+1)*scaleX, minY+y*scaleY, minY+(y+1)*scaleY)
			resImg.Set(x, y, averageColor)
		}
	}
	return resImg
}

func GetAverageColor(img image.Image, minX int, maxX int, minY int, maxY int) color.Color {
	var averageRed float64
	var averageGreen float64
	var averageBlue float64
	var averageAlpha float64
	scale := 1.0 / float64((maxX-minX)*(maxY-minY))

	for i := minX; i < maxX; i++ {
		for k := minY; k < maxY; k++ {
			r, g, b, a := img.At(i, k).RGBA()
			averageRed += float64(r) * scale
			averageGreen += float64(g) * scale
			averageBlue += float64(b) * scale
			averageAlpha += float64(a) * scale
		}
	}

	averageRed = math.Sqrt(averageRed)
	averageGreen = math.Sqrt(averageGreen)
	averageBlue = math.Sqrt(averageBlue)
	averageAlpha = math.Sqrt(averageAlpha)

	averageColor := color.RGBA{
		R: uint8(averageRed),
		G: uint8(averageGreen),
		B: uint8(averageBlue),
		A: uint8(averageAlpha)}

	return averageColor
}

func JpegBytes(img image.Image) []byte {
	var opt jpeg.Options
	opt.Quality = 80

	buff := bytes.NewBuffer(nil)
	err := jpeg.Encode(buff, img, &opt)
	if err != nil {
		log.Fatal(err)
	}

	return buff.Bytes()
}

func PNGBytes(img image.Image) []byte {

	buff := bytes.NewBuffer(nil)
	err := png.Encode(buff, img)
	if err != nil {
		log.Fatal(err)
	}

	return buff.Bytes()
}

func GifBytes(img image.Image) []byte {
	var opt gif.Options

	buff := bytes.NewBuffer(nil)
	err := gif.Encode(buff, img, &opt)
	if err != nil {
		log.Fatal(err)
	}

	return buff.Bytes()
}
