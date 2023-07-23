package ucmdb

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/jp-git1986/terraform-provider-ucmdb/utils"
	rest "github.com/jp-git1986/ucmdb-sdk/rest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUcmdbList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUcmdbReadList,

		Schema: map[string]*schema.Schema{
			"filter": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"names": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"items": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ucmdb_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceUcmdbReadList(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*rest.Client)

	var diags diag.Diagnostics

	tql := rest.TopologyQuery{}
	var nodes []rest.Node

	filters := d.Get("filter").(*schema.Set)
	filter_list := filters.List()
	for _, filter := range filter_list {
		f := filter.(map[string]interface{})
		ci_type := f["type"].(string)
		names := f["names"].(interface{})

		nodes = append(nodes, rest.Node{
			Type:            ci_type,
			QueryIdentifier: ci_type,
			Visible:         true,
			IncludeSubtypes: true,
			Layout:          []string{"name"},
			AttributeConditions: []rest.AttributeConditions{
				{
					Attribute: "name",
					Operator:  "in",
					Value:     names,
				},
			},
		})
	}

	tql.Nodes = nodes

	//out, err := conn.Ucmdb.GetIdByNameType(tql)
	td, err := conn.ExecuteQuery(tql)
	if err != nil {
		return diag.FromErr(err)
	}
	b, err := json.MarshalIndent(td, "", "  ")
	if err != nil {
		utils.LogMe("ERROR", "provider|dataSourceUcmdbReadList()|ExecuteQuery()|error to marshall ", err)
	}
	utils.LogMe("DEBUG", "provider|dataSourceUcmdbReadList()|ExecuteQuery()|response body", string(b))

	// will be used to create id for data source
	var ucmdb_ids []string
	cis := td.CIS

	items := make([]interface{}, 0, len(cis))

	for _, ci := range cis {
		item := make(map[string]interface{})
		ucmdb_ids = append(ucmdb_ids, ci.UcmdbId)
		item["ucmdb_id"] = ci.UcmdbId
		item["type"] = ci.Type
		item["name"] = ci.Properties["name"].(string)
		items = append(items, item)
	}

	if err := d.Set("items", items); err != nil {
		return diag.FromErr(err)
	}
	id := GenerateIdByHash(ucmdb_ids)
	d.SetId(id)

	return diags
}

func GenerateIdByHash(ids []string) string {
	var id string
	if len(ids) > 0 {
		id = strings.Join(ids, "")
	} else {

	}
	id = fmt.Sprintf("%d", schema.HashString(id))
	return id
}
