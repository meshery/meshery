import { Box, DialogTitle, Typography, styled } from '@sistent/sistent';
import { AddCircleOutlined as AddIcon } from '@/assets/icons';

export const ViewSwitchButton = styled(Box)(() => ({
  justifySelf: 'flex-end',
}));

export const CreateButton = styled(Box)(() => ({
  display: 'flex',
  justifyContent: 'flex-start',
  alignItems: 'center',
  whiteSpace: 'nowrap',
}));

export const AddIconStyled = styled(AddIcon)(() => ({
  paddingRight: '.35rem',
}));

export const SearchWrapper = styled(Box)(({ theme }) => ({
  justifySelf: 'flex-end',
  marginLeft: 'auto',
  paddingLeft: '1rem',
  display: 'flex',
  alignItems: 'center',
  '@media (max-width: 965px)': {
    width: 'max-content',
  },
  [theme.breakpoints.down('sm')]: {
    width: '100%',
    minWidth: 0,
    marginLeft: 0,
    paddingLeft: 0,
    flexWrap: 'wrap',
    gap: '0.5rem',
  },
}));

export const BtnText = styled('span')(() => ({
  display: 'block',
  '@media (max-width: 765px)': {
    display: 'none',
  },
}));

export const YamlDialogTitle = styled(DialogTitle)(() => ({
  display: 'flex',
  justifyContent: 'space-between',
  alignItems: 'end',
}));

export const YamlDialogTitleText = styled(Typography)(() => ({
  flexGrow: 1,
}));
