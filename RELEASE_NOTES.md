## What's Changed

### Change

- Limit the client → upstream reasoning replay (added in v0.1.15) to the upstream families whose Chat Completions endpoints actually require it: the `deepseek*` models. Replayed reasoning items and Anthropic `thinking` blocks still become `reasoning_content` there, and still never do for any other Chat Completions-routed model (`glm*`, `kimi*`, `mimo*`, `hy*`, `longcat*`), which keep the exact wire shape of earlier releases.
- The upstream → client direction is deliberately left ungated: it only forwards a chain of thought the upstream actually streamed, so it cannot fail a request, and it keeps working for every model that emits the field.

### Why

Replaying a field an endpoint never asked for is the only direction that can fail a request. Sending `reasoning_content` to a Chat Completions upstream that has no use for it adds no value, while a stricter endpoint could reject the unknown field; the response direction has no such failure mode and only affects what the client sees.

### Cost

If a non-DeepSeek Chat Completions model turns out to need the replay (the endpoint enforcing the same "must be passed back" rule), that model keeps failing until its prefix is added to `reasoningReplayFamilies` in `internal/adapter/chatcompletions/request.go` — a one-line change plus a test.

### Upgrade Notes

- Update from this fork's plugin-store source (or replace the plugin binary) and reload CLIProxyAPI.
- DeepSeek-routed models behave exactly as in v0.1.15. Every other Chat Completions-routed model behaves as it did before v0.1.15 for the request direction, while still surfacing whatever reasoning the upstream streams.
- The Messages and Responses routes are untouched, as are the v0.1.14 tool-call id normalization and the v0.1.13 verbatim error reporting.

**Full Changelog**: https://github.com/xspeed1989/opencode-go-cliproxyapi/compare/v0.1.15...v0.1.16
