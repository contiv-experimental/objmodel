/***
Copyright 2014 Cisco Systems Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"testing"

	log "github.com/Sirupsen/logrus"
)

// Simple test to parse json schema
func TestParseJsonSchema(t *testing.T) {
	inputStr := `{
		"objects": [
			{
				"name": "tenant",
				"type": "object",
				"key": "string",
				"properties": {
					"name": {
						"type": "string",
						"description": "Tenant Name"
					}
				},
				"link-sets": {
					"networks": {
						"ref": "network"
					},
					"apps": {
						"ref": "app"
					},
					"endpoint-groups": {
						"ref": "endpoint-group"
					}
				}
			},
			{
				"name": "network",
				"type": "object",
				"key": "string",
				"properties": {
					"name": {
						"type": "string"
					},
					"isPublic": {
						"type": "bool"
					},
					"isPrivate": {
						"type": "bool"
					},
					"encap": {
						"type": "string"
					},
					"subnet": {
						"type": "string"
					}
				},
				"links": {
					"tenant": {
						"ref": "tenant"
					}
				}
			}
		]

	}`

	schema, err := ParseSchema(inputStr)
	if err != nil {
		t.Fatalf("Error parsing json schema. Err: %v", err)
	}

	log.Printf("Parsed json schema: %+v", schema)

	goStr, err := schema.GenerateGoStructs()
	if err != nil {
		t.Fatalf("Error generating go code. Err: %v", err)
	}

	log.Printf("Generated go code: \n\n%s", goStr)
}

type TenantNetworksLinkSet struct {
	Type	string
	Key		string
	network	*Network
}

type TenantEndPointLinkSet struct {
	Type 	string
	Key		string
	// endpoint 	*Endpoint
}
type TenantLinkSets struct {
	Networks	[]TenantNetworksLinkSet
	Endpoints	[]TenantEndPointLinkSet
}

type Tenant struct {
	Key		string
	Name	string
	LinkSets	TenantLinkSets
}

type NetworkTenantLink struct {
	Type	string
	Key		string
	tenant	*Tenant
}

type NetworkLinks struct {
	Tenants 	[]NetworkTenantLink
}

type Network struct {
	Key		string
	Name 	string
	Reachability string
	Encap		string
	Subnet		string
	Links		NetworkLinks
}

// Sample json objects
/*
{
	tenants: {
		"default": {
			key: "default",
			name: "default",
			link-sets {
				networks: [
					"default:privateNet": {
						type: "network",
						key: "default:privateNet"
					},
					"default:publicNet": {
						type: "network",
						key: "default:publicNet"
					}
				]
			}
		}
	},
	networks: {
		"default/privateNet": {
			key: "default/privateNet",
			name: "privateNet",
			tenantName: "default",
			reachability: "private",
			encap: "vxlan",
			subnet: "20.1.1.0/24"
			links: {
				tenant: {
					type: "tenant",
					key: "default"
				}
			}
		}
	}
}
*/
