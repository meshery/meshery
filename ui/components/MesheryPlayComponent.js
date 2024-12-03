import React, { useState, useEffect, useCallback } from 'react';
import { connect } from 'react-redux';
import NoSsr from '@material-ui/core/NoSsr';
import {
  withStyles,
  Button,
  Divider,
  MenuItem,
  TextField,
  Grid,
  Typography,
} from '@material-ui/core';
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

const MesheryPlayComponent = (props) => {
  const router = useRouter();
  const { meshAdapters, classes } = props;
  const [adapter, setAdapterState] = useState({});

  const handleRouteChange = useCallback(() => {
    const queryParam = router.query.adapter;
    if (queryParam) {
      const selectedAdapter = meshAdapters.find(
        ({ adapter_location }) => adapter_location === queryParam,
      );
      if (selectedAdapter) {
        setAdapterState({ adapter: selectedAdapter });
      }
    } else if (meshAdapters.size > 0) {
      setAdapterState({ adapter: meshAdapters[0] });
    }
  }, [router.query.adapter, meshAdapters]);

  useEffect(() => {
    // mount
    router.events.on('routeChangeComplete', handleRouteChange);

    // unmount
    return () => {
      router.events.off('routeChangeComplete', handleRouteChange);
    };
  }, [router, handleRouteChange]);

  useEffect(() => {
    if (meshAdapters?.size > 0) {
      handleRouteChange();
    }
  }, [meshAdapters, handleRouteChange]);

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
          setAdapterState({ adapter: selectedAdapter[0] });
          setAdapter({ selectedAdapter: selectedAdapter[0].name });
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
        <>
          <MesheryAdapterPlayComponent
            adapter={adapter}
            adapCount={adapCount}
            adapter_icon={imageIcon}
          />
        </>
      );
    }
    return '';
  };

  if (meshAdapters.size === 0) {
    return (
      <NoSsr>
        <>
          <div
            style={{
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
              justifyContent: 'center',
              gap: '1rem',
              height: '100%',
              width: '100%',
            }}
          >
            <Typography variant="h6">Adapter Unavailable</Typography>
            <Typography variant="subtitle">Connect Meshery Adapter(s) in Settings</Typography>
            <Button variant="contained" color="primary" size="large" onClick={handleConfigure}>
              <SettingsIcon className={classes.icon} />
              Configure Settings
            </Button>
          </div>
        </>
      </NoSsr>
    );
  }

  if (adapter && adapter !== '') {
    const indContent = renderIndividualAdapter();
    if (indContent !== '') {
      return indContent;
    } // else it will render all the available adapters
  }

  const imageIcon = pickImage(adapter);
  return (
    <NoSsr>
      <>
        <div className={classes.playRoot}>
          <Grid container>
            <Grid item xs={12} className={classes.paneSection}>
              <TextField
                select
                id="adapter_id"
                name="adapter_name"
                label="Select Meshery Adapter"
                data-cy="lifecycle-service-mesh-type"
                fullWidth
                value={adapter && adapter.adapter_location ? adapter.adapter_location : ''}
                margin="normal"
                variant="outlined"
                onChange={handleAdapterChange}
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
      </>
    </NoSsr>
  );
};

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
