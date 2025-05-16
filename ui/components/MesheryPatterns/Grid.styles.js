import { Paper, styled, Typography } from '@layer5/sistent';

export const GridNoPapperStyles = styled(Paper)({
  padding: '0.5rem',
  fontSize: '3rem',
});

export const GridNoContainerStyles = styled('div')({
  padding: '2rem',
  display: 'flex',
  justifyContent: 'center',
  alignItems: 'center',
  flexDirection: 'column',
});

export const GridNoTextStyles = styled(Typography)({
  fontSize: '2rem',
  marginBottom: '2rem',
});

export const GridPaginationStyles = styled('div')({
  display: 'flex',
  justifyContent: 'center',
  alignItems: 'center',
  marginTop: '2rem',
});

export const GridAddIconStyles = styled('div')({
  paddingLeft: '0.5rem',
  marginRight: '0.5rem',
});
