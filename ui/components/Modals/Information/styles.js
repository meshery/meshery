import { Box, styled, Button, Typography } from '@layer5/sistent';

export const CreatAtContainer = styled(Typography)(({ isBold }) => ({
  fontWeight: isBold ? 'bold' : '',
  whiteSpace: 'wrap',
  fontSize: '0.8rem',
}));

export const ActionContainer = styled(Box)({
  width: '100%',
  display: 'flex',
  gap: '1rem',
  justifyContent: 'end',
});

export const CopyLinkButton = styled(Button)(({ theme }) => ({
  color: theme.palette.background.constant.white,
}));

export const ResourceName = styled(Typography)(() => ({
  fontFamily: 'Qanelas Soft, sans-serif',
  textAlign: 'left',
  marginTop: '0.5rem',
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  maxWidth: '10rem',
}));
