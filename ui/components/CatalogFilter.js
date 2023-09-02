import VisibilityIcon from '@material-ui/icons/Visibility';
import VisibilityOffIcon from '@material-ui/icons/VisibilityOff';
import { Button } from "@material-ui/core";
import PropTypes from "prop-types";
import { iconMedium } from '../css/icons.styles';

function CatalogFilter({ catalogVisibility, handleCatalogVisibility, hideCatalog, classes }) {
  return (
    <>
      { !hideCatalog &&  // In application we  don't have catalog, hence this check
        <Button
          style={{
            marginBottom : '0.2rem',
            alignItems : 'center'
          }}
          size="large"
          onClick = {handleCatalogVisibility} variant="contained" color="primary">
          {catalogVisibility ? <VisibilityIcon style={iconMedium} /> : <VisibilityOffIcon style={iconMedium} />}
          <span className={classes.btnText} style={{ marginLeft : '4px' }}> Catalog</span>
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