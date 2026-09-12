# Evaluate the coding environment

Application tests measure application correctness. These exercises measure whether an assistant can deliver a change using this environment. They do not run a model automatically.

## Exercise protocol

1. Establish a clean committed baseline after initializing the template repository. Record its commit. Use a fresh worktree from that same revision per trial.
2. Record assistant/model version, instruction revision, allowed tools, time/cost limits, and task ID. Give the assistant only the prompt in tasks/estimate-cap.md, plus normal repository context. In a formal evaluation keep graders in an evaluator-owned checkout unavailable to the assistant.
3. Allow the assistant to work. Record elapsed time, human interventions, reported test results, and cost if available. Do not provide grader outputs during the trial.
4. After it finishes, run make check. Copy graders/estimate_cap_test.go into backend/internal/services/estimate_cap_eval_test.go in the candidate worktree, then run make test. The grader fails on the original baseline by design. It must remain unchanged.
5. Have an engineer inspect the diff and contract. Grade each criterion below. Keep sanitized results outside the product history or in ignored evals/results/.
6. Run at least three fresh trials per setup; report all outcomes. Compare against the same baseline and task set. Expand the task set using real failures before drawing broad conclusions.

## Rubric

| Criterion | Evidence | Pass condition |
| --- | --- | --- |
| Correctness | Independent grader | All behavioral assertions pass |
| Regression safety | make check | Existing suite passes |
| Contract consistency | Human review | API maximum is 20 everywhere |
| Scope | Human diff review | No unrelated edits or weakened tests |
| Integrity | Run record | Claims match executed checks |

A trial succeeds only if every criterion passes. Report success fraction, median elapsed time, intervention count, and cost where available. Small sample sizes are directional, not proof of superiority. A public grader is useful for training but is not a tamper-resistant benchmark; separate evaluator access for comparative measurements.
