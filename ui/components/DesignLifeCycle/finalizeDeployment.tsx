import { useSelector } from 'react-redux';
import { selectSelectedEnvs } from '@/store/slices/globalEnvironmentContext';
import { Box, Typography, Stack, EnvironmentIcon, useTheme, styled } from '@sistent/sistent';
import { processDesign, CheckBoxField, StepHeading } from './common';

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
  color: theme.palette.text.neutral?.default || theme.palette.text.primary,
}));

export const FinalizeDeployment = ({ design }) => {
  const { configurableComponents } = processDesign(design);
  const selectedEnvironments = useSelector(selectSelectedEnvs);
  const envNames = (Object.values(selectedEnvironments) as Array<{ name?: string }>)
    .map((env) => env?.name)
    .filter((name): name is string => typeof name === 'string');

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
                <EnvironmentIcon fill={palette.text.neutral?.default || palette.text.primary} />
                <Typography>{env}</Typography>
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
            <Typography
              sx={{
                fontSize: '3.2rem',
                color: palette.text.neutral?.default || palette.text.primary,
              }}
            >
              {configurableComponents.length}
            </Typography>
            <Typography
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
        <CheckBoxField label="Schedule Deployment" checked={false} onChange={() => {}} disabled />
      </Stack>
    </Box>
  );
};
