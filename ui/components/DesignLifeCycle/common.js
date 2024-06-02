import React from 'react';
import { Box, Checkbox, CircularProgress, Stack, Typography } from '@layer5/sistent';
import PatternIcon from '@/assets/icons/Pattern';
import { processDesign } from '@/utils/utils';

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
  return `/${component.metadata.svgWhite}`;
};

export const CheckBoxField = ({ label, checked, onChange }) => {
  return (
    <Stack spacing={2} direction="row" alignItems="center">
      <Checkbox value={checked} onChange={onChange} />
      <Typography variant="body2">{label}</Typography>
    </Stack>
  );
};

export { processDesign };
