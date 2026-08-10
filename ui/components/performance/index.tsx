import React, { useState, useEffect } from 'react';
import { ModalBody, ModalFooter, NoSsr, useHasPermission } from '@sistent/sistent';
import { URLValidator } from '../../utils/URLValidator';
import { isValidJSON } from '../../utils/validators';
import LoadTestTimerDialog from '../load-test-timer-dialog';
import fetchControlPlanes from '@/graphql/queries/ControlPlanesQuery';
import { ctxUrl, getK8sClusterIdsFromCtxId } from '../../utils/multi-ctx';
import { useNotification } from '../../utils/hooks/useNotification';
import { EVENT_TYPES } from '../../lib/event-types';
import { generateTestName, generateUUID } from './helper';

import { Keys } from '@meshery/schemas/permissions';
import DefaultError from '@/components/general/error-404/index';
import { api } from '../../rtk-query';
import { useGetUserPrefWithContextQuery } from '@/rtk-query/user';
import { useSavePerformanceProfileMutation } from '@/rtk-query/performance-profile';
import { useGetMeshQuery } from '@/rtk-query/mesh';
import { CenterTimer } from './style';
import { getMeshModels } from '@/api/meshmodel';
import { useDispatch, useSelector } from 'react-redux';
import { updateProgress } from '@/store/slices/mesheryUi';
import { updateLoadTest } from '@/store/slices/prefTest';
import { normalizeLoadTestPrefs } from '../../lib/load-test-prefs';
import PerformanceForm from './PerformanceForm';
import PerformanceFormActions from './PerformanceFormActions';
import PerformanceTestResults from './PerformanceTestResults';
import { generatePerformanceProfile } from './performance-helpers';
import { createPerformanceFormChangeHandler } from './performance-handlers';

// Re-export the helper so external/legacy callers that imported it from the
// entry point keep working after the Phase 5.a split.
export { generatePerformanceProfile } from './performance-helpers';

let eventStream = null;
const PERFORMANCE_RTK_TAG = 'Performance_Profile_performance';

