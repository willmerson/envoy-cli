package envfile

import "fmt"

// PatchOp represents a single patch operation.
type PatchOp struct {
	Op    string // "set", "delete", "rename"
	Key   string
	Value string
	NewKey string
}

// PatchResult holds the outcome of applying a patch.
type PatchResult struct {
	Applied []string
	Skipped []string
	Errors  []string
}

// Patch applies a list of PatchOps to entries and returns the modified slice.
func Patch(entries []Entry, ops []PatchOp) ([]Entry, PatchResult) {
	result := PatchResult{}
	out := make([]Entry, len(entries))
	copy(out, entries)

	for _, op := range ops {
		switch op.Op {
		case "set":
			found := false
			for i, e := range out {
				if e.Key == op.Key {
					out[i].Value = op.Value
					found = true
					break
				}
			}
			if !found {
				out = append(out, Entry{Key: op.Key, Value: op.Value})
			}
			result.Applied = append(result.Applied, fmt.Sprintf("set %s", op.Key))

		case "delete":
			newOut := out[:0]
			deleted := false
			for _, e := range out {
				if e.Key == op.Key {
					deleted = true
					continue
				}
				newOut = append(newOut, e)
			}
			out = newOut
			if deleted {
				result.Applied = append(result.Applied, fmt.Sprintf("delete %s", op.Key))
			} else {
				result.Skipped = append(result.Skipped, fmt.Sprintf("delete %s (not found)", op.Key))
			}

		case "rename":
			if op.NewKey == "" {
				result.Errors = append(result.Errors, fmt.Sprintf("rename %s: new_key is empty", op.Key))
				continue
			}
			found := false
			for i, e := range out {
				if e.Key == op.Key {
					out[i].Key = op.NewKey
					found = true
					break
				}
			}
			if found {
				result.Applied = append(result.Applied, fmt.Sprintf("rename %s -> %s", op.Key, op.NewKey))
			} else {
				result.Skipped = append(result.Skipped, fmt.Sprintf("rename %s (not found)", op.Key))
			}

		default:
			result.Errors = append(result.Errors, fmt.Sprintf("unknown op: %s", op.Op))
		}
	}

	return out, result
}
