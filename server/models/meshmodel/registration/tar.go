package registration

type Tar struct {
    path string
}

func (t Tar) PkgUnit() (packagingUnit, error){
	pkg := packagingUnit{}
	return pkg, nil
}