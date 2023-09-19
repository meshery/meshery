import React, { useState, useEffect, useCallback, useRef } from 'react';
import { connect } from 'react-redux';
import NoSsr from '@material-ui/core/NoSsr';
import { withStyles, Button, Divider, MenuItem, TextField, Grid } from '@material-ui/core';
import { blue } from '@material-ui/core/colors';
import PropTypes from 'prop-types';
import { useRouter } from 'next/router';
import SettingsIcon from '@material-ui/icons/Settings';
import MesheryAdapterPlayComponent from './MesheryAdapterPlayComponent';
import { bindActionCreators } from 'redux';
import { setAdapter } from '../lib/store';

const styles = (theme) => ({
  icon: {
    fontSize: 23,
    width: theme.spacing(2.5),
    marginRight: theme.spacing(0.5),
    alignSelf: 'flex-start',
  },
  playRoot: {
    padding: theme.spacing(0),
    marginBottom: theme.spacing(2),
  },
  buttons: {
    display: 'flex',
    justifyContent: 'flex-end',
  },
  button: {
    marginTop: theme.spacing(3),
    marginLeft: theme.spacing(1),
  },
  margin: { margin: theme.spacing(1) },
  alreadyConfigured: {
    textAlign: 'center',
    padding: theme.spacing(20),
  },
  colorSwitchBase: {
    color: blue[300],
    '&$colorChecked': {
      color: blue[500],
      '& + $colorBar': { backgroundColor: blue[500] },
    },
  },
  colorBar: {},
  colorChecked: {},
  uploadButton: {
    margin: theme.spacing(1),
    marginTop: theme.spacing(3),
  },
  fileLabel: { width: '100%' },
  editorContainer: { width: '100%' },
  deleteLabel: { paddingRight: theme.spacing(2) },
  alignRight: { textAlign: 'right' },
  expTitleIcon: {
    width: theme.spacing(3),
    display: 'inline',
    verticalAlign: 'middle',
  },
  expIstioTitleIcon: {
    width: theme.spacing(2),
    display: 'inline',
    verticalAlign: 'middle',
    marginLeft: theme.spacing(0.5),
    marginRight: theme.spacing(0.5),
  },
  expTitle: {
    display: 'inline',
    verticalAlign: 'middle',
    marginLeft: theme.spacing(1),
  },
  paneSection: {
    backgroundColor: theme.palette.secondary.elevatedComponents,
    padding: theme.spacing(2.5),
    borderRadius: 4,
  },
});

