package registration


import(
	gcrv1 "github.com/google/go-containerregistry/pkg/v1"

)

type OCIImage struct {
    _ gcrv1.Image
}

func (o OCIImage) PkgUnit() (packagingUnit, error){
	pkg := packagingUnit{}
	return pkg, nil
}