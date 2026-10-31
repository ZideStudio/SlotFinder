# 🎭 Playwright BDD tests

End-to-end tests written in Gherkin with `playwright-bdd`, with the API mocked by MSW.

## 📑 Table of Contents

- [Folder Organization](#folder-organization)
- [Where to put a step](#where-to-put-a-step)
- [Network mocking](#network-mocking)
- [Commands](#commands)
- [Artifacts](#artifacts)

## <span id="folder-organization">Folder Organization</span>

```
tests/
  common/                  — Reusable steps, not tied to a feature
    fixtures.ts            — `test`, MSW `network` fixture, Given/When/Then
    content.steps.ts       — Check headings and texts
    interactions.steps.ts  — Click buttons, fill fields
  features/
    <domain>/
      <domain>.feature     — Gherkin scenarios
      <domain>.steps.ts    — Steps specific to this domain
  .artifacts/              — Generated, ignored by git
```

`playwright.config.ts` stays at the root of `front/`.

## <span id="where-to-put-a-step">Where to put a step</span>

- Specific to one page or domain: `features/<domain>/<domain>.steps.ts`.
- Generic and reusable: `common/*.steps.ts`.
- Always import `Given`, `When`, `Then` from `common/fixtures.ts`.

## <span id="network-mocking">Network mocking</span>

- The auto `network` fixture uses `defaultHandlers` (signed in). Tag a feature or scenario `@unauthenticated` to start signed out.
- Override handlers in a step with `network.use(...)`.
- An API request without a handler makes the test fail.

## <span id="commands">Commands</span>

```bash
npm run test:playwright     # Generate and run all scenarios
npm run test:playwright:ui  # Playwright UI with BDD regeneration on file changes
```

## <span id="artifacts">Artifacts</span>

Generated files live in `tests/.artifacts/`:

- `bdd-gen/` — generated specs
- `report/` — HTML report
- `results/` — traces, screenshots, videos
