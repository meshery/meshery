import { styled } from '@layer5/sistent';
import { Typography } from '@layer5/sistent';

export const DashboardSection = styled('div')(({ theme }) => ({
  backgroundColor: theme.palette.mode === 'dark' ? '#202020' : theme.palette.common.white,
  padding: theme.spacing(2),
  borderRadius: '4px',
  height: '100%',
}));

export const ChartSectionWithColumn = styled('div')(({ theme }) => ({
  backgroundColor: theme.palette.mode === 'dark' ? '#202020' : theme.palette.common.white,
  padding: theme.spacing(2),
  borderRadius: '4px',
  height: '100%',
  display: 'flex',
  flexDirection: 'column',
}));

export const LoadingContainer = styled('div')({
  height: '17rem',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  paddingBottom: '2.5rem',
});

export const LegendSection = styled('div')({
  display: 'flex',
  gap: '2rem',
  flexWrap: 'wrap',
});
export const ConnectClusterWrapper = styled('div')({
  padding: '2rem',
  display: 'flex',
  justifyContent: 'center',
  alignItems: 'center',
  flexDirection: 'column',
});

export const ConnectClusterText = styled(Typography)({
  marginBottom: '0.5rem',
});

export const HoneycombRoot = styled('div')(({ theme }) => ({
  backgroundColor: theme.palette.mode === 'dark' ? '#202020' : theme.palette.common.white,
  padding: theme.spacing(2),
  borderRadius: 4,
  width: '100%',
}));

export const HoneycombCell = styled('li')(({ row, column }) => ({
  gridRow: `${row} / span 4`,
  gridColumn: `${column} / span 4`,
  position: 'relative',
  pointerEvents: 'none',
  transform: row % 2 ? 'translateX(24%)' : 'translateX(-26%)',
}));

export const HoneycombContainer = styled('ul')(({ columnSize, columns, rowSize }) => ({
  display: 'grid',
  gridTemplateColumns: `repeat(${columns * 4}, ${columnSize}px)`,
  justifyContent: 'center',
  gridAutoRows: `${rowSize}px`,
  padding: `0 ${columnSize}px`,
  listStyle: 'none',
}));

export const HexagonWrapper = styled('div')(({ theme }) => ({
  backgroundColor: theme.palette.mode === 'dark' ? '#363636' : '#e9eff1', // TODO: this is the honeycomb color add this token in sistent
  position: 'absolute',
  inset: '3.5px',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  clipPath: 'polygon(50% 0%, 100% 25%, 100% 75%, 50% 100%, 0% 75%, 0% 25%)',
  pointerEvents: 'auto',
  boxSizing: 'border-box',
}));

export const IconWrapper = styled('div')({
  marginTop: '20px',
  display: 'flex',
  flexDirection: 'column',
});

export const ResourceCount = styled(Typography)({
  margin: '0px auto',
});

export const SelectedHexagon = styled('div')({
  display: 'flex',
  height: '95%',
  justifyContent: 'center',
  alignItems: 'center',
  '&:hover': {
    cursor: 'pointer',
  },
});

export const SkeletonHexagon = styled('div')(({ theme }) => ({
  display: 'flex',
  height: '95%',
  backgroundColor: theme.palette.mode === 'dark' ? '#363636' : '#e9eff1', // TODO: this is the honeycomb color add this token in sistent
  justifyContent: 'center',
  alignItems: 'center',
  opacity: 0.5,
}));
