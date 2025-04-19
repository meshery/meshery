import React, { useEffect, useState, useRef } from 'react';
import { Grid, Chip, Button, TextField, Tooltip, Avatar, styled } from '@layer5/sistent';
import { NoSsr } from '@layer5/sistent';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { withRouter } from 'next/router';
import ReactSelectWrapper from './ReactSelectWrapper';
import { updateAdaptersInfo, updateProgress } from '../lib/store';
import changeAdapterState from './graphql/mutations/AdapterStatusMutation';
import { useNotification } from '../utils/hooks/useNotification';
import { EVENT_TYPES } from '../lib/event-types';
import BadgeAvatars from './CustomAvatar';
import { keys } from '@/utils/permission_constants';
import CAN from '@/utils/can';
import { iconMedium } from 'css/icons.styles';
import {
  useGetAdaptersQuery,
  useGetAvailableAdaptersQuery,
  useLazyPingAdapterQuery,
  useManageAdapterMutation,
} from '@/rtk-query/system';

const WrapperStyledDiv = styled('div')(({ theme }) => ({
  padding: theme.spacing(5),
  backgroundColor: theme.palette.background.card,
  borderBottomLeftRadius: theme.spacing(1),
  borderBottomRightRadius: theme.spacing(1),
  marginTop: theme.spacing(2),
  marginBottom: theme.spacing(2),
}));

const AdapterButtons = styled('div')({
  display: 'flex',
  justifyContent: 'flex-end',
  paddingTop: '2rem',
});

const AdapterButton = styled(Button)(({ theme }) => ({
  marginLeft: theme.spacing(1),
}));

const AlignRight = styled('div')(({ theme }) => ({
  textAlign: 'right',
  marginBottom: theme.spacing(2),
}));

