import {
  Button,
  Typography,
  FormGroup,
  TextField,
  InputAdornment,
  Tooltip,
  makeStyles,
  Checkbox,
  Grid,
  Divider,
} from '@material-ui/core';
import React, { useEffect, useState } from 'react';
import { useRef } from 'react';
import AddIconCircleBorder from '../assets/icons/AddIconCircleBorder';
import PromptComponent from './PromptComponent';
import CloudUploadIcon from '@material-ui/icons/CloudUpload';
import { promisifiedDataFetch } from '../lib/data-fetch';
import { updateProgress } from '../lib/store';
import { useNotification } from '../utils/hooks/useNotification';
import { EVENT_TYPES } from '../lib/event-types';
import { CONNECTION_STATES } from '../utils/Enum';
import { ConnectionChip } from './ConnectionChip';
import useKubernetesHook from './hooks/useKubernetesHook';
import { K8sEmptyState } from './EmptyState/K8sContextEmptyState';

// const _contexts = {
//   inserted_contexts: [
//     {
//       id: '1',
//       name: 'ctx-1',
//       auth: {
//         name: 'admin@k3d-k3s-default',
//         user: {},
//       },
//       cluster: {
//         cluster: {
//           'insecure-skip-tls-verify': true,
//           server: 'https://ctx-1-halibut-59567v944vrh547-8000.app.github.dev',
//         },
//         name: 'k3d-k3s-default',
//       },
//       server: 'https://ctx-1-halibut-59567v944vrh547-8000.app.github.dev',
//       meshery_instance_id: 'daa10439-7204-43d5-a425-a1bda9d25add',
//       kubernetes_server_id: '53f76c0e-2679-4b9d-a306-123d37a33c4c',
//       version: 'v1.27.4+k3s1',
//     },
//     {
//       id: '2',
//       name: 'ctx-2',
//       auth: {
//         name: 'admin@k3d-k3s-default',
//         user: {},
//       },
//       cluster: {
//         cluster: {
//           'insecure-skip-tls-verify': true,
//           server: 'https://ctx-2-halibut-59567v944vrh547-8000.app.github.dev',
//         },
//         name: 'k3d-k3s-default',
//       },
//       server: 'https://ctx-2-halibut-59567v944vrh547-8000.app.github.dev',
//       meshery_instance_id: 'daa10439-7204-43d5-a425-a1bda9d25add',
//       kubernetes_server_id: '53f76c0e-2679-4b9d-a306-123d37a33c4c',
//       version: 'v1.27.4+k3s1',
//     },
//     {
//       id: '3',
//       name: 'ctx-3',
//       auth: {
//         name: 'admin@k3d-k3s-default',
//         user: {},
//       },
//       cluster: {
//         cluster: {
//           'insecure-skip-tls-verify': true,
//           server: 'https://ctx-3-halibut-59567v944vrh547-8000.app.github.dev',
//         },
//         name: 'k3d-k3s-default',
//       },
//       server: 'https://ctx-3-halibut-59567v944vrh547-8000.app.github.dev',
//       meshery_instance_id: 'daa10439-7204-43d5-a425-a1bda9d25add',
//       kubernetes_server_id: '53f76c0e-2679-4b9d-a306-123d37a33c4c',
//       version: 'v1.27.4+k3s1',
//     },
//   ],
//   updated_contexts: [
//     {
//       id: '4',
//       name: 'ctx-4',
//       auth: {
//         name: 'admin@k3d-k3s-default',
//         user: {},
//       },
//       cluster: {
//         cluster: {
//           'insecure-skip-tls-verify': true,
//           server: 'https://ctx4-halibut-59567v944vrh547-8000.app.github.dev',
//         },
//         name: 'k3d-k3s-default',
//       },
//       server: 'https://ctx-4-halibut-59567v944vrh547-8000.app.github.dev',
//       meshery_instance_id: 'daa10439-7204-43d5-a425-a1bda9d25add',
//       kubernetes_server_id: '53f76c0e-2679-4b9d-a306-123d37a33c4c',
//       version: 'v1.27.4+k3s1',
//     },
//   ],
//   errored_contexts: [],
// };
const styles = makeStyles((theme) => ({
  ctxIcon: {
    display: 'inline',
    verticalAlign: 'text-top',
    width: theme.spacing(2.5),
    marginLeft: theme.spacing(0.5),
  },
  chip: {
    height: '50px',
    fontSize: '15px',
    position: 'relative',
    top: theme.spacing(0.5),
    [theme.breakpoints.down('md')]: { fontSize: '12px' },
  },
}));
// Add links to docs
const MesherySettingsEnvButtons = () => {
  let k8sfileElementVal = '';
  let formData = new FormData();
  const ref = useRef(null);
  const { notify } = useNotification();

  let contextsRef = useRef();

  const handleConfigSnackbars = (ctxs) => {
    updateProgress({ showProgress: false });
    for (let ctx of ctxs.errored_contexts) {
      const msg = `Failed to add cluster "${ctx.name}" at ${ctx.server}`;
      notify({ message: msg, event_type: EVENT_TYPES.ERROR, details: ctx.error.toString() });
    }
  };

  const handleError = (msg) => (error) => {
    updateProgress({ showProgress: false });
    notify({
      message: `${msg}: ${error}`,
      event_type: EVENT_TYPES.ERROR,
      details: error.toString(),
    });
  };

  const handleChange = () => {
    const field = document.getElementById('k8sfile');
    const textField = document.getElementById('k8sfileLabelText');
    if (field instanceof HTMLInputElement) {
      if (field.files.length < 1) return;
      const name = field.files[0].name;
      const formdata = new FormData();
      formdata.append('k8sfile', field.files[0]);
      textField.value = name;
      formData = formdata;
    }
  };

  const uploadK8SConfig = async () => {
    return await promisifiedDataFetch('/api/system/kubernetes', {
      method: 'POST',
      body: formData,
    });
  };

  const handleChangeConnectionStatus = async () => {
    const body = {};
    contextsRef.current.contextsToConnect.forEach((id) => {
      body[id] = CONNECTION_STATES.CONNECTED;
    });

    contextsRef.current.registered_contexts.map((context) => {
      if (!contextsRef.current.contextsToConnect.includes(context.connection_id)) {
        body[context.connection_id] = CONNECTION_STATES.IGNORED;
      }
    });

    return await promisifiedDataFetch('/api/integrations/connections/kubernetes/status', {
      method: 'PUT',
      body: JSON.stringify(body),
    });
  };

  const showUploadedContexts = async () => {
    const modal = ref.current;
    const registeredContexts = contextsRef.current.registered_contexts;
    const connectedContexts = contextsRef.current.connected_contexts;
    const ignoredContexts = contextsRef.current.ignored_contexts;

    let response = await modal.show({
      title: 'Select one or more contexts to manage',
      subtitle: (
        <>
          <ShowDiscoveredContexts
            registeredContexts={registeredContexts}
            connectedContexts={connectedContexts}
            ignoredContexts={ignoredContexts}
            allContextsRef={contextsRef}
          />
        </>
      ),
      options: ['CONNECT', 'CANCEL'],
    });

    if (response === 'CONNECT') {
      handleChangeConnectionStatus().then((res) => {
        console.log(res);
      });
    }
  };

  const handleClick = async () => {
    const modal = ref.current;
    let response = await modal.show({
      title: 'Add Kubernetes Cluster(s)',
      subtitle: (
        <>
          <div style={{ overflow: 'hidden' }}>
            <Typography variant="h6">Upload your kubeconfig</Typography>
            <Typography variant="body2">commonly found at ~/.kube/config</Typography>
            <FormGroup>
              <input
                id="k8sfile"
                type="file"
                value={k8sfileElementVal}
                onChange={handleChange}
                style={{ display: 'none' }}
              />

              <TextField
                id="k8sfileLabelText"
                name="k8sfileLabelText"
                style={{ cursor: 'pointer' }}
                placeholder="Upload kubeconfig"
                variant="outlined"
                fullWidth
                onClick={() => {
                  document.querySelector('#k8sfile')?.click();
                }}
                margin="normal"
                InputProps={{
                  readOnly: true,
                  endAdornment: (
                    <InputAdornment position="end">
                      <CloudUploadIcon />
                    </InputAdornment>
                  ),
                }}
              />
            </FormGroup>
          </div>
        </>
      ),
      options: ['IMPORT', 'CANCEL'],
    });

    if (response === 'IMPORT') {
      if (formData.get('k8sfile') === null) {
        handleError('No file selected')('Please select a valid kube config');
        return;
      }

      const inputFile = formData.get('k8sfile').name;
      const invalidExtensions = /^.*\.(jpg|gif|jpeg|pdf|png|svg)$/i;

      if (invalidExtensions.test(inputFile)) {
        handleError('Invalid file selected')('Please select a valid kube config');
        return;
      }

      uploadK8SConfig()
        .then((obj) => {
          contextsRef.current = obj;
          showUploadedContexts();
          handleConfigSnackbars(obj);
        })
        .catch((err) => {
          handleError('failed to upload kubernetes config')(err);
        });
      formData.delete('k8sfile');
    }
  };

  return (
    <div>
      <Button
        type="submit"
        variant="contained"
        color="primary"
        size="large"
        onClick={handleClick}
        style={{
          padding: '8px',
          borderRadius: 5,
          marginRight: '2rem',
        }}
        data-cy="btnResetDatabase"
      >
        <AddIconCircleBorder style={{ width: '20px', height: '20px' }} />
        <Typography
          style={{
            paddingLeft: '4px',
            marginRight: '4px',
          }}
        >
          {' '}
          Add Cluster
        </Typography>
      </Button>
      <PromptComponent ref={ref} />
    </div>
  );
};

