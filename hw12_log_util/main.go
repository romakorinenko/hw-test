package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	file, level, output, err := getParams()
	if err != nil {
		log.Println(err)
		return
	}

	result, err := analyzeLogFile(file, level)
	if err != nil {
		log.Println(err)
		return
	}

	if output == "" {
		log.Println(result)
	} else {
		if writeErr := writeResult(result, output); writeErr != nil {
			log.Println(writeErr)
		}
	}
}

func getParams() (string, string, string, error) {
	var file string
	var level string
	var output string
	flag.StringVar(&file, "file", "", "path to log file")
	flag.StringVar(&level, "level", "", "received log level")
	flag.StringVar(&output, "output", "", "path to output file")
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		return "", "", "", fmt.Errorf("there is no .env file found")
	}

	if file == "" {
		logAnalyzerFile, ok := os.LookupEnv("LOG_ANALYZER_FILE")
		if ok {
			file = logAnalyzerFile
		} else {
			return "", "", "", fmt.Errorf("log file does not sign in .env or running command flag")
		}
	}
	if level == "" {
		logAnalyzerLevel, ok := os.LookupEnv("LOG_ANALYZER_LEVEL")
		if ok {
			level = strings.ToUpper(logAnalyzerLevel)
		} else {
			return "", "", "", fmt.Errorf("log level does not sign in .env or running command flag")
		}
	}
	if output == "" {
		logAnalyzerOutput, ok := os.LookupEnv("LOG_ANALYZER_OUTPUT")
		if ok {
			output = logAnalyzerOutput
		}
	}
	return file, level, output, nil
}

func analyzeLogFile(fileName, level string) (string, error) {
	var logFile *os.File
	defer func() {
		_ = logFile.Close()
	}()

	logFile, err := os.Open(fileName)
	if err != nil {
		return "", fmt.Errorf("cannot open file: %s, error: %w", fileName, err)
	}

	scanner := bufio.NewScanner(logFile)

	var logLevelLineCounter int
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), level) {
			logLevelLineCounter++
		}
	}

	return fmt.Sprintf("файл %s содержит %d строк с уровнем %s", fileName, logLevelLineCounter, level), nil
}

func writeResult(result, output string) error {
	var fileOutput *os.File
	defer func() {
		_ = fileOutput.Close()
	}()

	fileOutput, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("unable to create output file: %w", err)
	}

	_, err = fileOutput.WriteString(result)
	if err != nil {
		return fmt.Errorf("cannot write result in file: %w", err)
	}

	return nil
}
