package main

import (
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	flag "github.com/spf13/pflag"
	"io"

	"os"
)

var (
	version       string
	showVersion   bool
	inputFile     string
	outputFile    string
	encryptionKey string
)

func init() {
	flag.BoolVarP(&showVersion, "version", "v", false, "Output version information and exit")
	flag.StringVarP(&inputFile, "backup-file", "i", "", "Binlogic CloudBackup file path")
	flag.StringVarP(&outputFile, "output-file", "o", "", "Outfile to save backup decrypted and uncompressed")
	flag.StringVarP(&encryptionKey, "encryption-key", "e", "", "Encruption key to decrypt backup file")
}

func showError(message error) {
	fmt.Fprintf(os.Stderr, "%s\n", message)
	os.Exit(0)
}
func parseArgs(args []string) error {
	flag.CommandLine.Parse(args[1:])

	if showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	if len(inputFile) == 0 {
		return fmt.Errorf("Backup file is mandatory")
	}

	if len(outputFile) == 0 {
		return fmt.Errorf("Output file is mandatory")
	}

	if len(encryptionKey) == 0 {
		return fmt.Errorf("Encryption key is mandatory")
	}

	return nil
}

func main() {
	err := parseArgs(os.Args)
	if err != nil {
		showError(err)
	}

	err = prepareBackupFile(inputFile, encryptionKey, outputFile)
	if err != nil {
		showError(err)
	}

	fmt.Println("Process completed successfully")
	os.Exit(0)
}

func validateInputFile(filename string) error {
	finfo, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("Backup file doesn't exists")
		}
	}

	if finfo.IsDir() {
		return fmt.Errorf("Expecting a file but %s is a directory", filename)
	}

	return nil
}

func validateOutputFile(filename string) error {
	_, err := os.Stat(filename)
	if err == nil {
		return fmt.Errorf("Output file exists %s. Please remove it", filename)
	}

	return nil
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

func prepareBackupFile(filename, encryptionKey, output string) error {
	err := validateInputFile(filename)
	if err != nil {
		return fmt.Errorf("%s", err)
	}

	err = validateOutputFile(output)
	if err != nil {
		return fmt.Errorf("%s", err)
	}

	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("Fail to open %s: %s", filename, err)
	}
	defer f.Close()

	zipReader, err := zlib.NewReader(f)
	if err != nil {
		return fmt.Errorf("Fail opening zlib reader: %s", err)
	}

	cipherReader, err := getCipherReader(encryptionKey, zipReader)
	if err != nil {
		return fmt.Errorf("Fail to decrypt file %s using %s. Be sure encryption key is correct", filename, encryptionKey)
	}

	of, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("Fail to create %s: %s", output, err)
	}
	defer of.Close()

	_, err = io.Copy(of, cipherReader)
	if err != nil {
		return fmt.Errorf("Fail while preparing backup file: %s", err)
	}

	return nil
}
