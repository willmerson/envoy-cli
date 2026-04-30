# `envoy protect`

Mark one or more keys in your `.env` file as **protected** by inserting a
`#PROTECTED: <KEY>` sentinel comment directly above each matched entry.

This is a lightweight, human-readable convention that tooling (and teammates)
can respect when reviewing changes.

## Usage

```bash
envoy protect [flags]
```

## Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--key` | `-k` | Key name to protect. Repeatable; also accepts comma-separated values. |
| `--prefix` | `-p` | Protect all keys sharing this prefix. Repeatable. |
| `--dry-run` | | Preview which keys would be protected without writing changes. |
| `--env` | `-e` | Path to the `.env` file (default: `.env`). |

At least one `--key` or `--prefix` must be supplied.

## Examples

### Protect a single key

```bash
envoy protect --key DB_PASS
```

**Before:**
```
DB_PASS=supersecret
```

**After:**
```
#PROTECTED: DB_PASS
DB_PASS=supersecret
```

### Protect all keys sharing a prefix

```bash
envoy protect --prefix AWS_
```

### Protect multiple explicit keys

```bash
envoy protect --key SECRET_KEY --key API_TOKEN
# or comma-separated
envoy protect --key SECRET_KEY,API_TOKEN
```

### Dry-run preview

```bash
envoy protect --prefix DB_ --dry-run
```

Outputs the list of keys that *would* be protected without modifying the file.

## Notes

- Running `protect` on an already-protected key is a no-op; the sentinel will
  not be duplicated.
- The protection marker is a comment convention only — it does not enforce
  read-only access at the OS level. Use it in combination with code review
  policies or pre-commit hooks.
- To remove a protection marker, delete the `#PROTECTED: <KEY>` comment line
  manually or with `envoy prune --commented`.
