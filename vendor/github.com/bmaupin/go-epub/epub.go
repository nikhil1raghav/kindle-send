/*
Package epub generates valid EPUB 3.0 files with additional EPUB 2.0 table of
contents (as seen here: https://github.com/bmaupin/epub-samples) for maximum
compatibility.

Basic usage:

	// Create a new EPUB
	e := epub.NewEpub("My title")

	// Set the author
	e.SetAuthor("Hingle McCringleberry")

	// Add a section
	section1Body := `<h1>Section 1</h1>
	<p>This is a paragraph.</p>`
	e.AddSection(section1Body, "Section 1", "", "")

	// Write the EPUB
	err = e.Write("My EPUB.epub")
	if err != nil {
		// handle error
	}

*/
package epub

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	// TODO: Eventually this should include the major version (e.g. github.com/gofrs/uuid/v3) but that would break
	// compatibility with Go < 1.9 (https://github.com/golang/go/wiki/Modules#semantic-import-versioning)
	"github.com/gofrs/uuid"
	"github.com/vincent-petithory/dataurl"
)

// FilenameAlreadyUsedError is thrown by AddCSS, AddFont, AddImage, or AddSection
// if the same filename is used more than once.
type FilenameAlreadyUsedError struct {
	Filename string // Filename that caused the error
}

func (e *FilenameAlreadyUsedError) Error() string {
	return fmt.Sprintf("Filename already used: %s", e.Filename)
}

// FileRetrievalError is thrown by AddCSS, AddFont, AddImage, or Write if there was a
// problem retrieving the source file that was provided.
type FileRetrievalError struct {
	Source string // The source of the file whose retrieval failed
	Err    error  // The underlying error that was thrown
}

func (e *FileRetrievalError) Error() string {
	return fmt.Sprintf("Error retrieving %q from source: %+v", e.Source, e.Err)
}

// ParentDoesNotExistError is thrown by AddSubSection if the parent with the
// previously defined internal filename does not exist.
type ParentDoesNotExistError struct {
	Filename string // Filename that caused the error
}

func (e *ParentDoesNotExistError) Error() string {
	return fmt.Sprintf("Parent with the internal filename %s does not exist", e.Filename)
}

// Folder names used for resources inside the EPUB
const (
	CSSFolderName   = "css"
	FontFolderName  = "fonts"
	ImageFolderName = "images"
	VideoFolderName = "videos"
)

const (
	cssFileFormat          = "css%04d%s"
	defaultCoverBody       = `<img src="%s" alt="Cover Image" />`
	defaultCoverCSSContent = `body {
  background-color: #FFFFFF;
  margin-bottom: 0px;
  margin-left: 0px;
  margin-right: 0px;
  margin-top: 0px;
  text-align: center;
}
img {
  max-height: 100%;
  max-width: 100%;
}
`
	defaultCoverCSSFilename   = "cover.css"
	defaultCoverCSSSource     = "cover.css"
	defaultCoverImgFormat     = "cover%s"
	defaultCoverXhtmlFilename = "cover.xhtml"
	defaultEpubLang           = "en"
	fontFileFormat            = "font%04d%s"
	imageFileFormat           = "image%04d%s"
	videoFileFormat           = "video%04d%s"
	sectionFileFormat         = "section%04d.xhtml"
	urnUUIDPrefix             = "urn:uuid:"
)

// Epub implements an EPUB file.
type Epub struct {
	sync.Mutex
	*http.Client
	author string
	cover  *epubCover
	// The key is the css filename, the value is the css source
	css map[string]string
	// The key is the font filename, the value is the font source
	fonts      map[string]string
	identifier string
	// The key is the image filename, the value is the image source
	images map[string]string
	// The key is the video filename, the value is the video source
	videos map[string]string
	// Language
	lang string
	// Description
	desc string
	// Page progression direction
	ppd string
	// The package file (package.opf)
	pkg      *pkg
	sections []epubSection
	title    string
	// Table of contents
	toc *toc
}

type epubCover struct {
	cssFilename   string
	cssTempFile   string
	imageFilename string
	xhtmlFilename string
}

type epubSection struct {
	filename string
	xhtml    *xhtml
	children *[]epubSection
}

// NewEpub returns a new Epub.
func NewEpub(title string) *Epub {
	e := &Epub{}
	e.cover = &epubCover{
		cssFilename:   "",
		cssTempFile:   "",
		imageFilename: "",
		xhtmlFilename: "",
	}
	e.Client = http.DefaultClient
	e.css = make(map[string]string)
	e.fonts = make(map[string]string)
	e.images = make(map[string]string)
	e.videos = make(map[string]string)
	e.pkg = newPackage()
	e.toc = newToc()
	// Set minimal required attributes
	e.SetIdentifier(urnUUIDPrefix + uuid.Must(uuid.NewV4()).String())
	e.SetLang(defaultEpubLang)
	e.SetTitle(title)

	return e
}

