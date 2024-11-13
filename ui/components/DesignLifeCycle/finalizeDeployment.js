import { useSelectorRtk } from '@/store/hooks';
import { selectSelectedEnvs } from '@/store/slices/globalEnvironmentContext';
import { useSelector } from 'react-redux';
const { Box, Typography, Stack, EnvironmentIcon, useTheme, styled } = require('@layer5/sistent');
const { processDesign, CheckBoxField, StepHeading } = require('./common');
const _ = require('lodash');

const StyledSummaryItem = styled(Box)(({ theme }) => ({
  borderRadius: '0.5rem',
  padding: '1rem',
  backgroundColor: theme.palette.background.default,
  flexGrow: 1,
}));

const StyledEnvironment = styled(Box)(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  gap: '0.5rem',
  color: theme.palette.text.neutral.default,
}));

const isVisualizerEnabled = (capabilitiesRegistry) => {
  const navigatorExtension = _.get(capabilitiesRegistry, 'extensions.navigator') || [];
  return navigatorExtension.some((ext) => ext.title === 'MeshMap');
};

export const FinalizeDeployment = ({ design, openInVisualizer, setOpenInVisualizer }) => {
  const { configurableComponents } = processDesign(design);
  const selectedEnvironments = useSelectorRtk(selectSelectedEnvs);
  const envNames = Object.values(selectedEnvironments).map((env) => env.name);

  const capabilitiesRegistry = useSelector((state) => state.get('capabilitiesRegistry'));
  const visualizerEnabled = isVisualizerEnabled(capabilitiesRegistry);
  const theme = useTheme();
  const palette = theme.palette;
  return (
    <Box>
      <StepHeading>Finalize Deployment</StepHeading>
      <Box mt={2} display="flex" justifyContent="space-between" flexWrap={'wrap'} gap={2}>
        <StyledSummaryItem>
          <Typography color={palette.text.disabled} variant="textB2SemiBold">
            Environments
          </Typography>
          <Stack gap={1} mt={1}>
            {envNames.map((env) => (
              <StyledEnvironment key={env}>
                <EnvironmentIcon fill={palette.text.neutral.default} />
                <Typography variant="textB1Regular" key={env}>
                  {env}
                </Typography>
              </StyledEnvironment>
            ))}
          </Stack>
        </StyledSummaryItem>
        <StyledSummaryItem
          display="flex"
          flexDirection="column"
          alignItems="center"
          justifyContent={'center'}
        >
          <Box display="flex" alignItems="baseline">
            <Typography variant="textH1Bold" color={palette.text.neutral.default}>
              {configurableComponents.length}
            </Typography>
            <Typography
              variant="textB2SemiBold"
              color={palette.text.disabled}
              style={{
                textTransform: 'lowercase',
              }}
            >
              component(s)
            </Typography>
          </Box>
        </StyledSummaryItem>
      </Box>
      <Stack mt={3} gap={1}>
        {visualizerEnabled && (
          <CheckBoxField
            label="Open in Visualizer"
            checked={openInVisualizer}
            helpText="Opens the deployed design in visualizer"
            onChange={() => setOpenInVisualizer(!openInVisualizer)}
          />
        )}
        <CheckBoxField label="Schedule Deployment" checked={false} onChange={() => {}} disabled />
      </Stack>
    </Box>
  );
};
