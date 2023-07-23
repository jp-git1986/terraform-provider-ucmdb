package ucmdb

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/jp-git1986/terraform-provider-ucmdb/utils"
	rest "github.com/jp-git1986/ucmdb-sdk-new/rest"
)

func resourceDataModelCi() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDataModelCiCreate,
		ReadContext:   resourceDataModelCiRead,
		UpdateContext: resourceDataModelCiUpdate,
		DeleteContext: resourceDataModelCiDelete,

		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				Description:  "CI type.",
				ValidateFunc: validation.StringIsNotEmpty,
				Required:     true,
			},
			"properties": {
				Type:        schema.TypeSet,
				Description: "CI map of properties.",
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceDataModelCiCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*rest.Client)

	ci_type := d.Get("type").(string)
	ci_prop_set := d.Get("properties").(*schema.Set)
	ci_props := (ci_prop_set.List())[0].(map[string]interface{})

	td := &rest.TopologyData{
		CIS: []rest.DataInConfigurationItem{
			{
				UcmdbId:    ci_type,
				Type:       ci_type,
				Properties: ci_props,
			},
		},
	}

	dmc, err := conn.Ucmdb.CreateDataModel(td)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(dmc.AddedCis) == 1 {
		d.SetId(dmc.AddedCis[0])
	} else if len(dmc.AddedCis) == 0 && len(dmc.UpdatedCis) == 1 {
		utils.LogMe("WARNING", "create resource ci", fmt.Sprintf("resource ci updated: %s", dmc.UpdatedCis[0]))
	} else if len(dmc.AddedCis) == 0 && len(dmc.IgnoredCis) == 1 {
		utils.LogMe("WARNING", "create resource ci", fmt.Sprintf("resource ci ignored: %s", dmc.IgnoredCis[0]))
	} else {
		utils.LogMe("WARNING", "create resource ci", fmt.Sprintf("resource ci updated: %s", dmc.UpdatedCis[0]))
		return diag.FromErr(fmt.Errorf("create resource ci failed %s", dmc))
	}

	return resourceDataModelCiRead(ctx, d, meta)
}

func resourceDataModelCiRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*rest.Client)

	var diags diag.Diagnostics

	ucmdbid := d.Id()
	dci, err := conn.Ucmdb.GetConfigurationItem(ucmdbid)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("type", dci.Type)

	props := make(map[string]interface{})
	props["name"] = dci.Properties["name"].(string)
	props["description"] = dci.Properties["description"].(string)
	if err := d.Set("properties", []interface{}{props}); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceDataModelCiUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*rest.Client)
	id := d.Id()

	if d.HasChanges("type", "properties") {

		ci_type := d.Get("type").(string)
		ci_prop_set := d.Get("properties").(*schema.Set)
		ci_props := (ci_prop_set.List())[0].(map[string]interface{})

		dic := rest.DataInConfigurationItem{
			UcmdbId:    id,
			Type:       ci_type,
			Properties: ci_props,
		}

		dmc, err := conn.Ucmdb.UpdateConfigurationItem(id, dic)
		utils.LogMe("DEBUG", "resourceCiUpdate|UpdateConfigurationItem()", dmc)
		if err != nil {
			utils.LogMe("ERROR", "resourceCiUpdate|UpdateConfigurationItem()", err)
			return diag.FromErr(err)
		}
		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceDataModelCiRead(ctx, d, meta)
}

func resourceDataModelCiDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*rest.Client)

	var diags diag.Diagnostics

	id := d.Id()
	dmc, err := conn.Ucmdb.DeleteConfigurationItem(id)
	_ = dmc
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func anyUpdate(d *schema.ResourceData) bool {
	status := false
	keys := []string{"type", "properties"}
	for i := range keys {
		status = status || d.HasChange(keys[i])
	}
	return status
}
