package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wayfair/terraform-provider-utils/log"
)

const (
	OrganizationEndpointPrefix = "organizations"
)

// -----------------------------------------------------------------------------
// Struct Definition and Helpers
// -----------------------------------------------------------------------------

// The ForemanOrganization API model represents a puppet organization
type ForemanOrganization struct {
	// Inherits the base object's attributes
	ForemanObject
}

// -----------------------------------------------------------------------------
// CRUD Implementation
// -----------------------------------------------------------------------------

// CreateOrganization creates a new ForemanOrganization with the attributes of
// the supplied ForemanOrganization reference and returns the created
// ForemanOrganization reference.  The returned reference will have its ID and
// other API default values set by this function.
func (c *Client) CreateOrganization(e *ForemanOrganization) (*ForemanOrganization, error) {
	log.Tracef("foreman/api/organization.go#Create")

	reqEndpoint := fmt.Sprintf("/%s", OrganizationEndpointPrefix)

	organizationJSONBytes, jsonEncErr := json.Marshal(e)
	if jsonEncErr != nil {
		return nil, jsonEncErr
	}

	log.Debugf("organizationJSONBytes: [%s]", organizationJSONBytes)

	req, reqErr := c.NewRequest(
		http.MethodPost,
		reqEndpoint,
		bytes.NewBuffer(organizationJSONBytes),
	)
	if reqErr != nil {
		return nil, reqErr
	}

	var createdOrganization ForemanOrganization
	sendErr := c.SendAndParse(req, &createdOrganization)
	if sendErr != nil {
		return nil, sendErr
	}

	log.Debugf("createdOrganization: [%+v]", createdOrganization)

	return &createdOrganization, nil
}

// ReadOrganization reads the attributes of a ForemanOrganization identified by
// the supplied ID and returns a ForemanOrganization reference.
func (c *Client) ReadOrganization(id int) (*ForemanOrganization, error) {
	log.Tracef("foreman/api/organization.go#Read")

	reqEndpoint := fmt.Sprintf("/%s/%d", OrganizationEndpointPrefix, id)

	req, reqErr := c.NewRequest(
		http.MethodGet,
		reqEndpoint,
		nil,
	)
	if reqErr != nil {
		return nil, reqErr
	}

	var readOrganization ForemanOrganization
	sendErr := c.SendAndParse(req, &readOrganization)
	if sendErr != nil {
		return nil, sendErr
	}

	log.Debugf("readOrganization: [%+v]", readOrganization)

	return &readOrganization, nil
}

// UpdateOrganization updates a ForemanOrganization's attributes.  The
// organization with the ID of the supplied ForemanOrganization will be updated.
// A new ForemanOrganization reference is returned with the attributes from the
// result of the update operation.
func (c *Client) UpdateOrganization(e *ForemanOrganization) (*ForemanOrganization, error) {
	log.Tracef("foreman/api/organization.go#Update")

	reqEndpoint := fmt.Sprintf("/%s/%d", OrganizationEndpointPrefix, e.Id)

	organizationJSONBytes, jsonEncErr := json.Marshal(e)
	if jsonEncErr != nil {
		return nil, jsonEncErr
	}

	log.Debugf("organizationJSONBytes: [%s]", organizationJSONBytes)

	req, reqErr := c.NewRequest(
		http.MethodPut,
		reqEndpoint,
		bytes.NewBuffer(organizationJSONBytes),
	)
	if reqErr != nil {
		return nil, reqErr
	}

	var updatedOrganization ForemanOrganization
	sendErr := c.SendAndParse(req, &updatedOrganization)
	if sendErr != nil {
		return nil, sendErr
	}

	log.Debugf("updatedOrganization: [%+v]", updatedOrganization)

	return &updatedOrganization, nil
}

// DeleteOrganization deletes the ForemanOrganization identified by the supplied
// ID
func (c *Client) DeleteOrganization(id int) error {
	log.Tracef("foreman/api/organization.go#Delete")

	reqEndpoint := fmt.Sprintf("/%s/%d", OrganizationEndpointPrefix, id)

	req, reqErr := c.NewRequest(
		http.MethodDelete,
		reqEndpoint,
		nil,
	)
	if reqErr != nil {
		return reqErr
	}

	return c.SendAndParse(req, nil)
}

// -----------------------------------------------------------------------------
// Query Implementation
// -----------------------------------------------------------------------------

// QueryOrganization queries for a ForemanOrganization based on the attributes of
// the supplied ForemanOrganization reference and returns a QueryResponse struct
// containing query/response metadata and the matching organizations.
func (c *Client) QueryOrganization(e *ForemanOrganization) (QueryResponse, error) {
	log.Tracef("foreman/api/organization.go#Search")

	queryResponse := QueryResponse{}

	reqEndpoint := fmt.Sprintf("/%s", OrganizationEndpointPrefix)
	req, reqErr := c.NewRequest(
		http.MethodGet,
		reqEndpoint,
		nil,
	)
	if reqErr != nil {
		return queryResponse, reqErr
	}

	// dynamically build the query based on the attributes
	reqQuery := req.URL.Query()
	name := `"` + e.Name + `"`
	reqQuery.Set("search", "name="+name)

	req.URL.RawQuery = reqQuery.Encode()
	sendErr := c.SendAndParse(req, &queryResponse)
	if sendErr != nil {
		return queryResponse, sendErr
	}

	log.Debugf("queryResponse: [%+v]", queryResponse)

	// Results will be Unmarshaled into a []map[string]interface{}
	//
	// Encode back to JSON, then Unmarshal into []ForemanOrganization for
	// the results
	results := []ForemanOrganization{}
	resultsBytes, jsonEncErr := json.Marshal(queryResponse.Results)
	if jsonEncErr != nil {
		return queryResponse, jsonEncErr
	}
	jsonDecErr := json.Unmarshal(resultsBytes, &results)
	if jsonDecErr != nil {
		return queryResponse, jsonDecErr
	}
	// convert the search results from []ForemanOrganization to []interface
	// and set the search results on the query
	iArr := make([]interface{}, len(results))
	for idx, val := range results {
		iArr[idx] = val
	}
	queryResponse.Results = iArr

	return queryResponse, nil
}
