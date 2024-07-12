package registration

type RegisterableEntity interface {
	PkgUnit() (PackagingUnit, error)
}