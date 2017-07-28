package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(downloadCommand)
}

const getUrl = "http://i18n.ml.com/apps/%s/translations"

var downloadCommand = &cobra.Command{
	Use:   "download",
	Short: "Download the message bundle.",
	Long:  "Download the message bundle with all translations from Babel",
	Run: func(cmd *cobra.Command, args []string) {

		resp, err := http.Get(fmt.Sprintf(getUrl, flag(cmd, "app")))
		assert(err, "failed to get the messages bundle")
		defer resp.Body.Close()

		bundle, err := os.Create(flag(cmd, "bundle"))
		assert(err, "failed to create the messages output file")
		defer bundle.Close()

		io.Copy(bundle, resp.Body)
	},
}
