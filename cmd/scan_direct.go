package cmd

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Ayanrajpoot10/bugscanx-go/pkg/queuescanner"
	"github.com/spf13/cobra"
)

var scanDirectCmd = &cobra.Command{
	Use:   "direct",
	Short: "Scan using direct connection",
	Run:   scanDirectRun,
}

var (
	scanDirectFlagFilename string
	scanDirectFlagHttps    bool
	scanDirectFlagTimeout  int
	scanDirectFlagOutput   string
	scanDirectFlagMethod   string
)

func init() {
	scanCmd.AddCommand(scanDirectCmd)

	scanDirectCmd.Flags().StringVarP(&scanDirectFlagFilename, "filename", "f", "", "domain list filename")
	scanDirectCmd.Flags().StringVarP(&scanDirectFlagMethod, "method", "m", "HEAD", "HTTP method (e.g. GET, POST)")
	scanDirectCmd.Flags().IntVar(&scanDirectFlagTimeout, "timeout", 3, "connect timeout")
	scanDirectCmd.Flags().StringVarP(&scanDirectFlagOutput, "output", "o", "", "output result")
	scanDirectCmd.Flags().BoolVar(&scanDirectFlagHttps, "https", false, "use https")
	
	scanDirectCmd.MarkFlagRequired("filename")
}

type scanDirectRequest struct {
	Domain string
	Https  bool
	Method string
}

type scanDirectResponse struct {
	Request    *scanDirectRequest
	NetIPList  []net.IP
	StatusCode int
	Server     string
	Location   string
}

var httpClient = &http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	},
	Timeout: 10 * time.Second,
}

func scanDirect(c *queuescanner.Ctx, p *queuescanner.QueueScannerScanParams) {
	req := p.Data.(*scanDirectRequest)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	netIPList, err := net.DefaultResolver.LookupIP(ctxTimeout, "ip4", req.Domain)
	if err != nil {
		return
	}

	httpScheme := "http"
	if req.Https {
		httpScheme = "https"
	}

	httpReq, err := http.NewRequestWithContext(ctxTimeout, req.Method, httpScheme+"://"+req.Domain, nil)
	if err != nil {
		return
	}

	httpRes, err := httpClient.Do(httpReq)
	if err != nil {
		return
	}

	if httpRes.StatusCode == 302 {
		if location := httpRes.Header.Get("Location"); location == "https://jio.com/BalanceExhaust" {
			return
		}
	}

	res := &scanDirectResponse{
		Request:    req,
		NetIPList:  netIPList,
		StatusCode: httpRes.StatusCode,
		Server:     httpRes.Header.Get("Server"),
		Location:   httpRes.Header.Get("Location"),
	}
	c.ScanSuccess(res, nil)

	c.Log(fmt.Sprintf("%-15s  %-3d  %-15s    %s", netIPList[0], httpRes.StatusCode, httpRes.Header.Get("Server"), req.Domain))

}

func scanDirectRun(cmd *cobra.Command, args []string) {
	domainSet := make(map[string]struct{})

	domainListFile, err := os.Open(scanDirectFlagFilename)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer domainListFile.Close()

	scanner := bufio.NewScanner(domainListFile)
	for scanner.Scan() {
		domain := strings.TrimSpace(scanner.Text())
		if domain != "" {
			domainSet[domain] = struct{}{}
		}
	}

	domains := make([]string, 0, len(domainSet))
	for domain := range domainSet {
		domains = append(domains, domain)
	}

	queueScanner := queuescanner.NewQueueScanner(scanFlagThreads, scanDirect)

	fmt.Printf("%-15s  %-3s  %-15s    %-20s\n", "IP", "CODE", "SERVER", "HOST")
	fmt.Printf("%-15s  %-3s  %-15s    %-20s\n", "----", "----", "------", "----")

	scanParams := make([]*queuescanner.QueueScannerScanParams, 0, len(domains))
	for _, domain := range domains {
		scanParams = append(scanParams, &queuescanner.QueueScannerScanParams{
			Name: domain,
			Data: &scanDirectRequest{
				Domain: domain,
				Https:  scanDirectFlagHttps,
				Method: scanDirectFlagMethod,
			},
		})
	}

	queueScanner.Add(scanParams...)

	queueScanner.Start(func(c *queuescanner.Ctx) {
		successList := c.ScanSuccessList.Load()
		if len(*successList) == 0 {
			return
		}

		c.Log("")

		if scanDirectFlagOutput != "" {
			var outputContent strings.Builder
			for _, success := range *successList {
				res, ok := success.(*scanDirectResponse)
				if !ok {
					continue
				}
				outputContent.WriteString(res.Request.Domain)
				outputContent.WriteByte('\n')
			}

			err := os.WriteFile(scanDirectFlagOutput, []byte(outputContent.String()), 0644)
			if err != nil {
				fmt.Printf("Failed to write output file: %v\n", err)
			}
		}
	})
}