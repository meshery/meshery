import { Box, Card, Chip, styled } from '@sistent/sistent';
import { alpha } from '@/theme';
import { CompatibilityStatus } from './compatibilityAnalysisUtils';

export const RootContainer = styled(Box)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  gap: theme.spacing(2.5),
  padding: theme.spacing(2),
  width: '100%',
  backgroundColor: theme.palette.background.default,
  color: theme.palette.text.primary,
}));

export const HeaderContainer = styled(Box)(({ theme }) => ({
  display: 'flex',
  justifyContent: 'space-between',
  alignItems: 'center',
  flexWrap: 'wrap',
  gap: theme.spacing(1),
  borderBottom: `1px solid ${theme.palette.divider}`,
  paddingBottom: theme.spacing(1.5),
}));

export const MetricsGrid = styled(Box)(({ theme }) => ({
  display: 'grid',
  gridTemplateColumns: 'repeat(auto-fit, minmax(180px, 1fr))',
  gap: theme.spacing(2),
}));

export const MetricCard = styled(Card, {
  shouldForwardProp: (prop) => prop !== 'accentColor',
})<{ accentColor?: string }>(({ theme, accentColor }) => ({
  display: 'flex',
  alignItems: 'center',
  padding: theme.spacing(2),
  gap: theme.spacing(1.5),
  backgroundColor: theme.palette.background.paper,
  border: `1px solid ${accentColor ? alpha(accentColor, 0.4) : theme.palette.divider}`,
  borderRadius: theme.shape.borderRadius,
  transition: 'transform 0.2s ease, box-shadow 0.2s ease',
  '&:hover': {
    transform: 'translateY(-2px)',
    boxShadow: `0 4px 12px ${alpha(accentColor || theme.palette.common.black, 0.15)}`,
  },
}));

export const MatrixCard = styled(Card)(({ theme }) => ({
  backgroundColor: theme.palette.background.paper,
  border: `1px solid ${theme.palette.divider}`,
  borderRadius: theme.shape.borderRadius,
  overflow: 'hidden',
}));

export const TwoColumnGrid = styled(Box)(({ theme }) => ({
  display: 'grid',
  gridTemplateColumns: '2fr 1.2fr',
  gap: theme.spacing(2.5),
  alignItems: 'start',
  [theme.breakpoints.down('md')]: {
    gridTemplateColumns: '1fr',
  },
}));

export const BottomGrid = styled(Box)(({ theme }) => ({
  display: 'grid',
  gridTemplateColumns: '1fr 1fr',
  gap: theme.spacing(2.5),
  [theme.breakpoints.down('sm')]: {
    gridTemplateColumns: '1fr',
  },
}));

export const DetailCard = styled(Card)(({ theme }) => ({
  backgroundColor: theme.palette.background.paper,
  border: `1px solid ${theme.palette.divider}`,
  borderRadius: theme.shape.borderRadius,
  padding: theme.spacing(2.5),
  display: 'flex',
  flexDirection: 'column',
  gap: theme.spacing(2),
}));

export const DetailSection = styled(Box)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  gap: theme.spacing(0.75),
}));

export const CodeBox = styled(Box)(({ theme }) => ({
  backgroundColor: alpha(theme.palette.background.default, 0.8),
  border: `1px solid ${theme.palette.divider}`,
  borderRadius: theme.shape.borderRadius,
  padding: theme.spacing(1.5),
  fontFamily: 'monospace',
  fontSize: '0.85rem',
  color: theme.palette.error.main,
  wordBreak: 'break-word',
  maxHeight: '180px',
  overflowY: 'auto',
}));

export const InsightBox = styled(Box, {
  shouldForwardProp: (prop) => prop !== 'tone',
})<{ tone?: 'fix' | 'dependency' }>(({ theme, tone }) => {
  const isFix = tone === 'fix';
  const baseColor = isFix ? theme.palette.primary.main : theme.palette.warning.main;
  return {
    backgroundColor: alpha(baseColor, 0.08),
    borderLeft: `4px solid ${baseColor}`,
    borderRadius: `0 ${theme.shape.borderRadius}px ${theme.shape.borderRadius}px 0`,
    padding: theme.spacing(1.5),
    fontSize: '0.9rem',
    lineHeight: 1.5,
  };
});

export const StatusChipStyled = styled(Chip, {
  shouldForwardProp: (prop) => prop !== 'status' && prop !== 'isSelected',
})<{ status: CompatibilityStatus; isSelected?: boolean }>(({ theme, status, isSelected }) => {
  let color = theme.palette.success.main;
  let bg = alpha(theme.palette.success.main, 0.15);

  if (status === 'failed') {
    color = theme.palette.error.main;
    bg = alpha(theme.palette.error.main, 0.15);
  } else if (status === 'warning') {
    color = theme.palette.warning.main;
    bg = alpha(theme.palette.warning.main, 0.15);
  } else if (status === 'incomplete') {
    color = theme.palette.text.secondary;
    bg = alpha(theme.palette.action.disabled, 0.15);
  }

  return {
    color,
    backgroundColor: bg,
    border: `1px solid ${isSelected ? color : 'transparent'}`,
    fontWeight: 600,
    cursor: 'pointer',
    transition: 'all 0.2s ease',
    '&:hover': {
      backgroundColor: alpha(color, 0.25),
      transform: 'scale(1.02)',
    },
  };
});
