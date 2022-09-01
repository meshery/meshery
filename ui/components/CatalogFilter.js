import VisibilityIcon from '@material-ui/icons/Visibility';
import VisibilityOffIcon from '@material-ui/icons/VisibilityOff';
import { Button } from "@material-ui/core";

function CatalogFilter({ catalogVisibility, handleCatalogVisibility, hideCatalog }) {
  return (
    <>
      { !hideCatalog &&  // In application we  don't have catalog, hence this check
        <Button onClick = {handleCatalogVisibility} variant="contained" color="primary" endIcon={catalogVisibility ? <VisibilityIcon /> : <VisibilityOffIcon />}>
         CATALOG
        </Button>
      }
    </>
  )
}

export default CatalogFilter