# group

Partition `.env` entries by a common key prefix and display them in labelled sections.

## Usage

```
envoy group [flags]
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--env` | `-e` | `.env` | Path to the .env file |
| `--separator` | `-s` | `_` | Character used to split key segments |
| `--depth` | `-d` | `1` | Number of prefix segments to form the group label |
| `--ungrouped` | `-u` | `other` | Label for entries whose key has no separator |
| `--output` | `-o` | `summary` | Output format: `summary` or `keys` |

## Examples

### Default grouping (depth 1)

```
$ envoy group
[APP] (2 keys)
  APP_ENV=production
  APP_NAME=myapp
[DB] (2 keys)
  DB_HOST=localhost
  DB_PORT=5432
[other] (1 keys)
  STANDALONE=yes
```

### Depth-2 grouping for nested namespaces

```
$ envoy group --depth 2
[AWS_EC2] (1 keys)
  AWS_EC2_AMI=ami-123
[AWS_S3] (2 keys)
  AWS_S3_BUCKET=my-bucket
  AWS_S3_REGION=us-east-1
```

### List group names only

```
$ envoy group --output keys
APP (2)
DB (2)
other (1)
```

### Custom separator

```
$ envoy group --separator -
[DB] (2 keys)
  DB-HOST=localhost
  DB-PORT=5432
```

## Notes

- Groups are always displayed in alphabetical order.
- Entries whose key contains no separator character are placed in the `--ungrouped` bucket.
- Combine with `--depth 2` to group deeply-namespaced keys such as `AWS_S3_*` vs `AWS_EC2_*`.
