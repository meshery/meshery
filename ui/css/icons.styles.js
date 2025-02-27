// icon styles, setting general height and width properties to solves scaling and consistency problems

import { Typography, styled } from '@layer5/sistent';

export const iconSmall = {
  height: 20,
  width: 20,
};

export const iconMedium = {
  height: 24,
  width: 24,
};

export const iconLarge = {
  height: 32,
  width: 32,
};

export const iconXLarge = {
  height: 40,
  width: 40,
};

export const CardContainer = styled('div')(({ theme }) => ({
  padding: theme.spacing(4),
  borderRadius: theme.spacing(1),
  transformStyle: 'preserve-3d',
  boxShadow: '0 4px 8px 0 rgba(0,0,0,0.2)',
  backgroundColor: theme.palette.background.card,
  minHeight: '250px',
  position: 'relative',
}));

export const ImageWrapper = styled('img')(() => ({
  paddingRight: '1rem',
  height: '80px',
  width: '80px',
  flexShrink: 0,
}));

export const FrontSideDescription = styled(Typography)({
  paddingTop: '1rem',
  paddingBottom: '1rem',
  textAlign: 'left',
  display: 'flex',
  flexDirection: 'row',
});
