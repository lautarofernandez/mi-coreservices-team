package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/mercadolibre/go-meli-toolkit/babel/babel/scanner"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(scanCommand)

}

var scanCommand = &cobra.Command{
	Use:   "scan",
	Short: "Scan project files",
	Long:  `Scan project files to find translation keys`,
	Run: func(cmd *cobra.Command, args []string) {
		scanner := scanner.NewFileScanner()

		err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.Contains(path, "vendor/") {
				return scanner.Scan(path)
			}
			return nil
		})

		assert(err, "filed to scan the project files")
		scanner.Save(flag(cmd, "messages"))
	},
}
