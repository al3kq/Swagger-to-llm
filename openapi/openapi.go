package openapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

// =====================================================
// Simplified API Spec Structures (Enhanced)
// =====================================================

// APIDocument is the root of our simplified API spec.
type APIDocument struct {
	Title       string      `json:"title" yaml:"title"`
	Version     string      `json:"version" yaml:"version"`
	Description string      `json:"description,omitempty" yaml:"description,omitempty"`
	Endpoints   []Endpoint  `json:"endpoints" yaml:"endpoints"`
	Servers     []string    `json:"servers" yaml:"servers"`
	Components  *Components `json:"components" yaml:"components"`
}

// Endpoint represents a simplified API endpoint.
type Endpoint struct {
	Path        string               `json:"path" yaml:"path"`
	Method      string               `json:"method" yaml:"method"`
	Summary     string               `json:"summary" yaml:"summary"`
	Description string               `json:"description,omitempty" yaml:"description,omitempty"`
	Parameters  []*Parameter         `json:"parameters" yaml:"parameters"`
	RequestBody *RequestBody         `json:"requestBody" yaml:"requestBody"`
	Responses   map[string]*Response `json:"responses" yaml:"responses"`
}

// Parameter represents a simplified parameter.
type Parameter struct {
	Name        string  `json:"name" yaml:"name"`
	In          string  `json:"in" yaml:"in"`
	Required    bool    `json:"required" yaml:"required"`
	// New field: capture the type directly if present.
	Type        string  `json:"type,omitempty" yaml:"type,omitempty"`
	Schema      *Schema `json:"schema" yaml:"schema"`
	Description string  `json:"description,omitempty" yaml:"description,omitempty"`
	Ref         string  `json:"$ref,omitempty" yaml:"$ref,omitempty"`
}

// RequestBody represents a simplified request body.
type RequestBody struct {
	Description string                `json:"description" yaml:"description"`
	Content     map[string]*MediaType `json:"content" yaml:"content"`
	Ref         string                `json:"$ref,omitempty" yaml:"$ref,omitempty"`
}

// Response represents a simplified response.
type Response struct {
	Description string                `json:"description" yaml:"description"`
	Content     map[string]*MediaType `json:"content" yaml:"content"`
	Ref         string                `json:"$ref,omitempty" yaml:"$ref,omitempty"`
}

// MediaType holds the media type object (only schema is used here).
type MediaType struct {
	Schema *Schema `json:"schema" yaml:"schema"`
}

// Schema represents a simplified schema.
type Schema struct {
	Type string `json:"type" yaml:"type"`
	Ref  string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
}

// Components holds reusable objects.
type Components struct {
	Schemas       map[string]*Schema      `json:"schemas" yaml:"schemas"`
	Parameters    map[string]*Parameter   `json:"parameters" yaml:"parameters"`
	RequestBodies map[string]*RequestBody `json:"requestBodies" yaml:"requestBodies"`
	Responses     map[string]*Response    `json:"responses" yaml:"responses"`
}

// =====================================================
// Swagger 2.0 Structures
// =====================================================

// SwaggerSpec represents a Swagger 2.0 specification.
type SwaggerSpec struct {
	Swagger  string              `yaml:"swagger" json:"swagger"`
	Info     SwaggerInfo         `yaml:"info" json:"info"`
	BasePath string              `yaml:"basePath" json:"basePath"`
	Paths    map[string]PathItem `yaml:"paths" json:"paths"`
	// Additional fields (host, schemes, definitions, etc.) can be added as needed.
}

// SwaggerInfo holds API info for Swagger.
type SwaggerInfo struct {
	Title       string `yaml:"title" json:"title"`
	Description string `yaml:"description" json:"description"`
	Version     string `yaml:"version" json:"version"`
}

// PathItem represents the available operations for a single path.
type PathItem struct {
	Get     *Operation `yaml:"get" json:"get"`
	Put     *Operation `yaml:"put" json:"put"`
	Post    *Operation `yaml:"post" json:"post"`
	Delete  *Operation `yaml:"delete" json:"delete"`
	Options *Operation `yaml:"options" json:"options"`
	Head    *Operation `yaml:"head" json:"head"`
	Patch   *Operation `yaml:"patch" json:"patch"`
}

// Operation represents a Swagger operation.
type Operation struct {
	Summary     string              `yaml:"summary" json:"summary"`
	Description string              `yaml:"description" json:"description"`
	OperationID string              `yaml:"operationId" json:"operationId"`
	Parameters  []Parameter         `yaml:"parameters" json:"parameters"`
	Responses   map[string]Response `yaml:"responses" json:"responses"`
}

// =====================================================
// Parsing and Conversion Functions
// =====================================================