const ShowDiscoveredContexts = ({
  registeredContexts,
  connectedContexts,
  ignoredContexts,
  allContextsRef,
}) => {
  const [contextsToUpdate, setContextsToUpdate] = useState([]);
  const ping = useKubernetesHook();

  useEffect(() => {
    allContextsRef.current = { ...allContextsRef.current, contextsToConnect: contextsToUpdate };
  }, [contextsToUpdate]);

  const setContextViewer = (id) => {
    if (contextsToUpdate?.includes(id)) {
      const filteredContexts = contextsToUpdate.filter((cid) => cid !== id);
      setContextsToUpdate(filteredContexts);
    } else {
      setContextsToUpdate([...contextsToUpdate, id]);
    }
  };
  return registeredContexts.length === 0 &&
    connectedContexts.length == 0 &&
    ignoredContexts.length == 0 ? (
    // Since user cannot chage status should we show this msg only when no new contexts are found or only show when no new and already context found.
    <K8sEmptyState />
  ) : (
    <Grid
      direction="column"
      justifyContent="center"
      alignItems="center"
      container
      spacing={2}
      columns={1}
    >
      {registeredContexts.length > 0 && (
        <Grid item xs={8}>
          <K8sConnectionItems
            contexts={registeredContexts}
            setContextViewer={setContextViewer}
            ping={ping}
            contextsToUpdate={contextsToUpdate}
            title="Registered Contexts"
          />
        </Grid>
      )}
      {connectedContexts.length > 0 && (
        <>
          <Divider variant="middle" />
          <Grid item xs={8}>
            <K8sConnectionItems
              contexts={connectedContexts}
              setContextViewer={setContextViewer}
              ping={ping}
              contextsToUpdate={contextsToUpdate}
              disableSelection
              title="Connected Contexts"
            />
          </Grid>
        </>
      )}
      {ignoredContexts.length > 0 && (
        <>
          <Divider variant="middle" />
          <Grid item xs={8}>
            <K8sConnectionItems
              contexts={ignoredContexts}
              setContextViewer={setContextViewer}
              ping={ping}
              contextsToUpdate={contextsToUpdate}
              disableSelection
              title="Ignored Contexts"
            />
          </Grid>
        </>
      )}
    </Grid>
  );
};

const K8sConnectionItems = ({
  title,
  contexts,
  setContextViewer,
  ping,
  contextsToUpdate,
  disableSelection,
}) => {
  const classes = styles();
  return (
    <>
      <Typography variant="subtitle" className="">
        {title}
      </Typography>
      {contexts.map((context) => (
        <div id={context.connection_id} key={context.connection_id} className={classes.chip}>
          <Tooltip title={`Server: ${context.server}`}>
            <div
              style={{
                display: 'flex',
                justifyContent: 'flex-wrap',
                alignItems: 'center',
              }}
            >
              {!disableSelection && (
                <Checkbox
                  checked={
                    contextsToUpdate?.includes(context.connection_id) //add support for selecting all
                  }
                  onChange={() => setContextViewer(context.connection_id)}
                  color="primary"
                />
              )}
              <ConnectionChip
                title={context.name}
                handlePing={() => {
                  ping(context.name, context.server, context.connection_id);
                }}
                icon={<img src="/static/img/kubernetes.svg" className={classes.ctxIcon} />}
              />
            </div>
          </Tooltip>
        </div>
      ))}
    </>
  );
};
export default MesherySettingsEnvButtons;
