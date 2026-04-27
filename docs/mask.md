# mask

Mask sensitive values in a `.env` file, replacing them with a configurable placeholder.

## Usage

```
envoy mask [flags]
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--key` | — | Exact key name to mask (repeatable) |
| `--pattern` | — | Substring pattern matched against key names (repeatable, case-insensitive) |
| `--placeholder` | `****` | Replacement string for masked values |
| `--show-last` | `0` | Reveal the last N characters of the original value |
| `--output`, `-o` | source file | Write masked output to this path instead of overwriting the source |
| `--env` | `.env` | Path to the source `.env` file |

## Examples

### Mask a single key

```bash
envoy mask --key API_KEY --output .env.masked
```

### Mask all keys containing "SECRET" or "PASSWORD"

```bash
envoy mask --pattern SECRET --pattern PASSWORD -o .env.safe
```

### Reveal last 4 characters

Useful for confirming the correct credential is in place without exposing the full value:

```bash
envoy mask --key DATABASE_URL --show-last 4
```

Input:
```
DATABASE_URL=postgres://user:hunter2@host/db
```

Output:
```
DATABASE_URL=****/db
```

### Custom placeholder

```bash
envoy mask --key TOKEN --placeholder "[HIDDEN]"
```

## Notes

- Key matching for `--key` is case-insensitive.
- Pattern matching for `--pattern` checks whether the pattern appears anywhere in the key name (case-insensitive).
- When `--output` is omitted the source file is overwritten in place.
- The original file is never modified when `--output` points to a different path.
