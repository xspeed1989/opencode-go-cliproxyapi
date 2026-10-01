package plugin

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

// Reproduce the real failing path: a Responses client requests a DeepSeek
// model whose upstream protocol is Chat Completions. Catalog levels are
// deliberately contradictory; the exact client effort must still reach
// the mock upstream, for both streaming and non-streaming execution.
func TestDeepSeekReasoningEffortPassthrough(t *testing.T) {
	for _, stream := range []bool{false, true} {
		for _, effort := range []string{"xhigh", "max", "none", "auto", "ultra", " Custom-Tier ", " "} {
			t.Run(fmt.Sprintf("stream=%t/effort=%q", stream, effort), func(t *testing.T) {
				m, f, st, yamlText := newIntegrationManager(t)
				st.setCatalogBody(`{"data":[{"id":"deepseek-v4.1-flash","supported_reasoning_levels":[{"effort":"low"}]}]}`)
				resp := mustHandle(t, m, "plugin.reconfigure", lifecycleRequestBody(yamlText))
				if env := decodeEnv(t, resp); !env.OK {
					t.Fatalf("reconfigure failed: %s", resp)
				}

				body, err := json.Marshal(map[string]any{
					"model":     "opencode-go/deepseek-v4.1-flash",
					"input":     "hi",
					"stream":    stream,
					"reasoning": map[string]any{"effort": effort},
				})
				if err != nil {
					t.Fatalf("marshal request: %v", err)
				}
				if stream {
					resp, err = m.HandleCall("executor.execute_stream",
						execStreamReqBodyForKey("opencode-go/deepseek-v4.1-flash", "openai-response", body, "down-effort", "sk-test-1"))
					if err != nil {
						t.Fatalf("execute_stream: %v", err)
					}
					if env := decodeEnv(t, resp); !env.OK {
						t.Fatalf("plugin rejected %q: %s", effort, resp)
					}
					m.bridge.WaitForInFlight(5 * time.Second)
					assertCleanStreamClose(t, f)
				} else {
					if env := mustExecute(t, m, "opencode-go/deepseek-v4.1-flash", "openai-response", body); !env.OK {
						t.Fatalf("plugin rejected %q: %+v", effort, env.Error)
					}
				}

				path, _, upstream := st.chatRequest(t)
				if path != "/v1/chat/completions" || upstream["model"] != "deepseek-v4.1-flash" {
					t.Fatalf("wrong upstream route/model: %s %v", path, upstream["model"])
				}
				if got := upstream["reasoning_effort"]; got != effort {
					t.Fatalf("upstream effort = %q, want exact client value %q", got, effort)
				}
				if got, _ := upstream["stream"].(bool); got != stream {
					t.Fatalf("upstream stream = %v, want %v", got, stream)
				}
			})
		}
	}
}

// Removing the local gate must not hide a real upstream rejection. The
// mock records the request and returns an explicit 400 with a distinct
// message, proving that level acceptance remains the upstream's decision.
func TestDeepSeekUpstreamEffortRejectionPreserved(t *testing.T) {
	m, _, st, yamlText := newIntegrationManager(t)
	st.setCatalogBody(`{"data":[{"id":"deepseek-v4.1-flash"}]}`)
	resp := mustHandle(t, m, "plugin.reconfigure", lifecycleRequestBody(yamlText))
	if env := decodeEnv(t, resp); !env.OK {
		t.Fatalf("reconfigure failed: %s", resp)
	}
	st.setChatPlan(func(int) (int, string, string) {
		return 400, "", `{"error":{"message":"upstream rejected max"}}`
	})
	env := mustExecute(t, m, "opencode-go/deepseek-v4.1-flash", "openai-response",
		[]byte(`{"input":"hi","reasoning":{"effort":"max"}}`))
	if env.OK || env.Error == nil || env.Error.HTTPStatus != 400 {
		t.Fatalf("expected upstream 400, got %+v", env.Error)
	}
	if !strings.Contains(env.Error.Message, "upstream rejected max") {
		t.Fatalf("upstream error was replaced: %q", env.Error.Message)
	}
	if hits, _ := st.chatCalls(); hits != 1 {
		t.Fatalf("upstream must receive the request, got %d calls", hits)
	}
	_, _, upstream := st.chatRequest(t)
	if upstream["reasoning_effort"] != "max" {
		t.Fatalf("upstream effort changed: %v", upstream["reasoning_effort"])
	}
}
