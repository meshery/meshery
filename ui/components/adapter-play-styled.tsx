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
  width: '100%',
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

interface AdapterOperationsLayoutProps {
  hasAddons?: boolean;
}

export const AdapterOperationsLayout = styled('div', {
  shouldForwardProp: (prop) => prop !== 'hasAddons',
})<AdapterOperationsLayoutProps>(({ theme, hasAddons = false }) => ({
  display: 'grid',
  gridTemplateColumns: 'minmax(0, 1fr)',
  gap: theme.spacing(2),
  alignItems: 'stretch',
  width: '100%',
  minWidth: 0,
  [theme.breakpoints.up('lg')]: hasAddons
    ? {
        gridTemplateColumns: 'minmax(0, 5fr) minmax(16rem, 1fr)',
      }
    : {},
}));

export const AdapterCategoryGrid = styled('div')(({ theme }) => ({
  display: 'grid',
  gridTemplateColumns: 'repeat(auto-fit, minmax(16rem, 1fr))',
  gap: theme.spacing(2),
  alignItems: 'stretch',
  minWidth: 0,
}));

export const AdapterAddonPanel = styled('div')(() => ({
  display: 'flex',
  alignItems: 'stretch',
  minWidth: 0,
}));
