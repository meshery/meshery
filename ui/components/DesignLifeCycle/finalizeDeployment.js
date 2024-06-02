import { useSelectorRtk } from '@/store/hooks';
import { selectSelectedEnvs } from '@/store/slices/globalEnvironmentContext';
const { Box, Typography, Stack, EnvironmentIcon, useTheme, styled } = require('@layer5/sistent');
const { processDesign, CheckBoxField } = require('./common');

const StyledSummaryItem = styled(Box)(({ theme }) => ({
  borderRadius: '0.5rem',
  padding: '1rem',
  backgroundColor: theme.palette.background.constant.white,
  flexGrow: 1,
}));

const StyledContainer = styled(Box)(({ theme }) => ({
  backgroundColor: theme.palette?.background?.blur.light,
}));

const StyledEnvironment = styled(Box)(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  gap: '0.5rem',
  color: theme.palette.text.secondary,
}));

export const FinalizeDeployment = ({ design }) => {
  const { configurableComponents } = processDesign(design);
  const selectedEnvironments = useSelectorRtk(selectSelectedEnvs);
  const envNames = Object.values(selectedEnvironments).map((env) => env.name);
  const theme = useTheme();
  const palette = theme.palette;
  return (
    <StyledContainer>
      <Typography variant="textB2SemiBold">Deployment Summary</Typography>
      <Box mt={2} display="flex" justifyContent="space-between" flexWrap={'wrap'} gap={2}>
        <StyledSummaryItem>
          <Typography color={palette.text.disabled} variant="textB2SemiBold">
            Environments
          </Typography>
          <Stack gap={1} mt={1}>
            {envNames.map((env) => (
              <StyledEnvironment key={env}>
                <EnvironmentIcon />
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
            <Typography variant="textH1Bold" color={palette.text.secondary}>
              {configurableComponents.length}
            </Typography>
            <Typography variant="textB2SemiBold" color={palette.text.disabled}>
              component(s)
            </Typography>
          </Box>
        </StyledSummaryItem>
      </Box>
      <Stack mt={3} gap={1}>
        <CheckBoxField label="Open in Visualizer" checked={false} onChange={() => {}} />
        <CheckBoxField label="Schedule Deployment" checked={false} onChange={() => {}} />
      </Stack>
    </StyledContainer>
  );
};
