import React, { useEffect, useState, useRef } from 'react';
import PropTypes from 'prop-types';

import Grid from '@material-ui/core/Grid';
import { NoSsr, Chip, Button, TextField, Tooltip, Avatar, makeStyles } from '@material-ui/core';
import blue from '@material-ui/core/colors/blue';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { withRouter } from 'next/router';
import ReactSelectWrapper from './ReactSelectWrapper';
import { updateAdaptersInfo, updateProgress } from '../lib/store';
import dataFetch from '../lib/data-fetch';
import changeAdapterState from './graphql/mutations/AdapterStatusMutation';
import { useNotification } from '../utils/hooks/useNotification';
import { EVENT_TYPES } from '../lib/event-types';
import BadgeAvatars from './CustomAvatar';

const useStyles = makeStyles((theme) => ({
  wrapperClass: {
    padding: theme.spacing(5),
    backgroundColor: theme.palette.secondary.elevatedComponents,
    borderBottomLeftRadius: theme.spacing(1),
    borderBottomRightRadius: theme.spacing(1),
    marginTop: theme.spacing(2),
    marginBottom: theme.spacing(2),
  },
  buttons: {
    display: 'flex',
    justifyContent: 'flex-end',
    paddingTop: '2rem',
  },
  button: {
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
    '&$colorChecked + $colorBar': { backgroundColor: blue[500] },
  },
  colorBar: {},
  colorChecked: {},
  fileLabel: { width: '100%' },
  fileLabelText: {},
  inClusterLabel: { paddingRight: theme.spacing(2) },
  alignCenter: { textAlign: 'center' },
  alignRight: {
    textAlign: 'right',
    marginBottom: theme.spacing(2),
  },
  fileInputStyle: { opacity: '0.01' },
  // icon : { width : theme.spacing(2.5), },
  icon: {
    width: 20,
    height: 20,
  },
  istioIcon: { width: theme.spacing(1.5) },
  chip: {
    marginRight: theme.spacing(1),
    marginBottom: theme.spacing(1),
  },
}));
const STATUS = {
  DEPLOYED: 'DEPLOYED',
  UNDEPLOYED: 'UNDEPLOYED',
  DEPLOYING: 'DEPLOYING',
  UNDEPLOYING: 'UNDEPLOYING',
};

