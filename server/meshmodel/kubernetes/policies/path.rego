package path_builder

import future.keywords.in


ensureParentPathsExist(patches, object) = result {
	# Convert patches to a set
	paths := {p.path | p := patches[_]}

	# Compute all missing subpaths.
	#    Iterate over all paths and over all subpaths
	#    If subpath doesn't exist, add it to the set after making it a string
	missingPaths := {sprintf("/%s", [concat("/", prefixPath)]) |
		paths[path]
		pathArray := split(path, "/")
		pathArray[i] # walk over path
		i > 0 # skip initial element

		# array of all elements in path up to i
		prefixPath := [pathArray[j] | pathArray[j]; j <= i; j > 0] # j > 0: skip initial element
		walkPath := [toWalkElement(x) | x := prefixPath[_]]
		print("WALK PATH : ", walkPath)
		not inputPathExists(object, walkPath)
	}

	# Sort paths, to ensure they apply in correct order
	ordered_paths := sort(missingPaths)
	print("ORDERED PATHS", ordered_paths)

	# Return new patches prepended to original patches.
	#  Don't forget to prepend all paths with a /
	new_patches := [{"op": "add", "path": p, "value": value} |
		some i, p in ordered_paths
		print("CURRENT PATH: ", p)
		value := add_path(i, p, ordered_paths)
		print("VALUE", value)
	]

	result := array.concat(new_patches, patches)
}

# If next path exists and is a number, add an empty array.
add_path(currentPathIndex, currentPath, allPaths) := value {
	count(allPaths) > currentPathIndex + 1
	nextPath := allPaths[currentPathIndex + 1]
	print("NEXT PATH FOR 47: ", nextPath)
	re_match("[0-9]+$", nextPath)
	value = []
}

# If next path exists and is not a number, add an empty object.
add_path(currentPathIndex, currentPath, allPaths) := value {
	count(allPaths) > currentPathIndex + 1
	nextPath := allPaths[currentPathIndex + 1]
	print("NEXT PATH FOR 56: ", nextPath)
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
