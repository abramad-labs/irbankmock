# IR Bank Mock

A service for testing Iranian banks internet payment gateways.

## Support

This service currently only supports Saman Bank (SEP).

## Configuration

Visit [conf.go](./internal/conf/conf.go) to see which environment variables are supported.

## Deploy with Docker

Coming

## Sqlite Notice

We decided to use `github.com/glebarez/sqlite` instead of `gorm.io/driver/sqlite`. The trade-off 
was to lose a little bit of performance to gain a `cgo`-free package. 
This also enabled us to use `distroless/static` instead of `distroless/base`.

Read more:

- [cgo is not go](https://dave.cheney.net/2016/01/18/cgo-is-not-go)
- [SQLite in Go, with and without cgo](https://datastation.multiprocess.io/blog/2022-05-12-sqlite-in-go-with-and-without-cgo.html)