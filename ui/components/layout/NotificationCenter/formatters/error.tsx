import React from 'react';
import {
  Typography,
  Grid2,
  Box,
  Chip,
  List,
  ListItem,
  RenderMarkdown,
  styled,
} from '@sistent/sistent';
import { FormatStructuredData } from '../../../data-formatter';
import { canTruncateDescription } from '../notification';
import { normalizeEventError } from './event_error';
import type { ErrorCause, ErrorRemediation } from './event_error';

// The labels of the action policies that let a client act on a remediation.
// A manual remediation is the default and carries no label.
const ACTION_POLICY_LABEL = {
  automatic: 'Runs automatically',
  requiresApproval: 'Requires approval',
} as const;

const CauseCard = styled(Box, {
  shouldForwardProp: (prop) => prop !== 'matched',
})<{ matched?: boolean }>(({ theme, matched }) => ({
  padding: theme.spacing(1, 1.5),
  marginBottom: theme.spacing(1),
  borderRadius: theme.shape.borderRadius,
  border: `1px solid ${matched ? theme.palette.border.brand : theme.palette.border.default}`,
  backgroundColor: matched ? theme.palette.background.blur.light : 'transparent',
}));

const CauseHeading = styled(Box)(({ theme }) => ({
  display: 'flex',
  alignItems: 'flex-start',
  justifyContent: 'space-between',
  gap: theme.spacing(1),
}));

const MarkdownLine = ({ text }: { text: string }) => {
  const hashed = text.trim().startsWith('-');
  return (
    <ListItem
      sx={{
        display: hashed ? 'block' : 'list-item',
        lineHeight: '1.2rem',
        padding: '0',
        pb: 1,
        '& ul': { paddingInlineStart: hashed ? '20px' : '0px' },
        '& p': { lineHeight: '1.2rem', margin: 0 },
      }}
    >
      <RenderMarkdown content={text} />
    </ListItem>
  );
};

const ErrorDetailsList = ({ heading, value }: { heading: string; value: string[] }) => {
  const isUnorderedList = heading === 'Details' && value.length <= 1;
  return (
    <Box>
      <Typography variant="body1">
        <strong>{heading}</strong>
      </Typography>
      <List sx={{ listStyleType: isUnorderedList ? 'none' : 'decimal', pl: 3 }}>
        {value.map((line, idx) => (
          <MarkdownLine key={idx} text={line} />
        ))}
      </List>
    </Box>
  );
};

const RemediationItem = ({ remediation }: { remediation: ErrorRemediation }) => {
  const policyLabel =
    ACTION_POLICY_LABEL[remediation.actionPolicy as keyof typeof ACTION_POLICY_LABEL];
  return (
    <ListItem
      data-testid="error-remediation"
      sx={{ display: 'list-item', padding: 0, pb: 0.5, '& p': { margin: 0 } }}
    >
      <RenderMarkdown content={remediation.description} />
      {policyLabel && (
        <Chip
          size="small"
          variant="outlined"
          label={policyLabel}
          title={remediation.action}
          sx={{ ml: 1, height: 20 }}
        />
      )}
    </ListItem>
  );
};

const CauseItem = ({ cause }: { cause: ErrorCause }) => (
  <CauseCard matched={cause.matched} data-testid="error-cause" data-cause-id={cause.id}>
    <CauseHeading>
      <Box sx={{ '& p': { margin: 0, fontWeight: 600 } }}>
        <RenderMarkdown content={cause.description} />
      </Box>
      {cause.matched && (
        <Chip
          size="small"
          color="primary"
          label={
            cause.confidence === undefined
              ? 'Likely cause'
              : `Likely cause · ${Math.round(cause.confidence * 100)}%`
          }
          sx={{ flexShrink: 0, height: 20 }}
        />
      )}
    </CauseHeading>
    {cause.remediations.length > 0 && (
      <List sx={{ listStyleType: 'disc', pl: 3, pt: 0.5, pb: 0 }}>
        {cause.remediations.map((remediation) => (
          <RemediationItem key={remediation.id} remediation={remediation} />
        ))}
      </List>
    )}
  </CauseCard>
);

/**
 * Renders an error event's metadata.error. An error whose producer paired each
 * probable cause with its own remediations is shown cause by cause, the causes
 * the producer matched to this occurrence first. An unpaired error (MeshKit's
 * flat lists, including every row recorded before the pairing existed) is shown
 * as the two lists, never paired by position.
 */
export const ErrorMetadataFormatter = ({
  metadata,
  event,
}: {
  metadata: unknown;
  event?: { description?: unknown };
}) => {
  const error = normalizeEventError(metadata);
  if (!error) return null;

  // The notification's title already shows the description unless it had to
  // be truncated there.
  const description =
    typeof event?.description === 'string' && canTruncateDescription(event.description)
      ? event.description
      : null;
  const { longDescription, causes, probableCause, suggestedRemediation, referenceUrl, code } =
    error;

  return (
    <Grid2 size="grow">
      <div>
        {description && <FormatStructuredData data={description} />}
        <div style={{ marginTop: '1rem' }}>
          <FormatStructuredData
            data={{ Details: longDescription }}
            propertyFormatters={{
              Details: (value) => <ErrorDetailsList heading="Details" value={value} />,
            }}
          />
        </div>
      </div>
      {causes.length > 0 ? (
        <Box sx={{ mt: 1 }} data-testid="error-causes">
          <Typography variant="body1" sx={{ mb: 1 }}>
            <strong>Probable Causes</strong>
          </Typography>
          {causes.map((cause) => (
            <CauseItem key={cause.id} cause={cause} />
          ))}
        </Box>
      ) : (
        <Grid2 container spacing={1} style={{ marginTop: '0.5rem' }} size="grow">
          <Grid2 size={{ sm: probableCause.length > 0 ? 6 : 12 }}>
            <FormatStructuredData
              data={{ 'Probable Cause': probableCause }}
              propertyFormatters={{
                'Probable Cause': (value) => (
                  <ErrorDetailsList heading="Probable Cause" value={value} />
                ),
              }}
            />
          </Grid2>
          <Grid2 size={{ sm: suggestedRemediation.length > 0 ? 6 : 12 }}>
            <FormatStructuredData
              data={{ 'Suggested Remediation': suggestedRemediation }}
              propertyFormatters={{
                'Suggested Remediation': (value) => (
                  <ErrorDetailsList heading="Suggested Remediation" value={value} />
                ),
              }}
            />
          </Grid2>
        </Grid2>
      )}
      {referenceUrl && (
        <Typography variant="body2" sx={{ mt: 1 }}>
          <a href={referenceUrl} target="_blank" rel="noopener noreferrer">
            {code ? `${code} reference` : 'Error reference'}
          </a>
        </Typography>
      )}
    </Grid2>
  );
};
