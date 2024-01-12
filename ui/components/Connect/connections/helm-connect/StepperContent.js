import React from 'react';
import { Button, Grid, Box } from '@material-ui/core';
import { Alert, AlertTitle } from '@material-ui/lab';
import { ConnectAppContent, FinishContent, SelectRepositoryContent } from './constants';
import Router, { useRouter } from 'next/router';
import { useTheme } from '@material-ui/core';
import { useNotification } from '../../../../utils/hooks/useNotification';
import { EVENT_TYPES } from '../../../../lib/event-types';
import { generateSelectedHelmRepo, addStatusToCharts } from './utils';
import {
  useGetConnectionStatusQuery,
  useLazyGetConnectionDetailsQuery,
  useVerifyConnectionURLMutation,
  useConnectionMetaDataMutation,
  useConfigureConnectionMutation,
} from '../../../../rtk-query/connection';
import { useGetSchemaQuery } from '../../../../rtk-query/schema';
import StepperContent from '../../stepper/StepperContentWrapper';
import RJSFWrapper from '../../../MesheryMeshInterface/PatternService/RJSF_wrapper';
import debounce from 'lodash/debounce';
import CheckCircleIcon from '@material-ui/icons/CheckCircle';
import { selectCompSchema } from '../../../RJSFUtils/common';

export const ConnectApp = ({ handleNext, setSharedData }) => {
  const [isConnected, setIsConnected] = React.useState(false);
  const formRef = React.createRef();
  const [formData, setFormData] = React.useState({});
  const [isUrlValid, setIsUrlValid] = React.useState(false);
  const [extraErrors, setExtraErrors] = React.useState();
  const { notify } = useNotification();

  /**
   * RTK queries and mutations
   */
  const [triggerGetConnectionDetails] = useLazyGetConnectionDetailsQuery({
    connectionKind: 'helm',
    repoURL: '',
  }); //using lazy fetching
  const [verifyConnectionURL] = useVerifyConnectionURLMutation();
  const { data: schema, isSuccess: isSchemaFetched } = useGetSchemaQuery({
    schemaName: 'helmRepo',
  });
  const { data: connectionStatusData, isSuccess: isConnectionStatusDataFetched } =
    useGetConnectionStatusQuery('helm');

  const handleCallback = async () => {
    setExtraErrors({});
    if (formRef.current && formRef.current.validateForm() && isUrlValid) {
      // Check if we can fetch charts from this URL
      // URL is already verified if this function is running
      // If payload doesn't contain any chart don't allow user to go on next step
      handleNext();
      triggerGetConnectionDetails({
        connectionKind: 'helm',
        repoURL: formData.url,
      })
        .unwrap()
        .then((payload) => {
          if (payload) {
            setSharedData((prev) => {
              if (prev) {
                return {
                  ...prev,
                  repoURL: formData.url,
                  helmRepoChartsData: payload,
                };
              } else {
                return {
                  repoURL: formData.url,
                  helmRepoChartsData: payload,
                };
              }
            });
            handleNext();
          } else {
            setExtraErrors({
              url: {
                __errors: ['This Helm repo does not contain charts'],
              },
            });
          }
        })
        .catch((error) => {
          // TODO: use notify here
          console.error('Error coming when fetching connection details', error);
        });
    }
  };

  const cancelCallback = () => {
    Router.push('/management/connections');
  };

  const debouncedValidateUrl = debounce(async (url) => {
    verifyConnectionURL({
      connectionKind: 'helm',
      repoURL: url,
    })
      .unwrap()
      .then((payload) => {
        if (payload) {
          setExtraErrors({});
          setIsUrlValid(true);
        } else {
          setExtraErrors({
            url: {
              __errors: ['URL is not reachable, please provide valid URL'],
            },
          });
        }
      })
      .catch((error) => {
        notify({
          message: error.data
            ? JSON.stringify(error.data)
            : `${error.status} error coming, please check console for more info`,
          type: EVENT_TYPES.ERROR,
        });
      });
  }, 400);

  const handleChange = (data) => {
    setFormData(data);
    if (formData.url) {
      debouncedValidateUrl(formData.url);
    }
  };

  React.useEffect(() => {
    if (isConnectionStatusDataFetched) {
      if (connectionStatusData) {
        setIsConnected(true);
      }
    }
  }, [isConnectionStatusDataFetched]);

  return (
    <StepperContent
      {...ConnectAppContent}
      handleCallback={handleCallback}
      cancelCallback={cancelCallback}
      disabled={Object.keys(formData).length !== 3 || !isUrlValid}
      btnText={isConnected ? 'Update' : 'Connect'}
    >
      {isSchemaFetched && (
        <RJSFWrapper
          key="connect-helm-repo-rjsf-form"
          formData={formData}
          onChange={handleChange}
          uiSchema={schema.uiSchema}
          jsonSchema={schema.rjsfSchema}
          liveValidate={false}
          formRef={formRef}
          extraErrors={extraErrors}
        />
      )}
    </StepperContent>
  );
};

