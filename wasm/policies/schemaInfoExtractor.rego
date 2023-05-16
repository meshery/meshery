package schemaInfoExtractor

import input as schema

# list of all resources of a given type. given type must be defined in the resource_types variable above
resources = { resource_type: resources |
	some resource_type
#     print(resource_type)
	data.resource_types[resource_type]
    resources := { resource | 
    	walk_resources[resource]
        resource.type == resource_type
    }
}

walk_resources[resource] {	
    [path, value] := walk(schema)

    # Attempt to iterate over "resources" of the value, if the key doesn't
    # exist its OK, this iteration for walk will be undefined, and excluded
    # from the results.
    resource := {"type": path[count(path) - 1], "path": path, "value": getRefinedSchema(value)}
    # check if the resource type was contained in the set of desired resource types
#     print(resource.type)
	data.resource_types[resource.type]
}

getRefinedSchema(selectedschema) = refinedSchema {
	selectedschema.type == "array"
    refinedSchema := selectedschema.items
}

getRefinedSchema(selectedschema) = refinedSchema {
	not selectedschema.type == "array"
    refinedSchema := selectedschema
}
