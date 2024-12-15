package cmd

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/Ayanrajpoot10/bugscanx-go/pkg/queuescanner"
)

// dnsCmd represents the dns subcommand
var dnsCmd = &cobra.Command{
	Use:   "dns",
	Short: "Check DNS servers specified in a file",
	Long: `Check DNS servers specified in a file and output the results.
You can specify the number of concurrent workers.`,
	Run: dnsRun,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if dnsFlagFilename == "" {
			return fmt.Errorf("filename is required")
		}
		if dnsFlagWorkers <= 0 {
			return fmt.Errorf("workers must be greater than 0")
		}
		return nil
	},
}

var (
	dnsFlagFilename string
	dnsFlagWorkers  int
	dnsFlagOutput   string
)

func init() {
	scanCmd.AddCommand(dnsCmd)

	dnsCmd.Flags().StringVarP(&dnsFlagFilename, "filename", "f", "", "File containing DNS servers to check (required)")
	dnsCmd.Flags().IntVar(&dnsFlagWorkers, "workers", 10, "Number of concurrent workers (must be greater than 0)")
	dnsCmd.Flags().StringVarP(&dnsFlagOutput, "output", "o", "", "File to save the results (optional)")

	dnsCmd.MarkFlagRequired("filename")
}

func dnsRun(cmd *cobra.Command, args []string) {
	if dnsFlagFilename == "" {
		fmt.Println("Please specify a file using the -f flag.")
		return
	}

	dnsAddresses, err := readDNSFromFile(dnsFlagFilename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	scanner := queuescanner.NewQueueScanner(dnsFlagWorkers, checkDNS)
	for _, address := range dnsAddresses {
		scanner.Add(&queuescanner.QueueScannerScanParams{Name: address, Data: address})
	}

	scanner.Start(func(ctx *queuescanner.Ctx) {
		fmt.Printf("Success: %d, Failed: %d\n", len(ctx.ScanSuccessList), len(ctx.ScanFailedList))
		if dnsFlagOutput != "" {
			successList := make([]string, len(ctx.ScanSuccessList))
			for i, v := range ctx.ScanSuccessList {
				successList[i] = v.(string)
			}
			err := saveResultsToFile(dnsFlagOutput, successList)
			if err != nil {
				fmt.Printf("Error saving results to file: %v\n", err)
			}
		}
	})
}

// readDNSFromFile reads DNS addresses from a file
func readDNSFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var addresses []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		address := scanner.Text()
		if address != "" {
			addresses = append(addresses, address)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return addresses, nil
}

func buildDNSQuery(localIP string) ([]byte, error) {
	packetID := uint16(0x1234)
	flags := uint16(0x0100)
	questions := uint16(1)
	answerRRs := uint16(0)
	authorityRRs := uint16(0)
	additionalRRs := uint16(0)

	parts := strings.Split(localIP, ".")
	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}
	reverseIP := fmt.Sprintf("%s.in-addr.arpa", strings.Join(parts, "."))

	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, packetID)
	binary.Write(&buffer, binary.BigEndian, flags)
	binary.Write(&buffer, binary.BigEndian, questions)
	binary.Write(&buffer, binary.BigEndian, answerRRs)
	binary.Write(&buffer, binary.BigEndian, authorityRRs)
	binary.Write(&buffer, binary.BigEndian, additionalRRs)

	for _, part := range strings.Split(reverseIP, ".") {
		buffer.WriteByte(byte(len(part)))
		buffer.WriteString(part)
	}
	buffer.WriteByte(0)
	binary.Write(&buffer, binary.BigEndian, uint16(12)) // QTYPE=PTR
	binary.Write(&buffer, binary.BigEndian, uint16(1))  // QCLASS=IN

	return buffer.Bytes(), nil
}

func getLocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

func checkDNS(ctx *queuescanner.Ctx, params *queuescanner.QueueScannerScanParams) {
	address := params.Data.(string)
	localIP, err := getLocalIP()
	if err != nil {
		fmt.Printf("[!] Error getting local IP: %v\n", err)
		ctx.ScanFailed(address, func() {
			ctx.Log(fmt.Sprintf("DNS %s is not working: %v", address, err))
		})
		return
	}

	if checkDNSServer(address, localIP, 2*time.Second) {
		ctx.ScanSuccess(address, func() {
			ctx.Log(colorG1.Sprintf("DNS %s is respoding", address))
		})
	} else {
		ctx.ScanFailed(address, func() {
			ctx.Log(colorB1.Sprintf("DNS %s is not respoding", address))
		})
	}
}

func checkDNSServer(dnsServer, localIP string, timeout time.Duration) bool {
	conn, err := net.DialTimeout("udp", net.JoinHostPort(dnsServer, "53"), timeout)
	if err != nil {
		colorB1.Printf("[-] %s not working\n", dnsServer)
		return false
	}
	defer conn.Close()

	query, err := buildDNSQuery(localIP)
	if err != nil {
		colorB1.Printf("[!] Error building DNS query for %s: %v\n", dnsServer, err)
		return false
	}

	_, err = conn.Write(query)
	if err != nil {
		colorB1.Printf("[!] Error sending DNS query to %s: %v\n", dnsServer, err)
		return false
	}

	buffer := make([]byte, 512)
	conn.SetReadDeadline(time.Now().Add(timeout))
	_, err = conn.Read(buffer)
	if err != nil {
		colorB1.Printf("[-] %s not working\n", dnsServer)
		return false
	}

	colorG1.Printf("[+] %s is working\n", dnsServer)
	return true
}

// saveResultsToFile saves the successful DNS addresses to a file
func saveResultsToFile(filename string, results []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, result := range results {
		_, err := writer.WriteString(result + "\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}