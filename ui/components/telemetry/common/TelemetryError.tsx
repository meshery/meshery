import React from 'react';
import { Box, Typography } from '@sistent/sistent';

interface TelemetryErrorProps {
  error: unknown;
  fallback: string;
  variant?: 'body2' | 'caption';
}

const TelemetryError: React.FC<TelemetryErrorProps> = ({ error, fallback, variant = 'body2' }) => {
  const data = error && typeof error === 'object' && 'data' in error ? error.data : null;
  const message =
    typeof error === 'string'
      ? error
      : data && typeof data === 'object' && 'error' in data && typeof data.error === 'string'
        ? data.error
        : fallback;
  const remedies =
    data &&
    typeof data === 'object' &&
    'suggestedRemediation' in data &&
    Array.isArray(data.suggestedRemediation)
      ? data.suggestedRemediation.filter((step): step is string => typeof step === 'string')
      : [];

  return (
    <Box role="alert">
      <Typography variant={variant} color="error" sx={{ overflowWrap: 'anywhere' }}>
        {message}
      </Typography>
      {remedies.length > 0 && (
        <Box component="ul" sx={{ m: 0, pl: 2 }}>
          {remedies.map((step, index) => (
            <Typography key={`${index}-${step}`} component="li" variant={variant} color="error">
              {step}
            </Typography>
          ))}
        </Box>
      )}
    </Box>
  );
};

export default TelemetryError;
