# AI engineering template

Exploratory work created to investigate repository context, bounded agent workflows, and evidence-driven repair. This is a working demonstration, not a production-ready platform. See the [quick security review](docs/evidence/security-review.md) for scope and limitations.

A sample environment where an engineer can hand an AI a bounded problem, let it work independently, and review evidence that the result is correct.

The design has six parts: discoverable context, reproducible tools, small tasks, executable constraints, observable results, and human ownership. The repository is the shared memory. Tests and CI provide feedback. Engineers own product intent, tradeoffs, and release decisions.

## Try it

Requires Go 1.25+, Git, and Make. No API key, dependency installation, database, or containers required.

```sh
make doctor
make check
make dev
# In another terminal:
curl http://127.0.0.1:8080/api/v1/estimate?points=3
```

The sample estimates engineering capacity: 3 points × 4 hours = 12 hours. It is deliberately tiny so the engineering environment is easy to inspect. This is an illustrative rule, not a real forecasting model.

## Repository map

```text
AGENTS.md                         Short entry point for coding agents
Makefile                          Shared commands for engineers, agents, CI
backend/
  AGENTS.md                       Local implementation rules
  cmd/server/                     Process lifecycle and configuration
  internal/handlers/              HTTP boundary and contract tests
  internal/services/              Pure business rules and unit tests
  internal/architecture/          Executable import constraints
api/openapi.yaml                  Public API contract
scripts/doctor.sh                 Local prerequisite check
scripts/check.sh                  Formatting, architecture, tests, vet, build
.github/
  workflows/ci.yml                Same verification as local development
  pull_request_template.md        Review evidence and release impact
  ISSUE_TEMPLATE/task.md          Outcome-driven task intake
  CODEOWNERS.example              Ownership template to activate on adoption
docs/
  architecture.md                 Boundaries and extension points
  environment.md                  Workspaces, permissions, integrations, rollout
  workflows.md                    Engineer → agent → checks → reviewer loop
  decisions/0001-small-core.md     Rationale and tradeoffs
  tasks/                          Reusable brief and a worked task
  runbooks/local-development.md   Start, inspect, diagnose, stop
  evidence/                       Durable verification record template
evals/
  README.md                       Repeatable agent assessment protocol
  tasks/                          Held-out task prompts
  graders/                        Behavioral acceptance test for an exercise
```

## Start an AI-assisted change

Give your coding assistant this task:

> Read AGENTS.md and docs/tasks/example-estimate.md. Run the baseline checks, inspect the relevant implementation, and complete the task. Add meaningful tests and run make check. Report the behavior change, evidence, and remaining limitations.

That example is already implemented; use `evals/tasks/estimate-cap.md` for a fresh exercise. Keep each task in its own branch or worktree. Tool-specific adapters should point at AGENTS.md rather than duplicate it; configure your assistant to load it if it does not do so automatically.

## What makes this useful

| Need | Implementation |
| --- | --- |
| Agent finds the right context | Short root instructions link to focused documents |
| Agent can verify its work | One command, local and in CI |
| Architecture survives repeated edits | Import rules fail as tests |
| Review does not depend on chat history | Task acceptance criteria and evidence template |
| Work can resume in another session | Handoff records paths, decisions, and remaining work |
| Tool/model changes can be measured | Repeatable task, independent grader, review rubric |
| Engineers retain control | Workspace isolation and externally enforced permissions |

This is a working core, not a provisioned company platform. Hosting rules, secret stores, remote workspaces, issue trackers, telemetry backends, and actual AI runners must be connected by the adopting team. See [the environment design](docs/environment.md) for that layer and a rollout sequence.

The executable example is backend-only and uses the standard library to avoid dependency setup. A production extension follows the workspace defaults: Go + chi, React 18 + TypeScript + Vite, shadcn/ui + Tailwind, TanStack Query, and MongoDB when persistence is needed. Add these with approved dependency installation and appropriate checks; empty frontend/database scaffolds would not improve this example.

## Design references

The evaluation protocol draws on [Anthropic’s agent evaluation guidance](https://www.anthropic.com/engineering/demystifying-evals-for-ai-agents): assess task outcomes, retain execution evidence, and use multiple trials. The human review loop is consistent with [Humans and Agents in Software Engineering Loops](https://martinfowler.com/articles/exploring-gen-ai/humans-and-agents.html). The concrete repository layout and policies here are design choices for this sample.

## Beyond assisted coding: incident-driven repair

Run `make frontier-demo` for a [working repair-loop rehearsal](examples/repair-loop/README.md): reproduce a synthetic incident, reject an overfitted fix with independent tests, and emit a reviewable patch with evidence. The model proposals are recorded fixtures; the tests, rejection gates, and artifacts are real. The walkthrough explains how to extend it to a live agent and shadow verification.
