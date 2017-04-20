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
			expect: fmt.Errorf("Input file doesn't exists"),
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
			expect: fmt.Errorf("Input file is mandatory"),
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
		{
			name: "Missing output file",
			args: args{
				filename: ``,
			},
			expect: fmt.Errorf("Output file is mandatory"),
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
			expect: fmt.Errorf("Input file doesn't exists"),
		},
		{
			name: "Invalid encryption key",
			args: args{
				filename:      `tesfiles/mysqldump_sql.z`,
				encryptionKey: `InvalidKeyTest`,
				output:        `/tmp/mysqldump.sql`,
			},
			expect: fmt.Errorf("Fail to decrypt file tesfiles/mysqldump_sql.z using InvalidKeyTest: illegal base64 data at input byte 12"),
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
