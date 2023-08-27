package epubgen

import (
	"errors"
	"image/png"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/bmaupin/go-epub"
	"github.com/go-shiori/go-readability"
	"github.com/gosimple/slug"
	"github.com/nikhil1raghav/kindle-send/config"
	"github.com/nikhil1raghav/kindle-send/util"
	"golang.org/x/image/webp"
)

type epubImage struct {
	epubImgSrc  string
	localImgSrc string
}

type epubmaker struct {
	Epub      *epub.Epub
	tmpDir    string
	downloads map[string]epubImage
}

func NewEpubmaker(title string) *epubmaker {
	downloadMap := make(map[string]epubImage)
	return &epubmaker{
		Epub:      epub.NewEpub(title),
		downloads: downloadMap,
	}
}

func fetchReadable(url string) (readability.Article, error) {
	return readability.FromURL(url, 30*time.Second)
}

// Point remote image link to downloaded image
func (e *epubmaker) changeRefs(i int, img *goquery.Selection) {
	img.RemoveAttr("loading")
	img.RemoveAttr("srcset")
	imgSrc, exists := img.Attr("src")
	if exists {
		if _, ok := e.downloads[imgSrc]; ok {
			util.Green.Printf("Setting img src from %s to %s \n", imgSrc, e.downloads[imgSrc].epubImgSrc)
			img.SetAttr("src", e.downloads[imgSrc].epubImgSrc)
		}
	}
}

// Download images and add to epub zip
func (e *epubmaker) downloadImages(i int, img *goquery.Selection) {
	util.CyanBold.Println("Downloading Images")
	imgSrc, exists := img.Attr("src")

	if exists {

		//don't download same thing twice
		if _, ok := e.downloads[imgSrc]; ok {
			return
		}

		imageFileName, output, _ := e.processImage(imgSrc)

		var imgRef string
		var err error

		if len(output) > 0 {
			imgRef, err = e.Epub.AddImage(output, imageFileName)
		} else {
			imgRef, err = e.Epub.AddImage(imgSrc, imageFileName)
		}
		if err != nil {
			util.Red.Printf("Couldn't add image %s : %s\n", imgSrc, err)
			return
		} else {
			util.Green.Printf("Downloaded image %s\n", imgSrc)
			entry := epubImage{epubImgSrc: imgRef, localImgSrc: output}
			e.downloads[imgSrc] = entry
		}
	}
}

func (e *epubmaker) processImage(src string) (string, string, error) {
	getName := func(src string, ext string) string {
		return "img" + util.GetHash(src) + "." + ext
	}
	r, err := http.Head(src)
	if err != nil {
		return "", "", err
	}
	switch contentType := r.Header.Get("Content-type"); contentType {
	case "image/webp":
		fileName := getName(src, "webp")
		filePath := filepath.Join(e.tmpDir, fileName)
		f, _ := http.Get(src)
		img, _ := webp.Decode(f.Body)
		pngFile, _ := os.Create(filePath)
		png.Encode(pngFile, img)
		return fileName, filePath, nil
	case "image/jpeg":
		return getName(src, "jpg"), "", nil
	default:
		return getName(src, "png"), "", nil
	}
}

// Fetches images in article and then embeds them into epub
func (e *epubmaker) embedImages(wg *sync.WaitGroup, article *readability.Article) {
	util.Cyan.Println("Embedding images in ", article.Title)
	defer wg.Done()
	//TODO: Compress images before embedding to improve size
	doc := goquery.NewDocumentFromNode(article.Node)

	//download all images
	doc.Find("img").Each(e.downloadImages)

	//Change all refs, doing it in two phases to download repeated images only once
	doc.Find("img").Each(e.changeRefs)

	content, err := doc.Html()

	if err != nil {
		util.Red.Printf("Error converting modified %s to HTML, it will be transferred without images : %s \n", article.Title, err)
	} else {
		article.Content = content
	}
}

// TODO: Look for better formatting, this is bare bones
func prepare(article *readability.Article) string {
	return "<h1>" + article.Title + "</h1>" + article.Content
}

// Add articles to epub
func (e *epubmaker) addContent(articles *[]readability.Article) error {
	added := 0
	for _, article := range *articles {
		_, err := e.Epub.AddSection(prepare(&article), article.Title, "", "")
		if err != nil {
			util.Red.Printf("Couldn't add %s to epub : %s", article.Title, err)
		} else {
			added++
		}
	}
	util.Green.Printf("Added %d articles\n", added)
	if added == 0 {
		return errors.New("No article was added, epub creation failed")
	}
	return nil
}

// Generates a single epub from a slice of urls, returns file path
func Make(pageUrls []string, title string) (string, error) {
	//TODO: Parallelize fetching pages

	//Get readable article from urls
	readableArticles := make([]readability.Article, 0)
	for _, pageUrl := range pageUrls {
		article, err := fetchReadable(pageUrl)
		if err != nil {
			util.Red.Printf("Couldn't convert %s because %s", pageUrl, err)
			util.Magenta.Println("SKIPPING ", pageUrl)
			continue
		}
		util.Green.Printf("Fetched %s --> %s\n", pageUrl, article.Title)
		readableArticles = append(readableArticles, article)
	}

	if len(readableArticles) == 0 {
		return "", errors.New("No readable url given, exiting without creating epub")
	}

	if len(title) == 0 {
		title = readableArticles[0].Title
		util.Magenta.Printf("No title supplied, inheriting title of first readable article : %s \n", title)
	}

	book := NewEpubmaker(title)

	//get images and embed them
	var wg sync.WaitGroup

	// TODO: Error handling
	tmpDir, _ := os.MkdirTemp("", "kindle-send-")
	book.tmpDir = tmpDir

	defer os.RemoveAll(tmpDir)

	for i := 0; i < len(readableArticles); i++ {
		wg.Add(1)
		go book.embedImages(&wg, &readableArticles[i])
	}

	wg.Wait()

	err := book.addContent(&readableArticles)
	if err != nil {
		return "", err
	}
	var storeDir string
	if len(config.GetInstance().StorePath) == 0 {
		storeDir, err = os.Getwd()
		if err != nil {
			util.Red.Println("Error getting current directory, trying fallback")
			storeDir = "./"
		}
	} else {
		storeDir = config.GetInstance().StorePath
	}

	titleSlug := slug.Make(title)
	var filename string
	if len(titleSlug) == 0 {
		filename = "kindle-send-doc-" + util.GetHash(readableArticles[0].Content) + ".epub"
	} else {
		filename = titleSlug + ".epub"
	}
	filepath := path.Join(storeDir, filename)
	err = book.Epub.Write(filepath)
	if err != nil {
		return "", err
	}
	return filepath, nil
}
