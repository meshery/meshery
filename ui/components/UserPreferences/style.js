import { styled, Select, FormControl } from '@layer5/sistent';

export const StyledOrgSelect = styled(Select)(({ theme }) => ({
  '&.Mui-focused .MuiOutlinedInput-notchedOutline': {
    borderColor: 'transparent !important',
  },
  width: '85%',
  backgroundColor: theme.palette.primary.main,
  color: theme.palette.secondary.text,
  position: 'relative',
  overflow: 'hidden',
  '@media (max-width: 368px)': {
    width: '180px',
  },
  '& svg': {
    backgroundColor: theme.palette.common.white,
    fill: theme.palette.charcoal || '#333', // Fallback color
    float: 'right',
    top: '0',
    right: '0',
    height: '100%',
  },
}));

export const StyledSelectItem = styled('div')({
  width: '100%',
  '& svg': {
    backgroundColor: 'inherit',
  },
});

export const StyledOrgText = styled('span')({
  marginLeft: '10px',
  overflow: 'hidden',
  whiteSpace: 'nowrap',
  textOverflow: 'ellipsis',
});

export const StyledOrgIconWrapper = styled('div')({
  width: '32px',
  height: '100%',
  position: 'relative',
});

export const StyledFormControl = styled(FormControl)(({ theme }) => ({
  padding: theme.spacing(2.5),
  border: '1.5px solid #969696',
  display: 'flex',
  width: '70%',
}));

export const StyledFormContainer = styled('div')(({ theme }) => ({
  display: 'flex',
  flexWrap: 'wrap',
  justifyContent: 'space-evenly',
  padding: theme.spacing(6.25),
}));

// SX props object for components that don't need full styled component treatment
export const sxStyles = {
  formLabel: {
    fontSize: '20px',
  },
};
