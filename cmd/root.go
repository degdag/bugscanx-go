package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spf13/pflag"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bugscanx-go",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Set custom help function
	rootCmd.SetHelpFunc(customHelp)
}

var (
	colorD1 = color.New()
	colorB1 = color.New(color.FgHiBlack)
	colorW1 = color.New(color.FgWhite, color.Bold)
	colorG1 = color.New(color.FgGreen, color.Bold)
	colorC1 = color.New(color.FgCyan, color.Bold)
	colorC2 = color.New(color.FgHiCyan)
	colorY1 = color.New(color.FgYellow, color.Bold)
	colorY2 = color.New(color.FgHiYellow)
	colorM1 = color.New(color.FgMagenta, color.Bold)
	colorM2 = color.New(color.FgHiMagenta)
	colorR1 = color.New(color.FgRed, color.Bold)
)

// customHelp provides a custom help menu with colors and examples.
func customHelp(cmd *cobra.Command, args []string) {
	fmt.Println()

	colorY1.Println("Usage:")
	colorD1.Printf("  %s\n\n", cmd.UseLine())

	colorY1.Println("Available Commands:")
	for _, c := range cmd.Commands() {
		if !c.Hidden {
			colorG1.Printf("  %s\t", c.Name())
			colorD1.Printf("%s\n", c.Short)
		}
	}

	colorY1.Println("\nAvailable Modes:")
	colorG1.Printf("  %s\t", "direct")
	colorD1.Printf("%s\n", "Run the scan in direct mode.")

	colorG1.Printf("  %s\t", "proxy")
	colorD1.Printf("%s\n", "Run the scan in proxy mode.")

	colorG1.Printf("  %s\t", "cdn-ssl")
	colorD1.Printf("%s\n", "Run the scan in cdn-ssl mode.")

	colorG1.Printf("  %s\t", "SNI")
	colorD1.Printf("%s\n", "Run the scan in SNI mode.")

	colorG1.Printf("  %s\t", "ping")
	colorD1.Printf("%s\n", "Run the scan in ping mode.")

	fmt.Println()
	colorY1.Println("Examples:")
	colorD1.Println("  Scan using direct mode:")
	colorG1.Println("    bugscanx-go scan direct -f filename.txt")
	colorD1.Println("  Show detailed help for the scan command:")
	colorG1.Println("    bugscanx-go scan --help")

	fmt.Println()
	colorY1.Println("Flags:")
	cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
		if f.Name != "config" && f.Name != "toggle" {
			colorM1.Printf("  --%s\t", f.Name)
			colorD1.Printf("%s\n", f.Usage)
		}
	})

	fmt.Println()
	colorY1.Println("Use \"bugscanx-go [command] --help\" for more information about a command.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".bugscanx-go" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".bugscanx-go")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}