import { AppBar, styled } from '@material-ui/core';

const AppBarComponent = styled(AppBar)(({ theme }) => {
  return {
    marginBottom: 16,
    backgroundColor: theme.palette.secondary.appBar,
    borderRadius: '8px',
    color: theme.palette.secondary.text,
  };
});

export default AppBarComponent;
