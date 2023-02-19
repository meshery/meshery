import VisibilityIcon from '@material-ui/icons/Visibility';
import VisibilityOffIcon from '@material-ui/icons/VisibilityOff';
import { Button } from "@material-ui/core";
import PropTypes from "prop-types";
import { iconMedium } from '../css/icons.styles';

function CatalogFilter({ catalogVisibility, handleCatalogVisibility, hideCatalog }) {
  return (
    <>
      { !hideCatalog &&  // In application we  don't have catalog, hence this check
        <Button onClick = {handleCatalogVisibility} variant="contained" color="primary" endIcon={catalogVisibility ? <VisibilityIcon style={iconMedium} /> : <VisibilityOffIcon style={iconMedium} />}>
         CATALOG
        </Button>
      }
    </>
  )
}

CatalogFilter.propTypes={
  catalogVisibility : PropTypes.bool.isRequired,
  handleCatalogVisibility : PropTypes.func.isRequired,
  hideCatalog : PropTypes.bool
}
export default CatalogFilter