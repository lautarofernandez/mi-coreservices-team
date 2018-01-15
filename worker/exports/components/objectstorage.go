package components

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/mercadolibre/go-meli-toolkit/restful/rest"
)

var (
	// NoRedirectHTTPClient is an http client that purposefully does not follow redirects.
	NoRedirectHTTPClient = &http.Client{
		Timeout: 2 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// UploadTimeout is the timeout used for file upload operations
	UploadTimeout = 10 * time.Minute

	// UploadsHTTPClient is an http client with a very high Timeout specialized in uploading big files
	UploadsHTTPClient = &http.Client{
		Timeout: UploadTimeout,
	}
)

//ObjectStorage is a Object Storage wrapper
type ObjectStorage struct {
	writeClient *rest.RequestBuilder
}

// NewObjectStorage returns a Object Storage wrapper
func NewObjectStorage(writeClient *rest.RequestBuilder) *ObjectStorage {
	return &ObjectStorage{
		writeClient: writeClient,
	}
}

// UploadFile uploads the file `filename` to S3 using the complete `objectstorageURL`
func (obj *ObjectStorage) UploadFile(f *os.File, filename string) error {
	s3URL, err := GetRedirectURL(obj.writeClient.BaseURL + filename)
	if err != nil {
		return err
	}

	stat, err := f.Stat()
	if err != nil {
		return fmt.Errorf("error executing stat on file: %v", err)
	}

	req, err := http.NewRequest("PUT", s3URL, f)
	if err != nil {
		return fmt.Errorf("error preparing request to %s: %v", s3URL, err)
	}

	// Manually specify content-length. If content-length is not specified then http.Client
	// will initialize a chunked request to S3, which will fail because AWS does not
	// support not knowing the file size beforehand.
	req.ContentLength = stat.Size()

	res, err := UploadsHTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error executing AWS PUT request: %v", err)
	}

	if res.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code received from S3: %d", res.StatusCode)
	}

	return nil
}

// GetRedirectURL makes a dummy request with empty body to object storage API. This will
// validate our access to the underlying S3 bucket, and return a signed AWS url
// that we can later use for uploading the resulting file.
func GetRedirectURL(objectstorageURL string) (string, error) {
	req, err := http.NewRequest("PUT", objectstorageURL, bytes.NewBuffer(nil))
	if err != nil {
		return "", fmt.Errorf("error creating dummy redirect request: %v", err)
	}

	res, err := NoRedirectHTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making dummy redirect request: %v", err)
	}

	if res.StatusCode != 307 {
		return "", fmt.Errorf("unexpected redirect status code received: %s", res.Status)
	}

	// Get new path for file
	redirectURL := res.Header.Get("location")
	if redirectURL == "" {
		return "", fmt.Errorf("could not read location header from response")
	}

	return redirectURL, nil
}