// AddCSS adds a CSS file to the EPUB and returns a relative path to the CSS
// file that can be used in EPUB sections in the format:
// ../CSSFolderName/internalFilename
//
// The CSS source should either be a URL, a path to a local file, or an embedded data URL; in any
// case, the CSS file will be retrieved and stored in the EPUB.
//
// The internal filename will be used when storing the CSS file in the EPUB
// and must be unique among all CSS files. If the same filename is used more
// than once, FilenameAlreadyUsedError will be returned. The internal filename is
// optional; if no filename is provided, one will be generated.
func (e *Epub) AddCSS(source string, internalFilename string) (string, error) {
	e.Lock()
	defer e.Unlock()
	return e.addCSS(source, internalFilename)
}

func (e *Epub) addCSS(source string, internalFilename string) (string, error) {
	return addMedia(e.Client, source, internalFilename, cssFileFormat, CSSFolderName, e.css)
}

// AddFont adds a font file to the EPUB and returns a relative path to the font
// file that can be used in EPUB sections in the format:
// ../FontFolderName/internalFilename
//
// The font source should either be a URL, a path to a local file, or an embedded data URL; in any
// case, the font file will be retrieved and stored in the EPUB.
//
// The internal filename will be used when storing the font file in the EPUB
// and must be unique among all font files. If the same filename is used more
// than once, FilenameAlreadyUsedError will be returned. The internal filename is
// optional; if no filename is provided, one will be generated.
func (e *Epub) AddFont(source string, internalFilename string) (string, error) {
	e.Lock()
	defer e.Unlock()
	return addMedia(e.Client, source, internalFilename, fontFileFormat, FontFolderName, e.fonts)
}

// AddImage adds an image to the EPUB and returns a relative path to the image
// file that can be used in EPUB sections in the format:
// ../ImageFolderName/internalFilename
//
// The image source should either be a URL, a path to a local file, or an embedded data URL; in any
// case, the image file will be retrieved and stored in the EPUB.
//
// The internal filename will be used when storing the image file in the EPUB
// and must be unique among all image files. If the same filename is used more
// than once, FilenameAlreadyUsedError will be returned. The internal filename is
// optional; if no filename is provided, one will be generated.
func (e *Epub) AddImage(source string, imageFilename string) (string, error) {
	e.Lock()
	defer e.Unlock()
	return addMedia(e.Client, source, imageFilename, imageFileFormat, ImageFolderName, e.images)
}

// AddVideo adds an video to the EPUB and returns a relative path to the video
// file that can be used in EPUB sections in the format:
// ../VideoFolderName/internalFilename
//
// The video source should either be a URL, a path to a local file, or an embedded data URL; in any
// case, the video file will be retrieved and stored in the EPUB.
//
// The internal filename will be used when storing the video file in the EPUB
// and must be unique among all video files. If the same filename is used more
// than once, FilenameAlreadyUsedError will be returned. The internal filename is
// optional; if no filename is provided, one will be generated.
func (e *Epub) AddVideo(source string, videoFilename string) (string, error) {
	e.Lock()
	defer e.Unlock()
	return addMedia(e.Client, source, videoFilename, videoFileFormat, VideoFolderName, e.videos)
}

// AddSection adds a new section (chapter, etc) to the EPUB and returns a
// relative path to the section that can be used from another section (for
// links).
//
// The body must be valid XHTML that will go between the <body> tags of the
// section XHTML file. The content will not be validated.
//
// The title will be used for the table of contents. The section will be shown
// in the table of contents in the same order it was added to the EPUB. The
// title is optional; if no title is provided, the section will not be added to
// the table of contents.
//
// The internal filename will be used when storing the section file in the EPUB
// and must be unique among all section files. If the same filename is used more
// than once, FilenameAlreadyUsedError will be returned. The internal filename is
// optional; if no filename is provided, one will be generated.
//
// The internal path to an already-added CSS file (as returned by AddCSS) to be
// used for the section is optional.
func (e *Epub) AddSection(body string, sectionTitle string, internalFilename string, internalCSSPath string) (string, error) {
	e.Lock()
	defer e.Unlock()
	return e.addSection("", body, sectionTitle, internalFilename, internalCSSPath)
}

