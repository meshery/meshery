import { Card, Chip, Grid, TableCell, styled } from '@sistent/sistent';

export const AdapterChip = styled(Chip)(({ theme }) => ({
  height: '50px',
  fontSize: '15px',
  position: 'relative',
  top: theme.spacing(0.5),
  [theme.breakpoints.down('md')]: {
    fontSize: '12px',
  },
}));

export const AdapterTableHeader = styled(TableCell)({
  fontWeight: 'bolder',
  fontSize: 18,
});

export const AdapterSmWrapper = styled('div')(({ theme }) => ({
  backgroundColor: theme.palette.background.card,
}));

export const SecondaryTable = styled('div')({
  borderRadius: 10,
  backgroundColor: '#f7f7f7',
});

export const PaneSection = styled('div')(({ theme }) => ({
  backgroundColor: theme.palette.background.tabs,
  padding: theme.spacing(3),
  borderRadius: 4,
}));

export const ChipNamespaceContainer = styled(Grid)(() => ({
  gap: '2rem',
  margin: '0px',
}));

export const InputWrapper = styled('div')(() => ({
  flex: '1',
  minWidth: '250px',
}));

export const AdapterCard = styled(Card)(() => ({
  height: '100%',
  display: 'flex',
  flexDirection: 'column',
}));
