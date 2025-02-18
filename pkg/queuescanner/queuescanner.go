package queuescanner

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	terminal "github.com/wayneashleyberry/terminal-dimensions"
)

type Ctx struct {
	ScanSuccessList atomic.Pointer[[]interface{}]
	ScanFailedList  atomic.Pointer[[]interface{}]
	ScanComplete    int32
	dataList        []*QueueScannerScanParams
	context.Context
}

func NewCtx() *Ctx {
	successList := []interface{}{}
	failedList := []interface{}{}
	ctx := &Ctx{}
	ctx.ScanSuccessList.Store(&successList)
	ctx.ScanFailedList.Store(&failedList)
	return ctx
}

func (c *Ctx) Log(a ...interface{}) {
	fmt.Printf("\r\033[2K%s\n", fmt.Sprint(a...))
}

func (c *Ctx) Logf(f string, a ...interface{}) {
	c.Log(fmt.Sprintf(f, a...))
}

func (c *Ctx) LogReplace(a ...string) {
	scanSuccess := len(*c.ScanSuccessList.Load())
	scanFailed := len(*c.ScanFailedList.Load())
	scanComplete := atomic.LoadInt32(&c.ScanComplete)
	totalTasks := len(c.dataList)

	var scanCompletePercentage float64
	if totalTasks > 0 {
		scanCompletePercentage = float64(scanComplete) / float64(totalTasks) * 100
	}

	s := fmt.Sprintf(
		"  %.2f%% - C: %d / %d - S: %d - F: %d - %s", scanCompletePercentage, scanComplete, totalTasks, scanSuccess, scanFailed, strings.Join(a, " "),
	)

	if termWidth, _, err := terminal.Dimensions(); err == nil {
		if w := int(termWidth) - 3; len(s) >= w {
			s = s[:w] + "..."
		}
	}

	fmt.Print("\r\033[2K", s, "\r")
}

func (c *Ctx) LogReplacef(f string, a ...interface{}) {
	c.LogReplace(fmt.Sprintf(f, a...))
}

func (c *Ctx) ScanSuccess(a interface{}, fn func()) {
	if fn != nil {
		fn()
	}
	newList := append(*c.ScanSuccessList.Load(), a)
	c.ScanSuccessList.Store(&newList)
}

func (c *Ctx) ScanFailed(a interface{}, fn func()) {
	if fn != nil {
		fn()
	}
	newList := append(*c.ScanFailedList.Load(), a)
	c.ScanFailedList.Store(&newList)
}

type QueueScannerScanParams struct {
	Name string
	Data interface{}
}
type QueueScannerScanFunc func(c *Ctx, a *QueueScannerScanParams)
type QueueScannerDoneFunc func(c *Ctx)

type QueueScanner struct {
	threads  int
	scanFunc QueueScannerScanFunc
	queue    chan *QueueScannerScanParams
	wg       sync.WaitGroup
	ctx      *Ctx
}

func NewQueueScanner(threads int, scanFunc QueueScannerScanFunc) *QueueScanner {
	t := &QueueScanner{
		threads:  threads,
		scanFunc: scanFunc,
		queue:    make(chan *QueueScannerScanParams, threads*2),
		ctx:      NewCtx(),
	}

	for i := 0; i < t.threads; i++ {
		t.wg.Add(1)
		go t.run()
	}

	return t
}

func (s *QueueScanner) run() {
	defer s.wg.Done()
	for a := range s.queue {
		s.ctx.LogReplace(a.Name)
		s.scanFunc(s.ctx, a)
		atomic.AddInt32(&s.ctx.ScanComplete, 1)
		s.ctx.LogReplace(a.Name)
	}
}

func (s *QueueScanner) Add(dataList ...*QueueScannerScanParams) {
	s.ctx.dataList = append(s.ctx.dataList, dataList...)
}

func (s *QueueScanner) Start(doneFunc QueueScannerDoneFunc) {
	go func() {
		for _, data := range s.ctx.dataList {
			s.queue <- data
		}
		close(s.queue)
	}()

	s.wg.Wait()

	if doneFunc != nil {
		doneFunc(s.ctx)
	}
}
