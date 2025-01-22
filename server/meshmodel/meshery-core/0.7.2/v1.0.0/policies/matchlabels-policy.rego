package eval
import rego.v1


import data.core_utils
import data.core_utils.new_uuid
import data.core_utils.from_component_id
import data.core_utils.component_declaration_by_id
import data.core_utils.to_component_id
import data.core_utils.get_component_configuration
import data.core_utils.pop_first
import data.core_utils.object_get_nested
import data.core_utils.component_alias
import data.feasibility_evaluation_utils.is_relationship_feasible_to
import data.feasibility_evaluation_utils.is_relationship_feasible_from

#subtype matchlabels

is_matchlabel_relationship(relationship) if {
	lower(relationship.kind) == "hierarchical"
	lower(relationship.type) == "sibling" 
}




# returns set {
# 	{
# 		field,
# 		value,
# 		component
# 	}
# }
identify_matchlabels(design_file, relationship) := all_match_labels if {

	field_pairs := {pair |
		is_matchlabel_relationship(relationship)
		# identify against pair of components
		some component in design_file.components
		some other_component in design_file.components 
		component.id != other_component.id

		# print("is alias rel", component)
		is_relationship_feasible_from(component, relationship)
		is_relationship_feasible_to(other_component,relationship)
		


		some field,value in component.configuration.metadata.labels 
		some field2,value2 in other_component.configuration.metadata.labels 

        field == field2
		value == value2 

		print("matching fields",field,value)
        components := {component,other_component}

		pair := {
			"field": field,
			"value":value,
			"components":components
		}
	}

	all_match_labels := field_pairs 

}

is_match_labels_policy_identifier(relationship_policy_identifier) if {	
   relationship_policy_identifier == {
		"kind":"hierarchical",
		"type":"sibling",
		"subtype":"matchlabels"
	}
}

identify_relationships(design_file, relationships_in_scope,relationship_policy_identifier) := eval_results if {

    print("identifier",relationship_policy_identifier)
    is_match_labels_policy_identifier(relationship_policy_identifier)
	print("evaluatine matchlabels")


	
	eval_results := {new_relationship |
		some relationship in relationships_in_scope


     	is_matchlabel_relationship(relationship)

			
        identified_map := identify_matchlabels(design_file,relationship)

		

		# print("identified_map",identified_map)

		some match_label in identified_map 

		from_selectors := [from_selector | 
          some component in match_label.components
		  from_selector := {
		    "id" : component.id,
			"kind": component.component.kind,
			"model": component.model,
			"patch":{
		      "patchStrategy": "replace",
			  "mutatorRef": [[match_label.field]]
			}
		  }
		]

		selector_declaration := {
			"allow": {
				"from": from_selectors,
				"to": from_selectors,
			},
			"deny":{
				"from":[],
				"to":[]
			}
		}

		new_relationship := json.patch(relationship,[
			{
				"op":"replace",
				"path": "/selectors",
				"value": [selector_declaration]
			},
			{
				"op":"replace",
				"path":"/id",
				"value": new_uuid({selector_declaration,relationship})
			},
			{
				"op":"replace",
				"path":"/status",
				"value":"pending"
			}
		])

		# print("new_relationship",new_relationship )

        
	}

	print("Identify match rels Eval results", count(eval_results))

}


# validate all relationships in the design file
validate_relationships_phase(design_file,relationship_policy_identifier) := {validated |
    is_match_labels_policy_identifier(relationship_policy_identifier)
	some rel in design_file.relationships
	is_matchlabel_relationship(rel)
	validated := json.patch(rel,[{
		    "op":"replace",
			"path":"/status",
			"value":"deleted"	
	}])
}


action_phase(design_file,relationship_policy_identifier) := result if {
  
    is_match_labels_policy_identifier(relationship_policy_identifier)
	print("action phase matchlabels")
    result := { action | 
		  some rel in design_file.relationships
     	  is_matchlabel_relationship(rel)
		  rel.status == "pending"
		  rel_to_add := json.patch(rel,[{
		    "op":"replace",
			"path":"/status",
			"value":"approved"	
		  }])

		  action := {
			"op":"add_relationship",
			"value":rel_to_add
		}
	}
	|{ action |
	   some rel in design_file.relationships 
       is_matchlabel_relationship(rel)
	   rel.status == "deleted"
	   action := {
		"op":"delete_relationship",
		"value":rel
	   }
	}
}



	

#Notes 
# in each eval phase check the identifier matches 
# whenerever going through all relationships in desing always check if it matches the kind,type,subtype for the policy 
# because the relationships are never scoped down by default
