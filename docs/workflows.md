# Working loop

1. Engineer writes the outcome, acceptance criteria, scope, and known constraints in a task brief.
2. Agent reads the instructions and relevant code, runs baseline checks, and identifies missing decisions. Small changes need no separate planning document.
3. Agent implements one cohesive change with behavioral tests. Use existing utilities and keep the diff focused.
4. Agent runs make check, inspects the diff, and records evidence. A failed baseline is reported separately from regressions.
5. Engineer reviews behavior and tradeoffs with the test results visible. A second AI review may assist but does not replace ownership.
6. Authorized release workflow publishes the reviewed change. Production observations feed the next test or task.

For larger tasks, keep a short living plan in the task brief. At handoff record what is complete, what remains, exact paths, commands and outcomes, and blocking decisions. Do not preserve a full chat transcript as project documentation.

When stuck, state the specific blocker and what evidence established it. When a requirement is ambiguous but work can proceed, document a reasonable assumption and continue independent work. Respect existing authorization; ask for a decision only when necessary.
