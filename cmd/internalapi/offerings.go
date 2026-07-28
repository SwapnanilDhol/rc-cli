package internalapi

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	rcinternal "revenuecat-cli/internal"
)

var offeringsCmd = &cobra.Command{
	Use:   "offerings",
	Short: "Manage offerings (Internal Dashboard API)",
	Long:  `Manage offerings using the internal dashboard API at https://app.revenuecat.com/internal/v1`,
}

var offeringsListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all offerings",
	RunE:    runOfferingsList,
}

var offeringsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an offering",
	RunE:  runOfferingsCreate,
}

var offeringsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get offering details",
	RunE:  runOfferingGet,
}

var offeringsUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update offering",
	Long: `PATCHes the offerings_with_packages endpoint, which can set display name, identifier,
metadata and packages in a single call.

METADATA IS REPLACED, NOT MERGED. --metadata overwrites the entire metadata object with
what you pass. To change one field without losing the rest, use --metadata-merge, which
GETs the current metadata, applies your keys on top, and writes the result back.
Pass a JSON null to delete a key under --metadata-merge.

Example — replace the whole metadata object:
  rc internal offerings update -o ofrngXXXX --metadata '{"title":"New Title"}'

Example — change only the title, preserving every other metadata field:
  rc internal offerings update -o ofrngXXXX --metadata-merge '{"title":"New Title"}'

Example — delete a metadata key:
  rc internal offerings update -o ofrngXXXX --metadata-merge '{"headerImageURL":null}'

Example — attach packages:
  rc internal offerings update -o ofrngXXXX \
    --name "Weekly, Yearly, Lifetime" \
    --packages '[{"identifier":"my.weekly","display_name":"Weekly","products":[{"product_id":"prodXXXX"}]}]'`,
	RunE: runOfferingsUpdate,
}

var offeringsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an offering",
	RunE:  runOfferingsDelete,
}

var offeringsDuplicateCmd = &cobra.Command{
	Use:   "duplicate",
	Short: "Duplicate an offering",
	RunE:  runOfferingsDuplicate,
}

var offeringsSetCurrentCmd = &cobra.Command{
	Use:   "set-current",
	Short: "Set an offering as the current/default",
	RunE:  runOfferingsSetCurrent,
}

var offeringsArchiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Archive an offering",
	RunE:  runOfferingsArchive,
}

func init() {
	offeringsCmd.AddCommand(offeringsListCmd, offeringsGetCmd, offeringsCreateCmd, offeringsUpdateCmd, offeringsDeleteCmd, offeringsDuplicateCmd, offeringsSetCurrentCmd, offeringsArchiveCmd)

	offeringsListCmd.Flags().String("platform", "", "Filter by platform (IOS, ANDROID, MACOS, WINDOWS, LINUX, tvOS)")

	offeringsCreateCmd.Flags().StringP("identifier", "i", "", "Offering identifier (slug)")
	offeringsCreateCmd.Flags().StringP("name", "n", "", "Display name")

	offeringsGetCmd.Flags().StringP("offering-id", "o", "", "Offering ID, identifier, or display name")

	offeringsUpdateCmd.Flags().StringP("offering-id", "o", "", "Offering ID, identifier, or display name (required)")
	offeringsUpdateCmd.Flags().StringP("name", "n", "", "Display name (display_name)")
	offeringsUpdateCmd.Flags().StringP("identifier", "i", "", "Offering identifier / slug")
	offeringsUpdateCmd.Flags().StringP("metadata", "m", "", `JSON object REPLACING all metadata, e.g. '{"tier":"pro"}'`)
	offeringsUpdateCmd.Flags().String("metadata-merge", "", `JSON object merged into existing metadata; a null value deletes that key`)
	offeringsUpdateCmd.Flags().String("packages", "", `JSON array of packages to attach, e.g. '[{"identifier":"my.weekly","display_name":"Weekly","products":[{"product_id":"prodXXXX"}]}]'`)

	offeringsDeleteCmd.Flags().StringP("offering-id", "o", "", "Offering ID, identifier, or display name")

	offeringsDuplicateCmd.Flags().StringP("offering-id", "o", "", "Offering to duplicate — ID, identifier, or display name (required)")
	offeringsDuplicateCmd.Flags().StringP("identifier", "i", "", "New offering identifier (required)")
	offeringsDuplicateCmd.Flags().StringP("name", "n", "", "New offering display name (required)")
	offeringsDuplicateCmd.Flags().Bool("packages-only", false, "Only duplicate packages")

	offeringsSetCurrentCmd.Flags().StringP("offering-id", "o", "", "Offering ID, identifier, or display name (required)")

	offeringsArchiveCmd.Flags().StringP("offering-id", "o", "", "Offering ID, identifier, or display name (required)")
}

