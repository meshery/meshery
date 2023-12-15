import { styled } from '@mui/system';
import { Select, MenuItem } from '@mui/material';

export const OrgSelect = styled(Select)(({ theme }) => ({
  '&.Mui-focused .MuiOutlinedInput-notchedOutline': {
    borderColor: 'transparent !important',
  },
  width: '85%',
  backgroundColor: theme.palette.white,
  color: theme.palette.darkShadeGray,
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
}));
export const SelectItem = styled(MenuItem)(() => ({
  width: '100%',
  '& svg': {
    backgroundColor: 'none !important',
  },
}));

export const Org = styled('span')(() => ({
  marginLeft: '10px',
  overflow: 'hidden',
  whiteSpace: 'nowrap',
  textOverflow: 'ellipsis',
}));

export const OrgIconWrapper = styled('div')(() => ({
  width: '32px',
  height: '100%',
  position: 'relative',
}));
