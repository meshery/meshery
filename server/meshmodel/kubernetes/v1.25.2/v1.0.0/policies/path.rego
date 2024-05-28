package path_builder

import future.keywords.in

ensureParentPathsExist(patches, object) = result {
	# Extract paths from the patches to a set of paths
	paths := {p.path | p := patches[_]}

	# Compute all missing subpaths to ensure correct patch.
	# Iterate all paths.
	# Add missing subpaths to the set as an array.

	missingPaths := {sprintf("/%s", [concat("/", prefixPath)]) |
		paths[path]
		path[i] # walk over path


		# If a path is missing, all its subpaths will be added.
		# Eg: a/b/c: If path b is missing all its subpaths will be added.
		# array of all elements in path up to i
		prefixPath := [path[j] | path[j]; j <= i]
		walkPath := [toWalkElement(x) | x := prefixPath[_]]
		not inputPathExists(object, walkPath)
	}

	# Sort the paths.
	ordered_paths := sort(missingPaths)

	new_patches := [{"op": "add", "path": p, "value": value} |
		some i, p in ordered_paths
		value := add_path(i, p, ordered_paths)
	]

	result := array.concat(new_patches, patches)
}

# If next path exists and is a number, add an empty array.
add_path(currentPathIndex, currentPath, allPaths) := value {
	count(allPaths) > currentPathIndex + 1
	nextPath := allPaths[currentPathIndex + 1]
	re_match("[0-9]+$", nextPath)
	value = []
}

# If next path exists and is not a number, add an empty object.
add_path(currentPathIndex, currentPath, allPaths) := value {
	count(allPaths) > currentPathIndex + 1
	nextPath := allPaths[currentPathIndex + 1]
	not re_match("[0-9]+$", nextPath)
	value = {}
}

# Check that the given @path exists as part of the input object.
inputPathExists(object, path) {
	walk(object, [path, _])
}

toWalkElement(str) = str {
	not re_match("^[0-9]+$", str)
}
