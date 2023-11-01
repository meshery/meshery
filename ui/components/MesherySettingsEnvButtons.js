import {
  Button,
  Typography,
  FormGroup,
  TextField,
  InputAdornment,
  Chip,
  Tooltip,
  makeStyles,
  Checkbox,
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
    for (let ctx of ctxs.updated_contexts) {
      const msg = `Cluster ${ctx.name} at ${ctx.server} already exists`;
      notify({ message: msg, event_type: EVENT_TYPES.INFO });
    }

    for (let ctx of ctxs.errored_contexts) {
      const msg = `Failed to add cluster ${ctx.name} at ${ctx.server}`;
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
    const contextsToIgnore = contextsRef.current.inserted_contexts.filter((context) => {
      if (!contextsRef.current.contextsToRegister.includes(context.connection_id)) {
        return context.connection_id;
      }
    });

    contextsRef.current.contextsToRegister.forEach((id) => {
      body[id] = 'CONNECTED';
    });

    contextsToIgnore.forEach((id) => {
      body[id] = 'IGNORED';
    });

    console.log('BODY: ', body);
    return await promisifiedDataFetch('/api/integrations/connections/kubernetes/status', {
      method: 'PUT',
      body: JSON.stringify(body),
    });
  };

  const showUploadedContexts = async () => {
    const modal = ref.current;
    console.log('ppp:8888 ', contextsRef.current);
    const alreadyAddedContexts = contextsRef.current.updated_contexts.map((context) => context);
    let response = await modal.show({
      title: 'Select one or more contexts to manage',
      subtitle: (
        <>
          <ShowDiscoveredContexts
            insertedContexts={contextsRef.current.inserted_contexts}
            contextsRef={contextsRef}
            alreadyAddedContexts={alreadyAddedContexts}
          />
        </>
      ),
      options: ['REGISTER', 'CANCEL'],
    });

    console.log('response: ', response, contextsRef.current);

    if (response === 'REGISTER') {
      console.log('response: ', contextsRef.current);
      handleChangeConnectionStatus().then((res) => {
        console.log('result: ', res);
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
          console.log('obj line 136: ', obj);
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

const ShowDiscoveredContexts = ({ insertedContexts, contextsRef, alreadyAddedContexts }) => {
  const classes = styles();
  const [contextsToUpdate, setContextsToUpdate] = useState([]);

  useEffect(() => {
    contextsRef.current = { ...contextsRef.current, contextsToRegister: contextsToUpdate };
    console.log('ppp', contextsToUpdate);
  }, [contextsToUpdate]);

  const { notify } = useNotification();

  const setContextViewer = (id) => {
    console.log('ppp.......lll,l,', id);
    if (contextsToUpdate?.includes(id)) {
      const filteredContexts = contextsToUpdate.filter((cid) => cid !== id);
      setContextsToUpdate(filteredContexts);
    } else {
      setContextsToUpdate([...contextsToUpdate, id]);
    }
  };
  return insertedContexts.length === 0 && alreadyAddedContexts.length == 0 ? ( // Since user cannot chage statius should we show this msg only when no new contexts are found or only show when no new and already context found.
    <Typography variant="subtitle1">No Context found</Typography>
  ) : (
    <>
      {alreadyAddedContexts.map((context) => (
        <div id={context.connection_id} key={context.connection_id}>
          <Tooltip title={`Server: ${context.server}`}>
            <div
              style={{
                display: 'flex',
                justifyContent: 'flex-wrap',
                alignItems: 'center',
              }}
            >
              <Checkbox
                checked={
                  true //add support for selecting all
                }
                onChange={() =>
                  notify({
                    message: 'Connection state cannot be changed during on going process',
                    event_type: EVENT_TYPES.INFO,
                  })
                }
                color="primary"
              />
              <Chip
                label={context.name}
                onClick={() => console.log('line 103: ', context)}
                icon={<img src="/static/img/kubernetes.svg" className={classes.ctxIcon} />}
                variant="outlined"
              />
            </div>
          </Tooltip>
        </div>
      ))}
      {insertedContexts.map((context) => (
        <div id={context.connection_id} key={context.connection_id}>
          <Tooltip title={`Server: ${context.server}`}>
            <div
              style={{
                display: 'flex',
                justifyContent: 'flex-wrap',
                alignItems: 'center',
              }}
            >
              <Checkbox
                checked={
                  contextsToUpdate?.includes(context.connection_id) //add support for selecting all
                }
                onChange={() => setContextViewer(context.connection_id)}
                color="primary"
              />
              <Chip
                label={context.name}
                onClick={() => console.log('line 117: ', context)}
                icon={<img src="/static/img/kubernetes.svg" className={classes.ctxIcon} />}
                variant="outlined"
              />
            </div>
          </Tooltip>
        </div>
      ))}
    </>
  );
};

export default MesherySettingsEnvButtons;
