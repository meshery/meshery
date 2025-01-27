import PropTypes from 'prop-types';
import { iconMedium } from '../css/icons.styles';
import { OutlinedVisibilityOnIcon, OutlinedVisibilityOffIcon, Button } from '@layer5/sistent';

function CatalogFilter({ catalogVisibility, handleCatalogVisibility, hideCatalog, classes }) {
  return (
    <>
      {!hideCatalog && ( // In application we  don't have catalog, hence this check
        <Button
          style={{
            // marginBottom : '0.2rem',
            alignItems: 'center',
            marginLeft: '-0.6rem',
          }}
          size="large"
          onClick={handleCatalogVisibility}
          variant="contained"
          color="primary"
        >
          {catalogVisibility ? (
            <OutlinedVisibilityOnIcon style={iconMedium} />
          ) : (
            <OutlinedVisibilityOffIcon style={iconMedium} />
          )}
          <span className={classes.btnText} style={{ marginLeft: '4px' }}>
            {' '}
            Catalog
          </span>
        </Button>
      )}
    </>
  );
}

CatalogFilter.propTypes = {
  catalogVisibility: PropTypes.bool.isRequired,
  handleCatalogVisibility: PropTypes.func.isRequired,
  hideCatalog: PropTypes.bool,
};
export default CatalogFilter;
