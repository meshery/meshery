import { styled } from '@mui/system';
import { Button, Chip, Grid, Paper, Typography } from '@mui/material';
// import theme from "../../../styles/themes/theme";

export const StyledChip = styled(Chip)({
  padding: '5px 6px !important',
  color: '#000000',
  fontSize: '14px',
  textTransform: 'uppercase',
  fontWeight: 400,
  height: 'unset',
  borderRadius: '100px',
  border: '0.5px solid #51636B',
  background: '#ffffff',
  maxWidth: '230px',
});

export const StyledPaper = styled(Paper)({
  width: 300,
  height: 280,
  overflow: 'auto',
  backgroundColor: '#EBEFF1',
  borderRadius: '10px',
  boxShadow: '0px 1px 4px 0px rgba(0, 0, 0, 0.25) inset',
  '@media (max-width: 843px)': {
    width: 260,
  },
  '@media (max-width: 768px)': {
    width: 300,
  },
  '@media (max-width: 375px)': {
    width: '100%',
  },
});

export const ListHeading = styled(Typography)({
  paddingBottom: '15px',
  color: '#525252',
  textAlign: 'center',
  fontSize: '1rem',
  letterSpacing: '0.15px',
});

export const TransferButton = styled(Button)({
  margin: '5px 0',
  padding: '7px 0',
  borderRadius: '10px',
  borderColor: '#455a64',
  boxShadow: 'none',
  '&:hover': {
    borderColor: '#00B39F',
  },
});

export const ListGrid = styled(Grid)({
  padding: '0 1rem',
  '@media (max-width: 768px)': {
    padding: '0',
    display: 'flex',
    alignItems: 'center',
    flexDirection: 'column',
    width: '100%',
  },
});

export const ButtonGrid = styled(Grid)({
  padding: '40px 1rem 0 1rem',
  '@media (max-width: 768px)': {
    padding: '1rem',
    transform: 'rotate(90deg)',
    height: '100px',
    marginLeft: '100px',
  },
});
