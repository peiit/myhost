# myhost
Auto bind host ip to your customer domain record

[![GitHub](https://img.shields.io/github/license/peiit/myhost)](https://github.com/peiit/myhost/bfe/blob/main/LICENSE)
[![Build Status](https://travis-ci.com/peiit/myhost.svg?branch=main)](https://travis-ci.com/peiit/myhost)

# Install
`go install github.com/peiit/myhost@latest`
# USAGE
`myhost --help`
```shell script
NAME:
   myhost - 将本机的ip直接设置到域名

USAGE:
   myhost [--ak 阿里云AK] [--sk 阿里云SK] [-d example.com] [-r mylocalip]

VERSION:
   0.0.1

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --accesskeyid value, --ak value      阿里云ak
   --accesskeysecret value, --sk value  阿里云sk
   --domainname value, -d value         域名
   --rr value, -r value                 主机记录
   --help, -h                           show help (default: false)
   --version, -v                        print the version (default: false)
```