// offeringRef reads -o/--offering-id and resolves it to an offering ID, accepting
// an ID, an identifier, or a display name.
func offeringRef(cmd *cobra.Command, c *Ctx) (string, error) {
	ref, _ := cmd.Flags().GetString("offering-id")
	if ref == "" {
		return "", fmt.Errorf("offering is required (--offering-id or -o); accepts an ID, identifier, or display name")
	}
	return ResolveOfferingID(c.Client, c.ProjectID, ref)
}

func runOfferingsList(cmd *cobra.Command, args []string) error {
	c, err := Dashboard("\n📦 Fetching offerings...")
	if err != nil {
		return err
	}

	params := map[string]string{}
	if platform, _ := cmd.Flags().GetString("platform"); platform != "" {
		params["platform"] = platform
	}
	resp, err := c.Client.GetWithParams(c.Path("/offerings"), params)
	if err != nil {
		return err
	}

	return c.Respond(resp, func() error {
		var offerings []rcinternal.Offering
		if err := json.Unmarshal(ToJSON(resp.Items), &offerings); err != nil {
			return fmt.Errorf("error parsing offerings: %w", err)
		}
		if len(offerings) == 0 {
			fmt.Println(YellowStyle.Render("No offerings found."))
			return nil
		}
		fmt.Println(InternalStyle.Render("\n📦 Offerings:\n"))
		for _, o := range offerings {
			status := GreenStyle.Render("Current")
			if !o.IsCurrent {
				status = "Inactive"
			}
			archLabel := ""
			if o.IsArchived {
				archLabel = " [ARCHIVED]"
			}
			fmt.Printf("  ID: %s\n", CyanStyle.Render(o.ID))
			fmt.Printf("  Identifier: %s\n", o.Identifier)
			fmt.Printf("  Name: %s\n", o.DisplayName)
			fmt.Printf("  Status: %s%s\n", status, archLabel)
			fmt.Printf("  Packages: %d\n", len(o.Packages))
			if o.Metadata != nil {
				fmt.Printf("  Metadata: %v\n", o.Metadata)
			}
			fmt.Println()
		}
		fmt.Println(GrayStyle.Render(fmt.Sprintf("Total: %d offerings", len(offerings))))
		return nil
	})
}

func runOfferingGet(cmd *cobra.Command, args []string) error {
	c, err := Dashboard("\n📦 Fetching offering...")
	if err != nil {
		return err
	}
	offeringID, err := offeringRef(cmd, c)
	if err != nil {
		return err
	}

	resp, err := c.Client.Get(c.Path("/offerings/%s", offeringID))
	if err != nil {
		return err
	}

	return c.Respond(resp, func() error {
		var offering rcinternal.Offering
		if err := json.Unmarshal(ToJSON(resp.Data), &offering); err != nil {
			return fmt.Errorf("error parsing offering: %w", err)
		}
		fmt.Println(InternalStyle.Render("\n📦 Offering Details:\n"))
		fmt.Printf("  ID: %s\n", CyanStyle.Render(offering.ID))
		fmt.Printf("  Identifier: %s\n", offering.Identifier)
		fmt.Printf("  Name: %s\n", offering.DisplayName)
		fmt.Printf("  Archived: %v\n", offering.IsArchived)
		fmt.Printf("  Current: %v\n", offering.IsCurrent)
		if offering.Metadata != nil {
			fmt.Printf("  Metadata: %v\n", offering.Metadata)
		}
		fmt.Printf("  Packages: %d\n", len(offering.Packages))
		if len(offering.Packages) > 0 {
			fmt.Println(InternalStyle.Render("\n  Packages:"))
			for _, pkg := range offering.Packages {
				fmt.Printf("    - %s (%s)\n", CyanStyle.Render(pkg.Identifier), pkg.DisplayName)
			}
		}
		return nil
	})
}

