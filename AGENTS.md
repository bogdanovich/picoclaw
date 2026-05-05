# Local Working Notes

This fork uses `bogdanovich/picoclaw:fork/main` as the active local/runtime
development branch. `bogdanovich/picoclaw:main` is kept as a clean mirror of
upstream for contribution hygiene.

Branch policy:

- `fork/main` is the public feature-rich fork branch with all changes we
  actually run locally.
- `main` tracks `sipeed/picoclaw:main` and should stay clean.
- `upstream-mirror` is an optional local backup mirror of `sipeed/picoclaw:main`.
- Rebase or rebuild `fork/main` onto the latest `upstream/main` periodically to
  stay current.
- Do not open upstream PRs directly from `fork/main`.
- For upstream PRs, create a clean topic branch from the latest `upstream/main`
  or `upstream-mirror`, then cherry-pick or manually port only the intended
  patch.
- Do not use a `[codex]` prefix in PR titles.
- Use conventional PR titles with a functional scope and colon, such as
  `feat(providers): add Gemini search`, `fix(telegram): handle media groups`,
  `fix(agents): preserve topic routing`, or `feat(tools): add update_plan`.
