import React from 'react';
import { Box, Checkbox, CircularProgress, Stack, Typography } from '@layer5/sistent';
import PatternIcon from '@/assets/icons/Pattern';
import { processDesign } from '@/utils/utils';
import { CustomTooltip } from '@layer5/sistent';
import { InfoCircleIcon, useTheme } from '@layer5/sistent';
import { IconButton } from '@layer5/sistent';

export const DEPLOYMENT_TYPE = {
  DEPLOY: 'deploy',
  UNDEPLOY: 'undeploy',
};

export const ComponentIcon = ({ iconSrc }) => {
  if (!iconSrc) {
    return <PatternIcon fill="#fff" height="1.5rem" width="1.5rem" />;
  }

  return (
    <div style={{ height: '1.5rem', width: '1.5rem' }}>
      <img
        src={iconSrc}
        style={{
          height: '100%',
          width: '100%',
          objectFit: 'contain',
        }}
      />
    </div>
  );
};

export const Loading = ({ message }) => {
  return (
    <Box
      display="flex"
      alignItems="center"
      flexDirection="column"
      spacing={1}
      justifyContent="center"
    >
      <CircularProgress />
      <Typography variant="body2" style={{ marginTo: '1rem' }}>
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
  helpText = 'Open in visualizer',
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
    >
      <Stack spacing={2} direction="row" alignItems="center">
        <Checkbox value={checked} onChange={onChange} disabled={disabled} />
        <Typography variant="body2" color={color}>
          {label}
        </Typography>
      </Stack>
      {helpText && (
        <CustomTooltip title={helpText} placement="top">
          <IconButton>
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
