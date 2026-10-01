## What's Changed

### Fix

- Carry the chain of thought end to end on the Chat Completions route instead of dropping it:
  - **Upstream → client**: `reasoning_content` streamed by the upstream is no longer discarded. A Responses client receives it as a reasoning item that leads the output (`response.output_item.added` → `response.reasoning_text.delta` → `response.reasoning_text.done` → `response.output_item.done`, item carrying one `reasoning_text` part and an explicit empty `summary`), and the terminal `response.completed` output leads with the same item. Non-streaming conversions carry it too.
  - **Client → upstream**: replayed reasoning items map to `reasoning_content` on the assistant message of their own turn (the message carrying that turn's `tool_calls`), and Anthropic `thinking` blocks do the same. `redacted_thinking` stays omitted — it is encrypted metadata with no plaintext to carry.
- Decode the Anthropic `thinking` wire field: thinking blocks carry their plaintext in `thinking`, not `text`, so the block decoder previously normalized it to an empty string and every Claude-source translator saw blank thinking.

### Why

The upstreams this plugin fronts require their chain of thought back once `tools` are in play, and the replayed reasoning is spliced into the model's context rather than merely validated. The Chat Completions adapter parsed upstream chunks into typed structs (`content`, `tool_calls`) with no reasoning field and omitted reasoning items when building requests, so the thinking was dropped in both directions. That does not fail loudly: the tool loop simply continues without the earlier plan, exclusions, and intermediate conclusions, which shows up as repeated work, dropped plans, and target drift across turns.

### Scope

- Chat Completions route only. The native Responses passthrough and the Chat Completions passthrough already forward reasoning verbatim in both directions and are unchanged.
- The Messages route's Responses-to-Anthropic synthesis still drops upstream thinking (no Anthropic `signature` exists upstream); that remains a known gap, deliberately not changed here.

### Upgrade Notes

- Update from this fork's plugin-store source (or replace the plugin binary) and reload CLIProxyAPI.
- Responses-speaking clients paired with Chat Completions-routed models (`deepseek-v4-pro`, `glm-*`, and friends) are the main beneficiaries: their tool turns now carry the plan forward.
- Extra reasoning content means more input tokens per replay, which is the documented behaviour of these upstreams; a client that prefers to drop it can do so in its own history.
- No other behaviour changed: the tool-call id normalization from v0.1.14 and the verbatim error reporting from v0.1.13 are untouched.

**Full Changelog**: https://github.com/xspeed1989/opencode-go-cliproxyapi/compare/v0.1.14...v0.1.15
