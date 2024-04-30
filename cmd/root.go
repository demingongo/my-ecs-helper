package cmd

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/demingongo/my-ecs-helper/app/mainapp"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "my-ecs-helper",
	Short: "AWS ECS helper",
	Long: `
    .----.
    |~_  |
  __|____|__
 |  ______--|
 '-/.::::.\-'
  '--------'
Using AWS EC2/ECS?
Tired of using the slow AWS console?

Here's my personal helper. It uses aws-cli under the hood.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("verbose") {
			log.SetLevel(log.DebugLevel)
		}
		mainapp.Run()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobraapp.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().Bool("dummy", false, "dummy run (no aws call)")

	viper.BindPFlag("dummy", rootCmd.PersistentFlags().Lookup("dummy"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.SetDefault("dummy", false)
	viper.SetDefault("verbose", false)

}
