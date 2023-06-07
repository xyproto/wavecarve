package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// DownloadImage downloads the image from the internet and saves it into a temporary file.
func DownloadImage(url string) (*os.File, error) {
	// Retrieve the url and decode the response body.
	res, err := http.Get(url)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to download image file from URI: %s, status %v", url, res.Status))
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to read response body: %s", err))
	}

	tmpfile, err := ioutil.TempFile("/tmp", "image")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to create temporary file: %v", err))
	}

	// Copy the image binary data into the temporary file.
	_, err = io.Copy(tmpfile, bytes.NewBuffer(data))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to copy the source URI into the destination file"))
	}
	ctype, err := DetectContentType(tmpfile.Name())
	if err != nil {
		return nil, err
	}
	if !strings.Contains(ctype.(string), "image") {
		return nil, errors.New("the downloaded file is a valid image type.")
	}
	return tmpfile, nil
}

// IsValidUrl tests a string to determine if it is a well-structured url or not.
func IsValidUrl(uri string) bool {
	_, err := url.ParseRequestURI(uri)
	if err != nil {
		return false
	}

	u, err := url.Parse(uri)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

// DetectContentType detects the file type by reading MIME type information of the file content.
func DetectContentType(fname string) (interface{}, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return nil, err
	}

	// Reset the read pointer if necessary.
	file.Seek(0, 0)

	// Always returns a valid content-type and "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return string(contentType), nil
}