const MeshAdapterConfigComponent = (props) => {
  const labelRef = useRef(null);

  const [meshAdapters, setMeshAdapters] = useState(props.meshAdapters);
  const [setAdapterURLs, setSetAdapterURLs] = useState([]);
  const [availableAdapters, setAvailableAdapters] = useState([]);
  const [ts, setTs] = useState(props.meshAdaptersts);
  const [meshLocationURLError, setMeshLocationURLError] = useState(false);
  const [selectedAvailableAdapterError, setSelectedAvailableAdapterError] = useState(false);
  const [adapterStates, setAdapterStates] = useState(props.meshAdapterStates);
  const [meshLocationURL, setMeshLocationURL] = useState();
  const [meshDeployURL, setMeshDeployURL] = useState();
  const [meshDeployURLError, setMeshDeployURLError] = useState();
  const [selectedAvailableAdapter, setSelectedAvailableAdapter] = useState();
  const classes = useStyles();
  const { notify } = useNotification();

  useEffect(() => {
    if (props.meshAdaptersts > ts) {
      setMeshAdapters(props.meshAdapters);
      setTs(props.meshAdaptersts);
    }
  }, [props.meshAdaptersts, ts]);

  useEffect(() => {
    fetchSetAdapterURLs();
    fetchAvailableAdapters();
    setAdapterStatesFunction();
  }, []);

  const fetchSetAdapterURLs = () => {
    updateProgress({ showProgress: true });

    dataFetch(
      '/api/system/adapters',
      {
        method: 'GET',
        credentials: 'include',
      },
      (result) => {
        updateProgress({ showProgress: false });
        if (typeof result !== 'undefined') {
          const options = result.map((res) => ({
            value: res.adapter_location,
            label: res.adapter_location,
          }));
          setSetAdapterURLs(options);
        }
      },
      handleError('Unable to fetch available adapters'),
    );
  };

  const fetchAvailableAdapters = () => {
    updateProgress({ showProgress: true });

    dataFetch(
      '/api/system/availableAdapters',
      {
        method: 'GET',
        credentials: 'include',
      },
      (result) => {
        updateProgress({ showProgress: false });
        if (typeof result !== 'undefined') {
          const options = result.map((res) => ({
            value: res.adapter_location,
            label: res.name,
          }));
          setAvailableAdapters(options);
        }
      },
      handleError('Unable to fetch available adapters'),
    );
  };

  const setAdapterStatesFunction = () => {
    const initialAdapterStates = {};

    meshAdapters.forEach((adapter) => {
      const label = adapter.name.toUpperCase();
      initialAdapterStates[label] = STATUS.UNDEPLOYED;
    });

    setAdapterStates(initialAdapterStates);
  };

  const getStatusColor = (status) => {
    if (status === STATUS.DEPLOYED) {
      return '#00B39F';
    } else if (status === STATUS.UNDEPLOYED) {
      return '#808080';
    } else if (status === STATUS.DEPLOYING) {
      return '#EBC017';
    } else if (status === STATUS.UNDEPLOYING) {
      return '#E75225';
    }
  };

  // const handleChange = (name) => (event) => {
  //   if (name === "meshLocationURL" && event.target.value !== "") {
  //     setMeshLocationURLError(false);
  //   }

  //   setStateValues(prevState => ({
  //     ...prevState,
  //     [name] : event.target.value
  //   }));
  // };
  const handleMeshLocURLChange = (newValue) => {
    // console.log(newValue);
    // console.log(`action: ${actionMeta.action}`);
    // console.groupEnd();

    if (typeof newValue !== 'undefined') {
      setMeshLocationURL(newValue);
      setMeshLocationURLError(false);
    }
  };

  const handleDeployPortChange = (newValue) => {
    if (typeof newValue !== 'undefined') {
      console.log('port change to ' + newValue.value);
      setMeshDeployURL(newValue.value);
      setMeshDeployURLError(false);
    }
  };

  const handleAvailableAdapterChange = (newValue) => {
    if (typeof newValue !== 'undefined') {
      // Trigger label animation manually
      labelRef.current.querySelector('label').classList.add('MuiInputLabel-shrink');
      setSelectedAvailableAdapter(newValue);
      setSelectedAvailableAdapterError(false);

      if (newValue !== null) {
        setMeshDeployURL(newValue.value);
        setMeshDeployURLError(false);
      }
    }
  };

  const handleSubmit = () => {
    if (!meshLocationURL || !meshLocationURL.value || meshLocationURL.value === '') {
      setMeshLocationURLError(true);
      return;
    }

    submitConfig();
  };
  const submitConfig = () => {
    const data = { meshLocationURL: meshLocationURL.value };

    const params = Object.keys(data)
      .map((key) => `${encodeURIComponent(key)}=${encodeURIComponent(data[key])}`)
      .join('&');

    updateProgress({ showProgress: true });

    dataFetch(
      '/api/system/adapter/manage',
      {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8' },
        body: params,
      },
      (result) => {
        updateProgress({ showProgress: false });
        if (typeof result !== 'undefined') {
          // self.setState({ meshAdapters : result, meshLocationURL : "" });
          setMeshAdapters(result);
          setMeshLocationURL('');

          notify({ message: 'Adapter was configured!', event_type: EVENT_TYPES.SUCCESS });
          updateAdaptersInfo({ meshAdapters: result });
          fetchSetAdapterURLs();
        }
      },
      handleError('Adapter was not configured due to an error'),
    );
  };
  const handleDelete = (adapterLoc) => {
    updateProgress({ showProgress: true });

    dataFetch(
      `/api/system/adapter/manage?adapter=${encodeURIComponent(adapterLoc)}`,
      {
        method: 'DELETE',
        credentials: 'include',
      },
      (result) => {
        updateProgress({ showProgress: false });
        if (typeof result !== 'undefined') {
          setMeshAdapters(result);
          notify({ message: 'Adapter was removed!', event_type: EVENT_TYPES.SUCCESS });
          updateAdaptersInfo({ meshAdapters: result });
        }
      },
      handleError('Adapter was not removed due to an error'),
    );
  };

  const handleClick = (adapterLoc) => {
    updateProgress({ showProgress: true });

    dataFetch(
      `/api/system/adapters?adapter=${encodeURIComponent(adapterLoc)}`,
      {
        credentials: 'include',
      },
      (result) => {
        updateProgress({ showProgress: false });
        if (typeof result !== 'undefined') {
          notify({ message: 'Adapter was pinged!', event_type: EVENT_TYPES.SUCCESS });
        }
      },
      handleError('error'),
    );
  };

  const handleAdapterDeploy = () => {
    if (
      !selectedAvailableAdapter ||
      !selectedAvailableAdapter.value ||
      selectedAvailableAdapter.value === ''
    ) {
      setSelectedAvailableAdapterError(true);
      return;
    }

    const adapterLabel = selectedAvailableAdapter.label.replace(/^meshery-/, '').toUpperCase();
    setAdapterStates((prevState) => ({
      ...prevState,
      [adapterLabel]: STATUS.DEPLOYING,
    }));

    if (!meshDeployURL || meshDeployURL === '') {
      console.log(meshDeployURL);
      setMeshDeployURLError(true);
      return;
    }

    updateProgress({ showProgress: true });

    const variables = {
      status: 'ENABLED',
      adapter: selectedAvailableAdapter.label,
      targetPort: meshDeployURL,
    };

    changeAdapterState((response, errors) => {
      updateProgress({ showProgress: false });

      if (errors !== undefined) {
        handleError('Unable to Deploy adapter');
        setAdapterStates((prevState) => ({
          ...prevState,
          [adapterLabel]: STATUS.UNDEPLOYED,
        }));
      }

      setAdapterStates((prevState) => ({
        ...prevState,
        [adapterLabel]: STATUS.DEPLOYED,
      }));
      notify({
        message: 'Adapter ' + response.adapterStatus.toLowerCase(),
        event_type: EVENT_TYPES.SUCCESS,
      });
    }, variables);
  };

  const handleAdapterUndeploy = () => {
    if (!meshLocationURL || !meshLocationURL.value || meshLocationURL.value === '') {
      setMeshLocationURLError(true);
      return;
    }

    updateProgress({ showProgress: true });

    const targetPort = (() => {
      if (!meshLocationURL.value) {
        return null;
      }

      if (meshLocationURL.value.includes(':')) {
        return meshLocationURL.value.split(':')[1];
      }

      return meshLocationURL.value;
    })();

    const adapterName = (() => {
      if (!meshLocationURL.value) {
        return null;
      }

      if (meshLocationURL.value.includes(':')) {
        return meshLocationURL.value.split(':')[0];
      }

      return meshLocationURL.value;
    })();

    const adapterLabel = (
      availableAdapters.find((adapter) => adapter.value === targetPort)?.label || ''
    )
      .replace(/^meshery-/, '')
      .toUpperCase();

    setAdapterStates((prevState) => ({
      ...prevState,
      [adapterLabel]: STATUS.UNDEPLOYING,
    }));

    const variables = {
      status: 'DISABLED',
      adapter: adapterName,
      targetPort: targetPort,
    };

    changeAdapterState((response, errors) => {
      updateProgress({ showProgress: false });

      if (errors !== undefined) {
        console.error(errors);
        handleError('Unable to Deploy adapter');
        setAdapterStates((prevState) => ({
          ...prevState,
          [adapterLabel]: STATUS.DEPLOYED,
        }));
      }

      notify({
        message: 'Adapter ' + response.adapterStatus.toLowerCase(),
        event_type: EVENT_TYPES.SUCCESS,
      });

      setAdapterStates((prevState) => ({
        ...prevState,
        [adapterLabel]: STATUS.UNDEPLOYED,
      }));
    }, variables);
  };

  const handleError = (msg, error) => {
    updateProgress({ showProgress: false });
    notify({ message: msg, event_type: EVENT_TYPES.ERROR, details: error ? error.toString() : '' });
  };

  const configureTemplate = () => {
    let showAdapters = '';
    if (meshAdapters.length > 0) {
      showAdapters = (
        <div className={classes.alignRight}>
          {meshAdapters.map((adapter) => {
            let image = '/static/img/meshery-logo.png';
            // let logoIcon = <img src={image} className={classes.icon} />;
            if (adapter.name) {
              image = '/static/img/' + adapter.name.toLowerCase() + '.svg';
              // logoIcon = <img src={image} className={classes.icon} />;
            }

            return (
              <Tooltip
                key={adapter.uniqueID}
                title={`Meshery Adapter for
                        ${adapter.name
                          .toLowerCase()
                          .split(' ')
                          .map((s) => s.charAt(0).toUpperCase() + s.substring(1))
                          .join(' ')} (${adapter.version})`}
              >
                <Chip
                  className={classes.chip}
                  label={adapter.adapter_location}
                  onDelete={handleDelete(adapter.adapter_location)}
                  onClick={handleClick(adapter.adapter_location)}
                  icon={
                    // logoIcon
                    <BadgeAvatars color={getStatusColor(adapterStates[adapter.name])}>
                      <Avatar alt={adapter.name} src={image} className={classes.icon} />
                    </BadgeAvatars>
                  }
                  variant="outlined"
                  data-cy="chipAdapterLocation"
                />
              </Tooltip>
            );
          })}
        </div>
      );
    }

    return (
      <NoSsr>
        <div className={classes.wrapperClass} data-cy="mesh-adapter-connections">
          {showAdapters}

          <Grid container spacing={1} alignItems="flex-end">
            <Grid item xs={12} data-cy="mesh-adapter-url">
              <ReactSelectWrapper
                onChange={handleMeshLocURLChange}
                options={setAdapterURLs}
                value={meshLocationURL}
                // placeholder={'Mesh Adapter URL'}
                label="Mesh Adapter URL"
                error={meshLocationURLError}
              />
            </Grid>
          </Grid>
          <React.Fragment>
            <div className={classes.buttons}>
              <Button
                type="submit"
                variant="contained"
                color="primary"
                size="large"
                onClick={handleAdapterUndeploy}
                className={classes.button}
              >
                Undeploy
              </Button>
              <Button
                type="submit"
                variant="contained"
                color="primary"
                size="large"
                onClick={handleSubmit}
                className={classes.button}
                data-cy="btnSubmitMeshAdapter"
              >
                Connect
              </Button>
            </div>
          </React.Fragment>
          <Grid container spacing={1} alignItems="flex-end" style={{ marginTop: '50px' }}>
            <Grid item xs={12}>
              <ReactSelectWrapper
                onChange={handleAvailableAdapterChange}
                options={availableAdapters}
                value={selectedAvailableAdapter}
                // placeholder={'Mesh Adapter URL'}
                label="Available Mesh Adapter"
                error={selectedAvailableAdapterError}
              />
            </Grid>
          </Grid>
          <Grid container spacing={1} alignItems="flex-end" justifyContent="flex-end">
            <div ref={labelRef}>
              <TextField
                id="deployPort"
                type="text"
                label="Enter Port"
                onChange={(e) => handleDeployPortChange(e.target)}
                value={meshDeployURL}
                error={meshDeployURLError}
              />
            </div>
            <React.Fragment>
              <div className={classes.buttons}>
                <Button
                  type="submit"
                  variant="contained"
                  color="primary"
                  size="large"
                  onClick={handleAdapterDeploy}
                  className={classes.button}
                >
                  Deploy
                </Button>
              </div>
            </React.Fragment>
          </Grid>
        </div>
      </NoSsr>
    );
  };
  return configureTemplate();
};

MeshAdapterConfigComponent.propTypes = { classes: PropTypes.object.isRequired };

const mapDispatchToProps = (dispatch) => ({
  updateAdaptersInfo: bindActionCreators(updateAdaptersInfo, dispatch),
  updateProgress: bindActionCreators(updateProgress, dispatch),
});

const mapStateToProps = (state) => {
  const meshAdapters = state.get('meshAdapters').toJS();
  const meshAdaptersts = state.get('meshAdaptersts');

  return { meshAdapters, meshAdaptersts };
};

export default connect(mapStateToProps, mapDispatchToProps)(withRouter(MeshAdapterConfigComponent));
