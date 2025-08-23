package main

import (
	"github.com/aws-private-5g/mcp-server/config"
	"github.com/aws-private-5g/mcp-server/models"
	tools_v1 "github.com/aws-private-5g/mcp-server/tools/v1"
	tools_ping "github.com/aws-private-5g/mcp-server/tools/ping"
	tools_tags "github.com/aws-private-5g/mcp-server/tools/tags"
)

func GetAll(cfg *config.APIConfig) []models.Tool {
	return []models.Tool{
		tools_v1.CreateListdeviceidentifiersTool(cfg),
		tools_v1.CreateListnetworkresourcesTool(cfg),
		tools_v1.CreateGetnetworkresourceTool(cfg),
		tools_v1.CreateAcknowledgeorderreceiptTool(cfg),
		tools_v1.CreateActivatedeviceidentifierTool(cfg),
		tools_v1.CreateActivatenetworksiteTool(cfg),
		tools_v1.CreateListnetworksitesTool(cfg),
		tools_v1.CreateDeletenetworksiteTool(cfg),
		tools_v1.CreateGetnetworksiteTool(cfg),
		tools_v1.CreateListnetworksTool(cfg),
		tools_v1.CreateCreatenetworkTool(cfg),
		tools_ping.CreatePingTool(cfg),
		tools_v1.CreateDeactivatedeviceidentifierTool(cfg),
		tools_v1.CreateStartnetworkresourceupdateTool(cfg),
		tools_tags.CreateListtagsforresourceTool(cfg),
		tools_tags.CreateTagresourceTool(cfg),
		tools_v1.CreateUpdatenetworksiteTool(cfg),
		tools_v1.CreateDeletenetworkTool(cfg),
		tools_v1.CreateGetnetworkTool(cfg),
		tools_v1.CreateGetorderTool(cfg),
		tools_tags.CreateUntagresourceTool(cfg),
		tools_v1.CreateCreatenetworksiteTool(cfg),
		tools_v1.CreateUpdatenetworksiteplanTool(cfg),
		tools_v1.CreateConfigureaccesspointTool(cfg),
		tools_v1.CreateGetdeviceidentifierTool(cfg),
		tools_v1.CreateListordersTool(cfg),
	}
}
