package main

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type DownloadTask struct {
	URL         string
	Filename    string
	FileSize    int64
	Downloaded  int64
	Progress    int
	Speed       string
	TimeElapsed string
	ETA         string
	Status      string
	Connections int
	startTime   time.Time
	cancel      chan struct{}
	partFile    string
}

type DownloadModel struct {
	walk.TableModelBase
	items []*DownloadTask
}

func (m *DownloadModel) RowCount() int {
	return len(m.items)
}

func (m *DownloadModel) Value(row, col int) interface{} {
	item := m.items[row]
	switch col {
	case 0:
		return item.Filename
	case 1:
		return formatSize(item.FileSize)
	case 2:
		if item.FileSize > 0 {
			return float64(item.Progress)
		}
		return float64(0)
	case 3:
		return item.Speed
	case 4:
		return item.TimeElapsed
	case 5:
		return item.ETA
	case 6:
		return item.Status
	case 7:
		return strconv.Itoa(item.Connections)
	}
	return ""
}

var (
	model          = &DownloadModel{}
	urlEntry       *walk.LineEdit
	connSpin       *walk.NumberEdit
	downloadBtn    *walk.PushButton
	cancelBtn      *walk.PushButton
	clearBtn       *walk.PushButton
	tasksView      *walk.TableView
	statusTotal    *walk.Label
	statusActive   *walk.Label
	statusSpeed    *walk.Label
	activeCount    int32
	totalDownloaded int64
)

func guiMain() {
	var mw *walk.MainWindow

	if _, err := (MainWindow{
		AssignTo: &mw,
		Title:    "Memedovski Fast Download Manager v2.0",
		MinSize:  Size{900, 600},
		Size:     Size{1000, 650},
		Layout:   VBox{Margins: Margins{10, 10, 10, 10}, Spacing: 8},
		Children: []Widget{
			GroupBox{
				Title:  "New Download",
				Layout: HBox{Margins: Margins{8, 8, 8, 8}, Spacing: 6},
				Children: []Widget{
					Label{Text: "URL:", Font: Font{PointSize: 9}},
					LineEdit{
						AssignTo: &urlEntry,
						MinSize:  Size{400, 0},
					},
					Label{Text: "Connections:", Font: Font{PointSize: 9}},
					NumberEdit{
						AssignTo: &connSpin,
						Value:    100.0,
						MinValue: 1.0,
						MaxValue: 500.0,
						MinSize:  Size{70, 0},
					},
					PushButton{
						AssignTo: &downloadBtn,
						Text:     "Download",
						MinSize:  Size{100, 30},
						OnClicked: func() {
							startDownload()
						},
					},
				},
			},
			TableView{
				AssignTo: &tasksView,
				Model:    model,
				Columns: []TableViewColumn{
					{Title: "Filename", Width: 250},
					{Title: "Size", Width: 90},
					{Title: "Progress", Width: 80},
					{Title: "Speed", Width: 90},
					{Title: "Time", Width: 70},
					{Title: "ETA", Width: 70},
					{Title: "Status", Width: 90},
					{Title: "Conns", Width: 50},
				},
			},
			Composite{
				Layout: HBox{Margins: Margins{0, 0, 0, 0}, Spacing: 8},
				Children: []Widget{
					PushButton{
						AssignTo: &cancelBtn,
						Text:     "Cancel Selected",
						Enabled:  false,
						OnClicked: func() {
							cancelSelected()
						},
					},
					PushButton{
						AssignTo: &clearBtn,
						Text:     "Clear Completed",
						OnClicked: func() {
							clearCompleted()
						},
					},
					HSpacer{},
					Label{
						Text:     "Total: ",
						Font:     Font{PointSize: 9},
					},
					Label{
						AssignTo: &statusTotal,
						Text:     "0 B",
						Font:     Font{PointSize: 9, Bold: true},
						MinSize:  Size{80, 0},
					},
					Label{
						Text:    "Active: ",
						Font:    Font{PointSize: 9},
					},
					Label{
						AssignTo: &statusActive,
						Text:     "0",
						Font:     Font{PointSize: 9, Bold: true},
						MinSize:  Size{30, 0},
					},
					Label{
						Text:    "Speed: ",
						Font:    Font{PointSize: 9},
					},
					Label{
						AssignTo: &statusSpeed,
						Text:     "0 B/s",
						Font:     Font{PointSize: 9, Bold: true},
						MinSize:  Size{80, 0},
					},
				},
			},
		},
	}.Run()); err != nil {
		fmt.Fprintf(os.Stderr, "GUI error: %v\n", err)
		os.Exit(1)
	}
}

