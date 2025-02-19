import { styled, Button, alpha } from '@layer5/sistent';
import { TreeItem, treeItemClasses } from '@mui/x-tree-view/TreeItem';

export const DisableButton = styled(Button)(({ theme }) => ({
  '&.MuiButtonBase-root:disabled': {
    cursor: 'not-allowed',
    pointerEvents: 'auto',
    backgroundColor: theme.palette.background.brand.disabled,
    color: theme.palette.text.disabled,
  },
}));

export const JustifyAndAlignCenter = styled(`div`)(({ style, theme }) => ({
  display: 'flex',
  justifyContent: 'center',
  alignItems: 'center',
  width: '100%',
  height: '100%',
  fontFamily: theme.typography.fontFamily,
  ...style,
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
    fontWeight: theme.typography.fontWeightRegular,
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
      fontWeight: theme.typography.fontWeightRegular,
    },
  },
  [`& .${treeItemClasses.group}`]: {
    // marginLeft: 34,
    paddingLeft: 36,
    borderLeft: `1px dashed ${alpha(lineColor, 0.4)}`,
    borderOpacity: 0.5,
  },
}));

export const StyledKeyValuePropertyDiv = styled('div')(() => ({
  display: 'flex',
  flexDirection: 'column',
  margin: '0.6rem 0',
}));

export const StyledKeyValueProperty = styled('p')(({ theme }) => ({
  padding: '0',
  margin: '0 0.5rem 0 0',
  fontSize: theme.typography.htmlFontSize,
  fontWeight: theme.typography.fontWeightBold,
  fontFamily: theme.typography.fontFamily,
}));

export const StyledKeyValueFormattedValue = styled('div')(({ theme }) => ({
  padding: '0',
  margin: '0',
  fontSize: theme.typography.htmlFontSize,
  fontFamily: theme.typography.fontFamily,
}));

export const StyledTreeItemNameDiv = styled('div')(({ theme }) => ({
  fontWeight: theme.typography.fontWeightRegular,
}));

export const StyledTreeItemDiv = styled('div')(() => ({
  display: 'flex',
  alignItems: 'center',
  gap: '1rem',
}));
