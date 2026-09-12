// A deterministic rehearsal of incident-driven repair, not a live AI runner.
package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"go/format"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const broken = `package estimate
func Estimate(points int) (int, bool) {
 if points < 1 || points >= 100 { return 0, false }
 return points * 4, true
}
`

// Recorded proposals make the control loop reproducible without an AI account.
var proposals = []string{
	`package estimate
func Estimate(points int) (int, bool) { return points * 4, true }
`,
	`package estimate
func Estimate(points int) (int, bool) {
 if points < 1 || points > 100 { return 0, false }
 return points * 4, true
}
`,
}

type Incident struct {
	Points        int `json:"points"`
	ExpectedHours int `json:"expectedHours"`
}

type Attempt struct {
	Number         int    `json:"number"`
	SHA256         string `json:"sha256"`
	ReplayPassed   bool   `json:"replayPassed"`
	ContractPassed bool   `json:"contractPassed"`
	Decision       string `json:"decision"`
}

type Report struct {
	Mode               string    `json:"mode"`
	Incident           Incident  `json:"incident"`
	BaselineReproduced bool      `json:"baselineReproduced"`
	Attempts           []Attempt `json:"attempts"`
	State              string    `json:"state"`
}

// The incident becomes a concrete behavioral regression test; it is never a command.
func regression(i Incident) string {
	return fmt.Sprintf(`package estimate
import "testing"
func TestIncidentReplay(t *testing.T) {
 hours, ok := Estimate(%d)
 if !ok || hours != %d { t.Fatalf("incident still present: hours=%%d accepted=%%v", hours, ok) }
}
`, i.Points, i.ExpectedHours)
}

// This independent contract is owned by the evaluator, not the proposer.
const contract = `package estimate
import "testing"
func TestIndependentContract(t *testing.T) {
 for points := -100; points <= 200; points++ {
  hours, ok := Estimate(points)
  wantOK := points >= 1 && points <= 100
  if ok != wantOK { t.Errorf("points=%d accepted=%v want=%v", points, ok, wantOK) }
  if wantOK && hours != points*4 { t.Errorf("points=%d hours=%d", points, hours) }
 }
}
`

func write(dir, name, content string) error {
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0600); err != nil {
		return fmt.Errorf("write %s: %w", name, err)
	}
	return nil
}

func prepare(dir, source string, incident Incident) error {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	formatted, err := format.Source([]byte(source))
	if err != nil {
		return fmt.Errorf("format candidate: %w", err)
	}
	for name, content := range map[string]string{
		"go.mod":           "module example.com/repair-candidate\n\ngo 1.25.0\n",
		"estimate.go":      string(formatted),
		"incident_test.go": regression(incident),
		"contract_test.go": contract,
	} {
		if err := write(dir, name, content); err != nil {
			return err
		}
	}
	return nil
}

func check(dir, testName string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", "test", "-count=1", "-run", "^"+testName+"$", "./...")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if saveErr := write(dir, testName+".log", string(output)); saveErr != nil {
		return false, saveErr
	}
	if ctx.Err() != nil {
		return false, fmt.Errorf("%s timed out: %w", testName, ctx.Err())
	}
	if err == nil {
		return true, nil
	}
	if _, ok := err.(*exec.ExitError); ok {
		return false, nil
	}
	return false, fmt.Errorf("run evaluator: %w", err)
}

func run(root string) (Report, error) {
	incident := Incident{Points: 100, ExpectedHours: 400} // Synthetic, sanitized observation.
	report := Report{Mode: "recorded-proposals", Incident: incident, State: "blocked"}
	baseline := filepath.Join(root, "baseline")
	if err := prepare(baseline, broken, incident); err != nil {
		return report, err
	}
	// First show that the original service passes its old, narrow test.
	if err := write(baseline, "existing_test.go", `package estimate
import "testing"
func TestExisting(t *testing.T) { h, ok := Estimate(3); if !ok || h != 12 { t.Fatal(h, ok) } }
`); err != nil {
		return report, err
	}
	healthy, err := check(baseline, "TestExisting")
	if err != nil {
		return report, err
	}
	if !healthy {
		return report, fmt.Errorf("baseline build/existing test failed; inspect baseline logs")
	}
	passed, err := check(baseline, "TestIncidentReplay")
	if err != nil {
		return report, err
	}
	if passed {
		return report, fmt.Errorf("incident did not reproduce")
	}
	report.BaselineReproduced = true
	fmt.Println("REPRODUCED: 100 points incorrectly rejected; existing test passed")
	for index, source := range proposals {
		dir := filepath.Join(root, fmt.Sprintf("attempt-%d", index+1))
		if err := prepare(dir, source, incident); err != nil {
			return report, err
		}
		replay, err := check(dir, "TestIncidentReplay")
		if err != nil {
			return report, err
		}
		full, err := check(dir, "TestIndependentContract")
		if err != nil {
			return report, err
		}
		candidate, err := os.ReadFile(filepath.Join(dir, "estimate.go"))
		if err != nil {
			return report, err
		}
		attempt := Attempt{Number: index + 1, SHA256: fmt.Sprintf("%x", sha256.Sum256(candidate)), ReplayPassed: replay, ContractPassed: full, Decision: "rejected"}
		if replay && full {
			attempt.Decision = "ready-for-review"
		}
		report.Attempts = append(report.Attempts, attempt)
		fmt.Printf("ATTEMPT %d: replay=%v contract=%v -> %s\n", attempt.Number, replay, full, attempt.Decision)
		if attempt.Decision == "ready-for-review" {
			before, err := os.ReadFile(filepath.Join(baseline, "estimate.go"))
			if err != nil {
				return report, err
			}
			// git diff --no-index needs no commits and emits a standard reviewable patch.
			cmd := exec.Command("git", "diff", "--no-index", "--", filepath.Join(baseline, "estimate.go"), filepath.Join(dir, "estimate.go"))
			diff, err := cmd.CombinedOutput()
			if err != nil {
				if exit, ok := err.(*exec.ExitError); !ok || exit.ExitCode() != 1 {
					return report, fmt.Errorf("create diff: %w", err)
				}
			}
			if string(before) == string(candidate) {
				return report, fmt.Errorf("candidate has no changes")
			}
			if err := write(root, "proposal.patch", string(diff)); err != nil {
				return report, err
			}
			report.State = "ready-for-review"
			break
		}
	}
	return report, nil
}

func main() {
	root, err := os.MkdirTemp(".", "repair-run-")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	absolute, err := filepath.Abs(root)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	report, runErr := run(absolute)
	if runErr != nil {
		report.State = "blocked"
	}
	data, err := json.MarshalIndent(report, "", "  ")
	if err == nil {
		err = write(absolute, "report.json", string(data)+"\n")
	}
	fmt.Println("Evidence:", absolute)
	if runErr != nil {
		fmt.Fprintln(os.Stderr, runErr)
		os.Exit(1)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if report.State != "ready-for-review" {
		fmt.Fprintln(os.Stderr, "proposal budget exhausted")
		os.Exit(1)
	}
	fmt.Println("STOP: patch prepared for human review; no source checkout or deployment changed")
}
