import React, { useEffect, useMemo, useRef, useState } from 'react';
import { Grid2, Chip, Button, TextField, Tooltip, Avatar, styled } from '@sistent/sistent';
import { NoSsr } from '@sistent/sistent';
import ReactSelectWrapper from './ReactSelectWrapper';

import changeAdapterState from '@/graphql/mutations/AdapterStatusMutation';
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
import { updateProgress } from '@/store/slices/mesheryUi';
import { useDispatch, useSelector } from 'react-redux';
import { updateAdaptersInfo } from '@/store/slices/adapter';

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

type SelectOption = {
  value: string;
  label: string;
};

type MeshAdapter = {
  name: string;
  adapterLocation: string;
  uniqueID?: string;
  version?: string;
};

const MeshAdapterConfigComponent = () => {
  const labelRef = useRef<HTMLDivElement | null>(null);
  const { meshAdapters: globalAdapters, meshAdaptersts: meshAdapterStates } = useSelector(
    (state: any) => state.adapter,
  );
  const [meshAdapters, setMeshAdapters] = useState<MeshAdapter[]>(globalAdapters ?? []);
  const [ts, setTs] = useState<number>(meshAdapterStates);
  const [meshLocationURLError, setMeshLocationURLError] = useState(false);
  const [selectedAvailableAdapterError, setSelectedAvailableAdapterError] = useState(false);
  const [adapterStates, setAdapterStates] = useState<Record<string, string>>({});
  const [meshLocationURL, setMeshLocationURL] = useState<SelectOption | null>(null);
  const [meshDeployURL, setMeshDeployURL] = useState('');
  const [meshDeployURLError, setMeshDeployURLError] = useState(false);
  const [selectedAvailableAdapter, setSelectedAvailableAdapter] = useState<SelectOption | null>(
    null,
  );
  const { notify } = useNotification();

  const { data: adapters } = useGetAdaptersQuery();
  const { data: availableAdaptersData } = useGetAvailableAdaptersQuery();
  const [pingAdapter] = useLazyPingAdapterQuery();
  const [manageAdapter] = useManageAdapterMutation();
  const dispatch = useDispatch();

  const setAdapterURLs: SelectOption[] =
    adapters?.map((res: any) => ({
      value: String(res.adapterLocation ?? ''),
      label: String(res.adapterLocation ?? ''),
    })) || [];

  const availableAdapters: SelectOption[] =
    availableAdaptersData?.map((res: any) => ({
      value: String(res.adapterLocation ?? ''),
      label: String(res.name ?? ''),
    })) || [];

  const availableAdaptersByPort = useMemo(() => {
    const map = new Map<string, SelectOption>();
    for (const option of availableAdapters) {
      const raw = option.value;
      const port = raw.includes(':') ? raw.split(':')[1] : raw;
      if (port) {
        map.set(port, option);
      }
    }
    return map;
  }, [availableAdapters]);

  useEffect(() => {
    if (meshAdapterStates > ts) {
      setMeshAdapters(globalAdapters ?? []);
      setTs(meshAdapterStates);
    }
  }, [globalAdapters, meshAdapterStates, ts]);

  useEffect(() => {
    setAdapterStatesFunction();
  }, [meshAdapters]);

  const setAdapterStatesFunction = () => {
    const initialAdapterStates: Record<string, string> = {};

    meshAdapters.forEach((adapter: any) => {
      const label = adapter.name.toUpperCase();
      initialAdapterStates[label] = STATUS.UNDEPLOYED;
    });

    setAdapterStates(initialAdapterStates);
  };

  const getStatusColor = (status: string | undefined) => {
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

  const isValidPort = (value: string) => {
    if (!/^\d{1,5}$/.test(value)) {
      return false;
    }
    const port = Number(value);
    return port >= 1 && port <= 65535;
  };

  const parseAdapterLocation = (raw: string) => {
    const trimmed = raw.trim();
    const parts = trimmed.split(':');
    if (parts.length !== 2) {
      return null;
    }
    const [adapterName, port] = parts;
    if (!adapterName || !port || !isValidPort(port)) {
      return null;
    }
    return { adapterName, port };
  };

  const notifyValidationError = (message: string) => {
    notify({ message, event_type: EVENT_TYPES.ERROR });
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
  const handleMeshLocURLChange = (newValue: SelectOption | null) => {
    if (typeof newValue !== 'undefined') {
      setMeshLocationURL(newValue);
      setMeshLocationURLError(false);
    }
  };
  const handleDeployPortChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setMeshDeployURL(event.target.value);
    if (event.target.value) {
      setMeshDeployURLError(false);
    }
  };

  const handleAvailableAdapterChange = (newValue: SelectOption | null) => {
    if (typeof newValue !== 'undefined') {
      // Trigger label animation manually
      labelRef.current?.querySelector('label')?.classList.add('MuiInputLabel-shrink');
      setSelectedAvailableAdapter(newValue);
      setSelectedAvailableAdapterError(false);

      if (newValue !== null) {
        const derivedPort = newValue.value.includes(':')
          ? newValue.value.split(':')[1]
          : newValue.value;
        setMeshDeployURL(derivedPort);
        setMeshDeployURLError(false);
      }
    }
  };

  const handleSubmit = () => {
    if (!meshLocationURL?.value) {
      setMeshLocationURLError(true);
      notifyValidationError('Mesh Adapter URL is required.');
      return;
    }

    // Adapter management expects a location in the form host/name + port.
    // Example: meshery-istio:10000
    if (!parseAdapterLocation(meshLocationURL.value)) {
      setMeshLocationURLError(true);
      notifyValidationError(
        'Invalid Mesh Adapter URL. Expected format: adapter:port (e.g. meshery-istio:10000).',
      );
      return;
    }

    submitConfig();
  };

  const submitConfig = async () => {
    updateProgress({ showProgress: true });

    try {
      const result = await manageAdapter({
        meshLocationURL: meshLocationURL?.value,
      }).unwrap();

      updateProgress({ showProgress: false });

      if (result) {
        setMeshAdapters(result);
        setMeshLocationURL(null);
        notify({ message: 'Adapter configured.', event_type: EVENT_TYPES.SUCCESS });
        dispatch(updateAdaptersInfo({ meshAdapters: result }));
      } else {
        notify({
          message: 'Adapter was not configured. No response received from server.',
          event_type: EVENT_TYPES.ERROR,
        });
      }
    } catch (error) {
      handleError('Adapter was not configured due to an error')(error);
    }
  };

  const handleDelete = (adapterLoc: string) => async () => {
    updateProgress({ showProgress: true });

    try {
      const result = await manageAdapter({
        method: 'DELETE',
        adapter: adapterLoc,
      }).unwrap();

      updateProgress({ showProgress: false });

      if (result) {
        setMeshAdapters(result);
        notify({ message: 'Adapter removed.', event_type: EVENT_TYPES.SUCCESS });
        dispatch(updateAdaptersInfo({ meshAdapters: result }));
      } else {
        notify({
          message: 'Adapter was not removed. No response received from server.',
          event_type: EVENT_TYPES.ERROR,
        });
      }
    } catch (error) {
      handleError('Adapter was not removed due to an error')(error);
    }
  };

  const handleClick = (adapterLoc: string) => async () => {
    updateProgress({ showProgress: true });

    try {
      const result = await pingAdapter(adapterLoc).unwrap();

      updateProgress({ showProgress: false });

      if (result) {
        notify({ message: 'Adapter was pinged.', event_type: EVENT_TYPES.SUCCESS });
      } else {
        notify({ message: 'Ping returned no result.', event_type: EVENT_TYPES.WARNING });
      }
    } catch (error) {
      handleError('Error pinging adapter')(error);
    }
  };

  const handleAdapterDeploy = () => {
    if (
      !selectedAvailableAdapter ||
      !selectedAvailableAdapter.label ||
      selectedAvailableAdapter.label.trim() === ''
    ) {
      setSelectedAvailableAdapterError(true);
      notifyValidationError('Please select an available adapter.');
      return;
    }

    const adapterLabel = selectedAvailableAdapter.label.replace(/^meshery-/, '').toUpperCase();
    setAdapterStates((prevState) => ({
      ...prevState,
      [adapterLabel]: STATUS.DEPLOYING,
    }));

    if (!meshDeployURL || meshDeployURL.trim() === '') {
      setMeshDeployURLError(true);
      notifyValidationError('Port is required.');
      setAdapterStates((prevState) => ({
        ...prevState,
        [adapterLabel]: STATUS.UNDEPLOYED,
      }));
      return;
    }

    if (!isValidPort(meshDeployURL.trim())) {
      setMeshDeployURLError(true);
      notifyValidationError('Invalid port. Enter a number between 1 and 65535.');
      setAdapterStates((prevState) => ({
        ...prevState,
        [adapterLabel]: STATUS.UNDEPLOYED,
      }));
      return;
    }

    updateProgress({ showProgress: true });

    const variables = {
      status: 'ENABLED',
      adapter: selectedAvailableAdapter.label,
      targetPort: meshDeployURL.trim(),
    };

    changeAdapterState(
      (response: any, errors: any) => {
        updateProgress({ showProgress: false });

        if (errors) {
          setAdapterStates((prevState) => ({
            ...prevState,
            [adapterLabel]: STATUS.UNDEPLOYED,
          }));
          notify({ message: 'Unable to deploy adapter.', event_type: EVENT_TYPES.ERROR });
          return;
        }

        setAdapterStates((prevState) => ({
          ...prevState,
          [adapterLabel]: STATUS.DEPLOYED,
        }));
        notify({
          message: `Adapter ${(response?.adapterStatus ?? 'deployed').toString().toLowerCase()}`,
          event_type: EVENT_TYPES.SUCCESS,
        });
      },
      variables,
      (error) => {
        updateProgress({ showProgress: false });
        setAdapterStates((prevState) => ({
          ...prevState,
          [adapterLabel]: STATUS.UNDEPLOYED,
        }));
        notify({ message: 'Unable to deploy adapter.', event_type: EVENT_TYPES.ERROR });
      },
    );
  };

  const handleAdapterUndeploy = () => {
    if (!meshLocationURL?.value) {
      setMeshLocationURLError(true);
      notifyValidationError('Mesh Adapter URL is required.');
      return;
    }

    const parsed = parseAdapterLocation(meshLocationURL.value);
    if (!parsed) {
      setMeshLocationURLError(true);
      notifyValidationError(
        'Invalid Mesh Adapter URL. Expected format: adapter:port (e.g. meshery-istio:10000).',
      );
      return;
    }

    updateProgress({ showProgress: true });

    const targetPort = parsed.port;
    const adapterName = parsed.adapterName;

    const adapterLabel = (availableAdaptersByPort.get(targetPort)?.label || adapterName || '')
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

    changeAdapterState(
      (response: any, errors: any) => {
        updateProgress({ showProgress: false });

        if (errors) {
          setAdapterStates((prevState) => ({
            ...prevState,
            [adapterLabel]: STATUS.DEPLOYED,
          }));
          notify({ message: 'Unable to undeploy adapter.', event_type: EVENT_TYPES.ERROR });
          return;
        }

        notify({
          message: `Adapter ${(response?.adapterStatus ?? 'undeployed').toString().toLowerCase()}`,
          event_type: EVENT_TYPES.SUCCESS,
        });

        setAdapterStates((prevState) => ({
          ...prevState,
          [adapterLabel]: STATUS.UNDEPLOYED,
        }));
      },
      variables,
      (error) => {
        updateProgress({ showProgress: false });
        setAdapterStates((prevState) => ({
          ...prevState,
          [adapterLabel]: STATUS.DEPLOYED,
        }));
        notify({ message: 'Unable to undeploy adapter.', event_type: EVENT_TYPES.ERROR });
      },
    );
  };

  const handleError = (msg: string) => (error: any) => {
    updateProgress({ showProgress: false });
    notify({ message: msg, event_type: EVENT_TYPES.ERROR });
  };

  const configureTemplate = () => {
    let showAdapters: React.ReactNode = null;
    if (meshAdapters.length > 0) {
      showAdapters = (
        <AlignRight>
          {meshAdapters.map((adapter) => {
            let image = '/static/img/meshery-logo/meshery-logo.png';
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
                  label={adapter.adapterLocation}
                  onDelete={handleDelete(adapter.adapterLocation)}
                  onClick={handleClick(adapter.adapterLocation)}
                  icon={
                    // logoIcon
                    <BadgeAvatars
                      color={getStatusColor(adapterStates[adapter.name?.toUpperCase?.()])}
                    >
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

          <Grid2 container spacing={1} size="grow" sx={{ alignItems: 'flex-end' }}>
            <Grid2 data-cy="mesh-adapter-url" size={{ xs: 12 }}>
              <ReactSelectWrapper
                onChange={handleMeshLocURLChange}
                options={setAdapterURLs}
                value={meshLocationURL}
                // placeholder={'Mesh Adapter URL'}
                label="Mesh Adapter URL"
                error={meshLocationURLError}
                helperText={
                  meshLocationURLError
                    ? 'Expected format: adapter:port (e.g. meshery-istio:10000)'
                    : undefined
                }
              />
            </Grid2>
          </Grid2>
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
          <Grid2
            container
            spacing={1}
            style={{ marginTop: '50px' }}
            size="grow"
            sx={{ alignItems: 'flex-end' }}
          >
            <Grid2 size={{ xs: 12 }}>
              <ReactSelectWrapper
                onChange={handleAvailableAdapterChange}
                options={availableAdapters}
                value={selectedAvailableAdapter}
                // placeholder={'Mesh Adapter URL'}
                label="Available Adapters"
                data-testid="adapters-available-label"
                error={selectedAvailableAdapterError}
                helperText={
                  selectedAvailableAdapterError ? 'Select an adapter to deploy.' : undefined
                }
              />
            </Grid2>
          </Grid2>
          <Grid2
            container
            spacing={1}
            size="grow"
            sx={{ alignItems: 'flex-end', justifyContent: 'flex-end' }}
          >
            <div ref={labelRef}>
              <TextField
                id="deployPort"
                type="text"
                label="Enter Port"
                variant="standard"
                onChange={handleDeployPortChange}
                value={meshDeployURL}
                error={meshDeployURLError}
                helperText={meshDeployURLError ? 'Enter a port between 1 and 65535.' : undefined}
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
          </Grid2>
        </WrapperStyledDiv>
      </NoSsr>
    );
  };

  return configureTemplate();
};

export default MeshAdapterConfigComponent;