function MesheryPlayComponent(props) {
  const { classes, meshAdapters, setAdapter } = props;
  const router = useRouter();
  const [adapter, setAdapterState] = useState(
    meshAdapters && meshAdapters.size > 0 ? meshAdapters[0] : {},
  );
  const prevMeshAdaptersRef = useRef();

  useEffect(() => {
    prevMeshAdaptersRef.current = meshAdapters;
  });

  const prevMeshAdapters = prevMeshAdaptersRef.current;

  const handleRouteChange = useCallback(() => {
    const queryParam = router?.query?.adapter;
    if (queryParam) {
      const selectedAdapter = meshAdapters.find(
        ({ adapter_location }) => adapter_location === queryParam,
      );
      if (selectedAdapter) {
        setAdapterState(selectedAdapter);
      }
    } else if (meshAdapters.size > 0) {
      setAdapterState(meshAdapters.get(0));
    }
  }, [router, meshAdapters]);

  useEffect(() => {
    router.events.on('routeChangeComplete', handleRouteChange);

    if (prevMeshAdapters?.size !== meshAdapters?.size && meshAdapters.size > 0) {
      handleRouteChange();
    }

    return () => {
      router.events.off('routeChangeComplete', handleRouteChange);
    };
  }, [router, meshAdapters]); // Dependencies for useEffect

  const handleConfigure = () => {
    router.push('/settings#service-mesh');
  };

  const pickImage = (adapter) => {
    let image = '/static/img/meshery-logo.png';
    let imageIcon = <img src={image} className={classes.expTitleIcon} />;
    if (adapter && adapter.name) {
      image = '/static/img/' + adapter.name.toLowerCase() + '.svg';
      imageIcon = <img src={image} className={classes.expTitleIcon} />;
    }
    return imageIcon;
  };

  const handleAdapterChange = () => {
    return (event) => {
      if (event.target.value !== '') {
        const selectedAdapter = meshAdapters.filter(
          ({ adapter_location }) => adapter_location === event.target.value,
        );
        if (selectedAdapter && selectedAdapter.size === 1) {
          setAdapterState(selectedAdapter.get(0));
          setAdapter({ selectedAdapter: selectedAdapter.get(0).name });
        }
      }
    };
  };

  const renderIndividualAdapter = () => {
    let adapCount = 0;
    let adapter;
    meshAdapters.forEach((adap) => {
      if (adap.adapter_location === props.adapter) {
        adapter = adap;
        meshAdapters.forEach((ad) => {
          if (ad.name == adap.name) adapCount += 1;
        });
      }
    });
    if (adapter) {
      const imageIcon = pickImage(adapter);
      return (
        <React.Fragment>
          <MesheryAdapterPlayComponent
            adapter={adapter}
            adapCount={adapCount}
            adapter_icon={imageIcon}
          />
        </React.Fragment>
      );
    }
    return '';
  };

  if (meshAdapters.size === 0) {
    return (
      <NoSsr>
        <React.Fragment>
          <div className={classes.alreadyConfigured}>
            <Button variant="contained" color="primary" size="large" onClick={handleConfigure}>
              <SettingsIcon className={classes.icon} />
              Configure Settings
            </Button>
          </div>
        </React.Fragment>
      </NoSsr>
    );
  }

  if (props.adapter && props.adapter !== '') {
    const indContent = renderIndividualAdapter();
    if (indContent !== '') {
      return indContent;
    } // else it will render all the available adapters
  }

  const imageIcon = pickImage(adapter);
  return (
    <NoSsr>
      <React.Fragment>
        <div className={classes.playRoot}>
          <Grid container>
            <Grid item xs={12} className={classes.paneSection}>
              <TextField
                select
                id="adapter_id"
                name="adapter_name"
                label="Select Service Mesh Type"
                fullWidth
                value={adapter && adapter.adapter_location ? adapter.adapter_location : ''}
                margin="normal"
                variant="outlined"
                onChange={handleAdapterChange()}
                SelectProps={{
                  MenuProps: {
                    anchorOrigin: {
                      vertical: 'bottom',
                      horizontal: 'left',
                    },
                    transformOrigin: {
                      vertical: 'top',
                      horizontal: 'left',
                    },
                    getContentAnchorEl: null,
                  },
                }}
              >
                {meshAdapters.map((ada) => (
                  <MenuItem
                    key={`${ada.adapter_location}_${new Date().getTime()}`}
                    value={ada.adapter_location}
                  >
                    {/* <ListItemIcon> */}
                    {pickImage(ada)}
                    {/* </ListItemIcon> */}
                    <span className={classes.expTitle}>{ada.adapter_location}</span>
                  </MenuItem>
                ))}
              </TextField>
            </Grid>
          </Grid>
        </div>
        <Divider variant="fullWidth" light />
        {adapter && adapter.adapter_location && (
          <MesheryAdapterPlayComponent adapter={adapter} adapter_icon={imageIcon} />
        )}
      </React.Fragment>
    </NoSsr>
  );
}

MesheryPlayComponent.propTypes = { classes: PropTypes.object.isRequired };

const mapDispatchToProps = (dispatch) => ({
  setAdapter: bindActionCreators(setAdapter, dispatch),
});

const mapStateToProps = (state) => {
  const k8sconfig = state.get('k8sConfig');
  const meshAdapters = state.get('meshAdapters');
  const meshAdaptersts = state.get('meshAdaptersts');
  const selectedAdapter = state.get('selectedAdapter');
  return { k8sconfig, meshAdapters, meshAdaptersts, selectedAdapter };
};

export default withStyles(styles)(
  connect(mapStateToProps, mapDispatchToProps)(MesheryPlayComponent),
);
