import {
  Button,
  Typography,
  FormGroup,
  TextField,
  InputAdornment,
  Tooltip,
  Grid2,
  Box,
  styled,
  PROMPT_VARIANTS,
  FormControl,
  useTheme,
} from '@sistent/sistent';
import React, { useEffect, useState } from 'react';
import { useRef } from 'react';
import AddIconCircleBorder from '../assets/icons/AddIconCircleBorder';
import _PromptComponent from './PromptComponent';
import CloudUploadIcon from '@mui/icons-material/CloudUpload';
import { useNotification } from '../utils/hooks/useNotification';
import { EVENT_TYPES } from '../lib/event-types';
import { CONNECTION_STATES } from '../utils/Enum';
import { TooltipWrappedConnectionChip, ConnectionStateChip } from './connections/ConnectionChip';
import useKubernetesHook from './hooks/useKubernetesHook';
import { keys } from '@/utils/permission_constants';
import useTestIDsGenerator from './hooks/useTestIDs';
import CAN from '@/utils/can';
import { useAddKubernetesConfigMutation } from '../rtk-query/connection';
import { updateProgress } from '@/store/slices/mesheryUi';

const styles = styled((theme) => ({
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
const MesherySettingsEnvButtons = () => {
  let k8sfileElementVal = '';
  let formData = new FormData();
  const ref = useRef(null);
  const { notify } = useNotification();
  const theme = useTheme(); 
  const [deploymentMode, setDeploymentMode] = useState('operator');
  const deploymentModeRef = useRef('operator'); 

  useEffect(() => {
    deploymentModeRef.current = deploymentMode; // -> Keep ref in sync
  }, [deploymentMode]);

  let contextsRef = useRef();

  const testIDs = useTestIDsGenerator('connection');

  const [addK8sConfig] = useAddKubernetesConfigMutation();

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

  const handleValue = (mode) => {
    setDeploymentMode(mode);
    deploymentModeRef.current = mode; // Update ref immediately

    const operatorBtn = document.querySelector('[data-deployment-mode="operator"]');
    const embeddedBtn = document.querySelector('[data-deployment-mode="embedded"]');

    if (operatorBtn && embeddedBtn) {
      if (mode === 'operator') {
        operatorBtn.style.backgroundColor = theme.palette.background.brand.default;
        embeddedBtn.style.backgroundColor = 'transparent';
      } else {
        embeddedBtn.style.backgroundColor = theme.palette.background.brand.default;
        operatorBtn.style.backgroundColor = 'transparent';
      }
    }
  };
  const uploadK8SConfig = async () => {
    return await addK8sConfig({ body: formData }).unwrap();
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
            dataTestid={testIDs('discoveredModal')}
          />
        </>
      ),
      variant: PROMPT_VARIANTS.SUCCESS,
      primaryOption: 'OK',
    });
  };

  const handleClick = async () => {
    const modal = ref.current;
    let response = await modal.show({
      title: 'Add Kubernetes Cluster(s)',
      subtitle: (
        <>
          <div style={{ overflow: 'hidden' }} data-testid={testIDs('addKubernetesModal')}>
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
                data-testid={testIDs('uploadKubeConfig')}
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

            <FormControl style={{ marginTop: '16px' }}>
              <Typography variant="h6" style={{ marginBottom: '8px' }}>
                Deployment Mode
              </Typography>
              <Typography variant="body2" style={{ marginBottom: '12px' }}>
                Choose how MeshSync should be deployed in your cluster
              </Typography>

              <Box style={{ display: 'flex', gap: '12px', flexDirection: 'column' }}>
                <Button
                  data-deployment-mode="operator"
                  variant="contained"
                  onClick={() => handleValue('operator')}
                  style={{
                    textAlign: 'left',
                    justifyContent: 'flex-start',
                  }}
                >
                  <Box>
                    <Typography variant="body1">
                      <strong>Operator Mode</strong>
                    </Typography>
                    <Typography variant="body2" style={{ color: 'inherit' }}>
                      Deploy MeshSync as a separate operator in the cluster
                    </Typography>
                  </Box>
                </Button>

                <Button
                  data-deployment-mode="embedded"
                  variant="outlined"
                  onClick={() => handleValue('embedded')}
                  style={{
                    textAlign: 'left',
                    justifyContent: 'flex-start',
                  }}
                >
                  <Box>
                    <Typography variant="body1">
                      <strong>Embedded Mode</strong>
                    </Typography>
                    <Typography variant="body2" style={{ color: 'inherit' }}>
                      Run MeshSync embedded within Meshery Server
                    </Typography>
                  </Box>
                </Button>
              </Box>
            </FormControl>
          </div>
        </>
      ),
      primaryOption: 'IMPORT',
      variant: PROMPT_VARIANTS.SUCCESS,
      showInfoIcon: `If your config has not been autodetected, you can manually upload your kubeconfig file (or any number of kubeconfig files).
      The deployment mode determines how MeshSync operates: Operator mode deploys MeshSync as a separate operator in your cluster, while Embedded mode runs MeshSync within Meshery Server. [See Managing Kubernetes Clusters for more information](https://docs.meshery.io/installation/kubernetes).`,
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

      try {
        const currentDeploymentMode = deploymentModeRef.current;
        formData.delete('meshsync_deployment_mode');
        formData.append('meshsync_deployment_mode', currentDeploymentMode);

        const obj = await uploadK8SConfig();
        contextsRef.current = obj;
        await showUploadedContexts(inputFileName);
        handleConfigSnackbars(obj);
      } catch (err) {
        handleError('failed to upload kubernetes config')(err);
      }
      formData.delete('k8sfile');
    }
  };

  return (
    <div>
      <>
        <Button
          type="submit"
          variant="contained"
          onClick={handleClick}
          style={{
            width: '100%',
            borderRadius: 5,
            padding: '8px',
          }}
          disabled={!CAN(keys.ADD_CLUSTER.action, keys.ADD_CLUSTER.subject)}
          data-cy="btnResetDatabase"
        >
          <AddIconCircleBorder style={{ width: '20px', height: '20px' }} />
          <Typography
            style={{
              paddingLeft: '4px',
              width: 'max-content',
              marginRight: '4px',
            }}
            data-testid={testIDs('addCluster')}
          >
            Add Cluster
          </Typography>
        </Button>
      </>
      <_PromptComponent ref={ref} />
    </div>
  );
};

const ShowDiscoveredContexts = ({
  registeredContexts,
  connectedContexts,
  ignoredContexts,
  dataTestid,
}) => {
  const ping = useKubernetesHook();

  return (
    <Grid2
      direction="column"
      justifyContent="center"
      alignItems="center"
      spacing={2}
      columns={1}
      data-testid={dataTestid}
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
        <Grid2 size={{ xs: 8 }}>
          <K8sConnectionItems
            contexts={ignoredContexts}
            ping={ping}
            status={CONNECTION_STATES.IGNORED}
          />
        </Grid2>
      )}
    </Grid2>
  );
};

const K8sConnectionItems = ({ status, contexts, ping }) => {
  const classes = styles();
  return (
    <Grid2 container spacing={2} size={'grow'}>
      {contexts.map((context) => (
        <Grid2
          direction="column"
          alignContent="center"
          alignItems="center"
          container
          size="grow"
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
                <TooltipWrappedConnectionChip
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
        </Grid2>
      ))}
    </Grid2>
  );
};
export default MesherySettingsEnvButtons;
