package envfile

import (
	"fmt"
	"sort"
	"time"
)

// RollbackResult describes the outcome of a rollback operation.
type RollbackResult struct {
	SnapshotName string
	Restored     int
	Removed      int
	Added        int
	Timestamp    time.Time
}

// Rollback restores entries to the state captured in the given snapshot.
// It returns a RollbackResult summarising what changed.
func Rollback(current []Entry, snapshot []Entry) ([]Entry, RollbackResult) {
	currentMap := make(map[string]string, len(current))
	for _, e := range current {
		currentMap[e.Key] = e.Value
	}

	snapshotMap := make(map[string]string, len(snapshot))
	for _, e := range snapshot {
		snapshotMap[e.Key] = e.Value
	}

	result := RollbackResult{Timestamp: time.Now()}

	// Build restored entries from snapshot order.
	restored := make([]Entry, 0, len(snapshot))
	for _, e := range snapshot {
		prev, exists := currentMap[e.Key]
		if !exists {
			result.Added++
		} else if prev != e.Value {
			result.Restored++
		}
		restored = append(restored, e)
	}

	// Count keys present in current but absent from snapshot (will be removed).
	for _, e := range current {
		if _, ok := snapshotMap[e.Key]; !ok {
			result.Removed++
		}
	}

	return restored, result
}

// FormatRollbackResult returns a human-readable summary of the rollback.
func FormatRollbackResult(r RollbackResult) string {
	return fmt.Sprintf(
		"Rollback complete: %d restored, %d added, %d removed (snapshot: %s)",
		r.Restored, r.Added, r.Removed, r.SnapshotName,
	)
}

// RollbackPlan describes what would change without applying the rollback.
type RollbackPlan struct {
	ToRestore []string
	ToAdd     []string
	ToRemove  []string
}

// PlanRollback computes what a rollback would do without modifying anything.
func PlanRollback(current []Entry, snapshot []Entry) RollbackPlan {
	currentMap := make(map[string]string, len(current))
	for _, e := range current {
		currentMap[e.Key] = e.Value
	}
	snapshotMap := make(map[string]string, len(snapshot))
	for _, e := range snapshot {
		snapshotMap[e.Key] = e.Value
	}

	plan := RollbackPlan{}
	for _, e := range snapshot {
		prev, exists := currentMap[e.Key]
		if !exists {
			plan.ToAdd = append(plan.ToAdd, e.Key)
		} else if prev != e.Value {
			plan.ToRestore = append(plan.ToRestore, e.Key)
		}
	}
	for _, e := range current {
		if _, ok := snapshotMap[e.Key]; !ok {
			plan.ToRemove = append(plan.ToRemove, e.Key)
		}
	}
	sort.Strings(plan.ToRestore)
	sort.Strings(plan.ToAdd)
	sort.Strings(plan.ToRemove)
	return plan
}
