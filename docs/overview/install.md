# Installing sqlc

sqlc is distributed as a single binary with zero dependencies.

## macOS

```
brew install sqlc
```

## Ubuntu

```
sudo snap install sqlc
```

## go get

```
go get github.com/kyleconroy/sqlc/cmd/sqlc
```

## Docker

```
docker pull kjconroy/sqlc
```

Run `sqlc` using `docker run`:

```
docker run --rm -v $(pwd):/src -w /src kjconroy/sqlc generate
```

## Downloads

Get pre-built binaries for *v1.10.0*:

- [Linux](https://github.com/kyleconroy/sqlc/releases/download/v1.10.0/sqlc_1.10.0_linux_amd64.tar.gz)
- [macOS](https://github.com/kyleconroy/sqlc/releases/download/v1.10.0/sqlc_1.10.0_darwin_amd64.zip)
- [Windows (MySQL only)](https://github.com/kyleconroy/sqlc/releases/download/v1.10.0/sqlc_1.10.0_windows_amd64.zip)
