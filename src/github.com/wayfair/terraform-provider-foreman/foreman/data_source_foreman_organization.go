package foreman

import (
	"fmt"

	"github.com/wayfair/terraform-provider-foreman/foreman/api"
	"github.com/wayfair/terraform-provider-utils/autodoc"
	"github.com/wayfair/terraform-provider-utils/helper"
	"github.com/wayfair/terraform-provider-utils/log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceForemanOrganization() *schema.Resource {
	// copy attributes from resource definition
	r := resourceForemanOrganization()
	ds := helper.DataSourceSchemaFromResourceSchema(r.Schema)

	// define searchable attributes for the data source
	ds["name"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		Description: fmt.Sprintf(
			"The name of the puppet branch, environment. "+
				"%s \"production\"",
			autodoc.MetaExample,
		),
	}

	return &schema.Resource{

		Read: dataSourceForemanOrganizationRead,

		// NOTE(ALL): See comments in the corresponding resource file
		Schema: ds,
	}
}

func dataSourceForemanOrganizationRead(d *schema.ResourceData, meta interface{}) error {
	log.Tracef("data_source_foreman_organization.go#Read")

	client := meta.(*api.Client)
	e := buildForemanOrganization(d)

	log.Debugf("ForemanOrganization: [%+v]", e)

	queryResponse, queryErr := client.QueryOrganization(e)
	if queryErr != nil {
		return queryErr
	}

	if queryResponse.Subtotal == 0 {
		return fmt.Errorf("Data source organization returned no results")
	} else if queryResponse.Subtotal > 1 {
		return fmt.Errorf("Data source organization returned more than 1 result")
	}

	var queryOrganization api.ForemanOrganization
	var ok bool
	if queryOrganization, ok = queryResponse.Results[0].(api.ForemanOrganization); !ok {
		return fmt.Errorf(
			"Data source results contain unexpected type. Expected "+
				"[api.ForemanOrganization], got [%T]",
			queryResponse.Results[0],
		)
	}
	e = &queryOrganization

	log.Debugf("ForemanOrganization: [%+v]", e)

	setResourceDataFromForemanOrganization(d, e)

	return nil
}
