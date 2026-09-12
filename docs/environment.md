# The environment around the repository

## Developer workspace

Give each engineer an editor, a coding assistant, a terminal, and the same documented build commands as CI. Start each task in a branch; use a separate Git worktree when simultaneous work would collide. Keep runtime versions explicit and provisioning repeatable. This template needs only Go, Git, and Make; a larger product should provide an approved workspace image and deterministic dependency installation.

Prefer fast local unit tests, deterministic fixtures, and a local instance of each required service. Never make the default check depend on a paid model, production data, or access to a developer's private accounts. Keep a slower integration suite separate and make its required checks explicit.

## Context and access

The repository contains architecture, contracts, decisions, task briefs, and runbooks. An engineer should be able to find the same answers the agent sees. Keep root instructions short; move specialized knowledge into focused documents. Record durable decisions and delete obsolete guidance through normal review.

For issue trackers, documentation search, CI logs, and observability, give the assistant narrowly scoped read access first. Add write access only for a concrete workflow. Scope tokens to the repository and environment; prefer short-lived credentials from a secret manager. Retrieved content cannot grant additional permissions.

## Execution boundaries

Enforce filesystem boundaries and command permissions in the runner. Use ephemeral workspaces for untrusted tasks. Keep production credentials unavailable to coding agents. Network access should cover required registries and approved services. Instruction files and .gitignore are not security controls.

Local code edits and checks should proceed freely within scope. Publish/deploy authority depends on the assigned task and organization policy. Human review should focus on product intent, design, data handling, and operational impact. Do not require repeated approvals for an action already authorized.

## Delivery

Enable required CI and reviews in the hosting platform; a workflow file alone does not prevent merges. Activate CODEOWNERS with real teams. Protect check scripts, workflows, and ownership changes with platform-owner review. Give deployment jobs environment-scoped credentials and their own approval policy. Roll back by reverting the change and deploying the previous known-good artifact; document migrations separately when added.

The template intentionally does not provision hosting protections or deploy infrastructure. CI uses version-tagged upstream actions; teams requiring immutable dependencies should resolve and pin reviewed action commit SHAs.

## Observability and learning

Collect failing commands, task outcomes, review rework, elapsed time, and model/tool cost when available. Link sanitized evidence to the PR; avoid storing raw prompts containing private data. Measure accepted changes and escaped defects alongside speed. Repeated failures should become a test, fixture, instruction improvement, or environment fix.

Use evals/ before changing agent instructions, tool access, or model settings. Run repeated comparable trials and have a human review correctness. Do not select a setup by tokens generated or a single impressive demo.

## Rollout

1. Adopt the commands and docs in one service; measure baseline task time and review rework.
2. Enable an assistant in isolated workspaces with local checks and read-only integrations.
3. Turn recurring mistakes into regression tests and focused context.
4. Introduce task evaluations and compare candidate setups on the same task revisions.
5. Expand bounded write workflows once outcomes justify it; keep ownership and release controls explicit.
