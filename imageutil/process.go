package imageutil

import (
	"errors"
	"fmt"
	"github.com/nikhil1raghav/kindle-send/config"
	"github.com/nikhil1raghav/kindle-send/util"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

//Process the image and send the filename to add in epub
func Process(imageUrl string)(string, error){
	cfg:=config.GetInstance()
	downloaded, err :=download(imageUrl)
	if err!=nil{
		return "", err
	}
	if cfg.Color{
		fmt.Println("Value of color is true")
	}
	if !cfg.Color{
		WriteToGray(downloaded)
	}
	return downloaded, nil
}

func download(imageUrl string)(string, error){

	response, err:=http.Get(imageUrl)
	if err!=nil{
		return "",err
	}
	defer response.Body.Close()
	if response.StatusCode!=200{
		log.Println("Got non 200 response")
		return "", errors.New("Non 200 response")
	}

	log.Println("Got response for ", imageUrl)

	file, err:=os.Create(path.Join(config.GetInstance().ImageDir, util.GetHash(imageUrl)))
	log.Println("Writing to file ", file.Name())
	if err!=nil{
		return "", err
	}
	defer file.Close()
	_, err = io.Copy(file, response.Body)
	if err!=nil{
		return "", err
	}
	log.Println("Written successfully to ", file.Name())
	return file.Name(), nil
}