func runOfferingsCreate(cmd *cobra.Command, args []string) error {
	identifier, _ := cmd.Flags().GetString("identifier")
	name, _ := cmd.Flags().GetString("name")
	if identifier == "" {
		return fmt.Errorf("identifier is required (--identifier or -i)")
	}
	if name == "" {
		return fmt.Errorf("name is required (--name or -n)")
	}

	c, err := Dashboard("\n📦 Creating offering...")
	if err != nil {
		return err
	}

	resp, err := c.Client.Post(c.Path("/offerings"), map[string]string{
		"identifier":   identifier,
		"display_name": name,
	})
	if err != nil {
		return err
	}

	return c.Respond(resp, func() error {
		var offering rcinternal.Offering
		if err := json.Unmarshal(ToJSON(resp.Data), &offering); err != nil {
			return fmt.Errorf("error parsing response: %w", err)
		}
		fmt.Println(GreenStyle.Render("\n✓ Offering created:"))
		fmt.Printf("  ID: %s\n", CyanStyle.Render(offering.ID))
		fmt.Printf("  Identifier: %s\n", offering.Identifier)
		fmt.Printf("  Name: %s\n", offering.DisplayName)
		return nil
	})
}

func runOfferingsUpdate(cmd *cobra.Command, args []string) error {
	displayName, _ := cmd.Flags().GetString("name")
	identifier, _ := cmd.Flags().GetString("identifier")
	metadataStr, _ := cmd.Flags().GetString("metadata")
	metadataMergeStr, _ := cmd.Flags().GetString("metadata-merge")
	packagesStr, _ := cmd.Flags().GetString("packages")

	if displayName == "" && identifier == "" && metadataStr == "" && metadataMergeStr == "" && packagesStr == "" {
		return fmt.Errorf("at least one of --name (-n), --identifier (-i), --metadata (-m), --metadata-merge, or --packages is required")
	}
	if metadataStr != "" && metadataMergeStr != "" {
		return fmt.Errorf("use only one of --metadata (replace) and --metadata-merge (merge)")
	}

	var metadata, metadataPatch map[string]interface{}
	if metadataStr != "" {
		if err := json.Unmarshal([]byte(metadataStr), &metadata); err != nil {
			return fmt.Errorf("--metadata must be a JSON object: %w", err)
		}
	}
	if metadataMergeStr != "" {
		if err := json.Unmarshal([]byte(metadataMergeStr), &metadataPatch); err != nil {
			return fmt.Errorf("--metadata-merge must be a JSON object: %w", err)
		}
	}
	var packages []interface{}
	if packagesStr != "" {
		if err := json.Unmarshal([]byte(packagesStr), &packages); err != nil {
			return fmt.Errorf("--packages must be a JSON array: %w", err)
		}
	}

	c, err := Dashboard("\n📦 Saving offering (internal PATCH)...")
	if err != nil {
		return err
	}
	offeringID, err := offeringRef(cmd, c)
	if err != nil {
		return err
	}

	// --metadata-merge is read-modify-write: fetch current metadata so untouched keys survive.
	if metadataPatch != nil {
		current, err := c.offeringMetadata(offeringID)
		if err != nil {
			return err
		}
		metadata = mergeMetadata(current, metadataPatch)
	}

	body := map[string]interface{}{}
	if displayName != "" {
		body["display_name"] = displayName
	}
	if identifier != "" {
		body["identifier"] = identifier
	}
	if metadata != nil {
		body["metadata"] = metadata
	}
	if packages != nil {
		body["packages"] = packages
	}

	resp, err := c.Client.Patch(c.Path("/offerings_with_packages/%s", offeringID), body)
	if err != nil {
		return err
	}

	return c.Respond(resp, func() error {
		fmt.Println(GreenStyle.Render("\n✓ Offering updated via internal API."))
		fmt.Println(GrayStyle.Render("Open RevenueCat → Product catalog → Offerings, select this offering, and refresh if needed to see display name / metadata / packages."))
		return nil
	})
}

func runOfferingsDelete(cmd *cobra.Command, args []string) error {
	c, err := Dashboard("\n📦 Deleting offering...")
	if err != nil {
		return err
	}
	offeringID, err := offeringRef(cmd, c)
	if err != nil {
		return err
	}

	resp, err := c.Client.Delete(c.Path("/offerings/%s", offeringID))
	if err != nil {
		return err
	}

	return c.Respond(resp, func() error {
		fmt.Println(GreenStyle.Render("\n✓ Offering deleted: " + offeringID))
		return nil
	})
}

