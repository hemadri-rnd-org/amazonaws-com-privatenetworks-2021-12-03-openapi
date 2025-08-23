package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aws-private-5g/mcp-server/config"
	"github.com/aws-private-5g/mcp-server/models"
	"github.com/mark3labs/mcp-go/mcp"
)

func GetnetworksiteHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
		}
		networkSiteArnVal, ok := args["networkSiteArn"]
		if !ok {
			return mcp.NewToolResultError("Missing required path parameter: networkSiteArn"), nil
		}
		networkSiteArn, ok := networkSiteArnVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid path parameter: networkSiteArn"), nil
		}
		url := fmt.Sprintf("%s/v1/network-sites/%s", cfg.BaseURL, networkSiteArn)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to create request", err), nil
		}
		// Set authentication based on auth type
		// Handle multiple authentication parameters
		if cfg.BearerToken != "" {
			req.Header.Set("X-Amz-Security-Token", cfg.BearerToken)
		}
		req.Header.Set("Accept", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Request failed", err), nil
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to read response body", err), nil
		}

		if resp.StatusCode >= 400 {
			return mcp.NewToolResultError(fmt.Sprintf("API error: %s", body)), nil
		}
		// Use properly typed response
		var result models.GetNetworkSiteResponse
		if err := json.Unmarshal(body, &result); err != nil {
			// Fallback to raw text if unmarshaling fails
			return mcp.NewToolResultText(string(body)), nil
		}

		prettyJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format JSON", err), nil
		}

		return mcp.NewToolResultText(string(prettyJSON)), nil
	}
}

func CreateGetnetworksiteTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("get_v1_network-sites_networkSiteArn",
		mcp.WithDescription("Gets the specified network site."),
		mcp.WithString("networkSiteArn", mcp.Required(), mcp.Description("The Amazon Resource Name (ARN) of the network site.")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    GetnetworksiteHandler(cfg),
	}
}
