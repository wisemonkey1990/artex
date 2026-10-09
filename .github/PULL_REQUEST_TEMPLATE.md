# Pull request

Before submitting, please read the contribution guidelines in [CONTRIBUTING.md](../CONTRIBUTING.md).

## Summary

<!-- Briefly describe what changed and why. -->

## Change type

<!-- Mark all applicable items with x. -->

- [ ] Bug fix (`fix`)
- [ ] Feature (`feat`)
- [ ] Documentation (`docs`)
- [ ] Localization (`i18n`)
- [ ] Refactoring or maintenance (`refactor` / `chore`)
- [ ] Other:

## Related issue

<!-- For example: Closes #123 -->

## Verification

<!-- List the commands you ran and their results. -->

- [ ] Go changes: `go build ./...`, `go vet ./agent/`, and `go test ./agent/` pass
- [ ] Web changes: `npm run check` and `npm run build` pass
- [ ] UI changes: screenshots attached

## Localization checklist

<!-- Leave this section blank if the change is unrelated to localization. -->

- [ ] Internal agent reasoning prompts (`agent/promptcatalog.go` and `agent_prompts`) were not translated.
- [ ] Only user-facing output was localized.
- [ ] Commands, payloads, code, URLs, and original logs remain unchanged.
- [ ] UI strings use keys in `web/messages/zh.json`.

## Authorized use

- [ ] This change follows the [authorized-use policy](../README.md#authorized-use-and-legal-notice). Validation used owned or authorized targets, or a local isolated environment.
- [ ] Test and scan artifacts are not included in the commit.
