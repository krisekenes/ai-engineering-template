package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRepairLoopRejectsOverfitAndProducesEvidence(t *testing.T) {
	root := t.TempDir()
	report, err := run(root)
	if err != nil {
		t.Fatal(err)
	}
	if !report.BaselineReproduced || report.State != "ready-for-review" || len(report.Attempts) != 2 {
		t.Fatalf("unexpected report: %+v", report)
	}
	first, second := report.Attempts[0], report.Attempts[1]
	if !first.ReplayPassed || first.ContractPassed || first.Decision != "rejected" {
		t.Fatalf("overfit escaped: %+v", first)
	}
	if !second.ReplayPassed || !second.ContractPassed || second.SHA256 == first.SHA256 {
		t.Fatalf("invalid winner: %+v", second)
	}
	patch, err := os.ReadFile(filepath.Join(root, "proposal.patch"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(patch), "-\tif points < 1 || points >= 100") || !strings.Contains(string(patch), "+\tif points < 1 || points > 100") {
		t.Fatalf("unexpected patch: %s", patch)
	}
}

func TestEvaluatorFailsClosedOnInvalidSource(t *testing.T) {
	dir := t.TempDir()
	if err := prepare(dir, broken, Incident{100, 400}); err != nil {
		t.Fatal(err)
	}
	if err := write(dir, "estimate.go", "this will not compile"); err != nil {
		t.Fatal(err)
	}
	passed, err := check(dir, "TestIndependentContract")
	if err != nil {
		t.Fatal(err)
	}
	if passed {
		t.Fatal("broken build passed the gate")
	}
}
