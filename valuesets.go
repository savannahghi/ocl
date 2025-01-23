package ocl

import (
	"context"
	"net/http"
)

// The expected input is a JSON file with the following structure:
/*
{
    "resourceType": "ValueSet",
    "id": "Test1",
    "url": "https://ocl.org/ValueSet/test1",
    "version": "v5.0",
    "contact": [
        {
            "name": "Jon Doe 1",
            "telecom": [
                {
                    "system": "email",
                    "value": "jondoe1@gmail.com",
                    "use": "work",
                    "rank": 1,
                    "period": {
                        "start": "2020-10-29T10:26:15-04:00",
                        "end": "2025-10-29T10:26:15-04:00"
                    }
                }
            ]
        }
    ],
    "jurisdiction": [
        {
            "coding": [
                {
                    "system": "http://unstats.un.org/unsd/methods/m49/m49.htm",
                    "code": "USA",
                    "display": "United States of America"
                }
            ]
        }
    ],
    "compose": {
        "include": [
            {
                "system": "https://ocl.org/CodeSystem/test1",
                "version": "v2.0",
                "concept": [
                    {
                        "code": "AGYW_PREV"
                    },
                    {
                        "code": "CXCA_SCRN"
                    }
                ]
            }
        ]
    }
}
*/

// A ValueSet resource instance specifies a set of codes drawn from one or more code systems,
//  intended for use in a particular context.
//	Value sets link between CodeSystem definitions and their use in coded elements.
type Valueset struct {
	ResourceType string         `json:"resourceType,omitempty"`
	ID           string         `json:"id,omitempty"`
	URL          string         `json:"url,omitempty"`
	Version      string         `json:"version,omitempty"`
	Contact      []Contact      `json:"contact,omitempty"`
	Jurisdiction []Jurisdiction `json:"jurisdiction,omitempty"`
	Compose      Compose        `json:"compose,omitempty"`
}

type Contact struct {
	Name    string    `json:"name,omitempty"`
	Telecom []Telecom `json:"telecom,omitempty"`
}

type Telecom struct {
	System string `json:"system,omitempty"`
	Value  string `json:"value,omitempty"`
	Use    string `json:"use,omitempty"`
	Rank   int    `json:"rank,omitempty"`
	Period Period `json:"period,omitempty"`
}

type Period struct {
	Start string `json:"start,omitempty"`
	End   string `json:"end,omitempty"`
}

type Jurisdiction struct {
	Coding []Coding `json:"coding,omitempty"`
}

type Coding struct {
	System  string `json:"system,omitempty"`
	Code    string `json:"code,omitempty"`
	Display string `json:"display,omitempty"`
}

type Compose struct {
	Include []Include `json:"include,omitempty"`
}

type Include struct {
	System  string    `json:"system,omitempty"`
	Version string    `json:"version,omitempty"`
	Concept []Concept `json:"concept,omitempty"`
}

// Creates a new valueset.
func (c *Client) CreateValueset(ctx context.Context, valueset *Valueset) (*Valueset, error) {
	var resp Valueset
    var valuesetPath = getValuesetPath()

	err := c.makeRequest(ctx, http.MethodPost, valuesetPath, nil, valueset, resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