func runOfferingsDuplicate(cmd *cobra.Command, args []string) error {
	identifier, _ := cmd.Flags().GetString("identifier")
	name, _ := cmd.Flags().GetString("name")
	packagesOnly, _ := cmd.Flags().GetBool("packages-only")
	if identifier == "" {
		return fmt.Errorf("identifier is required (--identifier or -i)")
	}
	if name == "" {
		return fmt.Errorf("name is required (--name or -n)")
	}

	c, err := Dashboard("\n📦 Duplicating offering...")
	if err != nil {
		return err
	}
	offeringID, err := offeringRef(cmd, c)
	if err != nil {
		return err
	}

	resp, err := c.Client.Post(c.Path("/offerings/%s/duplicate", offeringID), map[string]interface{}{
		"identifier":    identifier,
		"display_name":  name,
		"packages_only": packagesOnly,
	})
	if err != nil {
		return err
	}

	return c.Respond(resp, func() error {
		var offering map[string]interface{}
		if err := json.Unmarshal(ToJSON(resp.Data), &offering); err != nil {
			return fmt.Errorf("error parsing response: %w", err)
		}
		fmt.Println(GreenStyle.Render("\n✓ Offering duplicated:"))
		fmt.Printf("  ID: %s\n", CyanStyle.Render(GetStringValue(offering, "id", "?")))
		fmt.Printf("  Identifier: %s\n", GetStringValue(offering, "identifier", "?"))
		fmt.Printf("  Name: %s\n", GetStringValue(offering, "display_name", "?"))
		return nil
	})
}

func runOfferingsSetCurrent(cmd *cobra.Command, args []string) error {
	c, err := Dashboard("\n📦 Setting offering as current...")
	if err != nil {
		return err
	}
	offeringID, err := offeringRef(cmd, c)
	if err != nil {
		return err
	}

	resp, err := c.Client.Patch(c.Path("/offerings/%s", offeringID), map[string]interface{}{"is_current": true})
	if err != nil {
		return err
	}

	return c.Respond(resp, func() error {
		fmt.Println(GreenStyle.Render("\n✓ Offering set as current"))
		return nil
	})
}

func runOfferingsArchive(cmd *cobra.Command, args []string) error {
	c, err := Dashboard("\n📦 Archiving offering...")
	if err != nil {
		return err
	}
	offeringID, err := offeringRef(cmd, c)
	if err != nil {
		return err
	}

	resp, err := c.Client.Post(c.Path("/offerings/%s/actions/archive", offeringID), map[string]interface{}{})
	if err != nil {
		return err
	}

	return c.Respond(resp, func() error {
		fmt.Println(GreenStyle.Render("\n✓ Offering archived"))
		return nil
	})
}

// mergeMetadata applies patch on top of current and returns the result. A nil value
// in patch deletes that key. The merge is one level deep: a nested object in patch
// replaces the whole nested object, it is not merged recursively — matching how the
// dashboard treats metadata as an opaque document.
//
// Always returns a non-nil map, so the caller sends `metadata: {}` rather than
// omitting the field, when every key has been deleted.
func mergeMetadata(current, patch map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{}, len(current)+len(patch))
	for k, v := range current {
		merged[k] = v
	}
	for k, v := range patch {
		if v == nil {
			delete(merged, k)
			continue
		}
		merged[k] = v
	}
	return merged
}

// offeringMetadata returns the offering's current metadata object, or an empty map
// if it has none. Used by --metadata-merge to avoid clobbering existing keys.
func (c *Ctx) offeringMetadata(offeringID string) (map[string]interface{}, error) {
	resp, err := c.Client.Get(c.Path("/offerings/%s", offeringID))
	if err == nil {
		err = CheckResponse(resp)
	}
	if err != nil {
		return nil, fmt.Errorf("could not read current metadata: %w", err)
	}
	obj, ok := resp.Data.(map[string]interface{})
	if !ok {
		return map[string]interface{}{}, nil
	}
	meta, ok := obj["metadata"].(map[string]interface{})
	if !ok {
		return map[string]interface{}{}, nil
	}
	return meta, nil
}
