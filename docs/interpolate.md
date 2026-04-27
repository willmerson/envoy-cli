# interpolate

Expand `$VAR` and `${VAR}` references inside `.env` values using other keys defined in the same file.

## Usage

```
envoy interpolate [flags]
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--strict` | `false` | Exit with an error if any reference cannot be resolved |
| `--output`, `-o` | *(source file)* | Write the result to a different file instead of overwriting the source |
| `--env` | `.env` | Path to the source `.env` file |

## Examples

### Basic expansion

```dotenv
# .env
BASE_URL=https://api.example.com
HEALTH_URL=${BASE_URL}/health
CALLBACK=$BASE_URL/callback
```

```bash
envoy interpolate
```

Result:

```dotenv
BASE_URL=https://api.example.com
HEALTH_URL=https://api.example.com/health
CALLBACK=https://api.example.com/callback
```

### Write to a separate file

```bash
envoy interpolate --env .env.template --output .env.resolved
```

### Strict mode

Fails when a reference like `${MISSING}` cannot be resolved:

```bash
envoy interpolate --strict
# Error: unresolved references: MISSING
```

## Notes

- Self-referential or circular references are **not** detected; the original token is left unchanged.
- Only uppercase `$VAR` bare-style references are matched to avoid false positives with shell arithmetic or lowercase variable names.
- Comments and blank lines are preserved in the output.
