import React from 'react';
import { FormatDryRunResponse } from '../../../designs/lifecycle/DryRun';
import { formatDryRunResponse } from 'machines/validator/designValidator';
import { ValidationResults } from '@/components/designs/lifecycle/ValidateDesign';
import { Box, Typography } from '@sistent/sistent';
import { DesignEventErrorsFormatter } from './design_event_errors';
import { getDesignEventErrorSummary, getDesignEventIdentity } from './design_event_metadata';

export const DryRunResponse = ({ response }: { response: unknown }) => {
  return <FormatDryRunResponse dryRunErrors={formatDryRunResponse(response)} />;
};

type SchemaValidationEvent = {
  metadata?: {
    designName?: string;
    design_name?: string;
    status?: string;
    errorCount?: number;
    errors?: { code?: string; componentId?: string }[];
    truncated?: boolean;
    // Events persisted before #22764 embedded the per-component validation
    // results, components included.
    validationResult?: Record<string, { errors?: unknown[] }>;
    total_components?: number;
    configurable_components?: number;
  };
};

// SchemaValidationFormatter renders a pattern validate event. Since #22764 it
// carries status, errorCount, errors[] of {code, componentId} and truncated,
// like a deploy event; an errors[].code is the schema validator keyword
// (`required`, `type`, ...), not a MeshKit code.
export const SchemaValidationFormatter = ({ event }: { event: SchemaValidationEvent }) => {
  const metadata = event?.metadata;
  const designName = getDesignEventIdentity(metadata).designName || 'Unknown Design';
  const errorSummary = getDesignEventErrorSummary(metadata);

  if (errorSummary || !metadata?.validationResult) {
    return (
      <Box display="flex" flexDirection="column" gap={1} sx={{ mt: 1 }}>
        <Typography variant="textB1Regular">
          Design: <strong>{designName}</strong>
        </Typography>
        {errorSummary && <DesignEventErrorsFormatter summary={errorSummary} />}
      </Box>
    );
  }

  const results = metadata.validationResult;
  const totalComponents = metadata.total_components || 0;
  const configurableComponents = metadata.configurable_components || 0;

  const totalErrors = Object.values(results).reduce(
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
