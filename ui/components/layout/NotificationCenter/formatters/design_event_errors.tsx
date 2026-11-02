import React from 'react';
import { Box, Typography, styled } from '@sistent/sistent';
import { isMeshkitErrorCode, meshkitErrorDocLink } from './design_event_metadata';
import type { DesignEventErrorSummary } from './design_event_metadata';

const ErrorDocLink = styled('a')({
  color: 'inherit',
});

// Only a MeshKit code has an error code reference entry to link to; a
// JSON-schema keyword or Kubernetes cause type is shown as-is.
const DesignEventErrorCode = ({ code }: { code?: string }) => {
  if (!code) {
    return <>Unknown error</>;
  }
  if (!isMeshkitErrorCode(code)) {
    return <>{code}</>;
  }
  return (
    <ErrorDocLink href={meshkitErrorDocLink(code)} target="_blank" rel="noopener noreferrer">
      {code}
    </ErrorDocLink>
  );
};

// DesignEventErrorsFormatter renders the status, error count and errors of a
// pattern validate, deploy, undeploy or dry-run event (#22764). The errors
// list is capped server-side; truncated marks that errorCount exceeds it.
export const DesignEventErrorsFormatter = ({ summary }: { summary: DesignEventErrorSummary }) => {
  const { status, errorCount, errors, truncated } = summary;
  const countLabel = `${errorCount} ${errorCount === 1 ? 'error' : 'errors'}`;

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
      <Typography variant="textB1Regular">
        {status ? (
          <>
            Status: <strong>{status}</strong> - {countLabel}
          </>
        ) : (
          countLabel
        )}
      </Typography>
      {errors.map((error, index) => (
        <Typography
          variant="textB1Regular"
          key={`${error.componentId}-${error.code}-${index}`}
          data-testid="design-event-error"
        >
          <DesignEventErrorCode code={error.code} />
          {error.componentId && ` - component ${error.componentId}`}
        </Typography>
      ))}
      {truncated && (
        <Typography variant="textB1Regular">
          Showing the first {errors.length} of {errorCount} errors.
        </Typography>
      )}
    </Box>
  );
};
