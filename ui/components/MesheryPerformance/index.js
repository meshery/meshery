// @ts-nocheck
import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import Button from '@material-ui/core/Button';
import Typography from '@material-ui/core/Typography';
import { Autocomplete } from '@material-ui/lab';
import { withStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import { URLValidator } from '../../utils/URLValidator';
import {
  NoSsr,
  Tooltip,
  MenuItem,
  IconButton,
  CircularProgress,
  FormControl,
  FormLabel,
  RadioGroup,
  FormControlLabel,
  Radio,
  Divider,
  Link,
  ExpansionPanel,
  ExpansionPanelSummary,
  ExpansionPanelDetails,
} from '@material-ui/core';
import TextField from '@material-ui/core/TextField';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import GetAppIcon from '@material-ui/icons/GetApp';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import HelpOutlineOutlinedIcon from '@material-ui/icons/HelpOutlineOutlined';
import SaveOutlinedIcon from '@material-ui/icons/SaveOutlined';
import {
  updateLoadTestData,
  updateStaticPrometheusBoardConfig,
  updateLoadTestPref,
  updateProgress,
} from '../../lib/store';
import dataFetch from '../../lib/data-fetch';
import MesheryChart from '../MesheryChart';
import LoadTestTimerDialog from '../load-test-timer-dialog';
import GrafanaCustomCharts from '../telemetry/grafana/GrafanaCustomCharts';
import { durationOptions } from '../../lib/prePopulatedOptions';
import fetchControlPlanes from '../graphql/queries/ControlPlanesQuery';
import { ctxUrl, getK8sClusterIdsFromCtxId } from '../../utils/multi-ctx';
import { iconMedium } from '../../css/icons.styles';
import { withNotify } from '../../utils/hooks/useNotification';
import { EVENT_TYPES } from '../../lib/event-types';
import { generateTestName, generateUUID } from './helper';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import DefaultError from '@/components/General/error-404/index';

// =============================== HELPER FUNCTIONS ===========================

/**
 * generatePerformanceProfile takes in data and generate a performance
 * profile object from it
 * @param {*} data
 */
export function generatePerformanceProfile(data) {
  const {
    id,
    name,
    loadGenerator,
    additional_options,
    endpoint,
    serviceMesh,
    concurrentRequest,
    qps,
    duration,
    requestHeaders,
    requestCookies,
    requestBody,
    contentType,
    caCertificate,
  } = data;

  const performanceProfileName = generateTestName(name, serviceMesh);

  return {
    ...(id && { id }),
    name: performanceProfileName,
    load_generators: [loadGenerator],
    endpoints: [endpoint],
    service_mesh: serviceMesh,
    concurrent_request: concurrentRequest,
    qps,
    duration,
    request_headers: requestHeaders,
    request_body: requestBody,
    request_cookies: requestCookies,
    content_type: contentType,
    metadata: {
      additional_options: [additional_options],
      ca_certificate: {
        file: caCertificate.file,
        name: caCertificate.name,
      },
    },
  };
}

const styles = (theme) => ({
  title: {
    textAlign: 'center',
    minWidth: 400,
    padding: '10px',
    color: '#fff',
    backgroundColor:
      theme.palette.type === 'dark'
        ? theme.palette.secondary.headerColor
        : theme.palette.secondary.mainBackground,
  },
  wrapperClss: {
    padding: theme.spacing(10),
    position: 'relative',
    paddingTop: theme.spacing(5),
  },
  buttons: { display: 'flex', justifyContent: 'flex-end' },
  spacing: {
    marginTop: theme.spacing(3),
    marginLeft: theme.spacing(1),
  },
  button: {
    backgroundColor: theme.palette.type === 'dark' ? '#00B39F' : '#607d8b',
    '&:hover': {
      backgroundColor: theme.palette.type === 'dark' ? '#00B39F' : '#607d8b',
    },
    color: '#fff',
  },
  iconColor: {
    color: '#929292',
  },
  upload: {
    paddingLeft: '0.7rem',
    paddingTop: '8px',
  },
  expansionPanel: { boxShadow: 'none', border: '1px solid rgb(196,196,196)' },
  margin: { margin: theme.spacing(1) },
  chartTitle: { textAlign: 'center' },
  chartTitleGraf: {
    textAlign: 'center',
    // marginTop: theme.spacing(5),
  },
  chartContent: {
    // minHeight: window.innerHeight * 0.7,
  },
  centerTimer: {
    width: '100%',
    height: '100%',
    position: 'absolute',
    top: '0',
    left: '0',
    backgroundColor: 'rgba(0, 0, 0, 0.6)',
    zIndex: 1201,
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
  },
  paper: {
    backgroundColor: theme.palette.background.paper,
    border: '2px solid #000',
    boxShadow: theme.shadows[5],
    padding: theme.spacing(2, 4, 3),
  },
  smallIcons: {
    width: '15px',
    height: '18px',
    marginBottom: theme.spacing(1),
    marginLeft: theme.spacing(0.3),
  },
  radio: {
    '&.Mui-checked': {
      color:
        theme.palette.type === 'dark' ? theme.palette.secondary.focused : theme.palette.primary,
    },
  },
});
// =============================== PERFORMANCE COMPONENT =======================
const loadGenerators = ['fortio', 'wrk2', 'nighthawk'];

const infoFlags = <>Only .json files are supported.</>;

const infoCRTCertificates = <>Only .crt files are supported.</>;

const infoloadGenerators = (
  <>
    Which load generators does Meshery support?
    <ul>
      <li>
        fortio - Fortio load testing library, command line tool, advanced echo server and web UI in
        go (golang). Allows to specify a set query-per-second load and record latency histograms and
        other useful stats.{' '}
      </li>
      <li> wrk2 - A constant throughput, correct latency recording variant of wrk.</li>
      <li>
        {' '}
        nighthawk - Enables users to run distributed performance tests to better mimic real-world,
        distributed systems scenarios.
      </li>
    </ul>
    <Link
      style={{ textDecoration: 'underline' }}
      color="inherit"
      href="https://docs.meshery.io/functionality/performance-management"
    >
      {' '}
      Performance Management
    </Link>
  </>
);

let eventStream = null;
const MesheryPerformanceComponent = (props) => {
  const {
    testName = '',
    meshName = '',
    url = '',
    qps = '0',
    c = '0',
    t = '30s',
    result,
    staticPrometheusBoardConfig,
    performanceProfileID,
    profileName,
    loadGenerator,
    additional_options,
    headers,
    cookies,
    reqBody,
    contentType,
    metadata,
  } = props;
  const isJsonString = (str) => {
    try {
      JSON.parse(str);
    } catch (e) {
      return false;
    }
    return true;
  };
  // Create individual state variables for each property
  const [testNameState, setTestName] = useState(testName);
  const [meshNameState, setMeshName] = useState(meshName);
  const [urlState, setUrl] = useState(url);
  const [qpsState, setQps] = useState(qps);
  const [cState, setC] = useState(c);
  const [tState, setT] = useState(t);
  const [tValueState, setTValue] = useState(t);
  const [loadGeneratorState, setLoadGenerator] = useState(loadGenerator || 'fortio');
  const [additionalOptionsState, setAdditionalOptions] = useState(additional_options || '');
  const [resultState, setResult] = useState(result);
  const [headersState, setHeaders] = useState(headers || '');
  const [cookiesState, setCookies] = useState(cookies || '');
  const [reqBodyState, setReqBody] = useState(reqBody || '');
  const [contentTypeState, setContentType] = useState(contentType || '');
  const [caCertificateState, setCaCertificate] = useState({});
  const [profileNameState, setProfileName] = useState(profileName || '');
  const [performanceProfileIDState, setPerformanceProfileID] = useState(performanceProfileID || '');
  const [timerDialogOpenState, setTimerDialogOpen] = useState(false);
  const [blockRunTestState, setBlockRunTest] = useState(false);
  const [urlErrorState, setUrlError] = useState(false);
  const [tErrorState, setTError] = useState('');
  const [jsonErrorState, setJsonError] = useState(false);
  const [disableTestState, setDisableTest] = useState(
    !(URLValidator(urlState) || isJsonString(additionalOptionsState)),
  );
  const [testUUIDState, setTestUUID] = useState(generateUUID());
  const [selectedMeshState, setSelectedMesh] = useState('');
  const [availableAdaptersState, setAvailableAdapters] = useState([]);
  const [availableSMPMeshesState, setAvailableSMPMeshes] = useState([]);
  const [metadataState, setMetadata] = useState(metadata);
  const [staticPrometheusBoardConfigState, setStaticPrometheusBoardConfig] = useState(
    staticPrometheusBoardConfig,
  );
  const handleChange = (name) => (event) => {
    const { value } = event.target;
    if (name === 'caCertificate') {
      if (!event.target.files?.length) return;

      const file = event.target.files[0];
      const reader = new FileReader();
      reader.addEventListener('load', (evt) => {
        setCaCertificate({
          name: file.name,
          file: evt.target.result,
        });
      });
      reader.readAsText(file);
    }

    if (name === 'url' && value !== '') {
      let urlPattern = value;

      let val = URLValidator(urlPattern);
      if (!val) {
        setDisableTest(true);
        setUrlError(true);
      } else {
        setDisableTest(false);
        setUrlError(false);
      }
    } else setUrlError(false);

    if (name === 'additional_options') {
      // Check if the target event is an input element (typing) or a file input (upload)
      const isFileUpload = event.target.getAttribute('type') === 'file';

      if (isFileUpload) {
        // Handle file upload
        const file = event.target.files[0];
        if (file) {
          const reader = new FileReader();
          reader.onload = (event) => {
            try {
              const fileContent = event.target.result;
              // Validate JSON
              JSON.parse(fileContent);
              setAdditionalOptions(fileContent);
              setJsonError(false);
            } catch (error) {
              setAdditionalOptions(event.target.result);
              setJsonError(true);
            }
          };
          reader.readAsText(file);
        }
      } else {
        // Handle text input
        try {
          // empty text input exception
          if (value !== '') JSON.parse(value);
          setAdditionalOptions(value);
          setJsonError(false);
        } catch (error) {
          setAdditionalOptions(value);
          setJsonError(true);
        }
      }
    }
    switch (name) {
      case 'profileName':
        setProfileName(value);
        break;
      case 'meshName':
        setMeshName(value);
        break;
      case 'c':
        setC(value);
        break;
      case 'qps':
        setQps(value);
        break;
      case 'headers':
        setHeaders(value);
        break;
      case 'cookies':
        setCookies(value);
        break;
      case 'contentType':
        setContentType(value);
        break;
      case 'reqBody':
        setReqBody(value);
        break;
      case 'loadGenerator':
        setLoadGenerator(value);
        break;
      case 'url':
        setUrl(value);
        break;
      default:
        // Handle any other cases or do nothing if not matched
        break;
    }
  };

  const handleDurationChange = (event, newValue) => {
    setTValue(newValue);
    if (newValue !== null) {
      setTError('');
    }
  };

  const handleInputDurationChange = (event, newValue) => {
    setT(newValue);
  };

  const handleSubmit = () => {
    if (urlState === '') {
      setUrlError(true);
      return;
    }

    let err = false;
    let tNum = 0;
    try {
      tNum = parseInt(t.substring(0, tState.length - 1));
    } catch (ex) {
      err = true;
    }

    if (
      tState === '' ||
      tState === null ||
      !(
        tState.toLowerCase().endsWith('h') ||
        tState.toLowerCase().endsWith('m') ||
        tState.toLowerCase().endsWith('s')
      ) ||
      err ||
      tNum <= 0
    ) {
      setTError('error-autocomplete-value');
      return;
    }

    if (!performanceProfileIDState) {
      submitProfile(({ id }) => submitLoadTest(id));
      return;
    }
    submitLoadTest(performanceProfileIDState);
  };

  const submitProfile = (cb) => {
    const profile = generatePerformanceProfile({
      name: profileNameState,
      loadGenerator: loadGeneratorState,
      additional_options: additionalOptionsState,
      endpoint: urlState,
      serviceMesh: meshNameState,
      concurrentRequest: +cState || 0,
      qps: +qpsState || 0,
      duration: tState,
      requestHeaders: headersState,
      requestCookies: cookiesState,
      requestBody: reqBodyState,
      contentType: contentTypeState,
      caCertificate: caCertificateState,
      testName: testNameState,
      id: performanceProfileIDState,
    });

    handleProfileUpload(profile, true, cb);
  };

  const handleAbort = () => {
    setProfileName('');
    setLoadGenerator('');
    setAdditionalOptions('');
    setUrl('');
    setMeshName('');
    setC(0);
    setQps(0);
    setT('30s');
    setHeaders('');
    setCookies('');
    setReqBody('');
    setContentType('');
    setTestName('');
    setPerformanceProfileID('');
    setDisableTest(true);
    return;
  };

  const handleProfileUpload = (body, generateNotif, cb) => {
    if (generateNotif) props.updateProgress({ showProgress: true });

    dataFetch(
      '/api/user/performance/profiles',
      { method: 'POST', credentials: 'include', body: JSON.stringify(body) },
      (result) => {
        if (typeof result !== 'undefined') {
          props.updateProgress({ showProgress: false });
          setPerformanceProfileID(result.id);
          if (cb) cb(result);
          if (generateNotif) {
            const notify = props.notify;
            notify({
              message: 'Performance Profile Created!',
              event_type: EVENT_TYPES.SUCCESS,
            });
          }
        }
      },
      (err) => {
        console.error(err);
        props.updateProgress({ showProgress: false });
        const notify = props.notify;
        notify({
          message: 'Failed to create performance profile',
          event_type: EVENT_TYPES.ERROR,
          details: err.toString(),
        });
      },
    );
  };

  const submitLoadTest = (id) => {
    const computedTestName = generateTestName(testNameState, meshNameState);
    setTestName(computedTestName);

    const t1 = tState.substring(0, tState.length - 1);
    const dur = tState.substring(tState.length - 1, tState.length).toLowerCase();

    const data = {
      name: computedTestName,
      mesh: meshName === '' || meshName === 'None' ? '' : meshNameState, // to prevent None from getting to the DB
      url: urlState,
      qps: qpsState,
      c: cState,
      t: t1,
      dur,
      uuid: testUUIDState,
      loadGenerator: loadGeneratorState,
      additional_options: additionalOptionsState,
      headers: headersState,
      cookies: cookiesState,
      reqBody: reqBodyState,
      contentType: contentTypeState,
    };
    const params = Object.keys(data)
      .map((key) => `${encodeURIComponent(key)}=${encodeURIComponent(data[key])}`)
      .join('&');

    const runURL =
      ctxUrl(`/api/user/performance/profiles/${id}/run`, props?.selectedK8sContexts) + '&cert=true';
    startEventStream(`${runURL}${props?.selectedK8sContexts?.length > 0 ? '&' : '?'}${params}`);
    setBlockRunTest(true); // to block the button
  };

  function handleSuccess() {
    return (result) => {
      if (typeof result !== 'undefined' && typeof result.runner_results !== 'undefined') {
        const notify = props.notify;
        notify({
          message: 'fetched the data.',
          event_type: EVENT_TYPES.SUCCESS,
        });
        props.updateLoadTestData({
          loadTest: {
            testName: testNameState,
            meshName: meshNameState,
            url: urlState,
            qps: qpsState,
            c: cState,
            t: tState,
            loadGenerator: loadGeneratorState,
            result: resultState,
          },
        });
        setTestUUID(generateUUID());
        setResult(result);
      }
      closeEventStream();
      setBlockRunTest(false);
      setTimerDialogOpen(false);
    };
  }
  async function startEventStream(url) {
    closeEventStream();
    eventStream = new EventSource(url);
    eventStream.onmessage = handleEvents();
    eventStream.onerror = handleError(
      'Connection to the server got disconnected. Load test might be running in the background. Please check the results page in a few.',
    );
    const notify = props.notify;
    notify({
      message: 'Load test has been submitted',
      event_type: EVENT_TYPES.SUCCESS,
    });
  }

  function handleEvents() {
    const notify = props.notify;
    let track = 0;
    return (e) => {
      const data = JSON.parse(e.data);
      switch (data.status) {
        case 'info':
          notify({ message: data.message, event_type: EVENT_TYPES.INFO });
          if (track === 0) {
            setTimerDialogOpen(true);
            setResult({});
            track++;
          }
          break;
        case 'error':
          handleError('Load test did not run with msg')(data.message);
          break;
        case 'success':
          handleSuccess()(data.result);
          break;
      }
    };
  }

  function closeEventStream() {
    if (eventStream && eventStream.close) {
      eventStream.close();
      eventStream = null;
    }
  }
  useEffect(() => {
    getStaticPrometheusBoardConfig();
    scanForMeshes();
    getLoadTestPrefs();
    getSMPMeshes();
    if (props.runTestOnMount) handleSubmit();
  }, []);

  const getLoadTestPrefs = () => {
    dataFetch(
      ctxUrl('/api/user/prefs', props?.selectedK8sContexts),
      { credentials: 'same-origin', method: 'GET' },
      (result) => {
        if (typeof result !== 'undefined') {
          setQps(result.loadTestPrefs.qps);
          setC(result.loadTestPrefs.c);
          setT(result.loadTestPrefs.t);
          setLoadGenerator(result.loadTestPrefs.gen);
        }
      },
      () => {},
    ); //error is already captured from the handler, also we have a redux-store for same & hence it's not needed here.
  };

  const getStaticPrometheusBoardConfig = () => {
    if (
      (staticPrometheusBoardConfig &&
        staticPrometheusBoardConfig !== null &&
        Object.keys(props.staticPrometheusBoardConfig).length > 0) ||
      (staticPrometheusBoardConfigState &&
        staticPrometheusBoardConfigState !== null &&
        Object.keys(staticPrometheusBoardConfigState).length > 0)
    ) {
      return;
    }
    dataFetch(
      '/api/telemetry/metrics/static-board',
      { credentials: 'include' },
      (result) => {
        if (
          typeof result !== 'undefined' &&
          typeof result.cluster !== 'undefined' &&
          typeof result.node !== 'undefined' &&
          typeof result.cluster.panels !== 'undefined' &&
          result.cluster.panels.length > 0 &&
          typeof result.node.panels !== 'undefined' &&
          result.node.panels.length > 0
        ) {
          props.updateStaticPrometheusBoardConfig({
            staticPrometheusBoardConfig: result, // will contain both the cluster and node keys for the respective boards
          });
          setStaticPrometheusBoardConfig(result);
        }
      },
      handleError('unable to fetch pre-configured boards'),
    );
  };

  const getK8sClusterIds = () => {
    return getK8sClusterIdsFromCtxId(props.selectedK8sContexts, props.k8sconfig);
  };

  const scanForMeshes = () => {
    if (typeof props.k8sConfig === 'undefined' || !props.k8sConfig.clusterConfigured) {
      return;
    }
    /**
     * ALL_MESH indicates that we are interested in control plane
     * component of all of the service meshes supported by meshsync v2
     */

    const ALL_MESH = {
      type: 'ALL_MESH',
      k8sClusterIDs: getK8sClusterIds(),
    };

    fetchControlPlanes(ALL_MESH).subscribe({
      next: (res) => {
        let result = res?.controlPlanesState;
        if (typeof result !== 'undefined' && Object.keys(result).length > 0) {
          const adaptersList = [];
          result.forEach((mesh) => {
            if (mesh?.members.length > 0) {
              let name = mesh?.name;
              adaptersList.push(
                // Capatilize First Letter and replace undersocres
                name
                  .split(/ |_/i)
                  .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
                  .join(' '),
              );
            }
          });
          setAvailableAdapters(adaptersList);
          result.forEach((mesh) => {
            setSelectedMesh(mesh?.name);
          });
        }
      },
      error: (err) => console.error(err),
    });
  };

  const getSMPMeshes = () => {
    dataFetch(
      '/api/mesh',
      { credentials: 'include' },
      (result) => {
        if (result && Array.isArray(result.available_meshes)) {
          setAvailableSMPMeshes(result.available_meshes.sort((m1, m2) => m1.localeCompare(m2)));
        }
      },
      handleError('unable to fetch SMP meshes'),
    );
  };

  function handleError(msg) {
    return (error) => {
      setBlockRunTest(false);
      setTimerDialogOpen(false);
      closeEventStream();
      let finalMsg = msg;
      if (typeof error === 'string') {
        finalMsg = `${msg}: ${error}`;
      }
      const notify = props.notify;
      notify({
        message: finalMsg,
        event_type: EVENT_TYPES.ERROR,
        details: error.toString(),
      });
    };
  }

  const handleCertificateUpload = (event) => {
    const file = event.target.files[0];
    if (file) {
      const newMetadata = {
        ...metadataState,
        ca_certificate: {
          ...metadataState.ca_certificate,
          name: file.name,
        },
      };
      setMetadata(newMetadata);
    }
  };

  const handleTimerDialogClose = () => {
    setTimerDialogOpen(false);
  };
  const { classes, grafana, prometheus } = props;
  let localStaticPrometheusBoardConfig;
  if (
    props.staticPrometheusBoardConfig &&
    props.staticPrometheusBoardConfig != null &&
    Object.keys(props.staticPrometheusBoardConfig).length > 0
  ) {
    localStaticPrometheusBoardConfig = props.staticPrometheusBoardConfig;
  } else {
    localStaticPrometheusBoardConfig = staticPrometheusBoardConfigState;
  }
  let chartStyle = {};
  if (timerDialogOpenState) {
    chartStyle = { opacity: 0.3 };
  }
  let displayStaticCharts = null;
  let displayGCharts = null;
  let displayPromCharts = null;

  availableAdaptersState.forEach((item) => {
    const index = availableSMPMeshesState.indexOf(item);
    if (index !== -1) availableSMPMeshesState.splice(index, 1);
  });

  if (
    localStaticPrometheusBoardConfig &&
    localStaticPrometheusBoardConfig !== null &&
    Object.keys(localStaticPrometheusBoardConfig).length > 0 &&
    prometheus.prometheusURL !== ''
  ) {
    // only add testUUID to the board that should be persisted
    if (localStaticPrometheusBoardConfig.cluster) {
      localStaticPrometheusBoardConfig.cluster.testUUID = testUUIDState;
    }
    displayStaticCharts = (
      <React.Fragment>
        <Typography variant="h6" gutterBottom className={classes.chartTitle}>
          Node Metrics
        </Typography>
        <GrafanaCustomCharts
          boardPanelConfigs={[
            localStaticPrometheusBoardConfig.cluster,
            localStaticPrometheusBoardConfig.node,
          ]}
          prometheusURL={prometheus.prometheusURL}
        />
      </React.Fragment>
    );
  }
  if (prometheus.selectedPrometheusBoardsConfigs.length > 0) {
    displayPromCharts = (
      <React.Fragment>
        <Typography variant="h6" gutterBottom className={classes.chartTitleGraf}>
          Prometheus charts
        </Typography>
        <GrafanaCustomCharts
          boardPanelConfigs={prometheus.selectedPrometheusBoardsConfigs}
          prometheusURL={prometheus.prometheusURL}
        />
      </React.Fragment>
    );
  }
  if (grafana.selectedBoardsConfigs.length > 0) {
    displayGCharts = (
      <React.Fragment>
        <Typography variant="h6" gutterBottom className={classes.chartTitleGraf}>
          Grafana charts
        </Typography>
        <GrafanaCustomCharts
          boardPanelConfigs={grafana.selectedBoardsConfigs}
          grafanaURL={grafana.grafanaURL}
          grafanaAPIKey={grafana.grafanaAPIKey}
        />
      </React.Fragment>
    );
  }
  return (
    <NoSsr>
      {CAN(keys.VIEW_PERFORMANCE_PROFILES.action, keys.VIEW_PERFORMANCE_PROFILES.subject) ? (
        <>
          <React.Fragment>
            <div className={classes.wrapperClss} style={props.style || {}}>
              <Grid container spacing={1}>
                <Grid item xs={12} md={6}>
                  <TextField
                    id="profileName"
                    name="profileName"
                    label="Profile Name"
                    fullWidth
                    value={profileNameState}
                    margin="normal"
                    variant="outlined"
                    onChange={handleChange('profileName')}
                    inputProps={{
                      maxLength: 300,
                    }}
                    InputProps={{
                      endAdornment: (
                        <Tooltip title="Create a profile providing a name, if a profile name is not provided, a random one will be generated for you.">
                          <HelpOutlineOutlinedIcon className={classes.iconColor} />
                        </Tooltip>
                      ),
                    }}
                  />
                </Grid>

                <Grid item xs={12} md={6}>
                  <TextField
                    select
                    id="meshName"
                    name="meshName"
                    label="Service Mesh"
                    fullWidth
                    value={
                      meshNameState === '' && selectedMeshState !== ''
                        ? selectedMeshState
                        : meshNameState
                    }
                    margin="normal"
                    variant="outlined"
                    onChange={handleChange('meshName')}
                  >
                    {availableAdaptersState &&
                      availableAdaptersState.map((mesh) => (
                        <MenuItem key={`mh_-_${mesh}`} value={mesh.toLowerCase()}>
                          {mesh}
                        </MenuItem>
                      ))}
                    {availableAdaptersState && availableAdaptersState.length > 0 && <Divider />}
                    <MenuItem key="mh_-_none" value="None">
                      None
                    </MenuItem>
                    {availableSMPMeshesState &&
                      availableSMPMeshesState.map((mesh) => (
                        <MenuItem key={`mh_-_${mesh}`} value={mesh.toLowerCase()}>
                          {mesh}
                        </MenuItem>
                      ))}
                  </TextField>
                </Grid>
                <Grid item xs={12}>
                  <TextField
                    required
                    id="url"
                    name="url"
                    label="URL to test"
                    type="url"
                    fullWidth
                    value={urlState}
                    error={urlErrorState}
                    helperText={urlErrorState ? 'Please enter a valid URL along with protocol' : ''}
                    margin="normal"
                    variant="outlined"
                    onChange={handleChange('url')}
                    InputProps={{
                      endAdornment: (
                        <Tooltip title="The Endpoint where the load will be generated and the perfromance test will run against.">
                          <HelpOutlineOutlinedIcon className={classes.iconColor} />
                        </Tooltip>
                      ),
                    }}
                  />
                </Grid>
                <Grid item xs={12} md={4}>
                  <TextField
                    required
                    id="c"
                    name="c"
                    label="Concurrent requests"
                    type="number"
                    fullWidth
                    value={cState}
                    inputProps={{ min: '0', step: '1' }}
                    margin="normal"
                    variant="outlined"
                    onChange={handleChange('c')}
                    InputLabelProps={{ shrink: true }}
                    InputProps={{
                      endAdornment: (
                        <Tooltip title="Load Testing tool will create this many concurrent request against the endpoint.">
                          <HelpOutlineOutlinedIcon className={classes.iconColor} />
                        </Tooltip>
                      ),
                    }}
                  />
                </Grid>
                <Grid item xs={12} md={4}>
                  <TextField
                    required
                    id="qps"
                    name="qps"
                    label="Queries per second"
                    type="number"
                    fullWidth
                    value={qpsState}
                    inputProps={{ min: '0', step: '1' }}
                    margin="normal"
                    variant="outlined"
                    onChange={handleChange('qps')}
                    InputLabelProps={{ shrink: true }}
                    InputProps={{
                      endAdornment: (
                        <Tooltip title="The Number of queries/second. If not provided then the MAX number of queries/second will be requested">
                          <HelpOutlineOutlinedIcon className={classes.iconColor} />
                        </Tooltip>
                      ),
                    }}
                  />
                </Grid>
                <Grid item xs={12} md={4}>
                  <Tooltip
                    title={
                      "Please use 'h', 'm' or 's' suffix for hour, minute or second respectively."
                    }
                  >
                    <Autocomplete
                      required
                      id="t"
                      name="t"
                      freeSolo
                      label="Duration*"
                      fullWidth
                      variant="outlined"
                      className={classes.errorValue}
                      classes={{ root: tErrorState }}
                      value={tValueState}
                      inputValue={tState}
                      onChange={handleDurationChange}
                      onInputChange={handleInputDurationChange}
                      options={durationOptions}
                      style={{ marginTop: '16px', marginBottom: '8px' }}
                      renderInput={(params) => (
                        <TextField {...params} label="Duration*" variant="outlined" />
                      )}
                      InputProps={{
                        endAdornment: (
                          <Tooltip title="Default duration is 30 seconds">
                            <HelpOutlineOutlinedIcon className={classes.iconColor} />
                          </Tooltip>
                        ),
                      }}
                    />
                  </Tooltip>
                </Grid>
                <Grid item xs={12} md={12}>
                  <ExpansionPanel className={classes.expansionPanel}>
                    <ExpansionPanelSummary expanded={true} expandIcon={<ExpandMoreIcon />}>
                      <Typography align="center" color="textSecondary" variant="h6">
                        Advanced Options
                      </Typography>
                    </ExpansionPanelSummary>
                    <ExpansionPanelDetails>
                      <Grid container spacing={1}>
                        <Grid item xs={12}>
                          <TextField
                            id="headers"
                            name="headers"
                            label='Request Headers e.g. {"host":"bookinfo.meshery.io"}'
                            fullWidth
                            value={headersState}
                            multiline
                            margin="normal"
                            variant="outlined"
                            onChange={handleChange('headers')}
                          ></TextField>
                        </Grid>
                        <Grid item xs={12}>
                          <TextField
                            id="cookies"
                            name="cookies"
                            label='Request Cookies e.g. {"yummy_cookie":"choco_chip"}'
                            fullWidth
                            value={cookiesState}
                            multiline
                            margin="normal"
                            variant="outlined"
                            onChange={handleChange('cookies')}
                          ></TextField>
                        </Grid>
                        <Grid item xs={12}>
                          <TextField
                            id="contentType"
                            name="contentType"
                            label="Content Type e.g. application/json"
                            fullWidth
                            value={contentTypeState}
                            multiline
                            margin="normal"
                            variant="outlined"
                            onChange={handleChange('contentType')}
                          ></TextField>
                        </Grid>
                        <Grid item xs={12} md={12}>
                          <TextField
                            id="cookies"
                            name="cookies"
                            label='Request Body e.g. {"method":"post","url":"http://bookinfo.meshery.io/test"}'
                            fullWidth
                            value={reqBodyState}
                            multiline
                            margin="normal"
                            variant="outlined"
                            onChange={handleChange('reqBody')}
                          ></TextField>
                        </Grid>
                        <Grid container xs={12} md={12}>
                          <Grid item xs={6}>
                            <TextField
                              id="additional_options"
                              name="additional_options"
                              label="Additional Options e.g. { `requestPerSecond`: 20 }"
                              fullWidth
                              error={jsonErrorState}
                              helperText={jsonErrorState ? 'Please enter a valid JSON string' : ''}
                              value={
                                additionalOptionsState.length > 150
                                  ? `${additionalOptionsState.slice(0, 150)} .....`
                                  : additionalOptionsState
                              }
                              multiline
                              margin="normal"
                              variant="outlined"
                              size="small"
                              onChange={handleChange('additional_options')}
                            />
                          </Grid>
                          <Grid item xs={6}>
                            <label
                              htmlFor="upload-additional-options"
                              style={{ paddingLeft: '0' }}
                              className={classes.upload}
                              fullWidth
                            >
                              <Button
                                variant="outlined"
                                onChange={handleChange('additional_options')}
                                aria-label="Upload Button"
                                component="span"
                                className={classes.button}
                                style={{ margin: '0.5rem', marginTop: '1.15rem' }}
                              >
                                <input
                                  id="upload-additional-options"
                                  type="file"
                                  accept={'.json'}
                                  name="upload-button"
                                  hidden
                                  data-cy="additional-options-upload-button"
                                />
                                Browse
                              </Button>
                              <Tooltip title={infoFlags} interactive>
                                <HelpOutlineOutlinedIcon className={classes.smallIcons} />
                              </Tooltip>
                            </label>
                          </Grid>
                        </Grid>
                        <Grid container xs={12} md={12}>
                          <Grid item xs={6}>
                            <TextField
                              size="small"
                              variant="outlined"
                              margin="mormal"
                              fullWidth
                              label={
                                caCertificateState?.name || 'Upload SSL Certificate e.g. .crt file'
                              }
                              style={{ width: '100%', margin: '0.5rem 0' }}
                              value={metadataState?.ca_certificate.name}
                            />
                          </Grid>
                          <Grid item xs={6}>
                            <label
                              htmlFor="upload-cacertificate"
                              className={classes.upload}
                              style={{ paddingLeft: '0' }}
                            >
                              <Button
                                variant="outlined"
                                aria-label="Upload Button"
                                onChange={handleChange('caCertificate')}
                                component="span"
                                className={classes.button}
                                style={{ margin: '0.5rem' }}
                              >
                                <input
                                  id="upload-cacertificate"
                                  type="file"
                                  accept={'.crt'}
                                  name="upload-button"
                                  hidden
                                  data-cy="cacertificate-upload-button"
                                  onChange={handleCertificateUpload}
                                />
                                Browse
                              </Button>
                              <Tooltip title={infoCRTCertificates} interactive>
                                <HelpOutlineOutlinedIcon className={classes.smallIcons} />
                              </Tooltip>
                            </label>
                          </Grid>
                        </Grid>
                      </Grid>
                    </ExpansionPanelDetails>
                  </ExpansionPanel>
                </Grid>
                <Grid item xs={12} md={4}>
                  <FormControl component="loadGenerator" className={classes.margin}>
                    <FormLabel
                      component="loadGenerator"
                      style={{
                        display: 'flex',
                        alignItems: 'center',
                        flexWrap: 'wrap',
                      }}
                    >
                      Load generator
                      <Tooltip title={infoloadGenerators} interactive>
                        <HelpOutlineOutlinedIcon className={classes.smallIcons} />
                      </Tooltip>
                    </FormLabel>
                    <RadioGroup
                      aria-label="loadGenerator"
                      name="loadGenerator"
                      value={loadGeneratorState}
                      onChange={handleChange('loadGenerator')}
                      row
                    >
                      {loadGenerators.map((lg, index) => (
                        <FormControlLabel
                          key={index}
                          value={lg}
                          disabled={lg === 'wrk2'}
                          control={<Radio color="primary" className={classes.radio} />}
                          label={lg}
                        />
                      ))}
                    </RadioGroup>
                  </FormControl>
                </Grid>
              </Grid>
              <React.Fragment>
                <div className={classes.buttons}>
                  <Button
                    type="submit"
                    variant="contained"
                    color="primary"
                    size="large"
                    className={classes.spacing}
                    disabled={disableTestState}
                    onClick={() => handleAbort()}
                  >
                    Clear
                  </Button>
                  <Button
                    type="submit"
                    variant="contained"
                    color="primary"
                    size="large"
                    onClick={() => submitProfile()}
                    className={classes.spacing}
                    disabled={disableTestState}
                    startIcon={<SaveOutlinedIcon />}
                  >
                    Save Profile
                  </Button>
                  <Button
                    type="submit"
                    variant="contained"
                    color="primary"
                    size="large"
                    onClick={handleSubmit}
                    className={classes.spacing}
                    disabled={
                      blockRunTestState ||
                      disableTestState ||
                      !CAN(keys.RUN_TEST.action, keys.RUN_TEST.subject)
                    }
                  >
                    {blockRunTestState ? <CircularProgress size={30} /> : 'Run Test'}
                  </Button>
                </div>
              </React.Fragment>

              {timerDialogOpenState ? (
                <div className={classes.centerTimer}>
                  <LoadTestTimerDialog
                    open={timerDialogOpenState}
                    t={tState}
                    onClose={handleTimerDialogClose}
                    countDownComplete={handleTimerDialogClose}
                  />
                </div>
              ) : null}

              {result && result.runner_results && (
                <div>
                  <Typography
                    variant="h6"
                    gutterBottom
                    className={classes.chartTitle}
                    id="timerAnchor"
                  >
                    Test Results
                    <IconButton
                      key="download"
                      aria-label="download"
                      color="inherit"
                      // onClick={() => self.props.closeSnackbar(key) }
                      href={`/api/perf/profile/result/${encodeURIComponent(result.meshery_id)}`}
                    >
                      <GetAppIcon style={iconMedium} />
                    </IconButton>
                  </Typography>
                  <div className={classes.chartContent} style={chartStyle}>
                    <MesheryChart
                      rawdata={[result && result.runner_results ? result : {}]}
                      data={[result && result.runner_results ? result.runner_results : {}]}
                    />
                  </div>
                </div>
              )}
            </div>
          </React.Fragment>

          {displayStaticCharts}

          {displayPromCharts}

          {displayGCharts}
        </>
      ) : (
        <DefaultError />
      )}
    </NoSsr>
  );
};
MesheryPerformanceComponent.propTypes = {
  classes: PropTypes.object.isRequired,
};

const mapDispatchToProps = (dispatch) => ({
  updateLoadTestData: bindActionCreators(updateLoadTestData, dispatch),
  updateStaticPrometheusBoardConfig: bindActionCreators(
    updateStaticPrometheusBoardConfig,
    dispatch,
  ),
  updateLoadTestPref: bindActionCreators(updateLoadTestPref, dispatch),
  updateProgress: bindActionCreators(updateProgress, dispatch),
});

const mapStateToProps = (state) => {
  const grafana = state.get('grafana').toJS();
  const prometheus = state.get('prometheus').toJS();
  const k8sConfig = state.get('k8sConfig');
  const staticPrometheusBoardConfig = state.get('staticPrometheusBoardConfig').toJS();
  const selectedK8sContexts = state.get('selectedK8sContexts');

  return {
    grafana,
    prometheus,
    staticPrometheusBoardConfig,
    k8sConfig,
    selectedK8sContexts,
  };
};

export default withStyles(styles)(
  connect(mapStateToProps, mapDispatchToProps)(withNotify(MesheryPerformanceComponent)),
);
