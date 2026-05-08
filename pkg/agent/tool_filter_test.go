package agent

import (
	"context"
	"testing"

	"github.com/sipeed/picoclaw/pkg/config"
	"github.com/sipeed/picoclaw/pkg/tools"
)

type testFilterTool struct {
	name string
}

func (t *testFilterTool) Name() string { return t.name }

func (t *testFilterTool) Description() string { return "test tool" }

func (t *testFilterTool) Parameters() map[string]any {
	return map[string]any{
		"type":       "object",
		"properties": map[string]any{},
	}
}

func (t *testFilterTool) Execute(context.Context, map[string]any) *tools.ToolResult {
	return tools.SilentResult("ok")
}

func TestToolAllowedByConfig_DefaultAllow(t *testing.T) {
	if !toolAllowedByConfig("main", nil, "mcp_gpt_researcher_deep_research") {
		t.Fatal("expected nil filter config to allow tool")
	}
}

func TestToolAllowedByConfig_DenyPattern(t *testing.T) {
	cfg := &config.AgentToolsConfig{
		Deny: []string{"mcp_gpt_researcher_*"},
	}
	if toolAllowedByConfig("main", cfg, "mcp_gpt_researcher_deep_research") {
		t.Fatal("expected deny pattern to block tool")
	}
	if !toolAllowedByConfig("main", cfg, "web_fetch") {
		t.Fatal("expected unrelated tool to remain allowed")
	}
}

func TestToolAllowedByConfig_AllowThenDeny(t *testing.T) {
	cfg := &config.AgentToolsConfig{
		Allow: []string{"mcp_*", "web_fetch"},
		Deny:  []string{"mcp_gpt_researcher_*"},
	}
	if toolAllowedByConfig("main", cfg, "read_file") {
		t.Fatal("expected allow list to exclude unspecified tool")
	}
	if !toolAllowedByConfig("main", cfg, "web_fetch") {
		t.Fatal("expected allow list to permit web_fetch")
	}
	if toolAllowedByConfig("main", cfg, "mcp_gpt_researcher_deep_research") {
		t.Fatal("expected deny to override allow")
	}
	if !toolAllowedByConfig("main", cfg, "mcp_inventorydb_get_location") {
		t.Fatal("expected allowed MCP tool to remain visible")
	}
}

func TestRegisterToolIfAllowed(t *testing.T) {
	agent := &AgentInstance{
		ID:    "main",
		Tools: tools.NewToolRegistry(),
		ToolFilter: &config.AgentToolsConfig{
			Deny: []string{"mcp_gpt_researcher_*"},
		},
	}

	if registerToolIfAllowed(agent, &testFilterTool{name: "mcp_gpt_researcher_deep_research"}) {
		t.Fatal("expected denied tool to be skipped")
	}
	if _, ok := agent.Tools.Get("mcp_gpt_researcher_deep_research"); ok {
		t.Fatal("denied tool should not be registered")
	}

	if !registerToolIfAllowed(agent, &testFilterTool{name: "web_fetch"}) {
		t.Fatal("expected allowed tool to register")
	}
	if _, ok := agent.Tools.Get("web_fetch"); !ok {
		t.Fatal("allowed tool should be registered")
	}
}

func TestRegisterHiddenToolIfAllowed(t *testing.T) {
	agent := &AgentInstance{
		ID:    "main",
		Tools: tools.NewToolRegistry(),
		ToolFilter: &config.AgentToolsConfig{
			Deny: []string{"mcp_gpt_researcher_*"},
		},
	}

	if registerHiddenToolIfAllowed(agent, &testFilterTool{name: "mcp_gpt_researcher_deep_research"}) {
		t.Fatal("expected denied hidden tool to be skipped")
	}
	if len(agent.Tools.List()) != 0 {
		t.Fatalf("expected no tools after skipped hidden registration, got %v", agent.Tools.List())
	}

	if !registerHiddenToolIfAllowed(agent, &testFilterTool{name: "mcp_inventorydb_get_location"}) {
		t.Fatal("expected allowed hidden tool to register")
	}
	if len(agent.Tools.List()) != 1 || agent.Tools.List()[0] != "mcp_inventorydb_get_location" {
		t.Fatalf("unexpected hidden tool registry contents: %v", agent.Tools.List())
	}
}
