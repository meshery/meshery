import React from 'react';
import { Box, Checkbox, CircularProgress, Stack, styled, Typography } from '@sistent/sistent';
import PatternIcon from '@/assets/icons/Pattern';
import { processDesign } from '@/utils/utils';
import { CustomTooltip } from '@sistent/sistent';
import { InfoCircleIcon, useTheme } from '@sistent/sistent';
import { IconButton } from '@sistent/sistent';

export const DEPLOYMENT_TYPE = {
  DEPLOY: 'deploy',
  UNDEPLOY: 'undeploy',
};

const ModelName = styled(Typography)(() => ({
  fontSize: '10px',
  fontWeight: 'bold',
  textTransform: 'uppercase',
}));
export const ComponentIcon = ({ iconSrc, label }) => {
  const [imgError, setImgError] = React.useState(false);
  if (!iconSrc) {
    return <PatternIcon fill="#fff" height="1.5rem" width="1.5rem" />;
  }
  if (imgError) {
    return <ModelName title={label}>{label}</ModelName>;
  }

  return (
    <div style={{ height: '24px', width: '24px' }}>
      <img
        src={iconSrc}
        style={{
          height: '100%',
          width: '100%',
          objectFit: 'contain',
        }}
        onError={() => setImgError(true)}
      />
    </div>
  );
};

export const Loading = ({ message, 'data-testid': testId = 'loading' }) => {
  return (
    <Box
      display="flex"
      alignItems="center"
      flexDirection="column"
      spacing={1}
      justifyContent="center"
      data-testid={testId}
    >
      <CircularProgress data-testid={`${testId}-spinner`} />
      <Typography variant="body2" style={{ marginTop: '1rem' }} data-testid={`${testId}-message`}>
        {message}
      </Typography>
    </Box>
  );
};

export const getSvgWhiteForComponent = (component) => {
  return `/${component.styles.svgWhite}`;
};

export const CheckBoxField = ({
  label,
  checked,
  onChange,
  disabled = false,
  helpText = 'Open in Operator',
  'data-testid': testId = 'checkbox-field',
}) => {
  const theme = useTheme();
  const color = disabled ? theme.palette.text.disabled : theme.palette.text.neutral.default;

  return (
    <Stack
      spacing={2}
      direction="row"
      alignItems="center"
      justifyContent="space-between"
      style={{
        cursor: disabled && 'not-allowed',
      }}
      data-testid={testId}
    >
      <Stack
        spacing={2}
        direction="row"
        alignItems="center"
        data-testid={`${testId}-label-section`}
      >
        <Checkbox
          value={checked}
          onChange={onChange}
          disabled={disabled}
          slotProps={{
            input: { 'data-testid': `${testId}-checkbox` },
          }}
        />
        <Typography variant="body2" color={color} data-testid={`${testId}-label`}>
          {label}
        </Typography>
      </Stack>
      {helpText && (
        <CustomTooltip title={helpText} placement="top">
          <IconButton data-testid={`${testId}-help-icon`}>
            <InfoCircleIcon fill={color} />
          </IconButton>
        </CustomTooltip>
      )}
    </Stack>
  );
};

export const StepHeading = ({ children }) => {
  return <Typography variant="textB2SemiBold"> {children}</Typography>;
};

export { processDesign };
