# vaultpull

> CLI tool to sync secrets from HashiCorp Vault into local `.env` files with rotation support.

---

## Installation

```bash
go install github.com/yourusername/vaultpull@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/vaultpull/releases).

---

## Usage

Authenticate with your Vault instance and run `vaultpull` pointing at a secret path:

```bash
export VAULT_ADDR="https://vault.example.com"
export VAULT_TOKEN="s.xxxxxxxxxxxxxxxx"

vaultpull --path secret/data/myapp --output .env
```

This will fetch all key-value pairs at the specified path and write them to `.env`.

### Rotation Support

To rotate secrets and update your local file automatically:

```bash
vaultpull --path secret/data/myapp --output .env --rotate --interval 24h
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--path` | Vault secret path | *(required)* |
| `--output` | Output `.env` file path | `.env` |
| `--rotate` | Enable automatic rotation | `false` |
| `--interval` | Rotation check interval | `12h` |
| `--overwrite` | Overwrite existing values | `false` |

---

## Requirements

- Go 1.21+
- HashiCorp Vault 1.10+

---

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

---

## License

[MIT](LICENSE)