package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"reflect"
	"testing"
)

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)

	if err != nil {
		return nil, err
	}

	return b, nil
}

func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)

	return "/tmp/" + base64.URLEncoding.EncodeToString(b), err
}

func resetArguments() {
	inputFile = ""
	outputFile = ""
	encryptionKey = ""
}

func Test_validateInputFile(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name   string
		args   args
		expect error
	}{
		{
			name: "Valid input file",
			args: args{
				filename: `tesfiles/mysqldump_sql.z`,
			},
			expect: nil,
		},
		{
			name: "Missing input file",
			args: args{
				filename: `tesfiles/foo.z`,
			},
			expect: fmt.Errorf("Backup file doesn't exists"),
		},
		{
			name: "Directory as input file",
			args: args{
				filename: `tesfiles`,
			},
			expect: fmt.Errorf("Expecting a file but tesfiles is a directory"),
		},
		{
			name: "No input file",
			args: args{
				filename: ``,
			},
			expect: fmt.Errorf("Backup file doesn't exists"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validateInputFile(tt.args.filename)
			if !reflect.DeepEqual(result, tt.expect) {
				t.Errorf("Expecting %s got %s", tt.expect, result)
			}
		})
	}
}

func Test_validateOutputFile(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name   string
		args   args
		expect error
	}{
		{
			name: "Valid output file",
			args: args{
				filename: `tesfiles/backup.sql`,
			},
			expect: nil,
		},
		{
			name: "Directory as output file",
			args: args{
				filename: `tesfiles`,
			},
			expect: fmt.Errorf("Output file exists tesfiles. Please remove it"),
		},
		{
			name: "Output file exists",
			args: args{
				filename: `tesfiles/dummy.sql`,
			},
			expect: fmt.Errorf("Output file exists tesfiles/dummy.sql. Please remove it"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validateOutputFile(tt.args.filename)
			if !reflect.DeepEqual(result, tt.expect) {
				t.Errorf("Expecting %s got %s", tt.expect, result)
			}
		})
	}
}

func Test_prepareBackupFile(t *testing.T) {
	backupFileNAme, _ := GenerateRandomString(8)
	type args struct {
		filename      string
		encryptionKey string
		output        string
	}
	tests := []struct {
		name   string
		args   args
		expect error
	}{
		{
			name: "Valid encryption key",
			args: args{
				filename:      `tesfiles/mysqldump_sql.z`,
				encryptionKey: `kpySdc2vfHL_4WebUstA29fRFacKis8LZRbLqFFY0HM=`,
				output:        backupFileNAme,
			},
			expect: nil,
		},
		{
			name: "Missing input file",
			args: args{
				filename:      `tesfiles/missing.file`,
				encryptionKey: `kpySdc2vfHL_4WebUstA29fRFacKis8LZRbLqFFY0HM=`,
				output:        `/tmp/backup.json`,
			},
			expect: fmt.Errorf("Backup file doesn't exists"),
		},
		{
			name: "Invalid encryption key",
			args: args{
				filename:      `tesfiles/mysqldump_sql.z`,
				encryptionKey: `InvalidKeyTest`,
				output:        `/tmp/mysqldump.sql`,
			},
			expect: fmt.Errorf("Fail to decrypt file tesfiles/mysqldump_sql.z using InvalidKeyTest. Be sure encryption key is correct"),
		},
		{
			name: "Directory as input file",
			args: args{
				filename:      `tesfiles`,
				encryptionKey: `InvalidKeyTest`,
				output:        `/tmp/backup_dir.sql`,
			},
			expect: fmt.Errorf("Expecting a file but tesfiles is a directory"),
		},
		{
			name: "Plain input file",
			args: args{
				filename:      `tesfiles/dummy.sql`,
				encryptionKey: `InvalidKeyTest`,
				output:        `/tmp/backup_dir.sql`,
			},
			expect: fmt.Errorf("Fail opening zlib reader: unexpected EOF"),
		},
		{
			name: "Incorrect output file",
			args: args{
				filename:      `tesfiles/mysqldump_sql.z`,
				encryptionKey: `kpySdc2vfHL_4WebUstA29fRFacKis8LZRbLqFFY0HM=`,
				output:        `/missing/dir/file.sql`,
			},
			expect: fmt.Errorf("Fail to create /missing/dir/file.sql: open /missing/dir/file.sql: no such file or directory"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := prepareBackupFile(tt.args.filename, tt.args.encryptionKey, tt.args.output)
			if !reflect.DeepEqual(result, tt.expect) {
				t.Errorf("Expecting %s got %s", tt.expect, result)
			}
		})
	}
}

func Test_parseArgs(t *testing.T) {
	type args struct {
		opts []string
	}
	tests := []struct {
		name   string
		args   args
		expect []string
	}{
		{
			name: "Validate backup-file argument",
			args: args{
				opts: []string{"cmd", "--backup-file", "infile.sql"},
			},
			expect: []string{"infile.sql", "", ""},
		},
		{
			name: "Validate backup-file and output-file arguments",
			args: args{
				opts: []string{"cmd", "--backup-file", "infile.sql", "--output-file", "outfile.sql"},
			},
			expect: []string{"infile.sql", "outfile.sql", ""},
		},
		{
			name: "Validate backup-file and output-file and encryption-key arguments",
			args: args{
				opts: []string{"cmd", "--backup-file", "file.json.z", "--output-file", "file.json", "--encryption-key", "XXXXXXXXXX"},
			},
			expect: []string{"file.json.z", "file.json", "XXXXXXXXXX"},
		},
		{
			name: "Validate output-file argument",
			args: args{
				opts: []string{"cmd", "--output-file", "out.sql"},
			},
			expect: []string{"", "out.sql", ""},
		},
		{
			name: "Validate encryptionKey argument",
			args: args{
				opts: []string{"cmd", "-e", "SOMERANDOMSTRING"},
			},
			expect: []string{"", "", "SOMERANDOMSTRING"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetArguments()
			_ = parseArgs(tt.args.opts)
			result := []string{inputFile, outputFile, encryptionKey}
			if !reflect.DeepEqual(result, tt.expect) {
				t.Errorf("Expecting %s got %s", tt.expect, result)
			}
		})
	}
}
