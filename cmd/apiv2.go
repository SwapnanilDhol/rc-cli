package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"revenuecat-cli/api"
)

func initApiV2() {
	cmd := &cobra.Command{
		Use:     "api",
		Aliases: []string{"v2", "http"},
		Short:   "Call any RevenueCat public v2 REST endpoint",
		Long: `Send arbitrary HTTP requests to https://api.revenuecat.com/v2 using your configured API key.
This mirrors every operation in the official Developer API v2 reference:
https://www.revenuecat.com/docs/api-v2

Path substitution: write {project_id} or {{project_id}} in the path — it is replaced with
the project from rc config or --project-id. Other placeholders (e.g. customer_id) are left as-is; set them in the path string.

Examples:
  rc api GET '/projects/{project_id}/customers' -q limit=20
  rc api POST '/projects/{project_id}/paywalls' -d '{"name":"Example"}'
  rc api GET '/projects' -q limit=10`,
		Args: cobra.ExactArgs(2),
		RunE: runApiV2,
	}
	cmd.Flags().StringArrayP("query", "q", nil, "Query parameter key=value (repeatable)")
	cmd.Flags().StringP("data", "d", "", "Request body (JSON string)")
	cmd.Flags().String("data-file", "", "Read request body from file (use - for stdin)")
	cmd.Flags().Bool("pretty", true, "Pretty-print JSON response body when possible")
	cmd.Flags().Bool("substitute-project", true, "Replace {project_id} / {{project_id}} using default project")

	RootCmd.AddCommand(cmd)
}

func runApiV2(cmd *cobra.Command, args []string) error {
	method := strings.ToUpper(strings.TrimSpace(args[0]))
	path := strings.TrimSpace(args[1])
	switch method {
	case "GET", "POST", "DELETE", "PUT", "PATCH", "HEAD":
	default:
		return fmt.Errorf("unsupported HTTP method %q (use GET, POST, DELETE, PUT, PATCH)", method)
	}

	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}

	subProj, _ := cmd.Flags().GetBool("substitute-project")
	if subProj && cfg.ProjectID != "" {
		path = strings.ReplaceAll(path, "{{project_id}}", cfg.ProjectID)
		path = strings.ReplaceAll(path, "{project_id}", cfg.ProjectID)
	}
	if strings.Contains(path, "{project_id}") || strings.Contains(path, "{{project_id}}") {
		return fmt.Errorf("path still contains project_id placeholder; set project in rc config, use --project-id, or pass --substitute-project=false with a full path")
	}

	pairs, _ := cmd.Flags().GetStringArray("query")
	values := url.Values{}
	for _, p := range pairs {
		k, v, ok := strings.Cut(p, "=")
		if !ok {
			return fmt.Errorf("query %q must be key=value", p)
		}
		values.Add(k, v)
	}

	var body []byte
	dataStr, _ := cmd.Flags().GetString("data")
	dataFile, _ := cmd.Flags().GetString("data-file")
	if dataStr != "" && dataFile != "" {
		return fmt.Errorf("use only one of --data and --data-file")
	}
	if dataStr != "" {
		body = []byte(dataStr)
	}
	if dataFile != "" {
		if dataFile == "-" {
			body, err = io.ReadAll(os.Stdin)
		} else {
			body, err = os.ReadFile(dataFile)
		}
		if err != nil {
			return err
		}
	}

	pretty, _ := cmd.Flags().GetBool("pretty")
	status, raw, err := client.DoRaw(method, path, values, body)
	if err != nil {
		return err
	}

	out := raw
	if pretty && json.Valid(raw) {
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
