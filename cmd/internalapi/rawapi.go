package internalapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var RawAPICmd = &cobra.Command{
	Use:   "api METHOD PATH",
	Short: "Call any internal dashboard API endpoint directly",
	Long: `Send an arbitrary HTTP request to https://app.revenuecat.com/internal/v1 using
your dashboard session (rc login). Use this for endpoints without a typed command.

Path substitution: {project_id} in the path is replaced with the project from
'rc internal projects use' or --project-id.

Examples:
  rc internal api GET  '/developers/me/projects/{project_id}/offerings'
  rc internal api GET  '/developers/me/projects/{project_id}/offerings/ofrngXXXX'
  rc internal api PATCH '/developers/me/projects/{project_id}/offerings_with_packages/ofrngXXXX' \
      -d '{"metadata":{"title":"Hello"}}'`,
	Args: cobra.ExactArgs(2),
	RunE: runRawAPI,
}

func init() {
	RawAPICmd.Flags().StringP("data", "d", "", "Request body (JSON string)")
	RawAPICmd.Flags().String("data-file", "", "Read request body from file (use - for stdin)")
	RawAPICmd.Flags().Bool("substitute-project", true, "Replace {project_id} using the default project")
}

func runRawAPI(cmd *cobra.Command, args []string) error {
	method := strings.ToUpper(strings.TrimSpace(args[0]))
	path := strings.TrimSpace(args[1])
	switch method {
	case "GET", "POST", "PUT", "PATCH", "DELETE":
	default:
		return fmt.Errorf("unsupported HTTP method %q (use GET, POST, PUT, PATCH, DELETE)", method)
	}

	if sub, _ := cmd.Flags().GetBool("substitute-project"); sub &&
		(strings.Contains(path, "{project_id}") || strings.Contains(path, "{{project_id}}")) {
		projectID, err := GetProjectID()
		if err != nil {
			return err
		}
		path = strings.ReplaceAll(path, "{{project_id}}", projectID)
		path = strings.ReplaceAll(path, "{project_id}", projectID)
	}

	dataStr, _ := cmd.Flags().GetString("data")
	dataFile, _ := cmd.Flags().GetString("data-file")
	if dataStr != "" && dataFile != "" {
		return fmt.Errorf("use only one of --data and --data-file")
	}
	var body []byte
	switch {
	case dataStr != "":
		body = []byte(dataStr)
	case dataFile == "-":
		var err error
		if body, err = io.ReadAll(os.Stdin); err != nil {
			return err
		}
	case dataFile != "":
		var err error
		if body, err = os.ReadFile(dataFile); err != nil {
			return err
		}
	}
	if len(body) > 0 && !json.Valid(body) {
		return fmt.Errorf("request body is not valid JSON")
	}

	client, err := GetInternalClient()
	if err != nil {
		return err
	}

	status, raw, err := client.DoRaw(method, path, body)
	if err != nil {
		return err
	}

	out := raw
	if json.Valid(raw) {
		var buf bytes.Buffer
		if err := json.Indent(&buf, raw, "", "  "); err == nil {
			out = buf.Bytes()
		}
	}
	os.Stdout.Write(out)
	if len(out) > 0 && out[len(out)-1] != '\n' {
		fmt.Println()
	}

	if status < 200 || status >= 300 {
		return fmt.Errorf("HTTP %d", status)
	}
	return nil
}
