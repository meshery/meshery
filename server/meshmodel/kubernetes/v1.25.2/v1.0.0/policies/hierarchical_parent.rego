package meshmodel_policy
# https://play.openpolicyagent.org/p/I25U0udl0I

hierarchical_parent_relationship := [ result | 
  relationships := data.relationships
  # some res
  print("WWWEWEW")
  result := hierarchical_inventory_relationship with data.relationships as relationships
]

