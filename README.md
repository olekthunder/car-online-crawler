# *car-online.ru* crawler

Allows you to dump your car data to xlsx file

# Installation

Download the binary and example config [here](/releases/latest) or build it manually

After that extract the archive and create file `config.yaml`

Fill it like this
```yaml
# Token from example api on car-online.ru. Place yours here
api_token: 71442b58F1ccbaA6e406C4598d7E03
file_to_save: car-online.xlsx
date_from: 2017-01-02
date_to: 2018-09-19
timezone: Europe/Kiev

```
Then, save `config.yaml` to the directory contains executable and run it

# Build

Requires `go` toolchain. See [https://golang.org/doc/install](https://golang.org/doc/install)

```bash
go get -u github.com/olekthunder/car-online-crawler
cd $GOPATH   # Windows: cd %GOPATH%
cd github.com/olekthunder/car-online-crawler
go build
```