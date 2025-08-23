package envv

import (
	"bufio"
	"log/slog"
	"os"
	"strings"
)

// LoadDotEnv loads the .env file in the current directory.
func LoadDotEnv() {
	LoadFile(".env")
}

// LoadFile loads the given file and sets the env vars.
func LoadFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		slog.Error("failed to open file", "error", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// handle comments
		if strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			slog.Warn("invalid env value", "value", parts[0])
			continue
		}

		os.Setenv(parts[0], parts[1])
	}
}
