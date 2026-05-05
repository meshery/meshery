import React from 'react';
import { FormatDryRunResponse } from '../../DesignLifeCycle/DryRun';
import { formatDryRunResponse } from 'machines/validator/designValidator';
import { ValidationResults } from '@/components/DesignLifeCycle/ValidateDesign';

export const DryRunResponse = ({ response }: { response: unknown }) => {
  return <FormatDryRunResponse dryRunErrors={formatDryRunResponse(response)} />;
};

type SchemaValidationEvent = {
  metadata?: {
    validationResult?: Record<string, { errors?: unknown[] }>;
    design_name?: string;
    total_components?: number;
    configurable_components?: number;
  };
};

export const SchemaValidationFormatter = ({ event }: { event: SchemaValidationEvent }) => {
  const results = event?.metadata?.validationResult || {};
  const designName = event?.metadata?.design_name || 'Unknown Design';
  const totalComponents = event?.metadata?.total_components || 0;
  const configurableComponents = event?.metadata?.configurable_components || 0;

  const totalErrors = Object.values(results || {}).reduce(
    (acc: number, serviceResult) => acc + (serviceResult?.errors?.length || 0),
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
