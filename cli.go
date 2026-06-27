package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

func cliMain() {
	urlStr := os.Args[1]
	conns := 100
	if len(os.Args) >= 3 {
		if n, err := strconv.Atoi(os.Args[2]); err == nil && n > 0 {
			conns = n
		}
	}

	fileSize, rangesSupported, filename, err := fetchFileInfo(urlStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", T("Error"), err)
		os.Exit(1)
	}

	if fileSize > 0 && rangesSupported {
		fmt.Printf("%s: %s (%s), %d %s\n\n",
			T("File"), filename, formatSize(fileSize), conns, T("ParallelConns"))
		startTime := time.Now()
		parallelDownloadCLI(urlStr, filename, fileSize, conns, startTime)
		elapsed := time.Since(startTime)
		speed := formatSize(int64(float64(fileSize)/elapsed.Seconds())) + "/s"
		fmt.Printf("\n%s: %s (%s, %s)\n", T("Done"), filename, formatSize(fileSize), speed)
	} else if fileSize > 0 {
		fmt.Printf("%s: %s (%s), %s\n",
			T("File"), filename, formatSize(fileSize), T("SingleConnection"))
		startTime := time.Now()
		singleDownloadCLI(urlStr, filename, startTime)
		elapsed := time.Since(startTime)
		speed := formatSize(int64(float64(fileSize)/elapsed.Seconds())) + "/s"
		fmt.Printf("\n%s: %s (%s, %s)\n", T("Done"), filename, formatSize(fileSize), speed)
	} else {
		fmt.Fprintf(os.Stderr, "%s: %s\n", T("Error"), T("UnknownFileSize"))
		os.Exit(1)
	}
}

func parallelDownloadCLI(urlStr, filename string, fileSize int64, conns int, startTime time.Time) {
	partFile := filename + ".part"
	f, err := os.Create(partFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", T("Error"), err)
		os.Exit(1)
	}
	f.Truncate(fileSize)
	f.Close()

	f, err = os.OpenFile(partFile, os.O_RDWR, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", T("Error"), err)
		os.Exit(1)
	}
	defer f.Close()

	numChunks := conns * 3
	if int64(numChunks) > fileSize/65536 {
		numChunks = int(fileSize / 65536)
		if numChunks < conns {
			numChunks = conns
		}
	}

	chunkSize := fileSize / int64(numChunks)
	type workItem struct {
		start int64
		end   int64
	}

	work := make(chan workItem, numChunks)
	for i := 0; i < numChunks; i++ {
		s := int64(i) * chunkSize
		e := s + chunkSize - 1
		if i == numChunks-1 {
			e = fileSize - 1
		}
		work <- workItem{start: s, end: e}
	}
	close(work)

	var wg sync.WaitGroup
	errCh := make(chan error, conns)
	var dld atomic.Int64
	writes := make(chan writeReq, 1024)
	go writeWorker(f, writes, &dld, nil)

	for w := 0; w < conns; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for wi := range work {
				ci := chunkInfo{start: wi.start, end: wi.end}
				if err := downloadChunk(urlStr, ci, writes, nil); err != nil {
					errCh <- err
					return
				}
			}
		}()
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(writes)
		close(done)
	}()

	ticker := time.NewTicker(200 * time.Millisecond)
loop:
	for {
		select {
		case <-done:
			break loop
		case <-ticker.C:
			d := dld.Load()
			elapsed := time.Since(startTime)
			pct := float64(d) * 100.0 / float64(fileSize)
			fmt.Printf("\r%s %5.1f%%  %s / %s  %s  %s",
				progressBar(pct, 30), pct,
				formatSize(d), formatSize(fileSize),
				formatSpeed(d, elapsed), formatETA(d, fileSize, elapsed))
		}
	}
	ticker.Stop()

	select {
	case err := <-errCh:
		fmt.Fprintf(os.Stderr, "\n%s: %v\n", T("Error"), err)
		f.Close()
		os.Remove(partFile)
		os.Exit(1)
	default:
	}

	f.Close()
	os.Rename(partFile, filename)
}

func singleDownloadCLI(urlStr, filename string, startTime time.Time) {
	resp, err := sharedClient.Get(urlStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", T("Error"), err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	f, err := os.Create(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", T("Error"), err)
		os.Exit(1)
	}
	defer f.Close()

	var dld int64
	total := resp.ContentLength
	buf := make([]byte, 1024*1024)
	ticker := time.NewTicker(200 * time.Millisecond)
	done := make(chan struct{})

	go func() {
		for {
			n, rerr := resp.Body.Read(buf)
			if n > 0 {
				f.Write(buf[:n])
				dld += int64(n)
			}
			if rerr == io.EOF {
				close(done)
				return
			}
			if rerr != nil {
				close(done)
				return
			}
		}
	}()

	for {
		select {
		case <-done:
			ticker.Stop()
			return
		case <-ticker.C:
			elapsed := time.Since(startTime)
			s := formatSpeed(dld, elapsed)
			if total > 0 {
				pct := float64(dld) * 100.0 / float64(total)
				fmt.Printf("\r%s %5.1f%%  %s / %s  %s  %s",
					progressBar(pct, 30), pct,
					formatSize(dld), formatSize(total),
					s, formatETA(dld, total, elapsed))
			} else {
				fmt.Printf("\r%s: %s  %s", T("Downloaded"), formatSize(dld), s)
			}
		}
	}
}