const MesheryPerformanceComponent_ = (props) => {
  const hasPermission = useHasPermission(Keys.PerformanceManagementViewPerformanceProfiles);
  const {
    testName = '',
    meshName = '',
    url = '',
    qps = '0',
    c = '0',
    t = '30s',
    performanceProfileID,
    profileName,
    loadGenerator,
    additional_options,
    headers,
    cookies,
    reqBody,
    contentType,
    metadata,
    closeModal,
  } = props;
  // Create individual state variables for each property
  const [testNameState, setTestName] = useState(testName);
  const [meshNameState, setMeshName] = useState(meshName);
  const [meshModels, setMeshModels] = useState([]);
  const [urlState, setUrl] = useState(url);
  const [qpsState, setQps] = useState(qps);
  const [cState, setC] = useState(c);
  const [tState, setT] = useState(t);
  const [tValueState, setTValue] = useState(t);
  const [loadGeneratorState, setLoadGenerator] = useState(loadGenerator || 'fortio');
  const [additionalOptionsState, setAdditionalOptions] = useState(additional_options || '');
  const [testResult, setTestResult] = useState();
  const [testResultsOpen, setTestResultsOpen] = useState(false);
  const { selectedK8sContexts, k8sConfig } = useSelector((state) => state.ui);
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
    !(URLValidator(urlState) || isValidJSON(additionalOptionsState)),
  );
  const [testUUIDState, setTestUUID] = useState(generateUUID());
  const [selectedMeshState, setSelectedMesh] = useState('');
  const [availableAdaptersState, setAvailableAdapters] = useState([]);
  const [availableSMPMeshesState, setAvailableSMPMeshes] = useState([]);
  const [metadataState, setMetadata] = useState(metadata);
  const { notify } = useNotification();
  const dispatch = useDispatch();
  const { data: userData, isSuccess: isUserDataFetched } =
    useGetUserPrefWithContextQuery(selectedK8sContexts);

  const [savePerformanceProfile] = useSavePerformanceProfileMutation();
  const {
    data: smpMeshes,
    isSuccess: isSMPMeshesFetched,
    isError: isSMPMeshError,
  } = useGetMeshQuery();

  useEffect(() => {
    const fetchMeshModels = async () => {
      const { models } = await getMeshModels();
      const modelNames = models.map((model) => model.displayName);
      setMeshModels(modelNames);
    };
    fetchMeshModels();
  }, []);

  const handleChange = createPerformanceFormChangeHandler({
    setProfileName,
    setMeshName,
    setC,
    setQps,
    setHeaders,
    setCookies,
    setContentType,
    setReqBody,
    setLoadGenerator,
    setUrl,
    setCaCertificate,
    setAdditionalOptions,
    setJsonError,
    setDisableTest,
    setUrlError,
  });

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
    } catch {
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
      closeModal && closeModal();
      return;
    }

    if (!performanceProfileIDState) {
      submitProfile(({ id }) => submitLoadTest(id));
      closeModal && closeModal();
      return;
    }
    submitLoadTest(performanceProfileIDState);
    closeModal && closeModal();
  };

  const submitProfile = (cb) => {
    const profile = generatePerformanceProfile({
      name: profileNameState,
      loadGenerator: loadGeneratorState,
      additional_options: additionalOptionsState,
      endpoint: urlState,
      serviceMesh: meshNameState,
      concurrentRequest: +cState || 1,
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
    if (generateNotif) updateProgress({ showProgress: true });
    savePerformanceProfile({ body: body })
      .unwrap()
      .then((result) => {
        if (result) {
          updateProgress({ showProgress: false });
          setPerformanceProfileID(result.id);
          if (cb) cb(result);
          if (generateNotif) {
            notify({
              message: `Performance profile ${result.name} has been created`,
              event_type: EVENT_TYPES.SUCCESS,
            });
          }
        }
      })
      .catch((err) => {
        console.error(err);
        updateProgress({ showProgress: false });
        notify({
          message: 'Failed to create performance profile',
          event_type: EVENT_TYPES.ERROR,
          details: err.toString(),
        });
      });
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
      ctxUrl(`/api/user/performance/profiles/${id}/run`, selectedK8sContexts) + '&cert=true';
    startEventStream(`${runURL}${selectedK8sContexts?.length > 0 ? '&' : '?'}${params}`);
    setBlockRunTest(true); // to block the button
  };

  function handleSuccess() {
    return (result) => {
      if (
        typeof result !== 'undefined' &&
        typeof (result.runnerResults ?? result.runner_results) !== 'undefined'
      ) {
        notify({
          message: 'fetched the data.',
          event_type: EVENT_TYPES.SUCCESS,
          dataTestID: 'notify-fetch-data',
        });
        dispatch(
          updateLoadTest({
            loadTest: {
              testName: testNameState,
              meshName: meshNameState,
              url: urlState,
              qps: qpsState,
              c: cState,
              t: tState,
              loadGenerator: loadGeneratorState,
              result: result,
            },
          }),
        );
        setTestUUID(generateUUID());

        dispatch(api.util.invalidateTags([PERFORMANCE_RTK_TAG]));
        setTestResultsOpen(true);
        setTestResult(result);
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
    notify({
      message: 'Load test has been submitted',
      event_type: EVENT_TYPES.SUCCESS,
    });
  }

  function handleEvents() {
    let track = 0;
    return (e) => {
      const data = JSON.parse(e.data);

      switch (data.status) {
        case 'info':
          notify({ message: data.message, event_type: EVENT_TYPES.INFO });
          if (track === 0) {
            setTimerDialogOpen(true);
            // setResult({});
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
    scanForMeshes();
    getLoadTestPrefs();
    getSMPMeshes();
    if (props.runTestOnMount) handleSubmit();
  }, [userData, isUserDataFetched, smpMeshes]);

  const getLoadTestPrefs = () => {
    if (!isUserDataFetched || !userData) return;
    if (props.performanceProfileID) return;

    const loadTestPrefs = normalizeLoadTestPrefs(userData.loadTestPrefs);

    setQps(loadTestPrefs.qps);
    setC(loadTestPrefs.c);
    setT(loadTestPrefs.t);
    setLoadGenerator(loadTestPrefs.gen);
  };

  const getK8sClusterIds = () => {
    return getK8sClusterIdsFromCtxId(selectedK8sContexts, k8sConfig);
  };

  const scanForMeshes = () => {
    if (typeof k8sConfig === 'undefined' || !k8sConfig.clusterConfigured) {
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
    if (isSMPMeshesFetched && smpMeshes) {
      setAvailableSMPMeshes(
        [...(smpMeshes.availableMeshes || [])].sort((m1, m2) => m1.localeCompare(m2)),
      );
    } else if (isSMPMeshError) {
      handleError('unable to fetch SMP meshes');
    }
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
  let chartStyle = {};
  if (timerDialogOpenState) {
    chartStyle = { opacity: 0.3 };
  }

  availableAdaptersState.forEach((item) => {
    const index = availableSMPMeshesState.indexOf(item);
    if (index !== -1) availableSMPMeshesState.splice(index, 1);
  });

  if (testResultsOpen) {
    return (
      <PerformanceTestResults
        testResult={testResult}
        chartStyle={chartStyle}
        onBack={() => setTestResultsOpen(false)}
      />
    );
  }

  return (
    <NoSsr>
      {hasPermission ? (
        <>
          <React.Fragment>
            <ModalBody>
              <PerformanceForm
                profileName={profileNameState}
                meshName={meshNameState}
                selectedMesh={selectedMeshState}
                meshModels={meshModels}
                url={urlState}
                urlError={urlErrorState}
                c={cState}
                qps={qpsState}
                t={tState}
                tValue={tValueState}
                tError={tErrorState}
                headers={headersState}
                cookies={cookiesState}
                contentType={contentTypeState}
                reqBody={reqBodyState}
                additionalOptions={additionalOptionsState}
                jsonError={jsonErrorState}
                caCertificate={caCertificateState}
                metadata={metadataState}
                loadGenerator={loadGeneratorState}
                handleChange={handleChange}
                handleDurationChange={handleDurationChange}
                handleInputDurationChange={handleInputDurationChange}
                handleCertificateUpload={handleCertificateUpload}
              />
            </ModalBody>
            <ModalFooter variant="filled">
              <React.Fragment>
                <PerformanceFormActions
                  disableTest={disableTestState}
                  blockRunTest={blockRunTestState}
                  hasTestResult={!!testResult}
                  onAbort={handleAbort}
                  onShowResults={() => setTestResultsOpen(true)}
                  onSaveProfile={() => submitProfile()}
                  onRunTest={handleSubmit}
                />
              </React.Fragment>
            </ModalFooter>

            {timerDialogOpenState ? (
              <CenterTimer>
                <LoadTestTimerDialog
                  open={timerDialogOpenState}
                  t={tState}
                  onClose={handleTimerDialogClose}
                  countDownComplete={handleTimerDialogClose}
                />
              </CenterTimer>
            ) : null}
          </React.Fragment>
        </>
      ) : (
        <DefaultError />
      )}
    </NoSsr>
  );
};

export const MesheryPerformanceComponentWithStyles = MesheryPerformanceComponent_;

export const MesheryPerformanceComponent = (props) => {
  return <MesheryPerformanceComponentWithStyles {...props} />;
};

export default MesheryPerformanceComponent;
