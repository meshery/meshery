import { Grid, InputLabel, Paper, Typography, styled } from '@layer5/sistent';

export const AccordionRoot = styled('div')(() => ({
  width: '100%',
  marginBottom: '0rem',
}));

export const AccordionHeading = styled(Typography)(({ theme }) => ({
  fontSize: theme.typography.pxToRem(15),
  // fontWeight: theme.typography.fontWeightRegular,
}));

export const ArrayFieldWrapper = styled(Paper)(() => ({
  '& .MuiPaper-root': {
    backgroundColor: '#f4f4f4',
  },
}));

export const CustomInputLabel = styled(InputLabel)(({ theme }) => ({
  color: theme.palette.text.secondary,
  backgroundColor: theme.palette.background.default,
  padding: '0.2rem',
  height: '1rem',
  borderRadius: '3px',
}));

export const ObjectFieldGrid = styled(Grid)(() => ({
  padding: '.5rem',
  paddingTop: '0.7rem',
  width: '100%',
  margin: '0px',
}));
