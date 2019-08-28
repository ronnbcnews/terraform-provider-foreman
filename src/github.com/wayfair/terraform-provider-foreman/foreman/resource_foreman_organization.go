package foreman

import (
	"fmt"
	"strconv"

	"github.com/wayfair/terraform-provider-foreman/foreman/api"
	"github.com/wayfair/terraform-provider-utils/autodoc"
	"github.com/wayfair/terraform-provider-utils/log"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceForemanOrganization() *schema.Resource {
	return &schema.Resource{

		Create: resourceForemanOrganizationCreate,
		Read:   resourceForemanOrganizationRead,
		Update: resourceForemanOrganizationUpdate,
		Delete: resourceForemanOrganizationDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{

			autodoc.MetaAttribute: &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
				Description: fmt.Sprintf(
					"%s A puppet organization, branch.",
					autodoc.MetaSummary,
				),
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: fmt.Sprintf(
					"Name of the organization. Usually maps to the name of "+
						"a puppet branch. "+
						"%s \"production\"",
					autodoc.MetaExample,
				),
			},
		},
	}
}

// -----------------------------------------------------------------------------
// Conversion Helpers
// -----------------------------------------------------------------------------

// buildForemanOrganization constructs a ForemanOrganization reference from a
// resource data reference.  The struct's  members are populated from the data
// populated in the resource data.  Missing members will be left to the zero
// value for that member's type.
func buildForemanOrganization(d *schema.ResourceData) *api.ForemanOrganization {
	log.Tracef("resource_foreman_organization.go#buildForemanOrganization")

	organization := api.ForemanOrganization{}

	obj := buildForemanObject(d)
	organization.ForemanObject = *obj

	return &organization
}

// setResourceDataFromForemanOrganization sets a ResourceData's attributes from
// the attributes of the supplied ForemanOrganization reference
func setResourceDataFromForemanOrganization(d *schema.ResourceData, fe *api.ForemanOrganization) {
	log.Tracef("resource_foreman_organization.go#setResourceDataFromForemanOrganization")

	d.SetId(strconv.Itoa(fe.Id))
	d.Set("name", fe.Name)
}

// -----------------------------------------------------------------------------
// Resource CRUD Operations
// -----------------------------------------------------------------------------

func resourceForemanOrganizationCreate(d *schema.ResourceData, meta interface{}) error {
	log.Tracef("resource_foreman_organization.go#Create")
	return nil
}

func resourceForemanOrganizationRead(d *schema.ResourceData, meta interface{}) error {
	log.Tracef("resource_foreman_organization.go#Read")

	client := meta.(*api.Client)
	e := buildForemanOrganization(d)

	log.Debugf("ForemanOrganization: [%+v]", e)

	readOrganization, readErr := client.ReadOrganization(e.Id)
	if readErr != nil {
		return readErr
	}

	log.Debugf("Read ForemanOrganization: [%+v]", readOrganization)

	setResourceDataFromForemanOrganization(d, readOrganization)

	return nil
}

func resourceForemanOrganizationUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Tracef("resource_foreman_organization.go#Update")
	return nil
}

func resourceForemanOrganizationDelete(d *schema.ResourceData, meta interface{}) error {
	log.Tracef("resource_foreman_organization.go#Delete")

	// NOTE(ALL): d.SetId("") is automatically called by terraform assuming delete
	//   returns no errors

	return nil
}
