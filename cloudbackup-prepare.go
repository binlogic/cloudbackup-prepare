package main

import (
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"io"
	"os"
)

var (
	log           = logrus.New()
	logLevel      bool
	version       string
	showVersion   bool
	agentVersion  string
	inputFile     string
	encryptionKey string
)

func setupLogger(logger *logrus.Logger, verbose bool) {
	formatter := new(logrus.TextFormatter)
	formatter.FullTimestamp = true
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.ForceColors = true
	logger.Formatter = formatter
	logger.Out = os.Stderr

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
	flag.StringVarP(&encryptionKey, "encryption-key", "e", "", "Encruption key to decrypt backup file")
	flag.BoolVarP(&logLevel, "debug", "d", false, "Debug mode")
	flag.Parse()

	setupLogger(log, logLevel)

	if showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	log.Debugf("Agent Version: %s", agentVersion)
	log.Debugf("Backup File: %s", inputFile)

	validateInputFile(inputFile)

	prepareBackupFile(inputFile, encryptionKey)

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

func getCipherStream(key string) (cipher.Stream, error) {
	keyBytes, err := base64.URLEncoding.DecodeString(key)

	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(keyBytes)

	if err != nil {
		return nil, err
	}

	// If the key is unique for each ciphertext, then it's ok to use a zero IV.
	// TODO: an iv that is predictable is safer in this case ?
	var iv [aes.BlockSize]byte

	return cipher.NewOFB(block, iv[:]), nil
}

func getCipherReader(key string, wrappedReader io.Reader) (io.Reader, error) {
	if key == "" {
		return wrappedReader, nil
	}

	stream, err := getCipherStream(key)

	if err != nil {
		return nil, err
	}

	return &cipher.StreamReader{S: stream, R: wrappedReader}, nil
}

func prepareBackupFile(filename, encryptionKey string) {

	if encryptionKey == "" {
		log.Fatal("Encryption key is mandatory.")
	}

	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Fail to opening %s: %s", filename, err)
	}
	defer f.Close()

	zipReader, zerr := zlib.NewReader(f)
	if zerr != nil {
		log.Fatalf("%s opening zlib reader", zerr)
	}

	cipherReader, cerr := getCipherReader(encryptionKey, zipReader)
	if cerr != nil {
		log.Fatalf("Fail to decrypt file %s using %s: %s", filename, encryptionKey, zerr)
	}

	_, ferr := io.Copy(os.Stdout, cipherReader)
	if cerr != nil {
		log.Fatalf("Fail while preparing backup file: %s", ferr)
	}
}
