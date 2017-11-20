GKeepassXReader
===============

A simple command line interface for [KeePassX][1] database files. 
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
  -c, --chrs=CHRS        Copy selected characters from password [2,6,7..]
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

[Dep][3] is required to install the dependencies


```bash
dep ensure
make build
```

Search
======

Search by title or UUID
-----------------------

```bash
./gkeepassxreader search 'Sample Entry' --db Database.kdbx
Password (press enter for no password):
+----------------------------------+-----------+--------------+----------+--------------------------+-------+-------------------+
|               UUID               |   GROUP   |    TITLE     | USERNAME |           URL            | NOTES |     PASSWORD      |
+----------------------------------+-----------+--------------+----------+--------------------------+-------+-------------------+
| a8370aa88afd3c4593ce981eafb789c8 | Protected | Sample Entry |          | http://www.somesite.com/ | Notes | ProtectedPassword |
+----------------------------------+-----------+--------------+----------+--------------------------+-------+-------------------+

```

Copy password to clipboard
--------------------------

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

Select specific characters from password
----------------------------------------

```bash
./gkeepassxreader --db Database.kdbx search 'Sample Entry' --chrs 1,7,8
Password (press enter for no password): 
+----------------------------------+-----------+--------------+---------------------+--------------------------+-------+----------+
|               UUID               |   GROUP   |    TITLE     |      USERNAME       |           URL            | NOTES | PASSWORD |
+----------------------------------+-----------+--------------+---------------------+--------------------------+-------+----------+
| a8370aa88afd3c4593ce981eafb789c8 | Protected | Sample Entry | Protected User Name | http://www.somesite.com/ | Notes | Pte      |
+----------------------------------+-----------+--------------+---------------------+--------------------------+-------+----------+

```

List
====

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



Testing
-------

[Ginkgo][4] is used to run the tests


```bash
make test
```


[1]: https://www.keepassx.org/
[2]: https://golang.org/
[3]: https://github.com/golang/dep
[4]: http://onsi.github.io/ginkgo/
