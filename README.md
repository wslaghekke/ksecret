# ksecret [![Build Status](https://travis-ci.org/wslaghekke/ksecret.svg)](https://travis-ci.org/wslaghekke/ksecret)
Kubernetes secret editor

## Commands

### List secrets
Lists secrets in current namespace
```bash
$ ksecret
somesecret
configfilessecret
```

### Edit secret
Opens specified secret in your configured editor (based on env var EDITOR)
```bash
$ ksecret edit somesecret
```

### List secret keys
List the keys in specified secret
```bash
$ ksecret editfile configfilessecret
config.json
secret.key
```

### Edit secret file
Opens a single key from specified secret in configured editor with key as name to trigger vim syntax highlighting (based on name)
```bash
$ ksecret editfile config.json
```
