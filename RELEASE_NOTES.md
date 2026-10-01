## What's Changed

### Reasoning Effort Pass-through

- Forward explicit reasoning efforts unchanged, including `xhigh`, `max`, unknown values, case, and whitespace. Remove local capability validation, fallback allowlists, and level downgrades; the upstream owns acceptance and its errors remain visible.
- Translate only the field name across protocols: Chat Completions uses `reasoning_effort`, Responses uses `reasoning.effort`, and Anthropic Messages uses `output_config.effort`. Explicit effort values are no longer converted into guessed thinking budgets or used to change sampling parameters and output limits.
- Preserve explicit Anthropic `output_config.effort` when converting to OpenAI protocols. Numeric budgets retain the fixed threshold conversion only when no explicit effort is supplied.
- Read catalog `supported_reasoning_levels` alongside `thinking` metadata for display and registration, without using it to gate requests.
- Add offline native/cross-protocol, streaming/non-streaming, and DeepSeek executor regression tests. Contradictory catalog levels do not block forwarding, and real upstream 400 errors are still returned.
- Point the fork's plugin repository metadata to `xspeed1989/opencode-go-cliproxyapi`.

## Upgrade Notes

- Update to `v0.1.11` from this fork's plugin-store source or replace the plugin binary and reload it as supported by your CPA deployment.
- Keep existing plugin keys and configuration. No new configuration switch is required.
- For Messages routes, explicit efforts now use `output_config.effort` instead of inferred `thinking.budget_tokens`. Native Anthropic thinking controls remain unchanged.
- This removes the plugin's own rejection, not validation by Pi, CPA, or the actual upstream. It does not guarantee that every effort is supported or honored.

**Full Changelog**: https://github.com/xspeed1989/opencode-go-cliproxyapi/compare/v0.1.10...v0.1.11
