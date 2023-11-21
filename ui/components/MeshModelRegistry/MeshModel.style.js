import { styled } from '@mui/material/styles';
import { Button, alpha } from '@material-ui/core';
import { TreeItem, treeItemClasses } from '@mui/x-tree-view/TreeItem';

export const DisableButton = styled(Button)(({ theme }) => ({
  '&.MuiButtonBase-root:disabled': {
    cursor: 'not-allowed',
    pointerEvents: 'auto',
    backgroundColor: theme.palette.secondary.disableButtonBg,
    color: theme.palette.secondary.disableButton,
  },
}));

export const StyledTreeItemRoot = styled(TreeItem)(({ theme, root, lineColor }) => ({
  position: 'relative',
  '&:before': {
    pointerEvents: 'none',
    content: '""',
    position: 'absolute',
    width: 32,
    left: -34,
    top: 23,
    borderBottom: !root ? `1px dashed ${alpha(lineColor, 0.4)}` : 'none',
  },

  [`& .${treeItemClasses.content}`]: {
    fontWeight: theme.typography.fontWeightMedium,
    borderRadius: '0px 4px 4px 0px',
    '&.Mui-expanded': {
      fontWeight: theme.typography.fontWeightRegular,
    },
    '&:hover': {
      backgroundColor: `transparent`,
    },
    '&.Mui-focused, &.Mui-selected, &.Mui-selected.Mui-focused': {
      backgroundColor: `#00bfa030`,
      borderLeft: '3px solid #00bfa0',
    },
    [`& .${treeItemClasses.label}`]: {
      fontWeight: 'inherit',
    },
  },
  [`& .${treeItemClasses.group}`]: {
    // marginLeft: 34,
    paddingLeft: 36,
    borderLeft: `1px dashed ${alpha(lineColor, 0.4)}`,
    borderOpacity: 0.5,
  },
}));
