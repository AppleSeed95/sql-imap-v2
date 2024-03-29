go-imap-sql
[![Travis CI](https://img.shields.io/travis/com/foxcpp/go-imap-sql.svg?style=flat-square&logo=Linux)](https://travis-ci.com/foxcpp/go-imap-sql)
[![CodeCov](https://img.shields.io/codecov/c/github/foxcpp/go-imap-sql.svg?style=flat-square)](https://codecov.io/gh/foxcpp/go-imap-sql)
[![Reference](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/foxcpp/go-imap-sql)
[![stability-unstable](https://img.shields.io/badge/stability-unstable-yellow.svg?style=flat-square)](https://github.com/emersion/stability-badges#unstable)
=============

SQL-based storage backend for [go-imap] library.

Building
----------

Go 1.13 is required due to use of Go 1.13 error inspection features.

RDBMS support
---------------

go-imap-sql is known to work with (and constantly being tested against) following RDBMS:
- SQLite 3.25.0
- PostgreSQL 9.6

Following RDBMS have experimental support:
- CockroachDB 20.1.5

Following RDBMS were actively supported in the past, it's unknown whether they
still work with go-imap-sql:
- MariaDB 10.2 

IMAP Extensions Supported
---------------------------

- [CHILDREN]
- [APPEND-LIMIT]
- [MOVE]
- [SPECIAL-USE]
- [SORT]

Authentication
----------------

go-imap-sql does not implement any authentication. "password" argument of Login
method is not checked and the user account is created if it does not exist. You
are supposed to wrap it to implement your own authentication the way you need
it.

Usernames case-insensitivity
------------------------------

Usernames are always converted to lower-case before doing anything.
This means that if you type `imapsql-ctl ... users create FOXCPP`.  Account
with username `foxcpp` will be created. Also this means that you can use any
case in account settings in your IMAP client.

secure_delete
-------------

You may want to overwrite deleted messages and theirs meta-data with zeroes for
security/privacy reasons.
For MySQL, PostgreSQL - consult documentation (AFAIK, there is no such option).

For SQLite3, you should build go-imap-sql with `sqlite_secure_delete` build tag.
It will enable corresponding SQLite3 feature by default for all databases.

If you want to enable it per-database - you can use
`file:PATH?_secure_delete=ON` in DSN.

UIDVALIDITY
-------------

go-imap-sql never invalidates UIDs in an existing mailbox. If mailbox is
DELETE'd then UIDVALIDITY value changes.

Unlike many popular IMAP server implementations, go-imap-sql uses randomly
generated UIDVALIDITY values instead of timestamps.

This makes several things easier to implement with less edge cases. And answer
to the question you are already probably asked: To make go-imap-sql malfunction
you need to get Go's PRNG to generate two equal integers in range of [1,
2^32-1] just at right moment (seems unlikely enough to ignore it). Even then,
it will not cause much problems due to the way most client implementations
work.

go-imap-sql uses separate `math/rand.Rand` instance and seeds it with system
time on initialization (in `New`).

You can provide custom pre-seeded struct implementing `math/rand.Source` 
in `Opts` struct (`PRNG` field).

Maddy
-------

You can use go-imap-sql as part of the [maddy] mail server.

imapsql-ctl
-------------

For direct access to database you can use imapsql-ctl console utility. See more information in
separate README [here](cmd/imapsql-ctl).
```
go install github.com/foxcpp/go-imap-sql/cmd/imapsql-ctl
```

[CHILDREN]: https://tools.ietf.org/html/rfc3348
[APPEND-LIMIT]: https://tools.ietf.org/html/rfc7889
[UIDPLUS]: https://tools.ietf.org/html/rfc4315
[MOVE]: https://tools.ietf.org/html/rfc6851
[SPECIAL-USE]: https://tools.ietf.org/html/rfc6154
[SORT]: https://tools.ietf.org/html/rfc5256
[go-imap]: https://github.com/emersion/go-imap
[maddy]: https://github.com/emersion/maddy