// AddSubSection adds a nested section (chapter, etc) to an existing section.
// The method returns a relative path to the section that can be used from another
// section (for links).
//
// The parent filename must be a valid filename from another section already added.
//
// The body must be valid XHTML that will go between the <body> tags of the
// section XHTML file. The content will not be validated.
//
// The title will be used for the table of contents. The section will be shown
// as a nested entry of the parent section in the table of contents. The
// title is optional; if no title is provided, the section will not be added to
// the table of contents.
//
// The internal filename will be used when storing the section file in the EPUB
// and must be unique among all section files. If the same filename is used more
// than once, FilenameAlreadyUsedError will be returned. The internal filename is
// optional; if no filename is provided, one will be generated.
//
// The internal path to an already-added CSS file (as returned by AddCSS) to be
// used for the section is optional.
func (e *Epub) AddSubSection(parentFilename string, body string, sectionTitle string, internalFilename string, internalCSSPath string) (string, error) {
	e.Lock()
	defer e.Unlock()
	return e.addSection(parentFilename, body, sectionTitle, internalFilename, internalCSSPath)
}

func (e *Epub) addSection(parentFilename string, body string, sectionTitle string, internalFilename string, internalCSSPath string) (string, error) {
	parentIndex := -1

	// Generate a filename if one isn't provided
	if internalFilename == "" {
		index := 1
		for internalFilename == "" {
			internalFilename = fmt.Sprintf(sectionFileFormat, index)
			for item, section := range e.sections {
				if section.filename == parentFilename {
					parentIndex = item
				}
				if section.filename == internalFilename {
					internalFilename, index = "", index+1
					if parentFilename == "" || parentIndex != -1 {
						break
					}
				}
				// Check for nested sections with the same filename to avoid duplicate entries
				if section.children != nil {
					for _, subsection := range *section.children {
						if subsection.filename == internalFilename {
							internalFilename, index = "", index+1
						}
					}
				}
			}
		}
	} else {
		for item, section := range e.sections {
			if section.filename == parentFilename {
				parentIndex = item
			}
			if section.filename == internalFilename {
				return "", &FilenameAlreadyUsedError{Filename: internalFilename}
			}
			if section.children != nil {
				for _, subsection := range *section.children {
					if subsection.filename == internalFilename {
						return "", &FilenameAlreadyUsedError{Filename: internalFilename}
					}
				}
			}
		}
	}

	if parentFilename != "" && parentIndex == -1 {
		return "", &ParentDoesNotExistError{Filename: parentFilename}
	}

	x := newXhtml(body)
	x.setTitle(sectionTitle)
	x.setXmlnsEpub(xmlnsEpub)

	if internalCSSPath != "" {
		x.setCSS(internalCSSPath)
	}

	s := epubSection{
		filename: internalFilename,
		xhtml:    x,
		children: nil,
	}

	if parentIndex != -1 {
		if e.sections[parentIndex].children == nil {
			var section []epubSection
			e.sections[parentIndex].children = &section
		}
		(*e.sections[parentIndex].children) = append(*e.sections[parentIndex].children, s)
	} else {
		e.sections = append(e.sections, s)
	}

	return internalFilename, nil
}

// Author returns the author of the EPUB.
func (e *Epub) Author() string {
	return e.author
}

// Identifier returns the unique identifier of the EPUB.
func (e *Epub) Identifier() string {
	return e.identifier
}

// Lang returns the language of the EPUB.
func (e *Epub) Lang() string {
	return e.lang
}

// Description returns the description of the EPUB.
func (e *Epub) Description() string {
	return e.desc
}

// Ppd returns the page progression direction of the EPUB.
func (e *Epub) Ppd() string {
	return e.ppd
}

// SetAuthor sets the author of the EPUB.
func (e *Epub) SetAuthor(author string) {
	e.Lock()
	defer e.Unlock()
	e.author = author
	e.pkg.setAuthor(author)
}

