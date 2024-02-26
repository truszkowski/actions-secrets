package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func loadEnvVars(filePath string) (map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	values := map[string]string{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		// Skip empty lines and lines starting with '#'
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split line into key and value at the first '='
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // Skip lines that don't have a key=value pair
		}

		key, value := parts[0], parts[1]
		if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
			if err := json.Unmarshal([]byte(value), &value); err != nil {
				return nil, fmt.Errorf("error unmarshalling JSON value: %w", err)
			}
		}

		values[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return values, nil
}

func loadEnvVarsNames(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	var names []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		// Skip empty lines and lines starting with '#'
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split line into key and value at the first '='
		parts := strings.SplitN(line, "=", 2)
		names = append(names, parts[0])
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return names, nil
}
