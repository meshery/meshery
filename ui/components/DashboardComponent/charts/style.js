import { styled } from '@layer5/sistent';

export const DashboardSection = styled('div')(({ theme }) => ({
  backgroundColor: theme.palette.mode === 'dark' ? '#202020' : '#ffffff',
  padding: theme.spacing(2),
  borderRadius: '4px',
  height: '100%',
}));

export const LoadingContainer = styled('div')({
  height: '17rem',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  paddingBottom: '2.5rem',
});