// SetCover sets the cover page for the EPUB using the provided image source and
// optional CSS.
//
// The internal path to an already-added image file (as returned by AddImage) is
// required.
//
// The internal path to an already-added CSS file (as returned by AddCSS) to be
// used for the cover is optional. If the CSS path isn't provided, default CSS
// will be used.
func (e *Epub) SetCover(internalImagePath string, internalCSSPath string) {
	e.Lock()
	defer e.Unlock()
	// If a cover already exists
	if e.cover.xhtmlFilename != "" {
		// Remove the xhtml file
		for i, section := range e.sections {
			if section.filename == e.cover.xhtmlFilename {
				e.sections = append(e.sections[:i], e.sections[i+1:]...)
				break
			}
		}

		// Remove the image
		delete(e.images, e.cover.imageFilename)

		// Remove the CSS
		delete(e.css, e.cover.cssFilename)

		if e.cover.cssTempFile != "" {
			os.Remove(e.cover.cssTempFile)
		}
	}

	e.cover.imageFilename = filepath.Base(internalImagePath)
	e.pkg.setCover(e.cover.imageFilename)

	// Use default cover stylesheet if one isn't provided
	if internalCSSPath == "" {
		// Encode the default CSS
		e.cover.cssTempFile = dataurl.EncodeBytes([]byte(defaultCoverCSSContent))
		var err error
		internalCSSPath, err = e.addCSS(e.cover.cssTempFile, defaultCoverCSSFilename)
		// If that doesn't work, generate a filename
		if _, ok := err.(*FilenameAlreadyUsedError); ok {
			coverCSSFilename := fmt.Sprintf(
				cssFileFormat,
				len(e.css)+1,
				".css",
			)

			internalCSSPath, err = e.addCSS(e.cover.cssTempFile, coverCSSFilename)
			if _, ok := err.(*FilenameAlreadyUsedError); ok {
				// This shouldn't cause an error
				panic(fmt.Sprintf("Error adding default cover CSS file: %s", err))
			}
		}
		if err != nil {
			if _, ok := err.(*FilenameAlreadyUsedError); !ok {
				panic(fmt.Sprintf("DEBUG %+v", err))
			}
		}
	}
	e.cover.cssFilename = filepath.Base(internalCSSPath)

	coverBody := fmt.Sprintf(defaultCoverBody, internalImagePath)
	// Title won't be used since the cover won't be added to the TOC
	// First try to use the default cover filename
	coverPath, err := e.addSection("", coverBody, "", defaultCoverXhtmlFilename, internalCSSPath)
	// If that doesn't work, generate a filename
	if _, ok := err.(*FilenameAlreadyUsedError); ok {
		coverPath, err = e.addSection("", coverBody, "", "", internalCSSPath)
		if _, ok := err.(*FilenameAlreadyUsedError); ok {
			// This shouldn't cause an error since we're not specifying a filename
			panic(fmt.Sprintf("Error adding default cover XHTML file: %s", err))
		}
	}
	e.cover.xhtmlFilename = filepath.Base(coverPath)
}

// SetIdentifier sets the unique identifier of the EPUB, such as a UUID, DOI,
// ISBN or ISSN. If no identifier is set, a UUID will be automatically
// generated.
func (e *Epub) SetIdentifier(identifier string) {
	e.Lock()
	defer e.Unlock()
	e.identifier = identifier
	e.pkg.setIdentifier(identifier)
	e.toc.setIdentifier(identifier)
}

// SetLang sets the language of the EPUB.
func (e *Epub) SetLang(lang string) {
	e.Lock()
	defer e.Unlock()
	e.lang = lang
	e.pkg.setLang(lang)
}

// SetDescription sets the description of the EPUB.
func (e *Epub) SetDescription(desc string) {
	e.Lock()
	defer e.Unlock()
	e.desc = desc
	e.pkg.setDescription(desc)
}

// SetPpd sets the page progression direction of the EPUB.
func (e *Epub) SetPpd(direction string) {
	e.Lock()
	defer e.Unlock()
	e.ppd = direction
	e.pkg.setPpd(direction)
}

// SetTitle sets the title of the EPUB.
func (e *Epub) SetTitle(title string) {
	e.Lock()
	defer e.Unlock()
	e.title = title
	e.pkg.setTitle(title)
	e.toc.setTitle(title)
}

// Title returns the title of the EPUB.
func (e *Epub) Title() string {
	return e.title
}

// Add a media file to the EPUB and return the path relative to the EPUB section
// files
func addMedia(client *http.Client, source string, internalFilename string, mediaFileFormat string, mediaFolderName string, mediaMap map[string]string) (string, error) {
	err := grabber{client}.checkMedia(source)
	if err != nil {
		return "", &FileRetrievalError{
			Source: source,
			Err:    err,
		}
	}
	if internalFilename == "" {
		// If a filename isn't provided, use the filename from the source
		internalFilename = filepath.Base(source)
		_, ok := mediaMap[internalFilename]
		// if filename is too long, invalid or already used, try to generate a unique filename
		if len(internalFilename) > 255 || !fs.ValidPath(internalFilename) || ok {
			internalFilename = fmt.Sprintf(
				mediaFileFormat,
				len(mediaMap)+1,
				strings.ToLower(filepath.Ext(source)),
			)
		}
	}

	if _, ok := mediaMap[internalFilename]; ok {
		return "", &FilenameAlreadyUsedError{Filename: internalFilename}
	}

	mediaMap[internalFilename] = source

	return path.Join(
		"..",
		mediaFolderName,
		internalFilename,
	), nil
}
