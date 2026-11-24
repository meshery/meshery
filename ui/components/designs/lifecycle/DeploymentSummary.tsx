import { styled, useTheme } from '@/theme';
import { Box, Stack, Typography } from '@sistent/sistent';

import { FormatStructuredData, TextWithLinks } from '../../data-formatter';
import { SEVERITY_STYLE } from '../../layout/NotificationCenter/constants';
import { ComponentIcon } from './common';
import { Button } from '@sistent/sistent';
import { ExternalLinkIcon, componentIcon } from '@sistent/sistent';
import { ErrorMetadataFormatter } from '../../layout/NotificationCenter/formatters/error';
import { DesignEventErrorsFormatter } from '../../layout/NotificationCenter/formatters/design_event_errors';
import {
  getDesignEventErrorSummary,
  getDesignEventIdentity,
} from '../../layout/NotificationCenter/formatters/design_event_metadata';
import { DryRunResponse } from '../../layout/NotificationCenter/formatters/pattern_dryrun';

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

const DeploymentSummaryFormatter_ = ({ event }) => {
  const theme = useTheme();
  const eventStyle = SEVERITY_STYLE[event?.severity] || {};
  const metadata = event?.metadata || {};
  const errors = metadata.error;
  const errorAction = event?.action;
  const router = useRouter();
  const { designId, designName } = getDesignEventIdentity(metadata);
  const errorSummary = getDesignEventErrorSummary(metadata);
  // Events persisted before #22764 embedded the raw deployment or dry-run
  // response; it is rendered only when the event has no error summary.
  const legacySummary = errorSummary ? {} : metadata.summary || {};
  const componentsDetails = Object.values(legacySummary).flatMap((perComponentDetail) => {
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
        {is_operator_enabled && errorAction != 'register' && errorAction != 'dry-run' && (
          <Button
            variant="contained"
            color="primary"
            onClick={() => openViewScopedToDesignInOperator(designName, designId, router)}
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

      {errorSummary && <DesignEventErrorsFormatter summary={errorSummary} />}

      {legacySummary.dryRunResponse && <DryRunResponse response={legacySummary.dryRunResponse} />}

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
