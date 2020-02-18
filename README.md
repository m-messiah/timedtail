# timedtail

[![GitHub release](https://img.shields.io/github/release/m-messiah/timedtail.svg?style=flat-square)](https://github.com/m-messiah/timedtail/releases)
[![Travis](https://img.shields.io/travis/m-messiah/timedtail.svg?style=flat-square)](https://travis-ci.org/m-messiah/timedtail)
[![Maintainability](https://api.codeclimate.com/v1/badges/7c40d359fbd8e2915b4a/maintainability)](https://codeclimate.com/repos/5c600a617c358f028700122a/maintainability)

Tail logs by timestamps

## Usage

```
Usage: timedtail [options] <log files>...
  -b int
        From which datetime start to operate seconds of the <log file>
  -j int
        Max number of junk lines to read (default 500)
  -n int
        Delta time count from now() to look at from the end of the <log file> (default 300)
  -r string
        Regexp to pick timestamp from string ($1 must select timestamp)
  -t string
        Timestamp type (default "common")
  -utc
        Parse timestamps as they are in UTC timezone, not local.
```

## Examples

1. Show last minute from several nginx logs ```timedtail -n 60 /var/log/nginx/app1/access.log /var/log/nginx/app2/access.log```
2. Show five minutes before unixtime from postgresql.log (log is multilines, so we use junk lines setting for skip non-timed lines) ```timedtail -t postgres -n 300 -b 1549800882 -j 10000 /var/log/postgresql/postgresql-9.6-data.log```
3. Show last five seconds from custom.log ```timedtail -r '(\d\d\d\d-\d\d-\d\dT\d\d:\d\d:\d\d)' -n 5 /var/log/my_supper_app.log```

## Install

Just download required binary from GitHub Releases somewhere to $PATH

## Compile

If you want to compile by yourself use:
```
git clone https://github.com/m-messiah/timedtail.git && go get && go build
```
