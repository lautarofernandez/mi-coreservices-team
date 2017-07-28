package cmd

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(uploadCommand)
}

const postURL = "http://i18n.ml.com/apps/%s/sources?project_name=%s&force=false"

var uploadCommand = &cobra.Command{
	Use:   "upload",
	Short: "Upload the message files",
	Long:  "Upload the message files to Babel",
	Run: func(cmd *cobra.Command, args []string) {
		body := bytes.NewBuffer(nil)

		bodyWriter := multipart.NewWriter(body)
		defer bodyWriter.Close()

		form, err := bodyWriter.CreateFormFile("sources.zip", "sources.zip")
		assert(err, "failed to create the body form")

		zipfile := zip.NewWriter(form)

		source, err := os.Open(flag(cmd, "messages"))
		assert(err, "failed to read source file")

		file, err := zipfile.Create("messages.po")
		assert(err, "failed to create the zipped source file")

		io.Copy(file, source)

		zipfile.Close()
		bodyWriter.Close()

		url := fmt.Sprintf(postURL, flag(cmd, "app"), flag(cmd, "project"))
		resp, err := http.Post(url, bodyWriter.FormDataContentType(), body)
		assert(err, "failed to post to Babel")
		defer resp.Body.Close()

		io.Copy(os.Stdout, resp.Body)
	},
}
