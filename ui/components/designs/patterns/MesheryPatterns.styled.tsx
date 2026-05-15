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

export const SearchWrapper = styled(Box)(() => ({
  justifySelf: 'flex-end',
  marginLeft: 'auto',
  paddingLeft: '1rem',
  display: 'flex',
  '@media (max-width: 965px)': {
    width: 'max-content',
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
