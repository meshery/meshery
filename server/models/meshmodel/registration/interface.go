package registration

// Anything that can be parsed into a packagingUnit is a RegisterableEntity in Meshery server
type RegisterableEntity interface {
   /*
	This function returns two errors:
		1. `err` - this is a breaking error, which signifies that the given entity is invalid
		2. `parsingErrors` - these are a list of errors that are encountered while parsing the objects inside the entity
	*/
	PkgUnit() (packagingUnit, error, []error)
}