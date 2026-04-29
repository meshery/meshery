// @ts-nocheck
import React from 'react';
import { FormatDryRunResponse } from '../../DesignLifeCycle/DryRun';
import { formatDryRunResponse } from 'machines/validator/designValidator';
import { ValidationResults } from '@/components/DesignLifeCycle/ValidateDesign';

export const DryRunResponse = ({ response }) => {
  return <FormatDryRunResponse dryRunErrors={formatDryRunResponse(response)} />;
};

export const SchemaValidationFormatter = ({ event }) => {
  const results = event?.metadata?.validationResult || {};
  const designName = event?.metadata?.design_name || 'Unknown Design';
  const totalComponents = event?.metadata?.total_components || 0;
  const configurableComponents = event?.metadata?.configurable_components || 0;

  const totalErrors = Object.values(results || {}).reduce(
    (acc, serviceResult) => acc + (serviceResult?.errors?.length || 0),
    0,
  );

  const validationMachine = {
    send: () => {},
  };

  return (
    <ValidationResults
      validationResults={results}
      errorCount={totalErrors}
      compCount={configurableComponents}
      annotationCount={totalComponents - configurableComponents}
      design={designName}
      validationMachine={validationMachine}
    />
  );
};
