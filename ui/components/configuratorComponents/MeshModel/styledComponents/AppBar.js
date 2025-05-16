import { AppBar, styled } from '@layer5/sistent';

const AppBarComponent = styled(AppBar)(({ theme }) => {
  return {
    marginBottom: 16,
    backgroundColor: theme.palette.mode === 'dark' ? '#363636' : '#fff',
    borderRadius: '8px',
    color: theme.palette.text.default,
  };
});

export default AppBarComponent;
