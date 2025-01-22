# OCL Go SDK

A Go client library (SDK) for interacting with [Open Concept Lab (OCL) APIs](https://github.com/OpenConceptLab/oclapi2). This library aims to simplify the process of connecting and making requests to OCL’s API for terminology management (creating, retrieving and more)

## Installation
To start using the SDK in your Go project, run:

```bash
go get github.com/savannahghi/ocl
```
After that, import the package in your Go files:
```go
import "github.com/savannahghi/ocl"
```

## Getting Started
### Setting Up Environment Variables
Optionally, you can supply the base URL and token via environment variables. The following environment variables are recognized:
```
OCL_BASE_URL
OCL_TOKEN
```
This allows you to create a client with zero arguments:

```go
client, err := ocl.NewClientFromEnvVars()
if err != nil {
    // Handle error e.g. log or return
}
```

If you prefer to supply the values manually (or have them stored in code or a configuration file):

```go
client, err := ocl.NewClient("https://api.openconceptlab.org", "your-api-token")
if err != nil {
    // Handle error
}
```

## Usage Example

Below is a simple example showing how you might create a new Concept

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/savannahghi/ocl"
)

func main() {
    // Create a client from environment variables or explicitly
    client, err := ocl.NewClient("https://api.openconceptlab.org", "your-api-token")
    if err != nil {
        log.Fatalf("Failed to create OCL client: %v", err)
    }

    // Prepare the concept payload
    concept := &ocl.Concept{
        ID:           "malaria",
        ExternalID:   "MALARIA-ID-001",
        ConceptClass: "Diagnosis",
        Datatype:     "N/A",
        DisplayName:  "Malaria",
        Names: []ocl.Names{
            {
                Name:            "Malaria",
                Locale:          "en",
                LocalePreferred: true,
                NameType:        "Fully Specified",
            },
        },
        Descriptions: []ocl.Descriptions{
            {
                Description: "A dangerous disease caused by parasites transmitted through the bite of an infected mosquito.",
                Locale:      "en",
            },
        },
    }

    // Create concept
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    createdConcept, createErr := client.CreateConcept(ctx, concept)
    if createErr != nil {
        log.Fatalf("Failed to create concept: %v", createErr)
    }

    fmt.Printf("Created concept: %+v\n", createdConcept)
}
```

## How to Release
We use [release-please](https://github.com/googleapis/release-please) to manage our releases. This tool automates the process of creating a pull request with the next semantic version, updating the CHANGELOG, and tagging a release
****
