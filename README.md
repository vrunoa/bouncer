# bouncer

## Requirements

* docker

## Development

### Setup
```bash
go mod tidy
```

### Usage
```bash
go run cmd/bouncer/bouncer.go check --config-file ./test/config.yaml
```


## TODO
* Support multiple image checks
* Improve message on failure, use human-readable sizes
* More unit tests
