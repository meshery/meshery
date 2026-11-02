import { styled, useTheme } from '@/theme';
import { Box, Stack, Typography } from '@sistent/sistent';

import { FormatStructuredData, TextWithLinks } from '../../data-formatter';
import { SEVERITY_STYLE } from '../../layout/NotificationCenter/constants';
import { ComponentIcon } from './common';
import { Button } from '@sistent/sistent';
import { ExternalLinkIcon, componentIcon } from '@sistent/sistent';
import { ErrorMetadataFormatter } from '../../layout/NotificationCenter/formatters/error';

import { openViewScopedToDesignInOperator, useIsOperatorEnabled } from '@/utils/utils';
import { useRouter } from 'next/router';
import { capitalize } from 'lodash';

const StyledDetailBox = styled(Box)(() => ({
  display: 'flex',
}));

// deployment_type is deploy/undeploy
const DeploymentComponentFormatter = ({ componentDetail, deploymentType }) => {
  const theme = useTheme();

  return (
    <StyledDetailBox
      severityColor={
        componentDetail.Success ? theme.palette.success.main : theme.palette.error.main
      }
      bgOpacity={0.1}
      display="flex"
      gap={2}
      flexDirection="column"
    >
      <Stack direction="row" spacing={2} alignItems={'center'}>
        <ComponentIcon
          iconSrc={componentIcon({
            kind: componentDetail.Kind,
            model: componentDetail.Model,
            color: 'color',
          })}
          alt={componentDetail.Kind}
        />
        <Typography variant="textB1Regular">
          {componentDetail.Success
            ? `${capitalize(deploymentType)}ed ${componentDetail.Kind} "${componentDetail.CompName}"`
            : `Failed to ${deploymentType} ${componentDetail.Kind} "${componentDetail.CompName}"`}
        </Typography>
      </Stack>
      {componentDetail.Error && <ErrorMetadataFormatter metadata={componentDetail.Error} />}
      {componentDetail.metadata && <FormatStructuredData data={componentDetail.metadata} />}
    </StyledDetailBox>
  );
};

// DeploymentErrorsFormatter renders the status, error count and per-component
// errors that deploy, undeploy and dry-run events carry since #22764. The list
// is capped server-side; truncated marks that errorCount exceeds it.
const DeploymentErrorsFormatter = ({ metadata }) => {
  const errorList = Array.isArray(metadata?.errors) ? metadata.errors : [];
  const errorCount = metadata?.errorCount ?? errorList.length;

  return (
    <Box display="flex" flexDirection="column" gap={1}>
      <Typography variant="textB1Regular">
        Status: <strong>{metadata.status}</strong> - {errorCount}{' '}
        {errorCount === 1 ? 'error' : 'errors'}
      </Typography>
      {errorList.map((error, index) => (
        <Typography variant="textB1Regular" key={`${error?.componentId}-${error?.code}-${index}`}>
          {error?.code || 'Unknown error'} - component {error?.componentId || 'unknown'}
        </Typography>
      ))}
      {metadata.truncated && (
        <Typography variant="textB1Regular">
          Showing the first {errorList.length} of {errorCount} errors.
        </Typography>
      )}
    </Box>
  );
};

const DeploymentSummaryFormatter_ = ({ event }) => {
  const theme = useTheme();
  const eventStyle = SEVERITY_STYLE[event?.severity] || {};
  const metadata = event?.metadata || {};
  const errors = metadata.error;
  const errorAction = event?.action;
  const router = useRouter();
  // Events persisted before #22764 embedded the raw deployment summary.
  const componentsDetails = Object.values(metadata.summary || {}).flatMap((perComponentDetail) => {
    perComponentDetail = perComponentDetail?.flatMap ? perComponentDetail : [];
    return perComponentDetail.flatMap((perContextDetail) =>
      perContextDetail?.Summary?.map((detail) => ({
        ...detail,
        Location: perContextDetail.Location,
      })),
    );
  });

  const is_operator_enabled = useIsOperatorEnabled();

  return (
    <Box>
      <StyledDetailBox
        severityColor={eventStyle.color}
        bgOpacity={0.1}
        alignItems="center"
        justifyContent="space-between"
      >
        <TextWithLinks
          text={event?.description || ''}
          style={{
            color: theme.palette.text.default,
            textTransform: 'capitalize',
            fontWeight: 'bold',
          }}
        />
        {is_operator_enabled && errorAction != 'register' && (
          <Button
            variant="contained"
            color="primary"
            onClick={() =>
              openViewScopedToDesignInOperator(
                metadata.designName ?? metadata.design_name,
                metadata.designId ?? metadata.design_id,
                router,
              )
            }
            style={{ gap: '0.25rem' }}
          >
            Open In Operator <ExternalLinkIcon fill={theme.palette.background.constant.white} />
          </Button>
        )}
      </StyledDetailBox>
      {errors && (
        <StyledDetailBox severityColor={eventStyle.color} bgOpacity={0}>
          <ErrorMetadataFormatter metadata={errors} event={event} />
        </StyledDetailBox>
      )}

      {metadata.status && <DeploymentErrorsFormatter metadata={metadata} />}

      {componentsDetails.map((componentDetail) => (
        <DeploymentComponentFormatter
          componentDetail={componentDetail}
          deploymentType={event.action}
          key={componentDetail.CompName + componentDetail.Location}
        />
      ))}
    </Box>
  );
};

export const DeploymentSummaryFormatter = ({ event }) => (
  <DeploymentSummaryFormatter_ event={event} />
);
