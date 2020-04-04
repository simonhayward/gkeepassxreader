# GKeepassXReader

## Overview

A simple command line interface for [KeePassX][0] database files, to search and list entries.
GKeepassXReader currently supports the KeePass 2 (.kdbx) database format.

## From source

GKeepassXReader requires [Go 1.14][1] or later.

```bash
go get -v github.com/simonhayward/gkeepassxreader/...
```

This produces a binary in your `${GOPATH}/bin/gkeepassxreader`.

## Docker

This creates a container from [Go 1.14][1], downloads the source, builds and installs it.

```bash
docker run golang:1.14-buster go get -v github.com/simonhayward/gkeepassxreader/...
```

Create an image of the most recent container.
```bash
docker commit $(docker ps -lq) gkeepassxreader-image
```

Run a container mounting the directories holding your database files.

```bash
docker run --rm -it -v `pwd`:/db-path gkeepassxreader-image gkeepassxreader --db=/db-path/db.kdbx list
```

## Usage

### Search

```bash
usage: gkeepassxreader search [<flags>] <term>

Search for an entry

Flags:
      --help             Show context-sensitive help (also try --help-long and --help-man).
      --db=DB            Keepassx database
  -k, --keyfile=KEYFILE  Key file
  -d, --debug            Enable debug mode
  -h, --history          Include historical entries
      --version          Show application version.
  -c, --chrs=CHRS        Copy selected characters from password [2,6,7..]
  -x, --clipboard        Copy to clipboard

Args:
  <term>  Search by title or UUID


```

#### Search by title or UUID

```bash
./gkeepassxreader search 'Sample Entry' --db Database.kdbx
Password (press enter for no password):
+----------------------------------+-----------+--------------+----------+--------------------------+-------+-------------------+
|               UUID               |   GROUP   |    TITLE     | USERNAME |           URL            | NOTES |     PASSWORD      |
+----------------------------------+-----------+--------------+----------+--------------------------+-------+-------------------+
| a8370aa88afd3c4593ce981eafb789c8 | Protected | Sample Entry |          | http://www.somesite.com/ | Notes | ProtectedPassword |
+----------------------------------+-----------+--------------+----------+--------------------------+-------+-------------------+

```

#### Search by title or UUID and copy password to clipboard

```bash
./gkeepassxreader --db Database.kdbx search 'Sample Entry' -x
Password (press enter for no password):
password copied to clipboard
+----------------------------------+-----------+--------------+---------------------+--------------------------+-------+
|               UUID               |   GROUP   |    TITLE     |      USERNAME       |           URL            | NOTES |
+----------------------------------+-----------+--------------+---------------------+--------------------------+-------+
| a8370aa88afd3c4593ce981eafb789c8 | Protected | Sample Entry | Protected User Name | http://www.somesite.com/ | Notes |
+----------------------------------+-----------+--------------+---------------------+--------------------------+-------+
```

#### Search by title or UUID and select specific characters from password

```bash
./gkeepassxreader --db Database.kdbx search 'Sample Entry' --chrs 1,7,8
Password (press enter for no password):
+----------------------------------+-----------+--------------+---------------------+--------------------------+-------+----------+
|               UUID               |   GROUP   |    TITLE     |      USERNAME       |           URL            | NOTES | PASSWORD |
+----------------------------------+-----------+--------------+---------------------+--------------------------+-------+----------+
| a8370aa88afd3c4593ce981eafb789c8 | Protected | Sample Entry | Protected User Name | http://www.somesite.com/ | Notes | Pte      |
+----------------------------------+-----------+--------------+---------------------+--------------------------+-------+----------+

```

### List

```bash
usage: gkeepassxreader list

List entries

Flags:
      --help             Show context-sensitive help (also try --help-long and --help-man).
      --db=DB            Keepassx database
  -k, --keyfile=KEYFILE  Key file
  -d, --debug            Enable debug mode
  -h, --history          Include historical entries
      --version          Show application version.

```

```bash
./gkeepassxreader list --db Example.kdbx
Password (press enter for no password):
+----------------------------------+---------+-----------------+------------+-------------------------------------------+-------+
|               UUID               |  GROUP  |      TITLE      |  USERNAME  |                    URL                    | NOTES |
+----------------------------------+---------+-----------------+------------+-------------------------------------------+-------+
| 640c38611c3ea4489ced361f54e43dbe | example | Sample Entry    | User Name  | http://keepass.info/                      | Notes |
| db8e52f8c86d7d468ecd53d4c2fe0a31 | example | Sample Entry #2 | Michael321 | http://keepass.info/help/kb/testform.html |       |
+----------------------------------+---------+-----------------+------------+-------------------------------------------+-------+

```

## Testing

[Ginkgo][2] is used to run the tests

```bash
cd $GOPATH/src/github.com/simonhayward/gkeepassxreader
make test
```

[0]: https://www.keepassx.org/
[1]: https://golang.org/
[2]: http://onsi.github.io/ginkgo/
