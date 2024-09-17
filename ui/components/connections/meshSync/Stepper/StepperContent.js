import React, { useEffect, useRef, useState } from 'react';
import OutlinedInput from '@mui/material/OutlinedInput';
import InputLabel from '@mui/material/InputLabel';
import FormControl from '@mui/material/FormControl';
import { Checkbox, MenuItem, ListItemText, Select, Typography } from '@material-ui/core';

import {
  ConnectionDetailContent,
  FinishContent,
  CredentialDetailContent,
  SelectConnectionTypeContent,
} from './constants';
import StepperContent from './StepperContentWrapper';
import RJSFWrapper from '../../../MesheryMeshInterface/PatternService/RJSF_wrapper';
import dataFetch from '../../../../lib/data-fetch';
import { Box } from '@mui/material';
import { selectCompSchema } from '../../../RJSFUtils/common';
import { JsonParse, randomPatternNameGenerator } from '../../../../utils/utils';
import Notification from './Notification';

const CONNECTION_TYPES = ['Prometheus Connection', 'Grafana Connection'];

const schema = selectCompSchema(
  CONNECTION_TYPES,
  'Select one of the available Connection type',
  'Select type of Connection to register',
  'selectedConnectionType',
);
export const SelectConnection = ({ setSharedData, handleNext }) => {
  const formRef = useRef();

  const registerConnection = (componentName) => {
    dataFetch(
      '/api/integrations/connections/register',
      {
        method: 'POST',
        credentials: 'include',
        body: JSON.stringify({
          kind: componentName,
          status: 'initialize',
        }),
      },
      (result) => {
        let schemaObj = {
          connection: JsonParse(result?.connection?.schema),
          credential: JsonParse(result?.credential?.schema),
        };
        // formConnectionIdRef.current = result.id;
        setSharedData((prevState) => ({
          ...prevState,
          connection: result,
          schemas: schemaObj,
          kind: componentName.toLowerCase(),
        }));
        handleNext();
      },
    );
  };

  const handleCallback = () => {
    handleNext();
  };

  const handleChange = (data) => {
    if (data.selectedConnectionType) {
      const selectedConnectionType = data.selectedConnectionType;
      // The selectedConnectionType is the concatentaion of connectionType, ' ' and 'Connection' suffix.
      // Therefore, when initiating connection we are removing ' ' and suffix so that correct schema is retrieved.
      registerConnection(selectedConnectionType?.slice(0, selectedConnectionType.indexOf(' ')));
    }
  };

  return (
    <StepperContent {...SelectConnectionTypeContent} handleCallback={handleCallback}>
      <RJSFWrapper
        key="select-connection-type-rjsf-form"
        jsonSchema={schema}
        liveValidate={false}
        formRef={formRef}
        onChange={handleChange}
      />
    </StepperContent>
  );
};

export const ConnectionDetails = ({ sharedData, setSharedData, handleNext }) => {
  const formRef = React.createRef();
  const [selectedEndpoint, setSelectedEndpoint] = useState(null);
  // stores selected endpoint just before dropdown is closed
  const [prevSelectedEndpoint, setPrevSelectedEndpoint] = useState(null);

  useEffect(() => {
    ConnectionDetailContent.title = `Connecting to ${sharedData?.kind}`;
  }, [sharedData?.connection]);

  const handleCallback = () => {
    handleNext();
  };

  const cancelCallback = () => {
    sharedData.onClose();
  };

  const handleSelectEndpoint = (e) => {
    setSharedData((prevState) => ({
      ...prevState,
      componentForm: {
        name: randomPatternNameGenerator(),
        url: e.target.value,
      },
    }));

    setSelectedEndpoint(e.target.value);
    setPrevSelectedEndpoint(e.target.value);
  };

  const handleClose = () => {
    if (prevSelectedEndpoint === selectedEndpoint) {
      setSelectedEndpoint(null);
      setPrevSelectedEndpoint(null);
    }
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
      cancelCallback={cancelCallback}
    >
      {sharedData?.capabilities && (
        <FormControl fullWidth size="small">
          <InputLabel fontSize="inherit" id="endpoint-checkbox-label">
            Select from the discovered endpoints
          </InputLabel>
          <Select
            labelId="endpoint-checkbox-label"
            id="endpoint-checkbox"
            onChange={handleSelectEndpoint}
            value={selectedEndpoint}
            onClose={handleClose}
            input={<OutlinedInput label="Select discovered endpoint" />}
            renderValue={() => <div>{selectedEndpoint !== null ? selectedEndpoint : ''}</div>}
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
            {sharedData.capabilities?.urls?.map((endpoint, index) => (
              <MenuItem key={index} value={endpoint} name={endpoint}>
                <Checkbox checked={endpoint === selectedEndpoint} />
                <ListItemText primary={endpoint} />
              </MenuItem>
            ))}
          </Select>
        </FormControl>
      )}
      <p style={{ display: 'flex', justifyContent: 'center' }}>-OR-</p>
      <p>Enter the {sharedData.kind} service URL</p>
      <RJSFWrapper
        key="register-connection-rjsf-form"
        jsonSchema={sharedData?.schemas?.connection}
        liveValidate={true}
        formRef={formRef}
        disabled={selectedEndpoint !== null ? true : false}
        onChange={handleChange}
      />
    </StepperContent>
  );
};

