GKeepassXReader
===============

A simple command line interface for [KeePassX][1]database files written in [Go][2]. 
GKeepassXReader currently supports the KeePass 2 (.kdbx) password database format 
and not the older KeePass 1 (.kdb) databases.


Usage
-----

```bash
usage: gkeepassxreader search [<flags>] <term>

Search for an entry

Flags:
      --help             Show context-sensitive help (also try --help-long and --help-man).
      --db=DB            Keepassx database
  -k, --keyfile=KEYFILE  Key file
  -d, --debug            Enable debug mode
      --version          Show application version.
  -c, --chrs=CHRS        Copy characters from password [2,6,7..]
  -x, --clipboard        Copy to clipboard

Args:
  <term>  Search by title or UUID


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
./gkeepassxreader search 'Sample Entry' --db Database.kdbx
Password (press enter for no password): 
+----------------------------------+--------------+----------+--------------------------+-------+-------------------+
|               UUID               |    TITLE     | USERNAME |           URL            | NOTES |     PASSWORD      |
+----------------------------------+--------------+----------+--------------------------+-------+-------------------+
| a8370aa88afd3c4593ce981eafb789c8 | Sample Entry |          | http://www.somesite.com/ | Notes | ProtectedPassword |
+----------------------------------+--------------+----------+--------------------------+-------+-------------------+

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
