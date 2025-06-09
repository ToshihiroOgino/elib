package env

import (
	"bufio"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Env struct {
	Port   int
	DBFile string
}

var (
	once sync.Once
	env  Env
)

func trimSpace(s string) string {
	s = strings.TrimLeft(s, " \t")
	s = strings.TrimRight(s, " \t")
	return s
}

func load() {
	envPath := ".env"
	file, err := os.Open(envPath)
	if err != nil {
		slog.Error("failed to open .env file", "error", err)
		return
	}
	defer file.Close()

	envMap := make(map[string]string)
	reader := bufio.NewScanner(file)
	for reader.Scan() {
		line := reader.Text()
		line = trimSpace(line)
		if line == "" || line[0] == '#' {
			// Skip empty lines and comments
			continue
		}
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			slog.Error("invalid line in .env file", "line", line)
			continue
		}
		key := trimSpace(parts[0])
		value := trimSpace(parts[1])
		envMap[key] = value
	}

	port, err := strconv.Atoi(envMap["PORT"])
	if err != nil {
		slog.Error("failed to parse PORT from .env file", "error", err)
	}

	env = Env{
		Port:   port,
		DBFile: envMap["DB_FILE"],
	}
	slog.Debug("loaded environment variables", "env", env)
}

func loadOnce() {
	once.Do(func() {
		load()
	})
}

func Get() Env {
	loadOnce()
	return env
}