export const CredentialDetails = ({ sharedData, handleNext, handleRegistrationComplete }) => {
  const [existingCredentials, setExistingCredentials] = useState([]);
  const [selectedCredential, setSelectedCredential] = useState(null);
  const [prevSelectedCredential, setPrevSelectedCredential] = useState(null);
  const [formState, setFormState] = useState(null);
  const [skipCredentialVerification, setSkipCredentialVerification] = useState(false);
  const [disableVerify, setDisableVerify] = useState(true);
  const [isSuccess, setIsSuccess] = React.useState(null);
  const formRef = React.createRef();
  useEffect(() => {
    getchExistingCredential();
  }, []);

  useEffect(() => {
    CredentialDetailContent.title = `Credential for ${sharedData?.kind}`;
  }, [sharedData.kind]);

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
      (error) => {
        console.error('Error fetching existing credentials:', error);
      },
    );
  };

  const verifyConnection = () => {
    let credential = {};
    if (selectedCredential === null) {
      credential = formState;
    } else {
      credential = {
        secret: selectedCredential?.secret?.secret,
        name: selectedCredential?.name,
      };
      credential.id = selectedCredential?.id;
    }

    dataFetch(
      '/api/integrations/connections/register',
      {
        method: 'POST',
        credentials: 'include',
        body: JSON.stringify({
          skip_credential_verification: skipCredentialVerification,
          kind: sharedData?.kind, // this is "kind" column of the current row which is selected in the meshsync table. i.e. the entry against which registration process has been invoked.
          name: sharedData?.componentForm?.name, // This name is from the name field in schema
          type: sharedData?.connection?.connection?.model?.category?.name?.toLowerCase(),
          sub_type: sharedData?.connection?.connection?.model?.subCategory?.toLowerCase(),
          metadata: sharedData?.componentForm,
          credential_secret: credential,
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
      (error) => {
        console.error('Error verifying connection:', error);
        setIsSuccess(false);
      },
    );
  };

  const connectToConnection = () => {
    let credential = {};
    if (selectedCredential === null) {
      credential = formState;
    } else {
      credential = {
        name: selectedCredential?.name,
        secret: selectedCredential?.secret?.secret,
      };
      credential.id = selectedCredential?.id;
    }

    dataFetch(
      '/api/integrations/connections/register',
      {
        method: 'POST',
        credentials: 'include',
        body: JSON.stringify({
          kind: sharedData?.kind, // this is "kind" column of the current row which is selected in the meshsync table. i.e. the entry against which registration process has been invoked.
          name: sharedData?.componentForm?.name, // This name is from the name field in schema
          type: sharedData?.connection?.connection?.model?.category?.name?.toLowerCase(),
          sub_type: sharedData?.connection?.connection?.model?.subCategory?.toLowerCase(),
          metadata: sharedData?.componentForm,
          credential_secret: credential,
          id: sharedData?.connection?.id,
          status: 'connect',
        }),
      },
      (result) => {
        if (result !== undefined && result !== null && result === '') {
          setIsSuccess(true);
        } else {
          setIsSuccess(false);
        }
      },
      (error) => {
        console.error('Error connecting to connection:', error);
        setIsSuccess(false);
      },
    );
  };

  const handleCallback = () => {
    if (isSuccess === null || isSuccess === false) {
      verifyConnection();
    } else {
      handleNext();
      handleRegistrationComplete();
    }
  };

  const cancelCallback = () => {
    sharedData.onClose();
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
    if (selectedCredential !== null || (formState !== null && formState['secret']) !== undefined) {
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
      <Typography variant="body2" style={{ paddingLeft: '16px' }}>
        Select an existing credential to use for this connection
      </Typography>
      <FormControl sx={{ width: '100%' }} size="small">
        <InputLabel fontSize="20" id="credential-checkbox-label">
          Select existing credential
        </InputLabel>
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
              <MenuItem key={credential.id} value={credential.id} name={credential.name}>
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
        jsonSchema={sharedData?.schemas?.credential}
        liveValidate={true}
        formRef={formRef}
        disabled={selectedCredential !== null ? true : false}
        onChange={handleChange}
      />
      <Box
        style={{
          background: 'rgba(0, 211, 169, 0.05)',
          padding: '0.4rem',
          margin: '1rem 0',
        }}
      >
        <Typography style={{ fontSize: 'inherit' }}>
          <Checkbox
            id="bypass_verification"
            color="success"
            onChange={(e) => {
              setSkipCredentialVerification(e.target.checked);
              setDisableVerify(!e.target.checked);
            }}
          />
          <label fontSize="inherit" for="bypass_verification">
            Bypass connection verification
          </label>
        </Typography>
      </Box>
      {isSuccess !== null && (
        <Notification
          type={isSuccess ? 'success' : 'error'}
          message={`Credential for ${sharedData?.kind} ${
            isSuccess ? 'created' : 'verification failed'
          }`}
          retry={!isSuccess}
          onRetry={() => verifyConnection()}
        />
      )}
    </StepperContent>
  );
};

export const Finish = ({ sharedData }) => {
  const cancelCallback = () => {
    sharedData.onClose();
  };

  return (
    <StepperContent
      {...FinishContent}
      handleCallback={cancelCallback}
      disabled={false}
    ></StepperContent>
  );
};
