package meshmodel_policy

import rego.v1

hierarchical_parent_relationship := [result |
	relationships := data.relationships
	result := hierarchical_inventory_relationship with data.relationships as relationships
]
