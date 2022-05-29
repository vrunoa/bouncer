## bouncer

Bouncer says no to big docker images

## Install
```bash
sudo sh -c 'curl -L https://vrunoa.github.io/bouncer/install | bash -s -- -b /usr/local/bin'
```

## Usage
```bash
bouncer check --config-file ./config.yaml
```

**config**
```yaml
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
