package registration

// Anything that can be parsed into a packagingUnit is a RegisterableEntity in Meshery server
type RegisterableEntity interface {
   /*
		1. `err` - this is a breaking error, which signifies that the given entity is invalid
		2. Errors encountered while parsing items into meshmodel entites are stored in the registration logger
	*/
	PkgUnit() (packagingUnit, error)
}