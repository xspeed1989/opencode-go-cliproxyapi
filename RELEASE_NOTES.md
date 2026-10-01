## What's Changed

### Changes

- Remove error-message redaction. Error text now reaches the client verbatim, and an upstream error body is reported up to 2 KiB instead of an 80-character snippet. This is a deliberate operator decision: the previous redaction made upstream refusals undiagnosable, and the host renders an empty message as its generic `plugin call failed` placeholder.
- Applies everywhere messages are produced: adapter translation failures, upstream status classification, in-stream upstream errors, the stream pump close reasons, and the envelope edge. The 2 KiB bound is shared (`shared.ErrorSnippetLimit`), so an oversized body is still reported as a head plus a truncation marker rather than echoed whole.

### Security notice

- **Error output is no longer scrubbed.** A credential that an upstream, transport failure, or malformed payload echoes in its error text will now surface in the client-visible error message and in any log that records it. Treat client-visible errors and CPA logs as sensitive. Nothing else about the error path changed; only redaction was removed and the snippet bound widened.

## Upgrade Notes

- Replace the old plugin binary with the new release binary, or update from this fork's plugin-store source.
- Reload or restart CLIProxyAPI as your deployment requires.
- If you need redaction restored, do not deploy this release: it is a deliberate behaviour change, not a regression.

**Full Changelog**: https://github.com/xspeed1989/opencode-go-cliproxyapi/compare/v0.1.12...v0.1.13