// LoadAPISpec reads a YAML or JSON file and unmarshals it into an APIDocument.
// It supports both the simplified API spec format and Swagger 2.0.
func LoadAPISpec(path string) (*APIDocument, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		return &APIDocument{}, nil
	}

	// Unmarshal into a generic map to check for a "swagger" key.
	var raw map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	if _, isSwagger := raw["swagger"]; isSwagger {
		// Unmarshal into SwaggerSpec.
		var swaggerSpec SwaggerSpec
		if trimmed[0] == '{' {
			err = json.Unmarshal(data, &swaggerSpec)
		} else {
			err = yaml.Unmarshal(data, &swaggerSpec)
		}
		if err != nil {
			return nil, err
		}
		// Convert SwaggerSpec to APIDocument.
		doc := convertSwaggerToAPIDocument(swaggerSpec)
		return &doc, nil
	}

	// Otherwise, assume it's already in the simplified APIDocument format.
	var doc APIDocument
	if trimmed[0] == '{' {
		err = json.Unmarshal(data, &doc)
	} else {
		err = yaml.Unmarshal(data, &doc)
	}
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

// convertSwaggerToAPIDocument converts a SwaggerSpec into our simplified APIDocument.
func convertSwaggerToAPIDocument(sw SwaggerSpec) APIDocument {
	doc := APIDocument{
		Title:       sw.Info.Title,
		Version:     sw.Info.Version,
		Description: sw.Info.Description,
		Endpoints:   []Endpoint{},
		Servers:     []string{}, // Swagger 2.0 doesn't have a "servers" array.
	}

	for path, item := range sw.Paths {
		// For each HTTP method in the PathItem, create an Endpoint.
		if item.Get != nil {
			ep := createEndpointFromOperation(path, "GET", *item.Get)
			doc.Endpoints = append(doc.Endpoints, ep)
		}
		if item.Post != nil {
			ep := createEndpointFromOperation(path, "POST", *item.Post)
			doc.Endpoints = append(doc.Endpoints, ep)
		}
		if item.Put != nil {
			ep := createEndpointFromOperation(path, "PUT", *item.Put)
			doc.Endpoints = append(doc.Endpoints, ep)
		}
		if item.Delete != nil {
			ep := createEndpointFromOperation(path, "DELETE", *item.Delete)
			doc.Endpoints = append(doc.Endpoints, ep)
		}
		if item.Patch != nil {
			ep := createEndpointFromOperation(path, "PATCH", *item.Patch)
			doc.Endpoints = append(doc.Endpoints, ep)
		}
		if item.Head != nil {
			ep := createEndpointFromOperation(path, "HEAD", *item.Head)
			doc.Endpoints = append(doc.Endpoints, ep)
		}
		if item.Options != nil {
			ep := createEndpointFromOperation(path, "OPTIONS", *item.Options)
			doc.Endpoints = append(doc.Endpoints, ep)
		}
	}

	return doc
}

// createEndpointFromOperation creates an Endpoint from a given Operation.
func createEndpointFromOperation(path, method string, op Operation) Endpoint {
	return Endpoint{
		Path:        path,
		Method:      method,
		Summary:     op.Summary,
		Description: op.Description,
		Parameters:  convertParameters(op.Parameters),
		Responses:   convertResponses(op.Responses),
		// Swagger 2.0 does not have a separate RequestBody field (it uses parameters for body data).
	}
}

// convertParameters converts a slice of Parameter (from Swagger) to a slice of pointers to Parameter.
func convertParameters(params []Parameter) []*Parameter {
	var result []*Parameter
	for _, p := range params {
		paramCopy := p // create a copy so each pointer is unique
		result = append(result, &paramCopy)
	}
	return result
}

// convertResponses converts a map of Response (from Swagger) to a map of pointers to Response.
func convertResponses(responses map[string]Response) map[string]*Response {
	result := make(map[string]*Response)
	for code, r := range responses {
		respCopy := r
		result[code] = &respCopy
	}
	return result
}

func minifyText(text string) string {
	// Collapse all whitespace (including newlines) into a single space.
	return strings.Join(strings.Fields(text), " ")
}

// =====================================================
// Documentation Rendering (Enhanced)
// =====================================================

