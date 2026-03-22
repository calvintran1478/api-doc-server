package utils

import (
	"os"
	"fmt"
	"bufio"
	"strings"
)

func LoadEnv(filename string) {
	// Open env file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %s", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		equalsIndex := strings.Index(line, "=")
		if equalsIndex != -1 {
			os.Setenv(line[:equalsIndex], line[equalsIndex + 1:])
		}
	}
}
