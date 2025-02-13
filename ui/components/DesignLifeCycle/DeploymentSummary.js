import { NOTIFICATIONCOLORS } from '@/themes/index';
import { Box, Stack, Typography, styled, useTheme } from '@layer5/sistent';
import { alpha } from '@mui/material';
import { FormatStructuredData, TextWithLinks } from '../DataFormatter';
import { SEVERITY_STYLE } from '../NotificationCenter/constants';
import { ErrorMetadataFormatter } from '../NotificationCenter/metadata';
import { ComponentIcon } from './common';
import { Button } from '@layer5/sistent';
import { ExternalLinkIcon, componentIcon } from '@layer5/sistent';

import { openViewScopedToDesignInOperator, useIsOperatorEnabled } from '@/utils/utils';
import { useRouter } from 'next/router';

const StyledDetailBox = styled(Box)(({ theme, severityColor, bgOpacity }) => ({
  padding: theme.spacing(2),
  backgroundColor: alpha(severityColor, bgOpacity),
  border: `1px solid ${theme.palette.divider}`,
  display: 'flex',
}));

const DeployementComponentFormatter = ({ componentDetail }) => {
  return (
    <StyledDetailBox
      severityColor={
        componentDetail.Success ? NOTIFICATIONCOLORS.SUCCESS : NOTIFICATIONCOLORS.ERROR
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
            ? `Deployed ${componentDetail.Kind} "${componentDetail.CompName}"`
            : `Failed to deploy ${componentDetail.Kind} "${componentDetail.CompName}"`}
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
  const errors = event.metadata?.error;
  const router = useRouter();
  const componentsDetails = Object.values(event.metadata?.summary || {}).flatMap(
    (perComponentDetail) => {
      perComponentDetail = perComponentDetail?.flatMap ? perComponentDetail : [];
      return perComponentDetail.flatMap((perContextDetail) =>
        perContextDetail?.Summary?.map((detail) => ({
          ...detail,
          Location: perContextDetail.Location,
        })),
      );
    },
  );

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
          style={{ color: theme.palette.text.default, textTransform: 'capitalize' }}
        />
        {is_operator_enabled && (
          <Button
            variant="contained"
            color="primary"
            onClick={() =>
              openViewScopedToDesignInOperator(
                event?.metadata?.design_name,
                event?.metadata?.design_id,
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

      {componentsDetails.map((componentDetail) => (
        <DeployementComponentFormatter
          componentDetail={componentDetail}
          key={componentDetail.CompName + componentDetail.Location}
        />
      ))}
    </Box>
  );
};

export const DeploymentSummaryFormatter = ({ event }) => (
  <DeploymentSummaryFormatter_ event={event} />
);
