# Binlogic CloudBackup Prepare
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](https://github.com/cweill/gotests/blob/master/LICENSE) [![Build Status](https://travis-ci.org/binlogicinc/cloudbackup-prepare.svg?branch=master)](https://travis-ci.org/binlogicinc/cloudbackup-prepare) [![Coverage Status](https://coveralls.io/repos/github/binlogicinc/cloudbackup-prepare/badge.svg?branch=master)](https://coveralls.io/github/binlogicinc/cloudbackup-prepare?branch=master) [![codebeat badge](https://codebeat.co/badges/45b30f6f-a2f9-4eee-bd54-47eb12b586cc)](https://codebeat.co/projects/github-com-binlogicinc-cloudbackup-prepare-master) [![Go Report Card](https://goreportcard.com/badge/github.com/binlogicinc/cloudbackup-prepare)](https://goreportcard.com/report/github.com/binlogicinc/cloudbackup-prepare)

---

[CloudBackup](https://www.binlogic.io/) is a tool to orchestrate MySQL, MariaDB, MongoDB and PostgreSQL Backups in the Cloud. Each backup is encrypted with a unique key and compressed.  In case you need to use any backup in a environment
where cloudbackup is not present or for some reason you no longer want to use cloudbackup this project will help you
getting the backup ready yo use.

*cloudbackup-prepare* is a simple binary to decrypt and uncompress backups made with [CloudBackup](https://www.binlogic.io/)

## Installation

- If you already have go installed

```shell
  go get github.com/binlogicinc/cloudbackup-prepare
```

- Another alternative is to download the binary from the [Release Section ](https://github.com/binlogicinc/cloudbackup-prepare/releases).

## Usage

```console
Usage of cloudbackup-prepare:

      --agent-version string    Agent version used to take the backup (if empty it's assumed <= 1.2.0)
  -i, --backup-file string      Binlogic CloudBackup file path
  -e, --encryption-key string   Encruption key to decrypt backup file
  -o, --output-file string      Outfile to save backup decrypted and uncompressed
  -v, --version                 Output version information and exit
```

### Example

```shell
cloudbackup-prepare  -i tesfiles/mysqldump_sql.z  -e kpySdc2vfHL_4WebUstA29fRFacKis8LZRbLqFFY0HM= -o somefile.sql
```

`NOTE:` Since we cant validate  encryption key is valid to decrypt each file (otherwise you backup is not safe :( ), we only check key has a valid format. After process is complete, please check that the output file doesn't contain garbage. In that case
you are probably using an incorrect key.

where `tesfiles/mysqldump_sql.z` is the backup file, `kpySdc2vfHL_4WebUstA29fRFacKis8LZRbLqFFY0HM=` is the key used to encrypt this backup and `somefile.sql` is where you want to dump decreypted and uncompress backup. Problably is you are
processing a MongoDB backup you want to use `somefile.json` or if you are processing a xtrabackup file like `somefile.xbstream` but that is up to you.

## License

cloudbackup-prepare is distributed under Apache 2.0 License.