func startDownload() {
	urlStr := urlEntry.Text()
	if urlStr == "" {
		walk.MsgBox(nil, "Error", "Please enter a URL", walk.MsgBoxIconError)
		return
	}

	if _, err := url.ParseRequestURI(urlStr); err != nil {
		walk.MsgBox(nil, "Error", "Invalid URL", walk.MsgBoxIconError)
		return
	}

	conns := int(connSpin.Value())
	if conns < 1 {
		conns = 100
	}

	task := &DownloadTask{
		URL:         urlStr,
		Connections: conns,
		Status:      "Queued",
		cancel:      make(chan struct{}),
	}

	model.items = append(model.items, task)
	model.PublishRowsReset()
	downloadBtn.SetEnabled(false)
	urlEntry.SetEnabled(false)

	go runDownload(task)
}

func runDownload(task *DownloadTask) {
	fileSize, rangesSupported, filename, err := fetchFileInfo(task.URL)
	if err != nil {
		setTaskStatus(task, "Error: "+err.Error())
		return
	}

	task.Filename = filename

	if fileSize > 0 && rangesSupported {
		task.FileSize = fileSize
		updateModel()
		setTaskStatus(task, "Downloading")
		task.startTime = time.Now()
		atomic.AddInt32(&activeCount, 1)
		parallelDownloadGUI(task, fileSize)
		atomic.AddInt32(&activeCount, -1)
	} else if fileSize > 0 {
		task.FileSize = fileSize
		updateModel()
		setTaskStatus(task, "Downloading")
		task.startTime = time.Now()
		atomic.AddInt32(&activeCount, 1)
		singleDownloadGUI(task)
		atomic.AddInt32(&activeCount, -1)
	} else {
		setTaskStatus(task, "Error: Unknown file size")
	}
	updateStatusBar()
}

