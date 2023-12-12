import React, { useEffect, useState } from 'react';
import OutlinedInput from '@mui/material/OutlinedInput';
import InputLabel from '@mui/material/InputLabel';
import FormControl from '@mui/material/FormControl';
import {
  Checkbox,
  MenuItem,
  ListItemText,
  Select,
  Typography,
  Grid,
  Button,
} from '@material-ui/core';

import { ConnectionDetailContent, FinishContent, CredentialDetailContent } from './constants';
import StepperContent from './StepperContentWrapper';
import RJSFWrapper from '../../../MesheryMeshInterface/PatternService/RJSF_wrapper';
import dataFetch from '../../../../lib/data-fetch';
import { Box } from '@mui/material';

export const ConnectionDetails = ({ sharedData, setSharedData, handleNext }) => {
  const [schema, setSchema] = useState({ properties: {} });
  const [connectionObject, setConnectionObject] = useState(null);
  const formRef = React.createRef();

  useEffect(() => {
    sharedData !== null &&
      sharedData.connection === undefined &&
      registerConnection(sharedData?.model);
  }, [sharedData]);

  useEffect(() => {
    ConnectionDetailContent.title = `Connecting to ${connectionObject?.component?.metadata?.modelDisplayName}`;
  }, [connectionObject]);

  const registerConnection = (model) => {
    dataFetch(
      '/api/integrations/connections/register',
      {
        method: 'POST',
        credentials: 'include',
        body: JSON.stringify({
          application_data: model,
          status: 'initialize',
        }),
      },
      (result) => {
        setConnectionObject(result);
        setSharedData((prevState) => ({
          ...prevState,
          connection: result,
        }));
        const componentSchema = result?.component?.schema;
        if (componentSchema !== null || componentSchema !== undefined) {
          const schemaObj = JSON.parse(componentSchema);
          setSchema(schemaObj);
        }
      },
    );
  };

  const handleCallback = () => {
    handleNext();
  };

  const handleChange = (data) => {
    setSharedData((prevState) => ({
      ...prevState,
      componentForm: data,
    }));
  };
  const isDisabledNextButton =
    sharedData?.componentForm &&
    sharedData?.componentForm['name'] !== undefined &&
    sharedData?.componentForm &&
    sharedData?.componentForm['url'] !== undefined
      ? false
      : true;
  return (
    <StepperContent
      {...ConnectionDetailContent}
      handleCallback={handleCallback}
      disabled={isDisabledNextButton}
    >
      <RJSFWrapper
        key="register-connection-rjsf-form"
        jsonSchema={schema}
        liveValidate={false}
        formRef={formRef}
        onChange={handleChange}
      />
    </StepperContent>
  );
};

