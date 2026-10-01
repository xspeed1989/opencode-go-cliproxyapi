## What's Changed

### Bug Fixes

- Report the upstream reason when a streaming request is refused. A refusal with HTTP >= 400 before the first byte used to produce an error envelope with an empty message, which CLIProxyAPI renders as its generic `plugin call failed` (HTTP 400). That hid both the status and the upstream explanation, because `debugTrace` is compiled out of release builds. The stream open-failure path now reads a bounded prefix of the refused response body and classifies with that snippet, so the real cause is visible.
- Guarantee that every classified error carries a message. An empty message is replaced by a class-derived fallback naming the failure class and, when known, the upstream status (for example `upstream returned HTTP 400 (unsupported_protocol_or_parameter) without a message`). The guarantee applies in the constructors and again at the envelope edge, which also covers in-stream upstream failures that arrive without a message.
- Keep the added read bounded and safe: at most 8 KiB across at most 8 host reads per refused stream, and the snippet stays redacted and truncated, so an oversized upstream error body is never echoed back in full.

## Upgrade Notes

- Replace the old plugin binary with the new release binary, or update from this fork's plugin-store source.
- Reload or restart CLIProxyAPI as your deployment requires.
- A failed streaming request now surfaces the upstream status and a redacted snippet of its error body instead of `plugin call failed`. No new configuration switch, and no other behaviour changes.

**Full Changelog**: https://github.com/xspeed1989/opencode-go-cliproxyapi/compare/v0.1.11...v0.1.12
