import {
  Button,
  Typography,
  FormGroup,
  TextField,
  InputAdornment,
  Tooltip,
  makeStyles,
  Grid,
  Box,
} from '@material-ui/core';
import React from 'react';
import { useRef } from 'react';
import AddIconCircleBorder from '../assets/icons/AddIconCircleBorder';
import PromptComponent from './PromptComponent';
import CloudUploadIcon from '@material-ui/icons/CloudUpload';
import { promisifiedDataFetch } from '../lib/data-fetch';
import { updateProgress } from '../lib/store';
import { useNotification } from '../utils/hooks/useNotification';
import { EVENT_TYPES } from '../lib/event-types';
import { CONNECTION_STATES } from '../utils/Enum';
import { TootltipWrappedConnectionChip, ConnectionStateChip } from './connections/ConnectionChip';
import useKubernetesHook from './hooks/useKubernetesHook';

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

  const showUploadedContexts = async (inputFileName) => {
    const modal = ref.current;
    const registeredContexts = contextsRef.current.registered_contexts;
    const connectedContexts = contextsRef.current.connected_contexts;
    const ignoredContexts = contextsRef.current.ignored_contexts;
    if (
      registeredContexts.length === 0 &&
      connectedContexts.length == 0 &&
      ignoredContexts.length == 0
    ) {
      notify({
        message: `No reachable contexts found in the uploaded kube config "${inputFileName}". `,
        event_type: EVENT_TYPES.WARNING,
        showInNotificationCenter: true,
      });
      return;
    }
    await modal.show({
      title: `Available contexts in "${inputFileName}".`,
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
      options: ['OK'],
    });
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

      const inputFileName = formData.get('k8sfile').name;
      const invalidExtensions = /^.*\.(jpg|gif|jpeg|pdf|png|svg)$/i;

      if (invalidExtensions.test(inputFileName)) {
        handleError('Invalid file selected')('Please select a valid kube config');
        return;
      }

      uploadK8SConfig()
        .then((obj) => {
          contextsRef.current = obj;
          showUploadedContexts(inputFileName);
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

const ShowDiscoveredContexts = ({ registeredContexts, connectedContexts, ignoredContexts }) => {
  const ping = useKubernetesHook();

  return (
    <Grid
      direction="column"
      justifyContent="center"
      alignItems="center"
      spacing={2}
      columns={1}
      xs={16}
    >
      {registeredContexts.length > 0 && (
        <K8sConnectionItems
          contexts={registeredContexts}
          ping={ping}
          status={CONNECTION_STATES.REGISTERED}
        />
      )}
      {connectedContexts.length > 0 && (
        <K8sConnectionItems
          contexts={connectedContexts}
          ping={ping}
          status={CONNECTION_STATES.CONNECTED}
        />
      )}
      {ignoredContexts.length > 0 && (
        <Grid item xs={8}>
          <K8sConnectionItems
            contexts={ignoredContexts}
            ping={ping}
            status={CONNECTION_STATES.IGNORED}
          />
        </Grid>
      )}
    </Grid>
  );
};

const K8sConnectionItems = ({ status, contexts, ping }) => {
  const classes = styles();
  return (
    <Grid container xs={12} spacing={2}>
      {contexts.map((context) => (
        <Grid
          direction="column"
          alignContent="center"
          alignItems="center"
          container
          item
          xs={16}
          spacing={1}
          id={context.connection_id}
          key={context.connection_id}
          className={classes.chip}
        >
          <Box minWidth="25%" maxWidth="50%">
            <Tooltip title={`Server: ${context.server}`}>
              <div
                style={{
                  display: 'flex',
                  justifyContent: 'flex-wrap',
                  alignItems: 'center',
                }}
                whiteSpace="no-wrap"
                textOverflow="ellipsis"
              >
                <TootltipWrappedConnectionChip
                  title={context.name}
                  handlePing={() => {
                    ping(context.name, context.server, context.connection_id);
                  }}
                  iconSrc={'/static/img/kubernetes.svg'}
                />
              </div>
            </Tooltip>
          </Box>
          <Box minWidth="25%" maxWidth="50%">
            <ConnectionStateChip status={status} />
          </Box>
        </Grid>
      ))}
    </Grid>
  );
};
export default MesherySettingsEnvButtons;
