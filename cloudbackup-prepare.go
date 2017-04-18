package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"io"
	"os"
)

var (
	log          = logrus.New()
	logLevel     bool
	version      string
	showVersion  bool
	agentVersion string
	inputFile    string
	outputDir    string
)

func setupLogger(logger *logrus.Logger, verbose bool) {
	formatter := new(logrus.TextFormatter)
	formatter.FullTimestamp = true
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.ForceColors = true
	logger.Formatter = formatter
	logger.Out = os.Stdout

	if verbose {
		logger.Level = logrus.DebugLevel
	} else {
		logger.Level = logrus.InfoLevel
	}
}

func main() {

	flag.BoolVarP(&showVersion, "version", "v", false, "Output version information and exit")
	flag.StringVarP(&agentVersion, "agent-version", "a", "", "Binlogic CloudBackup agent version used to take backup")
	flag.StringVarP(&inputFile, "backup-file", "i", "", "Binlogic CloudBackup file path")
	flag.StringVarP(&outputDir, "output-dir", "o", "", "Output directory to place backup files.")
	flag.BoolVarP(&logLevel, "debug", "d", false, "Debug mode")
	flag.Parse()

	setupLogger(log, logLevel)

	if showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	log.Debugf("Agent Version: %s", agentVersion)
	log.Debugf("Backup File: %s", inputFile)
	log.Debugf("Output Directory: %s", outputDir)

	validateInputFile(inputFile)
	validateOutputDir(outputDir)

	os.Exit(0)

}

func validateInputFile(filename string) {

	if len(filename) == 0 {
		log.Fatal("Input file is mandatory.")
	}

	finfo, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatal("Input file doesn't exists!.")
		}
	}

	if finfo.IsDir() {
		log.Fatalf("Expecting a file but %s is a directory.", filename)
	}
	log.Debugf("File: %s exists.", filename)
}

func validateOutputDir(path string) {

	if len(path) == 0 {
		log.Fatal("Ourput dir is mandatory.")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	} else {
		isEmptyDir(path)
	}
}

func isEmptyDir(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Cannot access directory %s", path)
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return
	}
	log.Fatalf("Output directory %s should be empty!", path)
}
