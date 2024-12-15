package cmd

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/Ayanrajpoot10/bugscanx-go/pkg/queuescanner"
)

// pingCmd represents the ping subcommand
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Ping hosts specified in a file",
	Long: `Ping hosts specified in a file and output the results.
You can specify the timeout, number of threads, and output file for the results.`,
	Run: pingRun,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if pingFlagFilename == "" {
			return fmt.Errorf("filename is required")
		}
		if pingFlagTimeout <= 0 {
			return fmt.Errorf("timeout must be greater than 0")
		}
		if pingFlagThreads <= 0 {
			return fmt.Errorf("threads must be greater than 0")
		}
		return nil
	},
}

var (
	pingFlagFilename string
	pingFlagTimeout  int
	pingFlagThreads  int
	pingFlagOutput   string
	pingFlagTCP      bool
	pingFlagPort     int
)

func init() {
	scanCmd.AddCommand(pingCmd)

	pingCmd.Flags().StringVarP(&pingFlagFilename, "filename", "f", "", "File containing hosts to ping (required)")
	pingCmd.Flags().IntVar(&pingFlagTimeout, "timeout", 2, "Ping timeout in seconds (must be greater than 0)")
	pingCmd.Flags().IntVar(&pingFlagThreads, "threads", 64, "Number of threads (must be greater than 0)")
	pingCmd.Flags().StringVarP(&pingFlagOutput, "output", "o", "", "File to write results")
	pingCmd.Flags().BoolVar(&pingFlagTCP, "tcp", false, "Use TCP ping instead of ICMP")
	pingCmd.Flags().IntVar(&pingFlagPort, "port", 80, "Port to use for TCP ping")

	pingCmd.MarkFlagRequired("filename")
}

func pingRun(cmd *cobra.Command, args []string) {
	if pingFlagFilename == "" {
		fmt.Println("Please specify a file using the -f flag.")
		return
	}

	hosts, err := readHostsFromFile(pingFlagFilename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	colorM1.Printf("\n%-15s %-20s\n", "Status", "Host")
	colorW1.Printf("%-15s %-20s\n", "--------", "--------")

	scanner := queuescanner.NewQueueScanner(pingFlagThreads, pingHost)
	for _, host := range hosts {
		scanner.Add(&queuescanner.QueueScannerScanParams{Name: host, Data: host})
	}

	scanner.Start(func(ctx *queuescanner.Ctx) {
		fmt.Printf("Success: %d, Failed: %d\n", len(ctx.ScanSuccessList), len(ctx.ScanFailedList))
		if pingFlagOutput != "" {
			writeResultsToFile(pingFlagOutput, ctx)
		}
	})
}

// readHostsFromFile reads hosts from a file
func readHostsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var hosts []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		host := scanner.Text()
		if host != "" {
			hosts = append(hosts, host)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return hosts, nil
}

// pingHost is the scanning function used by QueueScanner
func pingHost(ctx *queuescanner.Ctx, params *queuescanner.QueueScannerScanParams) {
	host := params.Data.(string)
	if pingFlagTCP {
		pingHostTCP(ctx, host)
	} else {
		pingHostICMP(ctx, host)
	}
}

func pingHostICMP(ctx *queuescanner.Ctx, host string) {
	addr, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		ctx.ScanFailed(host, func() {
			ctx.Log(colorB1.Sprintf("%-15s%-20s", "Not Resolved:", host))
		})
		return
	}

	conn, err := net.DialTimeout("ip4:icmp", addr.String(), time.Duration(pingFlagTimeout)*time.Second)
	if err != nil {
		ctx.ScanFailed(host, func() {
			ctx.Log(colorR1.Sprintf("%-15s%-20s", "failed:", host))
		})
		return
	}
	defer conn.Close()

	ctx.ScanSuccess(host, func() {
		ctx.Log(colorG1.Sprintf("%-15s%-20s", "succeeded:", host))
	})
}

func pingHostTCP(ctx *queuescanner.Ctx, host string) {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, fmt.Sprintf("%d", pingFlagPort)), time.Duration(pingFlagTimeout)*time.Second)
	if err != nil {
		ctx.ScanFailed(host, func() {
			ctx.Log(colorR1.Sprintf("%-15s%-20s", "failed:", host))
		})
		return
	}
	defer conn.Close()

	ctx.ScanSuccess(host, func() {
		ctx.Log(colorG1.Sprintf("%-15s%-20s", "succeeded:", host))
	})
}

// writeResultsToFile writes ping results to the specified output file
func writeResultsToFile(filename string, ctx *queuescanner.Ctx) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for _, success := range ctx.ScanSuccessList {
		writer.WriteString(fmt.Sprintf("%v\n", success))
	}
}
