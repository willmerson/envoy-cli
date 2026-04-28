# archive

Save and restore named versions of your `.env` file.

## Usage

```
envoy archive <subcommand> [flags]
```

## Subcommands

### save

Captures the current state of the env file and stores it as a labelled archive.

```
envoy archive save --label <label> [--env <path>]
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--label` | `-l` | Label for the archive (required) |
| `--env` | `-e` | Path to the `.env` file (default: `.env`) |

**Example:**

```bash
envoy archive save --label "before-deploy"
# Archive saved: 1718123456789012345 (label: "before-deploy")
```

---

### list

Lists all saved archives, sorted newest first.

```
envoy archive list
```

**Example output:**

```
ID                       LABEL            CREATED              KEYS
1718123456789012345      before-deploy    2024-06-11 14:30:56  12
1718100000000000000      initial          2024-06-11 08:00:00  10
```

---

### restore

Restores the env file from a previously saved archive by ID.

```
envoy archive restore --id <id> [--env <path>]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--id` | Archive ID to restore (required) |
| `--env` | Path to the `.env` file (default: `.env`) |

**Example:**

```bash
envoy archive restore --id 1718123456789012345
# Restored archive "before-deploy" (1718123456789012345) → .env
```

## Storage

Archives are stored in `.envoy/archive/` relative to the working directory. Each archive is a JSON file named by its Unix nanosecond timestamp ID.

## Notes

- Archive IDs are unique Unix nanosecond timestamps.
- Archives are stored locally and not committed to version control (add `.envoy/` to `.gitignore`).
- Use `snapshot` for lightweight point-in-time saves without labels.
