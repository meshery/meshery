package registration

type Tar struct {
    _ string
}

func (t Tar) PkgUnit() (packagingUnit, error){
	pkg := packagingUnit{}
	return pkg, nil
}