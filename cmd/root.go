package cmd

import (

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bugscanx-go",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

var (
	colorD1 = color.New()
	colorB1 = color.New(color.FgHiBlack)
	colorG1 = color.New(color.FgGreen, color.Bold)
	colorY1 = color.New(color.FgYellow, color.Bold)
	colorC1 = color.New(color.FgHiCyan, color.Bold)
	colorM1 = color.New(color.FgHiMagenta, color.Bold)
)

func PrintBanner() {
	colorC1.Print("\nWelcome to BugScanX-Go ")
	colorY1.Print("Made by Ayan Rajpoot ")
	colorM1.Print("Telegram Channel: BugScanX\n")
	println()
}