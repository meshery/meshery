import React from "react";
import { Alert, AlertTitle, Button, Stack } from "@material-ui/core";
import { ConnectAppContent, FinishContent, SelectRepositoryContent } from "./constants";
import Router, { useRouter } from "next/router";
import { useTheme } from "@material-ui/core"
import { useNotification } from '../../../utils/hooks/useNotification';
import { EVENT_TYPES } from '../../../lib/event-types';
import { VerifyHelmRepoLink, modifiedArrayOfObjects, postHelmInstall } from "./utils";
import StepperContent from "../stepper/stepper-content-wrapper";
import RJSFWrapper from '../../MesheryMeshInterface/PatternService/RJSF_wrapper';
// import { useGetSchemaQuery, useUpsertApplicationMutation } from "@/api/api";
import DesignsIcon from "../../../assets/icons/DesignIcon";
import { promisifiedDataFetch } from '../../../lib/data-fetch';

export const ConnectApp = ({ handleNext }) => {
  const [isConnected, setIsConnected] = React.useState(false);
  const formRef = React.createRef();
  const formStateRef = React.createRef();
  const [schema, setSchema] = React.useState({});
  const [isSchemaFetched, setIsSchemaFetched] = React.useState(false)
  // const { data: schema, isSuccess: isSchemaFetched } = useGetSchemaQuery({
  //   name: "helmRepo"
  // });
  const [formData, setFormData] = React.useState({});

  const handleCallback = () => {
    if (formRef.current && formRef.current.validateForm()) {
      // Only allow user to go on next step if this URL get verified
      //   axios.post("/api/integrations/helm/data", formData).then(()=>handleNext())
      handleNext();
    }
  };
  console.log("formRef", formRef);

  const cancelCallback = () => {
    Router.push("/dashboard");
  };
  const handleChange = data => {
    formStateRef.current = data;
    setFormData(data);
  };

  React.useEffect(()=>{
    let isMounted = true;

    const fetchData = async () => {
      try {
        const response = await promisifiedDataFetch(`/api/schema/resource/helmRepo`, {
          method: 'GET',
          credentials: 'include',
        });

        // Check if the component is still mounted before setting the state
        if (isMounted) {
          setSchema(response);
          setIsSchemaFetched(true)
        }
      } catch (error) {
        console.error('Error fetching schema in connection wizard', error);
      }
    };

    fetchData();

    // Cleanup function
    return () => {
      isMounted = false;
    };
  }, [])

  return (
    <StepperContent
      {...ConnectAppContent}
      handleCallback={handleCallback}
      cancelCallback={cancelCallback}
      disabled={Object.keys(formData).length !== 3}
      btnText={isConnected ? "Update" : "Connect"}
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
          // setErrors={setErrors}
        />
      )}
    </StepperContent>
  );
};

export const SelectRepository = ({ handleNext }) => {
  const [availableRepos, setAvailableRepos] = React.useState([]);
  const [loading, setLoading] = React.useState(false);
  const [page, setPage] = React.useState(0);
  const { notify } = useNotification();
  const formRef = React.createRef();
  // const [upsertApplication] = useUpsertApplicationMutation();
  const [formState, setFormState] = React.useState({});

  const getHelmRepository = async () => {
    setLoading(true);
    try {
      //   const { data } = await axios.get(
      //     process.env.API_ENDPOINT_PREFIX + HELM_INTEGRATION_ENDPOINT
      //   );
      // dummy data
      let data = [
        {
          name: "oci chart",
          charData: "/ yaml data /"
        },
        {
          name: "oci backend",
          charData: "/ yaml data /"
        },
        {
          name: "Meshery chart",
          charData: "/ yaml data /"
        },
        {
          name: "Meshery cloud chart",
          charData: "/ yaml data /"
        }
      ];

      if (data.length === 0) {
        setHasMore(false);
        setLoading(false);
        return;
      }
      const helmRepoNames = data.map(helmRepo => helmRepo.name);
      let updatedRepositories = [...helmRepoNames];
      setAvailableRepos(updatedRepositories);
      setLoading(false);
    } catch (err) {
      setLoading(false);
      notify({messgae: err, type: EVENT_TYPES.ERROR});
    }
  };

  React.useEffect(() => {
    getHelmRepository();
  }, [page]);

  const handleCallback = async () => {
    if (formRef.current && formRef.current.validateForm()) {
      handleNext(); // Proceed to the next step
      // try {
      //   setLoading(true);
      //   const formData = formStateRef.current;
      //   const selectedRepos = formData.selectedHelmRepos;

      //   // Filter the selected repositories from the data array
      //   const selectedReposWithChartData = data.filter(repo =>
      //     selectedRepos.includes(repo.name.toLowerCase())
      //   );

      //   // write logic if we will create application of selected repo from frontend
      //   setLoading(false);

      // } catch (error) {
      //   setLoading(false);
      //   handleError(error);
      // }
    }
  };

  const schema = {
    properties: {
      selectedHelmRepos: {
        description: "Select one or more Helm charts from the available options",
        items: {
          enum: availableRepos,
          type: "string"
        },
        minItems: 1,
        title: "Select one or more helm chart",
        type: "array",
        uniqueItems: true,
        "x-rjsf-grid-area": 12
      }
    },
    required: ["selectedHelmRepos"],
    type: "object"
  };

  const handleChange = data => {
    setFormState(data);
  };

  return (
    <StepperContent
      {...SelectRepositoryContent}
      handleCallback={handleCallback}
      disabled={!formState}
    >
      <RJSFWrapper
        key="select-helm-repo-rjsf-form"
        formData={formState}
        onChange={handleChange}
        jsonSchema={schema}
        liveValidate={false}
        formRef={formRef}
        // setErrors={setErrors}
      />
    </StepperContent>
  );
};

