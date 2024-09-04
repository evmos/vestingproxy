# Vesting API

## Requirements:

- Go

## Usage

```sh
git clone http://github.com/evmos/vestingproxy
go run main.go
```

By default, it runs on port 8080, it can be changed by passing the port at the end of the command

```sh
go run main.go 9090
```

## Docker

```sh
docker build -t vestingproxy --progress plain .
docker run -p=8080:8080 vestingproxy
```
