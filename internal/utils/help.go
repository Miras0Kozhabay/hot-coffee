package utils

import "fmt"

const helpText = `Coffee Shop Management System

Usage:
  hot-coffee [--port <N>] [--dir <S>] 
  hot-coffee --help

Options:
  --help       Show this screen.
  --port N     Port number.
  --dir S      Path to the data directory.`

func PrintHelp() {
	fmt.Println(helpText)
}
