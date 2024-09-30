package relationship_evaluation_policy

import rego.v1

ensure_parent_paths_exist(patches, obj) := result if {
	# Extract paths from the patches to a set of paths
	paths := {p.path | some p in patches}

	# Compute all missing subpaths to ensure correct patch.
	# Iterate all paths.
	# Add missing subpaths to the set as an array.

	missing_paths := {sprintf("/%s", [concat("/", prefix_path)]) |
		some path in paths

		# walk over path
		some i, _ in path

		# If a path is missing, all its subpaths will be added.
		# Eg: a/b/c: If path b is missing all its subpaths will be added.
		# array of all elements in path up to i
		prefix_path := [resultant_path |
			some j, val in path
			j <= i

			resultant_path := val
		]

		not input_path_exists(obj, prefix_path)
	}

	# Sort the paths.
	ordered_paths := sort(missing_paths)

	new_patches := [{"op": "add", "path": p, "value": value} |
		some i, p in ordered_paths
		value := add_path(i, p, ordered_paths)
	]

	result := array.concat(new_patches, patches)
}

# If next path exists and is a number, add an empty array.
add_path(current_path_index, current_path, all_paths) := value if {
	count(all_paths) > current_path_index + 1
	next_path := all_paths[current_path_index + 1]
	regex.match(`[0-9]+$`, next_path)
	value = []
}

# If next path exists and is not a number, add an empty object.
add_path(current_path_index, current_path, all_paths) := value if {
	count(all_paths) > current_path_index + 1
	next_path := all_paths[current_path_index + 1]
	not regex.match(`[0-9]+$`, next_path)
	value = {}
}

# Check that the given @path exists as part of the input object.
input_path_exists(check_object, path) if {
	walk(check_object, [path, v])
}

walk_element(str) := str if {
	not regex.match(`^[0-9]+$`, str)
}
