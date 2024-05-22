import {
  Typography,
  Select,
  FormControlLabel,
  Box,
  TextField,
  Chip,
  Button,
  Divider,
  LeftArrowIcon,
} from '@layer5/sistent';
import { styled } from '@material-ui/core';

export const ErrorMain = styled('main')(({ theme }) => ({
  padding: '4rem 8rem',
  minHeight: '100vh',
  [theme.breakpoints.down('sm')]: {
    padding: '4rem 0rem',
  },
  display: 'flex',
  flexDirection: 'column',
  gap: '2rem',
}));

export const ErrorContainer = styled('div')(({ theme }) => ({
  backgroundColor: theme.palette.background.paper,
  boxShadow:
    '0px 2px 1px -1px rgb(0 0 0 / 20%), 0px 1px 1px 0px rgb(0 0 0 / 14%), 0px 1px 3px 1px rgb(0 0 0 / 12%)',
  borderRadius: '8px',
  paddingBlock: '2rem',
  display: 'flex',
  flexDirection: 'column',
  flexWrap: 'wrap',
}));

export const ImageContainer = styled('div')(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  gap: '2rem',
  marginBottom: '1rem',
  [theme.breakpoints.down('sm')]: {
    flexDirection: 'column',
  },
}));

export const ErrorSectionContainer = styled('div')(() => ({
  display: 'flex',
  flexWrap: 'wrap',
}));

export const ErrorSection = styled('div')(() => ({
  flex: '1',
  width: '100%',
  display: 'flex',
  flexDirection: 'column',
  paddingInline: '2rem',
  gap: '2rem',
  marginBlock: '2rem',
  minWidth: 'max-content',
}));

export const ErrorSectionContent = styled('div')(() => ({
  display: 'flex',
  flexDirection: 'column',
  justifyContent: 'space-between',
  gap: '1rem',
}));

export const ErrorContentContainer = styled('div')(() => ({
  textAlign: 'center',
  backgroundColor: '#fafafa',
  margin: '2rem',
  padding: '20px',
  borderRadius: '10px',
  boxShadow:
    ' 0 2px 2px 0 rgb(0 0 0 / 14%), 0 1px 5px 0 rgb(0 0 0 / 12%), 0 3px 1px -2px rgb(0 0 0 / 20%)',
}));

export const ErrorLink = styled('a')(({ theme }) => ({
  color: theme.palette.secondary.focused,
}));

export const ErrorMsg = styled('errormsg')(() => ({
  fontWeight: '600',
}));

export const SelectContainerWrapper = styled('div')(() => ({
  display: 'flex',
  flexWrap: 'wrap',
  justifyContent: 'space-evenly',
  padding: '50px',
}));

export const OrgName = styled('span')(() => ({
  marginLeft: '10px',
  overflow: 'hidden',
  whiteSpace: 'nowrap',
  textOverflow: 'ellipsis',
  fontSize: '1.1rem',
}));

export const OrgNameDisabled = styled('span')(({ theme }) => ({
  color: theme.palette.text.secondary,
  marginLeft: '10px',
  overflow: 'hidden',
  whiteSpace: 'nowrap',
  textOverflow: 'ellipsis',
  fontSize: '1.1rem',
}));

export const StyledSelect = styled(Select)(({ theme }) => ({
  '& .MuiOutlinedInput-notchedOutline': {
    borderColor: theme.palette.primary.main,
  },
  '&.MuiInputBase-root': {
    borderRadius: '10px',
    paddingBlock: '5px',
    paddingInline: '10px',
  },
  '@media (max-width: 368px)': {
    width: '180px',
  },
}));

export const StyledBox = styled(Box)(({ theme }) => ({
  display: 'flex',
  gap: '10px',
  alignItems: 'center',
  borderRadius: '10px',
  fontSize: '1.1rem',
  padding: '1.3rem',
  border: `1px solid ${theme.palette.action.disabled}`,
  color: theme.palette.text.secondary,
}));

export const FormField = styled('div')(() => ({
  display: 'flex',
  flexDirection: 'column',
  width: '100%',
}));

export const StyledFormControlLabel = styled(FormControlLabel)(() => ({
  '&.MuiFormControlLabel-root': {
    marginLeft: '0px',
    marginRight: '0px',
  },
}));

export const StyledTextField = styled(TextField)(({ theme }) => ({
  '& .MuiOutlinedInput-root': {
    color: theme.palette.text.primary,
    '& fieldset': {
      borderColor: theme.palette.primary.main,
      borderRadius: '10px',
    },
  },
}));

export const StyledTypography = styled(Typography)(() => ({
  fontSize: '1.3rem',
  margin: '0px',
  marginBlock: '0.5rem',
}));

export const StyledTypographyDisabled = styled(Typography)(({ theme }) => ({
  color: theme.palette.text.secondary,
  fontSize: '1.3rem',
  margin: '0px',
  marginBlock: '0.5rem',
}));

export const StyledChip = styled(Chip)(() => ({
  '&.MuiChip-root': {
    fontSize: '1rem',
    textTransform: 'capitalize',
  },
}));

export const StyledFormButton = styled(Button)(() => ({
  '&.MuiButton-root': {
    borderRadius: '10px',
    padding: '10px',
  },
}));

export const StyledButton = styled(Button)(() => ({
  '&.MuiButton-root': {
    borderRadius: '10px',
    padding: '1rem 2rem',
    alignSelf: 'center',
  },
}));

export const IconWrapper = styled('div')(({ theme }) => ({
  paddingInline: '2rem',
  paddingTop: '1rem',
  alignSelf: 'flex-end',
  color: theme.palette.info.main,
}));

export const Logo = styled('img')(() => ({
  height: 'clamp(100px, 20vw, 150px)',
}));

export const LogoText = styled('img')(() => ({
  height: 'clamp(30px, 5vw, 100px)',
}));

export const StyledDivider = styled(Divider)(({ theme }) => ({
  backgroundColor: theme.palette.action.disabled,
  [theme.breakpoints.down('md')]: {
    display: 'none',
  },
}));

export const CustomDownIcon = styled(LeftArrowIcon)(({ theme }) => ({
  rotate: '-90deg',
  marginInline: '1rem',
  '&.MuiSelect-icon': {
    fill: theme.palette.primary.main,
  },
}));
