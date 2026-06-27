package main

import "os"

func main() {
	if len(os.Args) >= 2 && os.Args[1] != "" && os.Args[1][0] != '-' {
		cliMain()
		return
	}
	guiMain()
}
