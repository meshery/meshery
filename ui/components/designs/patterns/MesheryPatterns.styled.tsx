import { DialogTitle, Typography, styled } from '@sistent/sistent';
import { AddCircleOutlined as AddIcon } from '@/assets/icons';

export const AddIconStyled = styled(AddIcon)(() => ({
  paddingRight: '.35rem',
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
