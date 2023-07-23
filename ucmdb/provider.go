package ucmdb

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	rest "github.com/jp-git1986/ucmdb-sdk/rest"
)

// this function returns a terraform ResourceProvider interface
func Provider() *schema.Provider {
	return &schema.Provider{
		// setting up shared configuration objects, e.g. addresses, secrets, access keys
		Schema: map[string]*schema.Schema{
			"target_env": {
				Type:         schema.TypeString,
				Description:  "A value which represents the UCMDB Target Environment. Valid values: CMS, OPSB, APM",
				Required:     true,
				ValidateFunc: StringInSlice([]string{"CMS", "OPSB", "APM"}, true),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"data_model_ci": resourceDataModelCi(),
			//"data_model":    resourceDataModel(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"ucmdb_list": dataSourceUcmdbList(),
		},

		// initialize shared configuration objects - the SDK client which makes API requests to OBM Downtime Service
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	i, err := readConfiguration(d.Get("target_env").(string))
	if err != nil {
		return nil, diag.FromErr(err)
	}

	conn, err := rest.NewClient(i["address"], i["user"], i["password"])
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return conn, diags
}

func StringInSlice(valid []string, ignoreCase bool) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (warnings []string, errors []error) {
		v, ok := i.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
			return warnings, errors
		}

		for _, str := range valid {
			if v == str || (ignoreCase && strings.ToLower(v) == strings.ToLower(str)) {
				return warnings, errors
			}
		}

		errors = append(errors, fmt.Errorf("expected %s to be one of %v, got %s", k, valid, v))
		return warnings, errors
	}
}

func readConfiguration(target_env string) (map[string]string, error) {
	address_env_name := fmt.Sprintf("UCMDB_%s_ADDRESS", strings.ToUpper(target_env))
	user_env_name := fmt.Sprintf("UCMDB_%s_API_USER", strings.ToUpper(target_env))
	password_env_name := fmt.Sprintf("UCMDB_%s_API_PASSWORD", strings.ToUpper(target_env))

	address, address_exists := os.LookupEnv(address_env_name)
	user, user_exists := os.LookupEnv(user_env_name)
	password, password_exists := os.LookupEnv(password_env_name)

	if !(address_exists && user_exists && password_exists) {
		return nil, fmt.Errorf("the following environment variables must be set: %s, %s, %s", address_env_name, user_env_name, password_env_name)
	}

	if strings.TrimSpace(address) == "" || strings.TrimSpace(user) == "" || strings.TrimSpace(password) == "" {
		return nil, fmt.Errorf("the following environment variables must not have empty values: %s, %s, %s", address_env_name, user_env_name, password_env_name)
	}

	config := make(map[string]string)
	config["address"] = address
	config["user"] = user
	config["password"] = password
	return config, nil
}
