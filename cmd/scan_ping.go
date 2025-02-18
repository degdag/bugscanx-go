package cmd

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/Ayanrajpoot10/bugscanx-go/pkg/queuescanner"
	"github.com/spf13/cobra"
)

var (
	pingFlagFilename string
	pingFlagTimeout  int
	pingFlagOutput   string
	pingFlagPort     int
)

var scanPingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Scan using TCP ping",
	Run:   pingRun,
}

func init() {
	scanCmd.AddCommand(scanPingCmd)
	scanPingCmd.Flags().StringVarP(&pingFlagFilename, "filename", "f", "", "domain list filename (required)")
	scanPingCmd.Flags().IntVarP(&pingFlagPort, "port", "p", 80, "Port for TCP ping")
	scanPingCmd.Flags().IntVar(&pingFlagTimeout, "timeout", 2, "Ping timeout")
	scanPingCmd.Flags().StringVarP(&pingFlagOutput, "output", "o", "", "output results")
	
	scanPingCmd.MarkFlagRequired("filename")
}

type pingRequest struct {
	Host string
}

func pingHostTCP(ctx *queuescanner.Ctx, params *queuescanner.QueueScannerScanParams) {
	req := params.Data.(*pingRequest)
	addr := net.JoinHostPort(req.Host, fmt.Sprintf("%d", pingFlagPort))
	conn, err := net.DialTimeout("tcp", addr, time.Duration(pingFlagTimeout)*time.Second)
	if err != nil {
		ctx.ScanFailed(req.Host, nil)
		return
	}
	defer conn.Close()
	ctx.ScanSuccess(req.Host, func() { ctx.Log(fmt.Sprintf("UP: %-20s", req.Host)) })
}

func pingRun(cmd *cobra.Command, args []string) {
	domainList := make(map[string]bool)
	domainListFile, err := os.Open(pingFlagFilename)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer domainListFile.Close()

	scanner := bufio.NewScanner(domainListFile)
	for scanner.Scan() {
		domain := scanner.Text()
		domainList[domain] = true
	}

	queueScanner := queuescanner.NewQueueScanner(scanFlagThreads, pingHostTCP)

	fmt.Printf("%-10s %-20s\n------     -----\n", "Status", "Host")

	for domain := range domainList {
		queueScanner.Add(&queuescanner.QueueScannerScanParams{
			Name: domain,
			Data: &pingRequest{Host: domain},
		})
	}

	queueScanner.Start(func(ctx *queuescanner.Ctx) {
		ctx.Log("")
		if pingFlagOutput != "" {
			outputFile, err := os.Create(pingFlagOutput)
			if err != nil {
				fmt.Printf("Failed to create output file: %v\n", err)
				return
			}
			defer outputFile.Close()

			for _, success := range *ctx.ScanSuccessList.Load() {
				_, err := outputFile.WriteString(success.(string) + "\n")
				if err != nil {
					fmt.Printf("Error writing to file: %v\n", err)
					break
				}
			}
		}
	})
}
