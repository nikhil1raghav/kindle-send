package epub

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"
	"github.com/vincent-petithory/dataurl"
)

// grabber is a top level structure that allows a custom http client.
// if onlyChecl is true, the methods will not perform actual grab to spare memory and bandwidth
type grabber struct {
	*http.Client
}

func (g grabber) checkMedia(mediaSource string) error {
	fetchErrors := make([]error, 0)
	for _, f := range []func(string, bool) (io.ReadCloser, error){
		g.localHandler,
		g.httpHandler,
		g.dataURLHandler,
	} {
		var err error
		source, err := f(mediaSource, true)
		if source != nil {
			source.Close()
		}
		if err == nil {
			return nil
		}
		fetchErrors = append(fetchErrors, err)
	}
	return &FileRetrievalError{Source: mediaSource, Err: fetchError(fetchErrors)}
}

// fetchMedia from mediaSource into mediaFolderPath as mediaFilename returning its type.
// the mediaSource can be a URL, a local path or an inline dataurl (as specified in RFC 2397)
func (g grabber) fetchMedia(mediaSource, mediaFolderPath, mediaFilename string) (mediaType string, err error) {

	mediaFilePath := filepath.Join(
		mediaFolderPath,
		mediaFilename,
	)
	// failfast, create the output file handler at the begining, if we cannot write the file, bail out
	w, err := filesystem.Create(mediaFilePath)
	if err != nil {
		return "", fmt.Errorf("unable to create file %s: %s", mediaFilePath, err)
	}
	defer w.Close()
	var source io.ReadCloser
	fetchErrors := make([]error, 0)
	for _, f := range []func(string, bool) (io.ReadCloser, error){
		g.localHandler,
		g.httpHandler,
		g.dataURLHandler,
	} {
		var err error
		source, err = f(mediaSource, false)
		if err != nil {
			fetchErrors = append(fetchErrors, err)
			continue
		}
		break
	}
	if source == nil {
		return "", &FileRetrievalError{Source: mediaSource, Err: fetchError(fetchErrors)}

	}
	defer source.Close()

	_, err = io.Copy(w, source)
	if err != nil {
		// There shouldn't be any problem with the writer, but the reader
		// might have an issue
		return "", &FileRetrievalError{Source: mediaSource, Err: err}
	}

	// Detect the mediaType
	r, err := filesystem.Open(mediaFilePath)
	if err != nil {
		return "", err
	}
	defer r.Close()
	mime, err := mimetype.DetectReader(r)
	if err != nil {
		panic(err)
	}

	// Is it CSS?
	mtype := mime.String()
	if mime.Is("text/plain") {
		if filepath.Ext(mediaSource) == ".css" || filepath.Ext(mediaFilename) == ".css" {
			mtype = "text/css"
		}
	}
	return mtype, nil
}

func (g grabber) httpHandler(mediaSource string, onlyCheck bool) (io.ReadCloser, error) {
	var resp *http.Response
	var err error
	if onlyCheck {
		resp, err = g.Head(mediaSource)
	} else {
		resp, err = g.Get(mediaSource)
	}
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > 400 {
		return nil, errors.New("cannot get file, bad return code")
	}
	return resp.Body, nil
}

func (g grabber) localHandler(mediaSource string, onlyCheck bool) (io.ReadCloser, error) {
	if onlyCheck {
		if _, err := os.Stat(mediaSource); os.IsNotExist(err) {
			return nil, err
		}
		return nil, nil
	}
	return os.Open(mediaSource)
}

func (g grabber) dataURLHandler(mediaSource string, onlyCheck bool) (io.ReadCloser, error) {
	if onlyCheck {
		_, err := dataurl.DecodeString(mediaSource)
		return nil, err
	}
	data, err := dataurl.DecodeString(mediaSource)
	if err != nil {
		return nil, err
	}
	return ioutil.NopCloser(bytes.NewReader(data.Data)), nil
}

type fetchError []error

func (f fetchError) Error() string {
	var message string
	for _, err := range f {
		message = fmt.Sprintf("%v\n %v", message, err.Error())
	}
	return message
}
