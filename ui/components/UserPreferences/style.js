const styles = (theme) => ({
  orgSelect: {
    '&.Mui-focused .MuiOutlinedInput-notchedOutline': {
      borderColor: 'transparent !important',
    },
    width: '85%',
    backgroundColor: theme.palette.primary,
    color: theme.palette.secondary.text,
    '@media (max-width: 368px)': {
      width: '180px',
    },
    position: 'relative',
    overflow: 'hidden',
    '& svg': {
      backgroundColor: theme.palette.white,
      fill: theme.palette.charcoal,
      float: 'right',
      top: '0',
      right: '0',
      height: '100%',
    },
  },
  selectItem: {
    width: '100%',
    '& svg': {
      backgroundColor: 'inherit',
    },
  },
  org: {
    marginLeft: '10px',
    overflow: 'hidden',
    whiteSpace: 'nowrap',
    textOverflow: 'ellipsis',
  },
  orgIconWrapper: {
    width: '32px',
    height: '100%',
    position: 'relative',
  },
  formControlWrapper: {
    padding: 20,
    border: '1.5px solid #969696',
    display: 'flex',
    width: '70%',
  },
  formLabelWrapper: {
    fontSize: 20,
  },
  formContainerWrapper: {
    display: 'flex',
    'flex-wrap': 'wrap',
    'justify-content': 'space-evenly',
    padding: 50,
  },
});

export default styles;
