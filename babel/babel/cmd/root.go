package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RootCmd = &cobra.Command{
	Use:   "babel",
	Short: "babel cli",
	Long:  `babel is command line tool to help manage tasks like update, upload and download translations from Babel`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().String("app", "", "github repository name")
	RootCmd.PersistentFlags().String("project", "", "Babel project name")
	RootCmd.PersistentFlags().String("messages", "./conf/messages.po", "messages filename")
	RootCmd.PersistentFlags().String("bundle", "./conf/all.zip", "message bundle with all translations")

	viper.BindPFlags(RootCmd.PersistentFlags())
}

func initConfig() {
	viper.SetConfigName(".babel")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("BABEL")
	viper.AutomaticEnv()
	viper.ReadInConfig()
}

func flag(cmd *cobra.Command, key string) string {
	value := viper.GetString(key)
	if len(value) == 0 {
		fmt.Printf("the --%s is required\n", key)
		cmd.Usage()
		os.Exit(-1)
	}
	return value
}

func assert(err error, message string) {
	if err != nil {
		fmt.Println(message, err)
		os.Exit(-1)
	}
}