export const SelectRepository = ({ handleNext, sharedData, setSharedData }) => {
  const [availableRepos, setAvailableRepos] = React.useState([]);
  const [loading, setLoading] = React.useState(false);
  const { notify } = useNotification();
  const formRef = React.createRef();
  const [formState, setFormState] = React.useState({});
  const [selectedCharts, setSelectedCharts] = React.useState([]);
  const [connectionMetaData] = useConnectionMetaDataMutation();

  const getHelmRepository = async () => {
    setLoading(true);
    try {
      const formData = [];
      sharedData.helmRepoChartsData.forEach((helmRepo) => {
        formData.push(helmRepo.Name);

        // Check if the Helm repo has sub charts (dependencies)
        if (helmRepo.Dependencies && helmRepo.Dependencies.length > 0) {
          helmRepo.Dependencies.forEach((dependency) => {
            formData.push(`${helmRepo.Name}/${dependency.name}`);
          });
        }
      });

      setAvailableRepos(formData);
      setLoading(false);
    } catch (error) {
      setLoading(false);
      notify({
        message: error,
        type: EVENT_TYPES.ERROR,
      });
    }
  };

  React.useEffect(() => {
    getHelmRepository();
  }, []);

  const handleCallback = async () => {
    if (formRef.current && formRef.current.validateForm()) {
      connectionMetaData({
        connectionKind: 'helm',
        body: {
          repoURL: sharedData.repoURL,
          charts: selectedCharts,
        },
      })
        .unwrap()
        .then((payload) => {
          // payload should be updated connection object
          if (payload) {
            setSharedData((prev) => {
              return {
                ...prev,
                selectedCharts: selectedCharts,
                selectedChartsTotalCount: formState.selectedHelmRepos.length, // Helps to not loop again in next step
              };
            });
            handleNext();
          }
        })
        .catch((error) => {
          notify({
            message: error.data
              ? JSON.stringify(error.data)
              : `${error.status} error coming, please check console for more info`,
            type: EVENT_TYPES.ERROR,
          });
        });
    }
  };

  const schema = selectCompSchema(
    availableRepos,
    'Select one or more Helm charts from the available options',
    'Select one or more of your Helm charts',
    'selectedHelmRepos',
    true,
  );

  const handleChange = (data) => {
    const { updatedData, arrayOfSelected } = generateSelectedHelmRepo(data, sharedData);
    setFormState(updatedData);
    setSelectedCharts(arrayOfSelected);
  };

  return (
    <StepperContent
      {...SelectRepositoryContent}
      handleCallback={handleCallback}
      disabled={!formState || loading}
    >
      <RJSFWrapper
        key="select-helm-repo-rjsf-form"
        formData={formState}
        onChange={handleChange}
        jsonSchema={schema}
        liveValidate={false}
        formRef={formRef}
      />
    </StepperContent>
  );
};

const HelmChartsStatus = ({ configuredRepository, handleInstall }) => {
  const theme = useTheme();
  const renderChart = (chart) => {
    const { name, status } = chart;
    const isSuccess = status === true;

    return (
      <Grid item key={name} style={{ width: '100%' }}>
        <Alert
          icon={false}
          severity={isSuccess ? 'success' : 'error'}
          style={{
            display: 'flex',
            alignItems: 'center',
            lineHeight: '0',
            color: '#000',
            backgroundColor: isSuccess ? theme.palette.keppelGreenLight : null,
          }}
          action={
            !isSuccess ? (
              <Button
                style={{
                  backgroundColor: theme.palette.error.main,
                  color: theme.palette.secondary.primaryTextModal,
                }}
                size="small"
                onClick={() => handleInstall(chart)}
              >
                Retry
              </Button>
            ) : (
              <Box style={{ color: '#00B39F', display: 'flex', alignItems: 'center' }}>
                <CheckCircleIcon size="small" style={{ marginRight: '0.3rem' }} />
                Installed
              </Box>
            )
          }
        >
          <AlertTitle style={{ fontSize: '1.4rem' }}>{name}</AlertTitle>
          {isSuccess ? (
            <div>
              <p>{`Design ${name} created`}</p>
            </div>
          ) : (
            <div>
              <p>{`Unable to import ${name}.yaml`}</p>
            </div>
          )}
        </Alert>
      </Grid>
    );
  };

  return (
    <Grid
      container
      spacing={2}
      style={{ width: '100%', marginBottom: '2rem', maxHeight: '21rem', overflow: 'auto' }}
    >
      {configuredRepository?.map((chart) => {
        return renderChart(chart);
      })}
    </Grid>
  );
};

export const Finish = ({ sharedData }) => {
  const router = useRouter();
  const { notify } = useNotification();
  const [unDigestedCharts, setUnDigestedCharts] = React.useState([]);
  const handleCallback = () => {
    router.push('/connections');
  };

  const [configureConnection] = useConfigureConnectionMutation();

  const configureCharts = () => {
    configureConnection({
      connectionKind: 'helm',
      body: sharedData.repoURL,
    })
      .unwrap()
      .then((payload) => {
        // payload charts are which didn't got digested
        if (payload) {
          setUnDigestedCharts(payload);
        }
      })
      .catch((error) => {
        console.error(error);
        notify({
          message: error.data
            ? JSON.stringify(error.data)
            : `${error.status} error coming, please check console for more info`,
          type: EVENT_TYPES.ERROR,
        });
      });
  };
  const updatedConfiguredRepository = addStatusToCharts(
    sharedData.selectedCharts,
    unDigestedCharts,
  );

  React.useEffect(() => {
    configureCharts();
  }, []);

  const handleInstall = () => {
    // for retry
  };

  return (
    <StepperContent
      {...FinishContent}
      title={`${
        unDigestedCharts
          ? sharedData.selectedChartsTotalCount - unDigestedCharts?.length
          : sharedData.selectedChartsTotalCount
      } of ${sharedData.selectedChartsTotalCount} charts imported`}
      handleCallback={handleCallback}
    >
      <HelmChartsStatus
        configuredRepository={updatedConfiguredRepository}
        handleInstall={handleInstall}
      />
    </StepperContent>
  );
};
