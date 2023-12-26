import { Typography, Card, Checkbox, styled } from '@material-ui/core';

/** Card Styles */
export const CardWrapper = styled(Card)({
  width: '100%',
  display: 'flex',
  flexDirection: 'column',
  padding: '20px',
  '&:hover': {
    cursor: 'pointer',
  },
});

export const BulkSelectCheckbox = styled(Checkbox)({
  padding: 0,
  marginRight: '0.5rem',
  height: '28px',
  '& .MuiSvgIcon-root': {
    borderColor: 'white',
  },
  color: 'white',
  '&:hover': {
    color: 'white',
    cursor: 'pointer',
  },
  '&.Mui-checked': {
    color: 'white',
  },
});

export const CardTitle = styled(Typography)({
  fontSize: '1.25rem',
  fontWeight: 800,
  '&:hover': {
    cursor: 'default',
  },
});

export const OrganizationName = styled(Typography)({
  fontSize: '0.9rem',
  display: 'flex',
  alignItems: 'end',
  padding: '0 5px',
  '&:hover': {
    cursor: 'default',
  },
});

export const StyledIconButton = styled('button')({
  background: 'transparent',
  border: 'none',
  '&:hover': {
    cursor: 'default',
  },
});

export const DateLabel = styled(Typography)({
  fontStyle: 'italic',
  fontSize: '12px',
  '&:hover': {
    cursor: 'default',
  },
});

export const EmptyDescription = styled(Typography)({
  fontSize: '0.9rem',
  textAlign: 'left',
  fontStyle: 'italic',
});

export const DescriptionLabel = styled(EmptyDescription)({
  height: 'fit-content',
  fontStyle: 'normal',
  '&:hover': {
    cursor: 'default',
  },
});