// RenderText produces LLM-readable documentation for the API.
func RenderText(doc *APIDocument) string {
	var sb strings.Builder

	// API Header
	sb.WriteString(fmt.Sprintf("API: %s (v%s)\n\n", doc.Title, doc.Version))
	sb.WriteString("DESCRIPTION:\n")
	if doc.Description != "" {
		sb.WriteString(doc.Description)
	} else {
		sb.WriteString("(None or your description here)")
	}
	sb.WriteString("\n\n")

	// Process each Endpoint.
	for _, ep := range doc.Endpoints {
		sb.WriteString(fmt.Sprintf("ENDPOINT: %s %s\n", strings.ToUpper(ep.Method), ep.Path))
		sb.WriteString(fmt.Sprintf("SUMMARY: %s\n", ep.Summary))
		// Truncate endpoint description if too long.
		desc := minifyText(ep.Description)
		if len(desc) > 20000 {
			desc = desc[:2000] + "..."
		}
		if desc == "" {
			sb.WriteString("DESCRIPTION: (None)\n")
		} else {
			sb.WriteString(fmt.Sprintf("DESCRIPTION: %s\n", desc))
		}

		// Parameters
		sb.WriteString("PARAMETERS:\n")
		if len(ep.Parameters) == 0 {
			sb.WriteString("  (None)\n")
		} else {
			for _, p := range ep.Parameters {
				// Use p.Type if present; otherwise, if a schema is provided, use that type.
				var pType string
				if p.Type != "" {
					pType = p.Type
				} else if p.Schema != nil {
					pType = p.Schema.Type
				} else {
					pType = "(unknown)"
				}
				sb.WriteString(fmt.Sprintf("  - %s (%s, %s, required=%t)", p.Name, pType, p.In, p.Required))
				if p.Description != "" {
					sb.WriteString(fmt.Sprintf(" : %s", p.Description))
				}
				sb.WriteString("\n")
			}
		}

		// Request Body
		sb.WriteString("REQUEST BODY: ")
		if ep.RequestBody != nil && ep.RequestBody.Description != "" {
			sb.WriteString(ep.RequestBody.Description)
		} else {
			sb.WriteString("None")
		}
		sb.WriteString("\n")

		// Responses
		sb.WriteString("RESPONSES:\n")
		if len(ep.Responses) == 0 {
			sb.WriteString("  (None)\n")
		} else {
			for code, resp := range ep.Responses {
				sb.WriteString(fmt.Sprintf("  - %s: %s\n", code, resp.Description))
			}
		}
		sb.WriteString("END\n")
	}
	return sb.String()
}

// =====================================================
// Existing Functions for Reference Resolution
// =====================================================

// ResolveReferences replaces $ref fields in the document with direct pointers to Components.
func ResolveReferences(doc *APIDocument) error {
	if doc.Components == nil {
		return nil
	}

	for i := range doc.Endpoints {
		ep := &doc.Endpoints[i]

		// Resolve parameters.
		for j, param := range ep.Parameters {
			if param == nil {
				continue
			}
			if param.Ref != "" {
				refName := extractNameFromRef(param.Ref, "parameters")
				if resolved, ok := doc.Components.Parameters[refName]; ok {
					ep.Parameters[j] = resolved
				} else {
					errMsg := fmt.Sprintf("unresolved parameter reference: %s", param.Ref)
					return fmt.Errorf(errMsg)
				}
			}
			if err := resolveSchema(&param.Schema, doc); err != nil {
				return err
			}
		}

		// Resolve requestBody.
		if ep.RequestBody != nil {
			if ep.RequestBody.Ref != "" {
				refName := extractNameFromRef(ep.RequestBody.Ref, "requestBodies")
				if resolved, ok := doc.Components.RequestBodies[refName]; ok {
					ep.RequestBody = resolved
				} else {
					errMsg := fmt.Sprintf("unresolved requestBody reference: %s", ep.RequestBody.Ref)
					return fmt.Errorf(errMsg)
				}
			}
			for _, mt := range ep.RequestBody.Content {
				if mt != nil && mt.Schema != nil {
					if err := resolveSchema(&mt.Schema, doc); err != nil {
						return err
					}
				}
			}
		}

		// Resolve responses.
		for code, resp := range ep.Responses {
			if resp == nil {
				continue
			}
			if resp.Ref != "" {
				refName := extractNameFromRef(resp.Ref, "responses")
				if resolved, ok := doc.Components.Responses[refName]; ok {
					ep.Responses[code] = resolved
				} else {
					errMsg := fmt.Sprintf("unresolved response reference: %s", resp.Ref)
					return fmt.Errorf(errMsg)
				}
			}
			for _, mt := range resp.Content {
				if mt != nil && mt.Schema != nil {
					if err := resolveSchema(&mt.Schema, doc); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

// resolveSchema replaces a Schema reference with a pointer to the component schema.
func resolveSchema(s **Schema, doc *APIDocument) error {
	if *s == nil {
		return nil
	}
	if (*s).Ref != "" {
		refName := extractNameFromRef((*s).Ref, "schemas")
		if resolved, ok := doc.Components.Schemas[refName]; ok {
			*s = resolved
		} else {
			errMsg := fmt.Sprintf("unresolved schema reference: %s", (*s).Ref)
			return fmt.Errorf(errMsg)
		}
	}
	return nil
}

// extractNameFromRef extracts the component name from a $ref string.
// E.g. "#/components/schemas/Pet" with componentType "schemas" returns "Pet".
func extractNameFromRef(ref, componentType string) string {
	prefix := "#/components/" + componentType + "/"
	return strings.TrimPrefix(ref, prefix)
}

// snippet is a helper function to safely print the first n bytes of a file.
func snippet(data []byte, n int) string {
	if len(data) <= n {
		return string(data)
	}
	return string(data[:n]) + "..."
}