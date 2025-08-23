package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"bytes"

	"github.com/aws-private-5g/mcp-server/config"
	"github.com/aws-private-5g/mcp-server/models"
	"github.com/mark3labs/mcp-go/mcp"
)

func StartnetworkresourceupdateHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
		}
		// Create properly typed request body using the generated schema
		var requestBody map[string]interface{}
		
		// Optimized: Single marshal/unmarshal with JSON tags handling field mapping
		if argsJSON, err := json.Marshal(args); err == nil {
			if err := json.Unmarshal(argsJSON, &requestBody); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to convert arguments to request type: %v", err)), nil
			}
		} else {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal arguments: %v", err)), nil
		}
		
		bodyBytes, err := json.Marshal(requestBody)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to encode request body", err), nil
		}
		url := fmt.Sprintf("%s/v1/network-resources/update", cfg.BaseURL)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
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
		var result models.StartNetworkResourceUpdateResponse
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

func CreateStartnetworkresourceupdateTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("post_v1_network-resources_update",
		mcp.WithDescription("<p>Use this action to do the following tasks:</p> <ul> <li> <p>Update the duration and renewal status of the commitment period for a radio unit. The update goes into effect immediately.</p> </li> <li> <p>Request a replacement for a network resource.</p> </li> <li> <p>Request that you return a network resource.</p> </li> </ul> <p>After you submit a request to replace or return a network resource, the status of the network resource changes to <code>CREATING_SHIPPING_LABEL</code>. The shipping label is available when the status of the network resource is <code>PENDING_RETURN</code>. After the network resource is successfully returned, its status changes to <code>DELETED</code>. For more information, see <a href="https://docs.aws.amazon.com/private-networks/latest/userguide/radio-units.html#return-radio-unit">Return a radio unit</a>.</p>"),
		mcp.WithString("updateType", mcp.Required(), mcp.Description("Input parameter: <p>The update type.</p> <ul> <li> <p> <code>REPLACE</code> - Submits a request to replace a defective radio unit. We provide a shipping label that you can use for the return process and we ship a replacement radio unit to you.</p> </li> <li> <p> <code>RETURN</code> - Submits a request to return a radio unit that you no longer need. We provide a shipping label that you can use for the return process.</p> </li> <li> <p> <code>COMMITMENT</code> - Submits a request to change or renew the commitment period. If you choose this value, then you must set <a href=\"https://docs.aws.amazon.com/private-networks/latest/APIReference/API_StartNetworkResourceUpdate.html#privatenetworks-StartNetworkResourceUpdate-request-commitmentConfiguration\"> <code>commitmentConfiguration</code> </a>.</p> </li> </ul>")),
		mcp.WithObject("commitmentConfiguration", mcp.Description("Input parameter: <p>Determines the duration and renewal status of the commitment period for a radio unit.</p> <p>For pricing, see <a href=\"http://aws.amazon.com/private5g/pricing\">Amazon Web Services Private 5G Pricing</a>.</p>")),
		mcp.WithString("networkResourceArn", mcp.Required(), mcp.Description("Input parameter: The Amazon Resource Name (ARN) of the network resource.")),
		mcp.WithString("returnReason", mcp.Description("Input parameter: The reason for the return. Providing a reason for a return is optional.")),
		mcp.WithObject("shippingAddress", mcp.Description("Input parameter: Information about an address.")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    StartnetworkresourceupdateHandler(cfg),
	}
}
