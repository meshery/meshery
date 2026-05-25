import { styled, Box } from '@sistent/sistent';

export const TableIconsContainer = styled('span')(({ disabled, theme }) => ({
  color: disabled ? theme.palette.icon.disabled : theme.palette.icon.default,
  pointerEvents: disabled ? 'none' : 'normal',
  display: 'flex',
  cursor: 'not-allowed',
  '& svg': {
    cursor: 'pointer',
  },
}));

export const IconWrapper = styled('div')(({ disabled = false }) => ({
  cursor: disabled ? 'not-allowed' : 'pointer',
  opacity: disabled ? '0.5' : '1',
  display: 'flex',
  '& svg': {
    cursor: disabled ? 'not-allowed' : 'pointer',
  },
}));

export const CreateButtonWrapper = styled('div')({
  display: 'flex',
  justifyContent: 'flex-start',
  alignItems: 'center',
  whiteSpace: 'nowrap',
});

export const BulkActionWrapper = styled(`div`)({
  width: '100%',
  padding: '0.8rem',
  justifyContent: 'space-between',
  marginTop: '0.18rem',
  marginBottom: '1rem',
  borderRadius: '.25rem',
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
});

export const UserCommonBox = styled(Box)(({ theme }) => ({
  backgroundColor: theme.palette.background.secondary,
  borderRadius: '.25rem',
  width: '100%',
  display: 'flex',
  flexDirection: 'row',
  justifyContent: 'space-between',
  alignItems: 'center',
}));