export const CredentialDetails = ({ sharedData, handleNext }) => {
  const [schema, setSchema] = useState({ properties: {} });
  const [connectionObject, setConnectionObject] = useState(null);
  const [existingCredentials, setExistingCredentials] = useState([]);
  const [selectedCredential, setSelectedCredential] = useState(null);
  const [prevSelectedCredential, setPrevSelectedCredential] = useState(null);
  const [formState, setFormState] = React.useState(null);
  const [skipCredentialVerification, setSkipCredentialVerification] = useState(false);
  const [disableVerify, setDisableVerify] = useState(true);
  const [isSuccess, setIsSuccess] = React.useState(null);
  const formRef = React.createRef();

  useEffect(() => {
    credentialSchema();
    getchExistingCredential();
  }, []);

  const credentialSchema = () => {
    setConnectionObject(sharedData?.connection);
    const credentialSchema = sharedData?.connection?.credential?.schema;
    if (credentialSchema !== null || credentialSchema !== undefined) {
      const schemaObj = JSON.parse(credentialSchema);
      setSchema(schemaObj);
    }
  };

  useEffect(() => {
    CredentialDetailContent.title = `Credential for ${connectionObject?.credential?.metadata?.modelDisplayName}`;
  }, [connectionObject]);

  const getchExistingCredential = () => {
    dataFetch(
      '/api/integrations/credentials',
      {
        method: 'GET',
        credentials: 'include',
      },
      (result) => {
        setExistingCredentials(result?.credentials);
      },
    );
  };

  const verifyConnection = () => {
    dataFetch(
      '/api/integrations/connections/register',
      {
        method: 'POST',
        credentials: 'include',
        body: JSON.stringify({
          skip_credential_verification: skipCredentialVerification,
          kind: sharedData?.kind, // this is "kind" column of the current row which is selected in the meshsync table. i.e. the entry against which registration process has been invoked.
          name: sharedData?.connection?.component?.displayName, // This name is from the name field in schema
          type: sharedData?.connection?.component?.model?.category?.name?.toLowerCase(),
          sub_type: sharedData?.connection?.component?.metadata?.subCategory.toLowerCase(),
          metadata: sharedData?.componentForm,
          credential_secret: selectedCredential !== null ? selectedCredential : formState,
          id: sharedData?.connection?.id,
          status: 'register',
        }),
      },
      (result) => {
        if (result === '') {
          setIsSuccess(true);
          connectToConnection();
        } else {
          setIsSuccess(false);
        }
      },
    );
  };

  const connectToConnection = () => {
    dataFetch(
      '/api/integrations/connections/register',
      {
        method: 'POST',
        credentials: 'include',
        body: JSON.stringify({
          kind: sharedData?.kind, // this is "kind" column of the current row which is selected in the meshsync table. i.e. the entry against which registration process has been invoked.
          name: sharedData?.connection?.component?.displayName, // This name is from the name field in schema
          type: sharedData?.connection?.component?.model?.category?.name?.toLowerCase(),
          sub_type: sharedData?.connection?.component?.metadata?.subCategory.toLowerCase(),
          metadata: sharedData?.componentForm,
          credential_secret: selectedCredential !== null ? selectedCredential : formState,
          id: sharedData?.connection?.id,
          status: 'connect',
        }),
      },
      (result) => {
        if (result === '') {
          setIsSuccess(true);
        } else {
          setIsSuccess(false);
        }
      },
    );
  };
  const cancelConnectionRegister = () => {
    dataFetch(
      '/api/integrations/connections/register',
      {
        method: 'DELETE',
        credentials: 'include',
        body: JSON.stringify({
          id: sharedData?.connection?.id,
        }),
      },
      (result) => {
        console.log(result);
      },
    );
  };

  const handleCallback = () => {
    if (isSuccess === null || isSuccess === false) {
      verifyConnection();
    } else {
      handleNext();
    }
  };

  const cancelCallback = () => {
    cancelConnectionRegister();
  };

  const handleSelectCredential = (e) => {
    const id = e.target.value;
    const credential = existingCredentials.find((credential) => credential.id === id);
    setSelectedCredential(credential);
    setPrevSelectedCredential(id);
  };

  const handleChange = (data) => {
    setFormState(data);
  };

  const handleClose = () => {
    if (prevSelectedCredential === selectedCredential?.id) {
      setSelectedCredential(null);
      setPrevSelectedCredential(null);
    }
  };

  useEffect(() => {
    if (
      selectedCredential !== null ||
      (formState !== null && formState['API key']) !== undefined ||
      (formState !== null && formState['Basic Auth']) !== undefined
    ) {
      setDisableVerify(false);
    } else {
      setDisableVerify(true);
    }
  }, [selectedCredential, formState]);

  return (
    <StepperContent
      {...CredentialDetailContent}
      handleCallback={handleCallback}
      cancelCallback={cancelCallback}
      disabled={disableVerify}
      btnText={isSuccess === null || isSuccess === false ? 'Verify Connection' : 'Next'}
    >
      <p className={{ paddingLeft: '16px' }}>
        Select an existing credential to use for this connection
      </p>
      <FormControl sx={{ width: '100%' }} size="small">
        <InputLabel id="credential-checkbox-label">Select existing credential</InputLabel>
        <Select
          labelId="credential-checkbox-label"
          id="credential-checkbox"
          onChange={handleSelectCredential}
          value={selectedCredential?.name}
          onClose={handleClose}
          input={<OutlinedInput label="Select existing credential" />}
          renderValue={() => (
            <div>{selectedCredential !== null ? selectedCredential.name : ''}</div>
          )}
          MenuProps={{
            anchorOrigin: {
              vertical: 'bottom',
              horizontal: 'left',
            },
            transformOrigin: {
              vertical: 'top',
              horizontal: 'left',
            },
            getContentAnchorEl: null,
            style: {
              maxHeight: 48 * 4.5 + 8,
              width: 250,
              zIndex: 10000,
            },
            PaperProps: {
              style: {
                zIndex: 10000,
              },
            },
          }}
        >
          {existingCredentials &&
            existingCredentials?.map((credential) => (
              <MenuItem value={credential.id} name={credential.name}>
                <Checkbox checked={selectedCredential?.id === credential.id} />
                <ListItemText primary={credential.name} />
              </MenuItem>
            ))}
        </Select>
      </FormControl>
      <p style={{ display: 'flex', justifyContent: 'center' }}>-OR-</p>
      <p>Configure a new credential to use for this connection</p>
      <RJSFWrapper
        key="register-connection-rjsf-form"
        jsonSchema={schema}
        liveValidate={false}
        formRef={formRef}
        disabled={selectedCredential !== null ? true : false}
        onChange={handleChange}
      />
      <Box
        style={{
          background: 'rgba(0, 211, 169, 0.05)',
          padding: '0.6rem',
          margin: '2rem 0',
        }}
      >
        <Typography style={{ fontSize: '1rem' }}>
          <Checkbox
            color="success"
            onChange={(e) => {
              setSkipCredentialVerification(e.target.checked);
              setDisableVerify(!e.target.checked);
            }}
          />
          <span>Force connection registration bypass verification</span>
        </Typography>
      </Box>
      {isSuccess === false && (
        <div
          style={{
            background: '#ff000010',
            borderRadius: '10px',
            padding: '10px',
            display: 'flex',
          }}
        >
          <Grid style={{ width: '80%' }}>
            <Typography>
              <b>Credential Invalid</b>
            </Typography>
            <Typography sx={{ color: '#00000020' }}>
              {`Unable to establish a connection using ${connectionObject?.credential?.metadata?.modelDisplayName}`}
            </Typography>
          </Grid>
          <Grid style={{ width: '20%', display: 'flex', justifyContent: 'flex-end' }}>
            <Button
              style={{
                backgroundColor: '#ff0000',
                padding: '10px 20px',
                borderRadius: '10px',
                border: '0',
                color: '#fff',
                height: '40px',
              }}
              onClick={() => verifyConnection()}
            >
              Retry
            </Button>
          </Grid>
        </div>
      )}
      {isSuccess === true && (
        <div
          style={{
            background: '#00B39F40',
            borderRadius: '10px',
            padding: '10px',
            display: 'flex',
          }}
        >
          <Grid style={{ width: '90%' }}>
            <Typography>
              <b>Credential Verified</b>
            </Typography>
            <Typography sx={{ color: '#00000020' }}>
              {`Credential ${connectionObject?.credential?.metadata?.modelDisplayName} created.`}
            </Typography>
          </Grid>
          <Grid style={{ width: '10%' }}></Grid>
        </div>
      )}
    </StepperContent>
  );
};

export const Finish = () => {
  const cancelCallback = () => {
    // Close Modal
  };

  return (
    <StepperContent
      {...FinishContent}
      cancelCallback={cancelCallback}
      disabled={false}
    ></StepperContent>
  );
};
