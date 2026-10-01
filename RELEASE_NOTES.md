## What's Changed

### Bug Fixes

- Close the Responses tool-call lifecycle in the Chat Completions and Messages adapters. Every announced `function_call` output item now emits `response.function_call_arguments.done` (with the complete accumulated arguments) followed by `response.output_item.done` before the terminal `response.completed`, and the aggregated message item emits its `response.output_item.done`. Responses clients that validate tool-call integrity rejected these otherwise-complete turns with an `unfinished tool call` error (observed with `opencode-go/deepseek-v4.1-flash` through pi).
- Both routes share the new emitter helpers so their synthesized streams stay byte-identical, pinned by the existing route-parity test plus a new lifecycle regression test that fails when an announced item is never closed.

## Upgrade Notes

- Replace the old plugin binary with the new release binary.
- Restart CLIProxyAPI after replacing the plugin.
- Hard-refresh Management Center if the plugin page looks stale.

**Full Changelog**: https://github.com/massiveits/opencode-go-cliproxyapi/compare/v0.1.9...v0.1.10
