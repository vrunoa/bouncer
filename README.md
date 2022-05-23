# bouncer

[WIP] Bouncer says no to big docker images

## Usage 

```
bouncer check --config-file ./config.yaml
```

**config**

```
---
apiVersion: v1alpha
kind: bouncer
image:
  name: bouncer:latest
  policy:
    deny:
      - desc: Image too big
        size: 20Mi
```

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
* Add test failure pipeline
* More unit tests, increase code coverage to 70%
* Fix PullImage stdout
* Refactor a lot
* Add bash installer
