# envoy-cli

A lightweight CLI for managing `.env` files across multiple environments and profiles.

---

## Installation

```bash
go install github.com/yourusername/envoy-cli@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envoy-cli.git
cd envoy-cli && go build -o envoy .
```

---

## Usage

```bash
# Initialize a new environment profile
envoy init --profile staging

# Set a variable in a profile
envoy set DATABASE_URL=postgres://localhost:5432/mydb --profile staging

# Get a variable from a profile
envoy get DATABASE_URL --profile staging

# List all variables in a profile
envoy list --profile staging

# Switch active profile
envoy use staging

# Export active profile to shell
eval $(envoy export)
```

Profiles are stored as `.env.<profile>` files in your project root, making them easy to version control or share with your team.

---

## Configuration

By default, `envoy-cli` looks for env files in the current directory. You can override this with the `--dir` flag or by setting `ENVOY_DIR` in your environment.

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any major changes.

---

## License

[MIT](LICENSE)