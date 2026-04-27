# envoy rollback

Restore a `.env` file to the state captured in a named snapshot.

## Usage

```
envoy rollback <snapshot-name> [flags]
```

## Flags

| Flag | Description |
|------|-------------|
| `--dry-run` | Preview changes without writing to disk |
| `-v`, `--verbose` | Show individual key changes after rollback |
| `--env` | Path to the `.env` file (default: `.env`) |

## Examples

### Restore from a snapshot

```bash
envoy rollback before-deploy
# Rollback complete: 3 restored, 0 added, 1 removed (snapshot: before-deploy)
```

### Preview without applying

```bash
envoy rollback before-deploy --dry-run
# Dry-run rollback plan:
#   restore: DB_HOST
#   restore: API_KEY
#   remove:  TEMP_FLAG
```

### Verbose output

```bash
envoy rollback before-deploy --verbose
# Rollback complete: 1 restored, 0 added, 0 removed (snapshot: before-deploy)
#   restored: DB_HOST
```

## Notes

- Snapshots are created with `envoy snapshot save <name>`.
- Keys present in the current file but absent from the snapshot are **removed**.
- Keys present in the snapshot but absent from the current file are **added**.
- Use `--dry-run` to audit changes before committing them.
