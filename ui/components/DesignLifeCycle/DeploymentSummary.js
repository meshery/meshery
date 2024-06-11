import { useGetComponentsByModelAndKindQuery } from '@/rtk-query/meshModel';
import { NOTIFICATIONCOLORS } from '@/themes/index';
import { Box, Stack, Typography, styled, useTheme } from '@layer5/sistent';
import { alpha } from '@mui/material';
import { FormatStructuredData, TextWithLinks } from '../DataFormatter';
import { SEVERITY_STYLE } from '../NotificationCenter/constants';
import { ErrorMetadataFormatter } from '../NotificationCenter/metadata';
import { ComponentIcon } from './common';
import { Button } from '@layer5/sistent';
import { ExternalLinkIcon } from '@layer5/sistent';
import { UsesSistent } from '../SistentWrapper';

const StyledDetailBox = styled(Box)(({ theme, severityColor, bgOpacity }) => ({
  padding: theme.spacing(2),
  backgroundColor: alpha(severityColor, bgOpacity),
  border: `1px solid ${theme.palette.divider}`,
  display: 'flex',
}));

const DeployementComponentFormatter = ({ componentDetail }) => {
  const { data } = useGetComponentsByModelAndKindQuery({
    model: componentDetail.Model || 'kubernetes',
    component: componentDetail.Kind,
  });
  const componentDef = data?.components?.[0];
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
        {componentDef?.metadata?.svgColor && (
          <ComponentIcon
            iconSrc={`/${componentDef?.metadata?.svgColor}`}
            alt={componentDetail.Kind}
          />
        )}
        <Typography variant="textB1Regular" style={{ textTransform: 'capitalize' }}>
          {componentDetail.Message}
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
        {event?.metadata?.view_link && (
          <Button
            variant="contained"
            color="primary"
            onClick={() => window.open(event.metadata?.view_link, '_blank')}
            style={{ gap: '0.25rem' }}
          >
            Open In Visualizer <ExternalLinkIcon fill={theme.palette.common.white} />
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
  <UsesSistent>
    <DeploymentSummaryFormatter_ event={event} />
  </UsesSistent>
);
