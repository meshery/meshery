import { Box, Button, styled } from '@layer5/sistent';

export const FallbackWrapper = styled(Box)(() => ({
  margin: '2rem',
}));

export const TryAgainButton = styled(Button)(({ theme }) => ({
  background: 'transparent',
  border: `1px solid ${theme.palette.border.brand}`,
  '&:hover': {
    border: `1px solid ${theme.palette.border.brand}`,
  },
}));

export const EditButton = styled(Button)(({ theme }) => ({
  backgroundImage: theme.palette.background.brand.prominent,
  backgroundColor: `${
    theme.palette.background.brand.prominent || theme.palette.background.brand.default
  }`,
  '@media (max-width: 768px)': {
    minWidth: '50px',
  },
}));

export const TextButton = styled('span')(({ style }) => ({
  marginLeft: '0.5rem',
  display: 'block',
  '@media (max-width: 853px)': {
    display: 'none',
  },
  ...style,
}));

export const ToolBarButtonContainer = styled('span')(() => ({
  margin: '0 1rem 0 0',
  '@media (max-width: 400px)': {
    margin: '0 0.25rem 0 0',
  },
}));
