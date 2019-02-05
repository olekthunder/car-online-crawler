# *car-online.ru* crawler

Allows you to dump your car data to xlsx file

# Installation

Download the binary and example config [here](/releases/latest) or
[build](https://github.com/olekthunder/car-online-crawler#build) it manually

After that extract the archive and create file `config.yaml`

Fill it like this
```yaml
api_token: 71442b58F1ccbaA6e406C4598d7E03  # Token from example api on car-online.ru
file_to_save: car-online.xlsx
date_from: 2017-01-02  # Date to get data from
date_to: 2018-09-19    # Date to get data until
timezone: Europe/Kiev

```
Then, save `config.yaml` to the directory that contains the executable and run it

# Build

Requires [go](https://golang.org/doc/install). Also I use [dep](https://golang.github.io/dep/docs/installation.html) as
package manager.

```bash
go get -u github.com/olekthunder/car-online-crawler
cd $GOPATH  # Windows: cd %GOPATH%
cd github.com/olekthunder/car-online-crawler
dep ensure  # To ensure all 3-rd party packages are installed
go build  # Now binary executable should appear inside your current directory
```