const AdapterChipStyled = styled(Chip)(({ theme }) => ({
  marginRight: theme.spacing(1),
  marginBottom: theme.spacing(1),
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
  const [ts, setTs] = useState(props.meshAdapterStates);
  const [meshLocationURLError, setMeshLocationURLError] = useState(false);
  const [selectedAvailableAdapterError, setSelectedAvailableAdapterError] = useState(false);
  const [adapterStates, setAdapterStates] = useState(props.meshAdapterStates);
  const [meshLocationURL, setMeshLocationURL] = useState();
  const [meshDeployURL, setMeshDeployURL] = useState('');
  const [meshDeployURLError, setMeshDeployURLError] = useState();
  const [selectedAvailableAdapter, setSelectedAvailableAdapter] = useState();
  const { notify } = useNotification();

  const { data: adapters } = useGetAdaptersQuery();
  const { data: availableAdaptersData } = useGetAvailableAdaptersQuery();
  const [pingAdapter] = useLazyPingAdapterQuery();
  const [manageAdapter] = useManageAdapterMutation();

  const setAdapterURLs =
    adapters?.map((res) => ({
      value: res.adapter_location,
      label: res.adapter_location,
    })) || [];

  const availableAdapters =
    availableAdaptersData?.map((res) => ({
      value: res.adapter_location,
      label: res.name,
    })) || [];

  useEffect(() => {
    if (props.meshAdapterStates > ts) {
      setMeshAdapters(props.meshAdapters);
      setTs(props.meshAdapterStates);
    }
  }, [props.meshAdapterStates, ts]);

  useEffect(() => {
    setAdapterStatesFunction();
  }, [meshAdapters]);

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
    if (typeof newValue !== 'undefined') {
      setMeshLocationURL(newValue);
      setMeshLocationURLError(false);
    }
  };
  const handleDeployPortChange = (newValue) => {
    if (typeof newValue !== 'undefined') {
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

  const submitConfig = async () => {
    props.updateProgress({ showProgress: true });

    try {
      const result = await manageAdapter({
        meshLocationURL: meshLocationURL.value,
      }).unwrap();

      props.updateProgress({ showProgress: false });

      if (result) {
        setMeshAdapters(result);
        setMeshLocationURL('');
        notify({ message: 'Adapter configured.', event_type: EVENT_TYPES.SUCCESS });
        props.updateAdaptersInfo({ meshAdapters: result });
      }
    } catch (error) {
      handleError('Adapter was not configured due to an error')(error);
    }
  };

  const handleDelete = (adapterLoc) => async () => {
    props.updateProgress({ showProgress: true });

    try {
      const result = await manageAdapter({
        method: 'DELETE',
        adapter: adapterLoc,
      }).unwrap();

      props.updateProgress({ showProgress: false });

      if (result) {
        setMeshAdapters(result);
        notify({ message: 'Adapter removed.', event_type: EVENT_TYPES.SUCCESS });
        props.updateAdaptersInfo({ meshAdapters: result });
      }
    } catch (error) {
      handleError('Adapter was not removed due to an error')(error);
    }
  };

  const handleClick = (adapterLoc) => async () => {
    props.updateProgress({ showProgress: true });

    try {
      const result = await pingAdapter(adapterLoc).unwrap();

      props.updateProgress({ showProgress: false });

      if (result) {
        notify({ message: 'Adapter was pinged.', event_type: EVENT_TYPES.SUCCESS });
      }
    } catch (error) {
      handleError('Error pinging adapter')(error);
    }
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
      setMeshDeployURLError(true);
      return;
    }

    props.updateProgress({ showProgress: true });

    const variables = {
      status: 'ENABLED',
      adapter: selectedAvailableAdapter.label,
      targetPort: meshDeployURL,
    };

    changeAdapterState((response, errors) => {
      props.updateProgress({ showProgress: false });

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

    props.updateProgress({ showProgress: true });

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
      props.updateProgress({ showProgress: false });

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

  const handleError = (msg) => (error) => {
    props.updateProgress({ showProgress: false });
    notify({ message: msg, event_type: EVENT_TYPES.ERROR, details: error.toString() });
  };

  const configureTemplate = () => {
    let showAdapters = '';
    if (meshAdapters.length > 0) {
      showAdapters = (
        <AlignRight>
          {meshAdapters.map((adapter) => {
            let image = '/static/img/meshery-logo.png';
            if (adapter.name) {
              image = '/static/img/' + adapter.name.toLowerCase() + '.svg';
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
                <AdapterChipStyled
                  label={adapter.adapter_location}
                  onDelete={handleDelete(adapter.adapter_location)}
                  onClick={handleClick(adapter.adapter_location)}
                  icon={
                    // logoIcon
                    <BadgeAvatars color={getStatusColor(adapterStates[adapter.name])}>
                      <Avatar alt={adapter.name} src={image} style={iconMedium} />
                    </BadgeAvatars>
                  }
                  variant="outlined"
                  data-cy="chipAdapterLocation"
                />
              </Tooltip>
            );
          })}
        </AlignRight>
      );
    }

    return (
      <NoSsr>
        <WrapperStyledDiv data-cy="mesh-adapter-connections">
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
            <AdapterButtons>
              <AdapterButton
                type="submit"
                variant="contained"
                data-testid="adapter-undeploy-button"
                color="primary"
                size="large"
                onClick={handleAdapterUndeploy}
                disabled={
                  !CAN(
                    keys.UNDEPLOY_CLOUD_NATIVE_INFRASTRUCTURE.action,
                    keys.UNDEPLOY_CLOUD_NATIVE_INFRASTRUCTURE.subject,
                  )
                }
              >
                Undeploy
              </AdapterButton>
              <AdapterButton
                type="submit"
                variant="contained"
                data-testid="adapter-connect-button"
                color="primary"
                size="large"
                onClick={handleSubmit}
                data-cy="btnSubmitMeshAdapter"
                disabled={!CAN(keys.CONNECT_ADAPTER.action, keys.CONNECT_ADAPTER.subject)}
              >
                Connect
              </AdapterButton>
            </AdapterButtons>
          </React.Fragment>
          <Grid container spacing={1} alignItems="flex-end" style={{ marginTop: '50px' }}>
            <Grid item xs={12}>
              <ReactSelectWrapper
                onChange={handleAvailableAdapterChange}
                options={availableAdapters}
                value={selectedAvailableAdapter}
                // placeholder={'Mesh Adapter URL'}
                label="Available Adapters"
                data-testid="adapters-available-label"
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
                variant="standard"
                onChange={(e) => handleDeployPortChange(e.target)}
                value={meshDeployURL}
                error={meshDeployURLError}
              />
            </div>
            <React.Fragment>
              <AdapterButtons>
                <AdapterButton
                  type="submit"
                  variant="contained"
                  data-testid="adapter-deploy-button"
                  color="primary"
                  size="large"
                  onClick={handleAdapterDeploy}
                  disabled={
                    !CAN(
                      keys.DEPLOY_CLOUD_NATIVE_INFRASTRUCTURE.action,
                      keys.DEPLOY_CLOUD_NATIVE_INFRASTRUCTURE.subject,
                    )
                  }
                >
                  Deploy
                </AdapterButton>
              </AdapterButtons>
            </React.Fragment>
          </Grid>
        </WrapperStyledDiv>
      </NoSsr>
    );
  };

  return configureTemplate();
};

const mapDispatchToProps = (dispatch) => ({
  updateAdaptersInfo: bindActionCreators(updateAdaptersInfo, dispatch),
  updateProgress: bindActionCreators(updateProgress, dispatch),
});

const mapStateToProps = (state) => {
  const meshAdapters = state.get('meshAdapters').toJS();
  const meshAdapterStates = state.get('meshAdaptersts');
  return { meshAdapters, meshAdapterStates };
};

export default connect(mapStateToProps, mapDispatchToProps)(withRouter(MeshAdapterConfigComponent));
