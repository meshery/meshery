import { alpha } from '@mui/system';
import { TextField, Typography, styled } from '@layer5/sistent';
export const Root = styled('div')(() => ({
  position: 'relative',
}));

export const InputField = styled(TextField)(({ theme }) => ({
  width: '100%',
  marginBottom: '.1rem',
  '& .MuiOutlinedInput-root': {
    borderRadius: '6px',
    backgroundColor: theme.palette.mode === 'dark' ? '#294957' : '#fff',
    '& fieldset': {
      borderRadius: '6px',
      border: `2px solid ${theme.palette.mode === 'dark' ? '#396679' : '#CCCCCC'}`,
    },
  },
}));
export const DropDown = styled('div')(({ theme }) => ({
  backgroundColor: theme.palette.mode === 'dark' ? '#294957' : '#fff',
  borderRadius: '6px',
  boxShadow:
    '0px 2px 4px 0px rgba(0, 0, 0, 0.20), 0px 1px 10px 0px rgba(0, 0, 0, 0.12), 0px 4px 5px 0px rgba(0, 0, 0, 0.14)',
  border: `2px solid ${theme.palette.background.brand.secondary}`,
  marginTop: '0.2rem',
  maxHeight: '20rem',
  overflowY: 'auto',
}));
export const Item = styled('div')(({ theme }) => ({
  display: 'flex',
  gap: '0.3rem',
  margin: '0.3rem',
  padding: '0.3rem',
  paddingInline: '3rem',
  borderRadius: '6px',
  cursor: 'pointer',
  '&:hover': {
    backgroundColor: alpha(theme.palette.text.brand, 0.25),
  },
}));
export const Label = styled(Typography)(({ theme }) => ({
  fontWeight: 500,
  color: theme.palette.secondary.icon,
}));

export const Description = styled(Typography)(({ theme }) => ({
  fontWeight: 400,
  color: theme.palette.secondary.number,
}));
