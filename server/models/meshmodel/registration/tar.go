package registration

type Tar struct {
    path string
}

func (t Tar) PkgUnit() (PackagingUnit, error){
	pkg := PackagingUnit{}
	return pkg, nil
}