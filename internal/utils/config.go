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

type Config struct {
	Port    int
	DataDir string
}

func Load() (*Config, error) {
	var (
		port    int
		dataDir string
	)

	flag.IntVar(&port, "port", 8080, "Port number")
	flag.StringVar(&dataDir, "dir", "data", "Path to the data directory")
	flag.Usage = printUsage()
	flag.Parse()
	// Railway предоставляет PORT через переменную окружения
	port := os.Getenv("PORT")
	if port == "" {
    	port = "8080"
	}

	dir := "data"

	if port < 1 || port > 65535 {
		return nil, fmt.Errorf("invalid port number %d, must be between 1 and 65535", port)
	}

	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create data directory '%s': %w", dataDir, err)
	}

	return &Config{
		Port:    port,
		DataDir: dataDir,
	}, nil
}