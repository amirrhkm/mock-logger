package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"time"
)

func main() {
	inputFile, err := os.Open("hub-gw-info.log")
	if err != nil {
		fmt.Println("Error opening input log file:", err)
		return
	}
	defer inputFile.Close()

	outputFile, err := os.Create("hub-monitoring.log")
	if err != nil {
		fmt.Println("Error creating output log file:", err)
		return
	}
	defer outputFile.Close()

	timestampRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\+\d{2}:\d{2}`)

	scanner := bufio.NewScanner(inputFile)

	var buffer string

	for scanner.Scan() {
		line := scanner.Text()

		if timestampRegex.MatchString(line) {
			if buffer != "" {
				time.Sleep(1 * time.Second)

				_, err := outputFile.WriteString(buffer + "\n")
				if err != nil {
					fmt.Println("Error writing to output log file:", err)
					return
				}

				err = outputFile.Sync()
				if err != nil {
					fmt.Println("Error flushing output log file:", err)
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

		_, err := outputFile.WriteString(buffer + "\n")
		if err != nil {
			fmt.Println("Error writing to output log file:", err)
			return
		}

		err = outputFile.Sync()
		if err != nil {
			fmt.Println("Error flushing output log file:", err)
			return
		}

		fmt.Println(buffer)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input log file:", err)
	}
}
