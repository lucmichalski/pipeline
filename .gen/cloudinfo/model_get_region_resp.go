/*
 * Product Info.
 *
 * The product info application uses the cloud provider APIs to asynchronously fetch and parse instance type attributes and prices, while storing the results in an in memory cache and making it available as structured data through a REST API.
 *
 * API version: 0.4.19
 * Contact: info@banzaicloud.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package cloudinfo

// GetRegionResp holds the detailed description of a specific region of a cloud provider
type GetRegionResp struct {
	Id    string   `json:"id,omitempty"`
	Name  string   `json:"name,omitempty"`
	Zones []string `json:"zones,omitempty"`
}
