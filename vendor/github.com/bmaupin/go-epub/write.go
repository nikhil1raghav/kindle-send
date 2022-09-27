package epub

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"unicode"
	"unicode/utf8"

	"github.com/gofrs/uuid"
)

// UnableToCreateEpubError is thrown by Write if it cannot create the destination EPUB file
type UnableToCreateEpubError struct {
	Path string // The path that was given to Write to create the EPUB
	Err  error  // The underlying error that was thrown
}

func (e *UnableToCreateEpubError) Error() string {
	return fmt.Sprintf("Error creating EPUB at %q: %+v", e.Path, e.Err)
}

const (
	containerFilename     = "container.xml"
	containerFileTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<container version="1.0" xmlns="urn:oasis:names:tc:opendocument:xmlns:container">
  <rootfiles>
    <rootfile full-path="%s/%s" media-type="application/oebps-package+xml" />
  </rootfiles>
</container>
`
	// This seems to be the standard based on the latest EPUB spec:
	// http://www.idpf.org/epub/31/spec/epub-ocf.html
	contentFolderName    = "EPUB"
	coverImageProperties = "cover-image"
	// Permissions for any new directories we create
	dirPermissions = 0755
	// Permissions for any new files we create
	filePermissions   = 0644
	mediaTypeCSS      = "text/css"
	mediaTypeEpub     = "application/epub+zip"
	mediaTypeJpeg     = "image/jpeg"
	mediaTypeNcx      = "application/x-dtbncx+xml"
	mediaTypeXhtml    = "application/xhtml+xml"
	metaInfFolderName = "META-INF"
	mimetypeFilename  = "mimetype"
	pkgFilename       = "package.opf"
	tempDirPrefix     = "go-epub"
	xhtmlFolderName   = "xhtml"
)

// WriteTo the dest io.Writer. The return value is the number of bytes written. Any error encountered during the write is also returned.
func (e *Epub) WriteTo(dst io.Writer) (int64, error) {
	e.Lock()
	defer e.Unlock()
	tempDir := uuid.Must(uuid.NewV4()).String()

	err := filesystem.Mkdir(tempDir, dirPermissions)
	if err != nil {
		panic(fmt.Sprintf("Error creating temp directory: %s", err))
	}
	defer func() {
		if err := filesystem.RemoveAll(tempDir); err != nil {
			panic(fmt.Sprintf("Error removing temp directory: %s", err))
		}
	}()
	writeMimetype(tempDir)
	createEpubFolders(tempDir)

	// Must be called after:
	// createEpubFolders()
	writeContainerFile(tempDir)

	// Must be called after:
	// createEpubFolders()
	err = e.writeCSSFiles(tempDir)
	if err != nil {
		return 0, err
	}

	// Must be called after:
	// createEpubFolders()
	err = e.writeFonts(tempDir)
	if err != nil {
		return 0, err
	}

	// Must be called after:
	// createEpubFolders()
	err = e.writeImages(tempDir)
	if err != nil {
		return 0, err
	}

	// Must be called after:
	// createEpubFolders()
	err = e.writeVideos(tempDir)
	if err != nil {
		return 0, err
	}

	// Must be called after:
	// createEpubFolders()
	e.writeSections(tempDir)

	// Must be called after:
	// createEpubFolders()
	// writeSections()
	e.writeToc(tempDir)

	// Must be called after:
	// createEpubFolders()
	// writeCSSFiles()
	// writeImages()
	// writeVideos()
	// writeSections()
	// writeToc()
	e.writePackageFile(tempDir)
	// Must be called last
	return e.writeEpub(tempDir, dst)
}

// Write writes the EPUB file. The destination path must be the full path to
// the resulting file, including filename and extension.
// The result is always writen to the local filesystem even if the underlying storage is in memory.
func (e *Epub) Write(destFilePath string) error {

	f, err := os.Create(destFilePath)
	if err != nil {
		return &UnableToCreateEpubError{
			Path: destFilePath,
			Err:  err,
		}
	}
	defer f.Close()
	_, err = e.WriteTo(f)
	return err
}

// Create the EPUB folder structure in a temp directory
func createEpubFolders(rootEpubDir string) {
	if err := filesystem.Mkdir(
		filepath.Join(
			rootEpubDir,
			contentFolderName,
		),
		dirPermissions); err != nil {
		// No reason this should happen if tempDir creation was successful
		panic(fmt.Sprintf("Error creating EPUB subdirectory: %s", err))
	}

	if err := filesystem.Mkdir(
		filepath.Join(
			rootEpubDir,
			contentFolderName,
			xhtmlFolderName,
		),
		dirPermissions); err != nil {
		panic(fmt.Sprintf("Error creating xhtml subdirectory: %s", err))
	}

	if err := filesystem.Mkdir(
		filepath.Join(
			rootEpubDir,
			metaInfFolderName,
		),
		dirPermissions); err != nil {
		panic(fmt.Sprintf("Error creating META-INF subdirectory: %s", err))
	}
}

// Write the contatiner file (container.xml), which mostly just points to the
// package file (package.opf)
//
// Sample: https://github.com/bmaupin/epub-samples/blob/master/minimal-v3plus2/META-INF/container.xml
// Spec: http://www.idpf.org/epub/301/spec/epub-ocf.html#sec-container-metainf-container.xml
func writeContainerFile(rootEpubDir string) {
	containerFilePath := filepath.Join(rootEpubDir, metaInfFolderName, containerFilename)
	if err := filesystem.WriteFile(
		containerFilePath,
		[]byte(
			fmt.Sprintf(
				containerFileTemplate,
				contentFolderName,
				pkgFilename,
			),
		),
		filePermissions,
	); err != nil {
		panic(fmt.Sprintf("Error writing container file: %s", err))
	}
}

// Write the CSS files to the temporary directory and add them to the package
// file
func (e *Epub) writeCSSFiles(rootEpubDir string) error {
	err := e.writeMedia(rootEpubDir, e.css, CSSFolderName)
	if err != nil {
		return err
	}

	// Clean up the cover temp file if one was created
	os.Remove(e.cover.cssTempFile)

	return nil
}

// writeCounter counts the number of bytes written to it.
type writeCounter struct {
	Total int64 // Total # of bytes written
}

// Write implements the io.Writer interface.
// Always completes and never returns an error.
func (wc *writeCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += int64(n)
	return n, nil
}

// Write the EPUB file itself by zipping up everything from a temp directory
// The return value is the number of bytes written. Any error encountered during the write is also returned.
func (e *Epub) writeEpub(rootEpubDir string, dst io.Writer) (int64, error) {
	counter := &writeCounter{}
	teeWriter := io.MultiWriter(counter, dst)

	z := zip.NewWriter(teeWriter)

	skipMimetypeFile := false

	// addFileToZip adds the file present at path to the zip archive. The path is relative to the rootEpubDir
	addFileToZip := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Get the path of the file relative to the folder we're zipping
		relativePath, err := filepath.Rel(rootEpubDir, path)
		if err != nil {
			// tempDir and path are both internal, so we shouldn't get here
			return err
		}
		relativePath = filepath.ToSlash(relativePath)

		// Only include regular files, not directories
		info, err := d.Info()
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}

		var w io.Writer
		if filepath.FromSlash(path) == filepath.Join(rootEpubDir, mimetypeFilename) {
			// Skip the mimetype file if it's already been written
			if skipMimetypeFile == true {
				return nil
			}
			// The mimetype file must be uncompressed according to the EPUB spec
			w, err = z.CreateHeader(&zip.FileHeader{
				Name:   relativePath,
				Method: zip.Store,
			})
		} else {
			w, err = z.Create(relativePath)
		}
		if err != nil {
			return fmt.Errorf("error creating zip writer: %w", err)
		}

		r, err := filesystem.Open(path)
		if err != nil {
			return fmt.Errorf("error opening file %v being added to EPUB: %w", path, err)
		}
		defer func() {
			if err := r.Close(); err != nil {
				panic(err)
			}
		}()

		_, err = io.Copy(w, r)
		if err != nil {
			return fmt.Errorf("error copying contents of file being added EPUB: %w", err)
		}
		return nil
	}

	// Add the mimetype file first
	mimetypeFilePath := filepath.Join(rootEpubDir, mimetypeFilename)
	mimetypeInfo, err := fs.Stat(filesystem, mimetypeFilePath)
	if err != nil {
		if err := z.Close(); err != nil {
			panic(err)
		}
		return counter.Total, fmt.Errorf("unable to get FileInfo for mimetype file: %w", err)
	}
	err = addFileToZip(mimetypeFilePath, fileInfoToDirEntry(mimetypeInfo), nil)
	if err != nil {
		if err := z.Close(); err != nil {
			panic(err)
		}
		return counter.Total, fmt.Errorf("unable to add mimetype file to EPUB: %w", err)
	}

	skipMimetypeFile = true

	err = fs.WalkDir(filesystem, rootEpubDir, addFileToZip)
	if err != nil {
		if err := z.Close(); err != nil {
			panic(err)
		}
		return counter.Total, fmt.Errorf("unable to add file to EPUB: %w", err)
	}

	err = z.Close()
	return counter.Total, err
}

// Get fonts from their source and save them in the temporary directory
func (e *Epub) writeFonts(rootEpubDir string) error {
	return e.writeMedia(rootEpubDir, e.fonts, FontFolderName)
}

// Get images from their source and save them in the temporary directory
func (e *Epub) writeImages(rootEpubDir string) error {
	return e.writeMedia(rootEpubDir, e.images, ImageFolderName)
}

// Get videos from their source and save them in the temporary directory
func (e *Epub) writeVideos(rootEpubDir string) error {
	return e.writeMedia(rootEpubDir, e.videos, VideoFolderName)
}

// Get media from their source and save them in the temporary directory
func (e *Epub) writeMedia(rootEpubDir string, mediaMap map[string]string, mediaFolderName string) error {
	if len(mediaMap) > 0 {
		mediaFolderPath := filepath.Join(rootEpubDir, contentFolderName, mediaFolderName)
		if err := filesystem.Mkdir(mediaFolderPath, dirPermissions); err != nil {
			return fmt.Errorf("unable to create directory: %s", err)
		}

		for mediaFilename, mediaSource := range mediaMap {
			mediaType, err := grabber{(e.Client)}.fetchMedia(mediaSource, mediaFolderPath, mediaFilename)
			if err != nil {
				return err
			}
			// The cover image has a special value for the properties attribute
			mediaProperties := ""
			if mediaFilename == e.cover.imageFilename {
				mediaProperties = coverImageProperties
			}

			// Add the file to the OPF manifest
			e.pkg.addToManifest(fixXMLId(mediaFilename), filepath.Join(mediaFolderName, mediaFilename), mediaType, mediaProperties)
		}
	}
	return nil
}

// fixXMLId takes a string and returns an XML id compatible string.
// https://www.w3.org/TR/REC-xml-names/#NT-NCName
// This means it must not contain a colon (:) or whitespace and it must not
// start with a digit, punctuation or diacritics
func fixXMLId(id string) string {
	if len(id) == 0 {
		panic("No id given")
	}
	fixedId := []rune{}
	for i := 0; len(id) > 0; i++ {
		r, size := utf8.DecodeRuneInString(id)
		if i == 0 {
			// The new id should be prefixed with 'id' if an invalid
			// starting character is found
			// this is not 100% accurate, but a better check than no check
			if unicode.IsNumber(r) || unicode.IsPunct(r) || unicode.IsSymbol(r) {
				fixedId = append(fixedId, []rune("id")...)
			}
		}
		if !unicode.IsSpace(r) && r != ':' {
			fixedId = append(fixedId, r)
		}
		id = id[size:]
	}
	return string(fixedId)
}

// Write the mimetype file
//
// Sample: https://github.com/bmaupin/epub-samples/blob/master/minimal-v3plus2/mimetype
// Spec: http://www.idpf.org/epub/301/spec/epub-ocf.html#sec-zip-container-mime
func writeMimetype(rootEpubDir string) {
	mimetypeFilePath := filepath.Join(rootEpubDir, mimetypeFilename)

	if err := filesystem.WriteFile(mimetypeFilePath, []byte(mediaTypeEpub), filePermissions); err != nil {
		panic(fmt.Sprintf("Error writing mimetype file: %s", err))
	}
}

func (e *Epub) writePackageFile(rootEpubDir string) {
	e.pkg.write(rootEpubDir)
}

// Write the section files to the temporary directory and add the sections to
// the TOC and package files
func (e *Epub) writeSections(rootEpubDir string) {
	var index int

	if len(e.sections) > 0 {
		// If a cover was set, add it to the package spine first so it shows up
		// first in the reading order
		if e.cover.xhtmlFilename != "" {
			e.pkg.addToSpine(e.cover.xhtmlFilename)
		}

		for _, section := range e.sections {
			// Set the title of the cover page XHTML to the title of the EPUB
			if section.filename == e.cover.xhtmlFilename {
				section.xhtml.setTitle(e.Title())
			}

			sectionFilePath := filepath.Join(rootEpubDir, contentFolderName, xhtmlFolderName, section.filename)
			section.xhtml.write(sectionFilePath)
			relativePath := filepath.Join(xhtmlFolderName, section.filename)

			// The cover page should have already been added to the spine first
			if section.filename != e.cover.xhtmlFilename {
				e.pkg.addToSpine(section.filename)
			}
			e.pkg.addToManifest(section.filename, relativePath, mediaTypeXhtml, "")

			// Don't add pages without titles or the cover to the TOC
			if section.xhtml.Title() != "" && section.filename != e.cover.xhtmlFilename {
				e.toc.addSection(index, section.xhtml.Title(), relativePath)

				// Add subsections
				if section.children != nil {
					for _, child := range *section.children {
						index += 1
						relativeSubPath := filepath.Join(xhtmlFolderName, child.filename)
						e.toc.addSubSection(relativePath, index, child.xhtml.Title(), relativeSubPath)

						subSectionFilePath := filepath.Join(rootEpubDir, contentFolderName, xhtmlFolderName, child.filename)
						child.xhtml.write(subSectionFilePath)

						// Add subsection to spine
						e.pkg.addToSpine(child.filename)
						e.pkg.addToManifest(child.filename, relativeSubPath, mediaTypeXhtml, "")
					}
				}
			}

			index += 1
		}
	}
}

// Write the TOC file to the temporary directory and add the TOC entries to the
// package file
func (e *Epub) writeToc(rootEpubDir string) {
	e.pkg.addToManifest(tocNavItemID, tocNavFilename, mediaTypeXhtml, tocNavItemProperties)
	e.pkg.addToManifest(tocNcxItemID, tocNcxFilename, mediaTypeNcx, "")

	e.toc.write(rootEpubDir)
}
