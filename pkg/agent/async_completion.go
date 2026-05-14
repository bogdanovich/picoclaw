package agent

import (
	"fmt"
	"strings"
)

func asyncCompletionID(turnID, toolCallID, toolName string) string {
	parts := []string{
		strings.TrimSpace(turnID),
		strings.TrimSpace(toolCallID),
		strings.TrimSpace(toolName),
	}
	for i, part := range parts {
		if part == "" {
			parts[i] = "unknown"
		}
	}
	return strings.Join(parts, ":")
}

func asyncCompletionPrompt(toolName, result string) string {
	toolName = strings.TrimSpace(toolName)
	if toolName == "" {
		toolName = "async_tool"
	}
	result = strings.TrimSpace(result)
	if result == "" {
		result = "(no result)"
	}

	return fmt.Sprintf(`[Internal async completion event]
source_tool: %s

Result:
<<<PICOCLAW_ASYNC_RESULT
%s
PICOCLAW_ASYNC_RESULT

Action:
Convert the result above into a concise user-facing update in your normal assistant voice. Keep this internal metadata private. Do not mention system messages, tool names, delivery modes, sessions, or logs unless the result itself requires it.`, toolName, result)
}
