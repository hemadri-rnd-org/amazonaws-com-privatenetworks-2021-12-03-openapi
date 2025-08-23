package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/aws-private-5g/mcp-server/config"
	"github.com/aws-private-5g/mcp-server/models"
	"github.com/mark3labs/mcp-go/mcp"
)

func DeletenetworkHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
		}
		networkArnVal, ok := args["networkArn"]
		if !ok {
			return mcp.NewToolResultError("Missing required path parameter: networkArn"), nil
		}
		networkArn, ok := networkArnVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid path parameter: networkArn"), nil
		}
		queryParams := make([]string, 0)
		if val, ok := args["clientToken"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("clientToken=%v", val))
		}
		queryString := ""
		if len(queryParams) > 0 {
			queryString = "?" + strings.Join(queryParams, "&")
		}
		url := fmt.Sprintf("%s/v1/networks/%s%s", cfg.BaseURL, networkArn, queryString)
		req, err := http.NewRequest("DELETE", url, nil)
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
		var result models.DeleteNetworkResponse
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

func CreateDeletenetworkTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("delete_v1_networks_networkArn",
		mcp.WithDescription("Deletes the specified network. You must delete network sites before you delete the network. For more information, see <a href="https://docs.aws.amazon.com/private-networks/latest/APIReference/API_DeleteNetworkSite.html">DeleteNetworkSite</a> in the <i>API Reference for Amazon Web Services Private 5G</i>."),
		mcp.WithString("clientToken", mcp.Description("Unique, case-sensitive identifier that you provide to ensure the idempotency of the request. For more information, see <a href=\"https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/Run_Instance_Idempotency.html\">How to ensure idempotency</a>.")),
		mcp.WithString("networkArn", mcp.Required(), mcp.Description("The Amazon Resource Name (ARN) of the network.")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    DeletenetworkHandler(cfg),
	}
}