const HelmChartsStatus = ({ configuredRepository, handleInstall }) => {
  const router = useRouter();
  const theme = useTheme();
  const meshMapRedirect = id => {
    let letHyperLink = `https://playground.meshery.io/extension/meshmap?application=${id}`;
    router.replace(letHyperLink);
  };
  return (
    <Stack
      style={{ width: "100%", marginBottom: "2rem", height: "21rem", overflow: "auto" }}
      spacing={2}
    >
      {configuredRepository?.map((item, index) => {
        const { name, status, id } = item;
        const isSuccess = status === "created";
        return (
          <Alert
            key={index}
            icon={false}
            severity={isSuccess ? "success" : "error"}
            // color: isSuccess ? theme.palette.success.dark : theme.palette.error.dark
            style={{
              display: "flex",
              alignItems: "center",
              lineHeight: "0",
              color: "#000",
              backgroundColor: isSuccess ? theme.palette.keppelGreenLight : null
            }}
            action={
              !isSuccess ? (
                <Button
                  style={{ backgroundColor: theme.palette.error.main, color: theme.palette.white }}
                  size="small"
                  onClick={() => handleInstall(item)}
                >
                  Retry
                </Button>
              ) : (
                <Button
                  style={{ color: theme.palette.keppelGreen }}
                  size="small"
                  onClick={() => meshMapRedirect(id)}
                >
                  <DesignsIcon
                    width="28"
                    height="28"
                    fill={theme.palette.keppelGreen}
                  />
                  <p style={{marginLeft: "0.2rem"}}>
                      Open Design
                  </p>
                </Button>
              )
            }
          >
            <AlertTitle style={{ fontSize: "1.4rem" }}>{name}</AlertTitle>
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
        );
      })}
    </Stack>
  );
};

export const Finish = ({
  configuredRepository,
  setConfiguredRepository,
  installationId
}) => {
  const router = useRouter();
  const { notify } = useNotification();
  const [chartCount, setChartCount] = React.useState({
    importedChartCount: 0,
    totalChartCount: configuredRepository.length
  });

  const handleCallback = () => {
    router.push("/connections");
  };

  let statusData = [
    {
      id: "uuid-uuid-uuid",
      name: "oci chart",
      status: "created",
      chart_yaml: "/ yaml data /"
    },
    {
      id: "uuid-uuid-uuid",
      name: "oci chart",
      status: "created",
      chart_yaml: "/ yaml data /"
    },
    {
      id: "uuid-uuid-uuid",
      name: "oci chart",
      status: "failed",
      chart_yaml: "/ yaml data /"
    }
  ];

  React.useEffect(() => {
    let importedChartCount = statusData.filter(item => item.status === "created").length;
    setChartCount({
      importedChartCount: importedChartCount,
      totalChartCount: statusData.length
    });
    setConfiguredRepository(statusData);
  }, []);

  //   React.useEffect(()=>{
  // write logic to update chart count
  //   }, [])

  const handleInstall = repository => () => {
    const body = {
      repositories: [repository]
    };
    postHelmInstall(JSON.stringify(body), installationId)
      .then(res => {
        // update configure repos based on new updated repo
        // setConfiguredRepository(modifiedArrayOfObjects(configuredRepository, res));
      })
      .catch(error => {
        notify({message: error, type: EVENT_TYPES.ERROR});
      });
  };

  return (
    <StepperContent
      {...FinishContent}
      title={`${chartCount.importedChartCount} of ${configuredRepository.length} charts imported`}
      handleCallback={handleCallback}
    >
      <HelmChartsStatus
        configuredRepository={configuredRepository}
        handleInstall={handleInstall}
      />
    </StepperContent>
  );
};
