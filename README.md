GKeepassXReader
===============

A simple command line interface for [KeePassX][1]database files. 
GKeepassXReader currently supports the KeePass 2 (.kdbx) password database format 
and not the older KeePass 1 (.kdb) databases.


Usage
-----

```bash
usage: gkeepassxreader --db=DB --search=SEARCH [<flags>]

Flags:
      --help             Show context-sensitive help (also try --help-long and --help-man).
      --db=DB            Keepassx database
  -s, --search=SEARCH    Search for title
  -c, --chrs=CHRS        Select characters from password [2,6,7..]
  -k, --keyfile=KEYFILE  Key file
  -d, --debug            Enable debug mode
  -x, --clipboard        Copy to clipboard
      --version          Show application version.


```

Download
--------

```bash
git clone git@github.com:simonhayward/gkeepassxreader.git

```

Build
-----

[Glide][3]is required to install the dependencies


```bash
glide install
make build
```

Run
---

```bash
./gkeepassxreader --db=database.kdbx --search=entry
```

Testing
-------

[Ginkgo][4]is used to run the tests


```bash
make test
```


[1]: https://www.keepassx.org/
[2]: https://golang.org/
[3]: https://github.com/Masterminds/glide
[4]: http://onsi.github.io/ginkgo/
