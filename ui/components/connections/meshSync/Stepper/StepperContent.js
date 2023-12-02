import React, { useEffect, useState } from 'react';
import { ConnectionDetailContent, FinishContent, CredentialDetailContent } from './constants';

import StepperContent from './StepperContentWrapper';
import RJSFWrapper from '../../../MesheryMeshInterface/PatternService/RJSF_wrapper';
import dataFetch from '../../../../lib/data-fetch';
import { FormControl, MenuItem, Select } from '@material-ui/core';

export const ConnectionDetails = ({ sharedData, setSharedData, handleNext }) => {
  const [schema, setSchema] = useState({ properties: {} });
  const [connectionObject, setConnectionObject] = useState(null);
  const formRef = React.createRef();

  useEffect(() => {
    sharedData !== null && registerConnection(sharedData?.model);
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

  return (
    <StepperContent {...ConnectionDetailContent} handleCallback={handleCallback} disabled={false}>
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
  const [schema, setSchema] = useState({});
  const [connectionObject, setConnectionObject] = useState(null);
  const [existingCredentials, setExistingCredentials] = useState([]);
  const [selectedCredential, setSelectedCredential] = useState(null);
  const [formState, setFormState] = React.useState(null);
  const [/*isSuccess,*/ setIsSuccess] = React.useState(null);
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
          Kind: sharedData?.connection?.component?.kind,
          Name: sharedData?.connection?.component?.displayName,
          Type: sharedData?.connection?.component?.metadata?.model,
          Subtype: sharedData?.connection?.component?.metadata?.subCategory,
          Metadata: sharedData?.componentForm,
          'Credential Secret': formState || selectedCredential,
          ID: sharedData?.connection?.id,
          Status: 'register',
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
    if (selectedCredential !== null) {
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
  };

  const handleChange = (data) => {
    setFormState(data);
  };

  return (
    <StepperContent
      {...CredentialDetailContent}
      handleCallback={handleCallback}
      cancelCallback={cancelCallback}
      disabled={selectedCredential === null || formState === null ? true : false}
      // btnText={selectedCredential === null ? 'Verify Connection' : 'Next'}
    >
      <p className={{ paddingLeft: '16px' }}>
        Select an existing credential to use for this connection
      </p>
      <FormControl fullWidth>
        <Select
          labelId="credential-label"
          id="credential"
          label="Select existing credential"
          defaultValue={'null'}
          onChange={(e) => handleSelectCredential(e)}
        >
          <MenuItem value={'null'} disabled>
            {'Select existing credential'}
          </MenuItem>
          {existingCredentials &&
            existingCredentials?.map((credential) => (
              <MenuItem value={credential.id} name={credential.name}>
                {credential.name}
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
    </StepperContent>
  );
};

export const Finish = () => {
  const handleCallback = () => {
    // MOve to step 2
  };

  const cancelCallback = () => {
    // Close Modal
  };

  return (
    <StepperContent
      {...FinishContent}
      handleCallback={handleCallback}
      cancelCallback={cancelCallback}
      disabled={false}
    >
      <p> Details</p>

      <p>[ METADATA FORMATTER OUTPUT HERE ]</p>
      <p>[ METADATA FORMATTER OUTPUT HERE ]</p>
      <p>[ METADATA FORMATTER OUTPUT HERE ]</p>
      <p>[ METADATA FORMATTER OUTPUT HERE ]</p>
    </StepperContent>
  );
};