func parallelDownloadGUI(task *DownloadTask, fileSize int64) {
	conns := task.Connections
	partFile := task.Filename + ".mfdm"
	task.partFile = partFile

	f, err := os.Create(partFile)
	if err != nil {
		setTaskStatus(task, "Error: "+err.Error())
		return
	}
	f.Truncate(fileSize)
	f.Close()

	f, err = os.OpenFile(partFile, os.O_RDWR, 0666)
	if err != nil {
		setTaskStatus(task, "Error: "+err.Error())
		return
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
	go writeWorker(f, writes, &dld, task.cancel)

	for w := 0; w < conns; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for wi := range work {
				select {
				case <-task.cancel:
					return
				default:
				}
				ci := chunkInfo{start: wi.start, end: wi.end}
				if err := downloadChunk(task.URL, ci, writes, task.cancel); err != nil {
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
			task.Downloaded = d
			elapsed := time.Since(task.startTime)
			task.Speed = formatSpeed(d, elapsed)
			task.TimeElapsed = formatDuration(elapsed)
			task.Progress = int(d * 100 / fileSize)
			remaining := fileSize - d
			if d > 0 {
				eta := time.Duration(float64(remaining)/float64(d)) * elapsed
				task.ETA = formatDuration(eta)
			} else {
				task.ETA = "-"
			}
			updateModel()
			updateStatusBar()
		case <-task.cancel:
			ticker.Stop()
			f.Close()
			os.Remove(partFile)
			setTaskStatus(task, "Cancelled")
			return
		}
	}
	ticker.Stop()

	select {
	case err := <-errCh:
		setTaskStatus(task, "Error: "+err.Error())
		f.Close()
		os.Remove(partFile)
		return
	default:
	}

	f.Close()
	os.Rename(partFile, task.Filename)

	task.Downloaded = fileSize
	task.Progress = 100
	task.Speed = ""
	task.ETA = "-"
	atomic.AddInt64(&totalDownloaded, fileSize)
	setTaskStatus(task, "Completed")
	elapsed := time.Since(task.startTime)
	task.TimeElapsed = formatDuration(elapsed)
	updateModel()
	enableUI()
}

func singleDownloadGUI(task *DownloadTask) {
	resp, err := sharedClient.Get(task.URL)
	if err != nil {
		setTaskStatus(task, "Error: "+err.Error())
		return
	}
	defer resp.Body.Close()

	f, err := os.Create(task.Filename)
	if err != nil {
		setTaskStatus(task, "Error: "+err.Error())
		return
	}
	defer f.Close()

	buf := make([]byte, 1024*1024)
	ticker := time.NewTicker(200 * time.Millisecond)
	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-task.cancel:
				close(done)
				return
			default:
				n, rerr := resp.Body.Read(buf)
				if n > 0 {
					f.Write(buf[:n])
					task.Downloaded += int64(n)
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
		}
	}()

	for {
		select {
		case <-done:
			ticker.Stop()
			if task.Status == "Cancelled" {
				f.Close()
				os.Remove(task.Filename)
				return
			}
			task.Progress = 100
			task.Speed = ""
			task.ETA = "-"
			atomic.AddInt64(&totalDownloaded, task.FileSize)
			setTaskStatus(task, "Completed")
			task.TimeElapsed = formatDuration(time.Since(task.startTime))
			updateModel()
			enableUI()
			return
		case <-task.cancel:
			ticker.Stop()
			f.Close()
			os.Remove(task.Filename)
			setTaskStatus(task, "Cancelled")
			enableUI()
			return
		case <-ticker.C:
			d := task.Downloaded
			elapsed := time.Since(task.startTime)
			task.Speed = formatSpeed(d, elapsed)
			task.TimeElapsed = formatDuration(elapsed)
			if task.FileSize > 0 {
				task.Progress = int(d * 100 / task.FileSize)
				remaining := task.FileSize - d
				if d > 0 {
					eta := time.Duration(float64(remaining)/float64(d)) * elapsed
					task.ETA = formatDuration(eta)
				}
			}
			updateModel()
			updateStatusBar()
		}
	}
}

func setTaskStatus(task *DownloadTask, status string) {
	task.Status = status
	updateModel()
	if status == "Completed" || strings.HasPrefix(status, "Error") || status == "Cancelled" {
		enableUI()
	}
}

func updateModel() {
	model.PublishRowsReset()
}

func cancelSelected() {
	idx := tasksView.CurrentIndex()
	if idx >= 0 && idx < len(model.items) {
		task := model.items[idx]
		if task.Status == "Downloading" || task.Status == "Queued" {
			select {
			case <-task.cancel:
			default:
				close(task.cancel)
			}
			setTaskStatus(task, "Cancelled")
			if task.partFile != "" {
				os.Remove(task.partFile)
			}
		}
	}
}

func clearCompleted() {
	var remaining []*DownloadTask
	for _, t := range model.items {
		if t.Status != "Completed" && t.Status != "Error" && t.Status != "Cancelled" {
			remaining = append(remaining, t)
		}
	}
	model.items = remaining
	model.PublishRowsReset()
}

func enableUI() {
	downloadBtn.SetEnabled(true)
	urlEntry.SetEnabled(true)
}

func updateStatusBar() {
	active := atomic.LoadInt32(&activeCount)
	total := atomic.LoadInt64(&totalDownloaded)
	statusActive.SetText(fmt.Sprintf("%d", active))
	statusTotal.SetText(formatSize(total))

	// Calculate total speed from active tasks
	var totalSpeed int64
	for _, t := range model.items {
		if t.Status == "Downloading" && t.Downloaded > 0 {
			speed := float64(t.Downloaded) / time.Since(t.startTime).Seconds()
			totalSpeed += int64(speed)
		}
	}
	if totalSpeed > 0 {
		statusSpeed.SetText(formatSize(totalSpeed) + "/s")
	} else {
		statusSpeed.SetText("0 B/s")
	}
}
