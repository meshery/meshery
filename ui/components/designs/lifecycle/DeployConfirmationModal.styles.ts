import { Button, DialogContentText } from '@sistent/sistent';
import { styled } from '@/theme';

export const ActionButton = styled(Button, {
  shouldForwardProp: (prop) => !['isUndeploy', 'isDisabled'].includes(prop as string),
})<{ isUndeploy?: boolean; isDisabled?: boolean }>(({ theme, isUndeploy, isDisabled }) => ({
  margin: theme.spacing(0.5),
  padding: theme.spacing(1),
  borderRadius: 5,
  minWidth: 100,
  ...(isUndeploy &&
    !isDisabled && {
      backgroundColor: theme.palette.error.main,
      '&:hover': {
        backgroundColor: theme.palette.error.dark,
        boxShadow: theme.shadows[8],
      },
    }),
  ...(!isUndeploy && {
    color: theme.palette.common.white,
    '&:hover': {
      boxShadow: theme.shadows[8],
    },
  }),
  ...(isDisabled && {
    '&.Mui-disabled': {
      cursor: 'not-allowed',
      pointerEvents: 'all !important',
    },
  }),
}));

export const ActionsRow = styled('div')({
  display: 'flex',
  justifyContent: 'space-evenly',
  width: '100%',
});

export const ContextsContainer = styled('div')({
  display: 'flex',
  flexWrap: 'wrap',
});

export const TabLabelWrapper = styled('span')(({ theme }) => ({
  [theme.breakpoints.up('sm')]: {
    fontSize: '1em',
  },
  [theme.breakpoints.between('xs', 'sm')]: {
    fontSize: '0.8em',
  },
  color: theme.palette.icon.default,
}));

export const TriangleContainer = styled('div')({
  position: 'relative',
  marginLeft: 2,
});

export const TriangleNumber = styled('div')(({ theme }) => ({
  position: 'absolute',
  bottom: 12,
  left: '37%',
  color: theme.palette.common.white,
  fontSize: '0.8rem',
}));

export const OctagonContainer = styled('div')({
  overflow: 'hidden',
  position: 'relative',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  height: 34,
  marginLeft: 2,
});

export const OctagonText = styled('div')(({ theme }) => ({
  position: 'absolute',
  bottom: 9.5,
  color: theme.palette.common.white,
  fontSize: '0.8rem',
}));

export const DialogSubtitle = styled(DialogContentText)({
  overflowWrap: 'anywhere',
  textAlign: 'center',
  padding: '5px',
});
