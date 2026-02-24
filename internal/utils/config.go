package utils

import (
	"flag"
	"fmt"
	"os"
)
const helpText = `Coffee Shop Management System

Usage:
  hot-coffee [--port <N>] [--dir <S>] 
  hot-coffee --help

Options:
  --help       Show this screen.
  --port N     Port number.
  --dir S      Path to the data directory.`

type Config struct {
	Port    string
	DataDir string
}

func Load() (*Config, error) {
	var (
		port    string
		dataDir string
	)

	flag.StringVar(&port, "port", "8080", "Port number")
	flag.StringVar(&dataDir, "dir", "data", "Path to the data directory")
	flag.Usage = func (){fmt.Println(helpText)}
	flag.Parse()
	// Railway предоставляет PORT через переменную окружения
	port = os.Getenv("PORT")
	if port == "" {
    	port = "8080"
	}

	dataDir = "data"
	portInt,_ := strconv.Atoi(port)
	if portInt < 1 || portInt > 65535 {
		return nil, fmt.Errorf("invalid port number %d, must be between 1 and 65535", portInt)
	}
	portStr := strconv.Itoa(portInt) 
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create data directory '%s': %w", dataDir, err)
	}

	return &Config{
		Port:    port,
		DataDir: dataDir,
	}, nil
}