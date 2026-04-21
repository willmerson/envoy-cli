# envoy patch

Apply a set of patch operations to a `.env` file using a JSON patch spec.

## Usage

```bash
envoy patch --patch-file ops.json [--env .env]
```

## Patch File Format

The patch file is a JSON file with an `ops` array. Each operation has an `op` field and relevant parameters.

### Supported Operations

| Op       | Fields              | Description                          |
|----------|---------------------|--------------------------------------|
| `set`    | `key`, `value`      | Set or add a key-value pair          |
| `delete` | `key`               | Remove a key from the file           |
| `rename` | `key`, `new_key`    | Rename an existing key               |

## Example

```json
{
  "ops": [
    { "op": "set",    "key": "PORT",  "value": "9090" },
    { "op": "delete", "key": "DEBUG" },
    { "op": "rename", "key": "HOST",  "new_key": "HOSTNAME" }
  ]
}
```

Given `.env`:
```
HOST=localhost
PORT=8080
DEBUG=true
```

After running `envoy patch --patch-file ops.json`:
```
HOSTNAME=localhost
PORT=9090
```

## Flags

| Flag            | Default | Description                        |
|-----------------|---------|------------------------------------|
| `--patch-file`  | —       | Path to JSON patch spec (required) |
| `--env`         | `.env`  | Target env file path               |

## Output

Each operation is reported as `applied`, `skipped`, or `error`:

```
applied: set PORT
applied: delete DEBUG
applied: rename HOST -> HOSTNAME
patch applied to .env
```
