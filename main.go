package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"time"
)

func writeToFile(outputPath string, buffer string) error {
	outputFile, err := os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening output file: %w", err)
	}
	defer outputFile.Close()

	_, err = outputFile.WriteString(buffer + "\n")
	if err != nil {
		return fmt.Errorf("error writing to output file: %w", err)
	}

	err = outputFile.Sync()
	if err != nil {
		return fmt.Errorf("error flushing output file: %w", err)
	}

	return nil
}

func main() {
	inputPath := flag.String("input", "hub-gw-info.log", "Path to input log file")
	outputPath := flag.String("output", "hub-monitoring.log", "Path to output log file")
	flag.Parse()

	inputFile, err := os.Open(*inputPath)
	if err != nil {
		fmt.Println("Error opening input log file:", err)
		return
	}
	defer inputFile.Close()

	timestampRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\+\d{2}:\d{2}`)

	scanner := bufio.NewScanner(inputFile)

	var buffer string

	for scanner.Scan() {
		line := scanner.Text()

		if timestampRegex.MatchString(line) {
			if buffer != "" {
				time.Sleep(1 * time.Second)

				if err := writeToFile(*outputPath, buffer); err != nil {
					fmt.Println(err)
					return
				}

				fmt.Println(buffer)
			}

			buffer = line
		} else {
			buffer += "\n" + line
		}
	}

	if buffer != "" {
		time.Sleep(1 * time.Second)

		if err := writeToFile(*outputPath, buffer); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(buffer)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input log file:", err)
	}
}
