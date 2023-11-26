import React from 'react';
import { ConnectionADetailContent, FinishContent, CredentialDetailContent } from './constants';

import StepperContent from '../../stepper/StepperContentWrapper';
import RJSFWrapper from '../../../MesheryMeshInterface/PatternService/RJSF_wrapper';

export const ConnectionDetails = () => {
  const formRef = React.createRef();
  const schema = {
    rjsfSchema: {
      properties: {
        name: {
          description: 'A short, memorable name of the registering connection.',
          minLength: 1,
          title: 'Name',
          type: 'string',
          'x-rjsf-grid-area': '12',
        },
      },
      required: ['name'],
      type: 'object',
    },
    uiSchema: {
      'ui:order': ['name'],
    },
  };

  const handleCallback = () => {
    // MOve to step 2
  };

  const cancelCallback = () => {
    // Close Modal
  };

  return (
    <StepperContent
      {...ConnectionADetailContent}
      handleCallback={handleCallback}
      cancelCallback={cancelCallback}
      disabled={false}
      // btnText={'Connect'}
    >
      <RJSFWrapper
        key="register-connection-rjsf-form"
        uiSchema={schema.uiSchema}
        jsonSchema={schema.rjsfSchema}
        liveValidate={false}
        formRef={formRef}
      />
    </StepperContent>
  );
};

export const CredentialDetails = () => {
  const formRef = React.createRef();
  const schemaExistingCredential = {
    rjsfSchema: {
      properties: {
        name: {
          description: 'A short, memorable name of the registering connection.',
          minLength: 1,
          title: 'Name',
          type: 'string',
          'x-rjsf-grid-area': '12',
        },
      },
      required: ['name'],
      type: 'object',
    },
    uiSchema: {
      'ui:order': ['name'],
    },
  };

  const schemaNewCredential = {
    rjsfSchema: {
      properties: {
        secret: {
          description: 'Secret',
          minLength: 1,
          title: 'Secret',
          type: 'string',
          'x-rjsf-grid-area': '12',
        },
        token: {
          description: 'Token',
          minLength: 1,
          title: 'Token',
          type: 'string',
          'x-rjsf-grid-area': '12',
        },
      },
      required: ['secret', 'token'],
      type: 'object',
    },
    uiSchema: {
      'ui:order': ['name'],
    },
  };

  const handleCallback = () => {
    // MOve to step 2
  };

  const cancelCallback = () => {
    // Close Modal
  };

  return (
    <StepperContent
      {...CredentialDetailContent}
      handleCallback={handleCallback}
      cancelCallback={cancelCallback}
      disabled={false}
    >
      <p className={{ paddingLeft: '16px' }}>
        Select an existing credential to use for this connection
      </p>
      <RJSFWrapper
        key="register-connection-rjsf-form"
        uiSchema={schemaExistingCredential.uiSchema}
        jsonSchema={schemaExistingCredential.rjsfSchema}
        liveValidate={false}
        formRef={formRef}
      />
      <p style={{ display: 'flex', justifyContent: 'center' }}>-OR-</p>
      <p style={{ paddingLeft: '16px' }}>Configure a new credential to use for this connection</p>
      <RJSFWrapper
        key="register-connection-rjsf-form"
        uiSchema={schemaNewCredential.uiSchema}
        jsonSchema={schemaNewCredential.rjsfSchema}
        liveValidate={false}
        formRef={formRef}
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
