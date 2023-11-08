import { styled, alpha } from '@mui/material/styles';
import { Button } from '@material-ui/core';
import { TreeItem, treeItemClasses } from '@mui/x-tree-view/TreeItem';

export const DisableButton = styled(Button)(({ theme }) => ({
  '&.MuiButtonBase-root:disabled': {
    cursor: 'not-allowed',
    pointerEvents: 'auto',
    backgroundColor: theme.palette.secondary.disableButtonBg,
    color: theme.palette.secondary.disableButton,
  },
}));

export const StyledTreeItemRoot = styled(TreeItem)(({ theme }) => ({
  [`& .${treeItemClasses.content}`]: {
    fontWeight: theme.typography.fontWeightMedium,
    borderRadius: '4px',
    '&.Mui-expanded': {
      fontWeight: theme.typography.fontWeightRegular,
    },
    '&:hover': {
      backgroundColor: `transparent`,
    },
    '&.Mui-focused, &.Mui-selected, &.Mui-selected.Mui-focused': {
      backgroundColor: `#00bfa030`,
      borderLeft: '2px solid #00bfa0',
    },
    [`& .${treeItemClasses.label}`]: {
      fontWeight: 'inherit',
    },
  },
  [`& .${treeItemClasses.group}`]: {
    paddingRight: '0',
    borderLeft: `1px solid ${
      theme.palette.type === 'dark' ? alpha('#00B39F', 0.4) : alpha(theme.palette.text.primary, 0.4)
    }`,
  },
}));
