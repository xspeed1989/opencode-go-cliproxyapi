## What's Changed

### Fix

- Normalize replayed tool-call ids on the Responses route. Every `function_call` / `function_call_output` id in an upstream-bound Responses request is rewritten into the upstream's own marked shape (`call_<NN>_ET_<suffix>`) unless it already carries it. Applies to all three source formats: native Responses passthrough, Chat Completions, and Anthropic Messages.

### Why

The upstream this plugin fronts runs DeepSeek in thinking mode. A replayed history is accepted only when, for every tool-call turn, one of these holds:

1. the turn's reasoning item is replayed (`include: ["reasoning.encrypted_content"]` history), or
2. the call id carries the upstream's own marker (`call_NN_ET_…`) — such calls are accepted without reasoning, or
3. the upstream still recognizes the id from its own live session state.

Otherwise the request is refused with HTTP 400 `The \`reasoning_text\` in the thinking mode must be passed back to the API.` Measured against the live endpoint with an otherwise identical body:

| replayed history | reasoning item | call id | result |
| --- | --- | --- | --- |
| any id, including `call_1` | present | — | 200 |
| reasoning stripped | absent | `call_00_ET_…` | 200 |
| reasoning stripped | absent | `call_1`, `call_00_ZZZ…`, `call_00_<random>` | 400 |

Clients that never store thinking (or store it only sometimes), histories produced while a model was served through another route (whose ids the endpoint cannot know), and sessions whose upstream state expired therefore failed permanently on every follow-up request — the reasoning the error asks for cannot be reconstructed after the fact.

Marking the ids makes condition 2 always true, so a reasoning-free replay is valid. The mapping is a pure function of the original id (already-marked ids pass through untouched), so a call and its output always normalize to the same value and repeated replays stay stable. Verified end-to-end against the live endpoint: a reasoning-free history with foreign ids returns 400 before the change and 200 after it.

### Upgrade Notes

- Update from this fork's plugin-store source (or replace the plugin binary) and reload CLIProxyAPI.
- Already-failing sessions are recoverable: the rewrite happens on every request, so a resumed transcript with previously unaccepted ids is healed — no need to start over.
- Ids seen by the upstream may differ from the ids the client sent; they remain opaque strings in every protocol that carries them.
- No other behavior changed: bodies without tool calls are still forwarded byte-identically.

**Full Changelog**: https://github.com/xspeed1989/opencode-go-cliproxyapi/compare/v0.1.13...v0.1.14
