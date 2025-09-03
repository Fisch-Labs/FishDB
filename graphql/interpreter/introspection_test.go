/*
 * FishDB
 *
 * // Copyright 2025 Fisch-labs
 */

package interpreter

import (
	"fmt"
	"testing"

	"github.com/Fisch-Labs/Toolkit/lang/graphql/parser"
)

func TestIntrospection(t *testing.T) {
	// Assuming songGraphGroups() is a helper that sets up the GraphQL schema.
	gm, _ := songGraphGroups()

	t.Run("Full introspection", func(t *testing.T) {
		// This is a standard, full introspection query.
		query := map[string]interface{}{
			"operationName": "IntrospectionQuery",
			"query": `
query IntrospectionQuery {
  __schema {
    queryType { name }
    mutationType { name }
    subscriptionType { name }
    types {
      ...FullType
    }
    directives {
      name
      description
      locations
      args {
        ...InputValue
      }
    }
  }
}

fragment FullType on __Type {
  kind
  name
  description
  fields(includeDeprecated: true) {
    name
    description
    args {
      ...InputValue
    }
    type {
      ...TypeRef
    }
    isDeprecated
    deprecationReason
  }
  inputFields {
    ...InputValue
  }
  interfaces {
    ...TypeRef
  }
  enumValues(includeDeprecated: true) {
    name
    description
    isDeprecated
    deprecationReason
  }
  possibleTypes {
    ...TypeRef
  }
}

fragment InputValue on __InputValue {
  name
  description
  type { ...TypeRef }
  defaultValue
}

fragment TypeRef on __Type {
  kind
  name
  ofType {
    kind
    name
    ofType {
      kind
      name
      ofType {
        kind
        name
        ofType {
          kind
          name
          ofType {
            kind
            name
            ofType {
              kind
              name
              ofType {
                kind
                name
              }
            }
          }
        }
      }
    }
  }
}
`,
		}

		res, err := runQuery("test", "main", query, gm, nil, false)

		// Check for initial query execution errors
		if err != nil {
			t.Errorf("runQuery failed: %v", err)
			return
		}

		data, ok := res["data"].(map[string]interface{})
		if !ok {
			t.Error("`data` key not found in response or is not a map")
			return
		}

		schema, ok := data["__schema"].(map[string]interface{})
		if !ok {
			t.Error("`__schema` key not found in data or is not a map")
			return
		}

		if _, ok := schema["types"]; !ok {
			t.Errorf("Unexpected result: `types` field not in schema. Error: %v", err)
			return
		}

		// Create runtime provider
		rtp := NewGraphQLRuntimeProvider("test", "main", gm,
			fmt.Sprint(query["operationName"]), make(map[string]interface{}), nil, true)

		// Parse the query and annotate the AST with runtime components
		ast, err := parser.ParseWithRuntime("test", fmt.Sprint(query["query"]), rtp)
		if err != nil {
			t.Errorf("ParseWithRuntime failed: %v", err)
			return
		}

		err = ast.Runtime.Validate()
		if err != nil {
			t.Errorf("AST validation failed: %v", err)
			return
		}

		// Evaluate the query
		// NOTE: This AST path is very specific and might break if the query changes.
		sr := ast.Children[0].Children[0].Children[0].Runtime.(*selectionSetRuntime)

		full := formatData(sr.ProcessFullIntrospection())
		filtered := formatData(sr.ProcessIntrospection())

		if full != filtered {
			t.Error("Full and filtered introspection are different")
			return
		}
	})

	// Now try out a reduced version
	t.Run("Reduced introspection", func(t *testing.T) {
		query := map[string]interface{}{
			"operationName": "IntrospectionQuery",
			"query": `
query IntrospectionQuery {
  __schema {
    queryType { name }
    mutationType { name }
    subscriptionType { name }
    directives {
      name
      description
      locations
      args {
        ...InputValue
        ...InputValue @skip(if: true)
        ... {
          name
        }
      }
      name1: name
    }
  }
}

fragment InputValue on __InputValue {
  name
  description
  type { ...TypeRef }
  defaultValue
}

fragment TypeRef on __Type {
  kind
  name
}
`,
		}

		res, err := runQuery("test", "main", query, gm, nil, false)

		expectedJSON := `{
  "data": {
    "__schema": {
      "directives": [
        {
          "args": [
            {
              "defaultValue": null,
              "description": "Skipped when true.",
              "name": "if",
              "type": {
                "kind": "NON_NULL",
                "name": null
              }
            }
          ],
          "description": "Directs the executor to skip this field or fragment when the ` + "`if`" + ` argument is true.",
          "locations": [
            "FIELD",
            "FRAGMENT_SPREAD",
            "INLINE_FRAGMENT"
          ],
          "name": "skip",
          "name1": "skip"
        },
        {
          "args": [
            {
              "defaultValue": null,
              "description": "Included when true.",
              "name": "if",
              "type": {
                "kind": "NON_NULL",
                "name": null
              }
            }
          ],
          "description": "Directs the executor to include this field or fragment only when the ` + "`if`" + ` argument is true.",
          "locations": [
            "FIELD",
            "FRAGMENT_SPREAD",
            "INLINE_FRAGMENT"
          ],
          "name": "include",
          "name1": "include"
        }
      ],
      "mutationType": {
        "name": "Mutation"
      },
      "queryType": {
        "name": "Query"
      },
      "subscriptionType": {
        "name": "Subscription"
      }
    }
  }
}`

		if formatData(res) != expectedJSON {
			t.Errorf("Unexpected result. Error: %v\nGot:\n%s\n\nExpected:\n%s", err, formatData(res), expectedJSON)
			return
		}
	})
}
