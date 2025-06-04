import { styled, Box, Typography, Link } from '@sistent/sistent';

export const StyledSummaryBox = styled(Box)(({ theme }) => ({
  backgroundColor: theme.palette.background.blur?.heavy,
  padding: '1rem',
  borderRadius: '8px',
  marginBottom: '1rem',
}));

export const StyledSummaryItem = styled(Box)(({ theme }) => ({
  borderRadius: '0.5rem',
  padding: '0.7rem',
  backgroundColor: theme.palette.background.hover,
  flexGrow: 1,
}));

export const SectionHeading = styled(Typography)(() => ({
  fontWeight: 'bold',
  marginTop: '1.5rem',
  marginBottom: '1rem',
}));

export const StyledColorBox = styled(Box)(({ color }) => ({
  width: '1.5rem',
  height: '1.5rem',
  borderRadius: '4px',
  backgroundColor: color,
  marginRight: '0.5rem',
  display: 'inline-block',
  verticalAlign: 'middle',
}));

export const StyledDocsRedirectLink = styled(Link)(({ theme }) => ({
  color: theme.palette.background.brand.default,
  textDecoration: 'underline',
}));
