package util

import (
	"image"
	"image/png"
	"log"
	"os"
)
func WriteToGray(filePath string){
	img, err:=loadImage(filePath)
	if err!=nil{
		log.Println(err)
		return
	}
	grayImage := rgbaToGray(img)
	os.Remove(filePath)
	log.Println("Removed old image")
	f,_:=os.Create(filePath)
	defer f.Close()
	png.Encode(f,grayImage)
	log.Println("Wrote the image :)")

}
func loadImage(filePath string)(image.Image, error){
	infile, err:=os.Open(filePath)
	if err!=nil{
		return nil, err
	}
	defer infile.Close()
	img, _, err:=image.Decode(infile)
	if err!=nil{
		return nil, err
	}
	return img, nil
}
func rgbaToGray(img image.Image) *image.Gray{
	var(
		bounds = img.Bounds()
		gray = image.NewGray(bounds)
	)
	for x:=0;x<bounds.Max.X;x++{
		for y:=0;y<bounds.Max.Y;y++{
			var rgba = img.At(x,y)
			gray.Set(x, y, rgba)
		}
	}
	return gray
}