import { Accordion, FormControl, IconButton, styled } from '@layer5/sistent';
import { HelpOutlineOutlined } from '@mui/icons-material';
import { Radio } from '@mui/material';
import { Calendar } from 'react-big-calendar';

export const CardButton = styled('div')(() => ({
  display: 'flex',
  justifyContent: 'flex-end',
  alignItems: 'center',
}));

export const BottomPart = styled('div')(() => ({
  display: 'flex',
  justifyContent: 'space-between',
  alignItems: 'center',
}));

export const ResultContainer = styled('div')(() => ({
  margin: '0 0 1rem',
  '& div': {
    display: 'flex',
    alignItems: 'center',
  },
}));

export const PaginationWrapper = styled('div')(() => ({
  display: 'flex',
  justifyContent: 'center',
  alignItems: 'center',
  marginTop: '2rem',
}));

export const ViewSwitchBUtton = styled('div')(() => ({
  justifySelf: 'flex-end',
  marginLeft: 'auto',
  paddingLeft: '1rem',
  display: 'flex',
}));

export const ProfileContainer = styled('div')(() => ({
  padding: '2rem',
  display: 'flex',
  justifyContent: 'center',
  alignItems: 'center',
  flexDirection: 'column',
}));

export const ButtonTextWrapper = styled('span')(() => ({
  display: 'block',
  '@media (max-width: 1450px)': {
    display: 'none',
  },
}));

export const ResultContainerWrap = styled('div')(() => ({
  display: 'flex',
  flexDirection: 'row',
  justifyContent: 'space-between',
  ['@media (max-width: 830px)']: {
    flexDirection: 'column',
  },
}));

export const CalendarComponent = styled(Calendar)(({ theme }) => ({
  '& .rbc-off-range-bg': {
    backgroundColor: theme.palette.mode === 'dark' ? '#a9a9a9' : '#e6e6e6',
    color: theme.palette.mode === 'dark' ? 'white' : 'black',
  },
  '& .rbc-off-range ': {
    color: theme.palette.text.secondary,
  },
  '& .rbc-btn-group': {
    color: theme.palette.text.secondary,
  },
  '& .rbc-toolbar button': {
    color: theme.palette.text.secondary,
  },
  '& .rbc-today': {
    backgroundColor: theme.palette.mode === 'dark' ? '#505050' : '#eaf6ff',
  },
  '& .rbc-day-slot .rbc-time-slot': {
    borderTop: `1px solid ${theme.palette.mode === 'dark' ? '#555555' : '#eaf6ff'}`,
  },
  '& .rbc-toolbar button.rbc-active': {
    backgroundColor: theme.palette.mode === 'dark' ? '#505050' : '#eaf6ff',
    color: theme.palette.mode === 'dark' ? 'white' : 'black',
  },
  '& .rbc-toolbar button:hover': {
    backgroundColor:
      theme.palette.mode === 'dark' ? 'rgba(0, 0, 0, 0.54)' : 'rgba(255, 255, 255, 0.54)',
  },
  '& .rbc-toolbar button.rbc-active:hover': {
    backgroundColor: theme.palette.mode === 'dark' ? '#909090' : '#eaf6ff',
  },
  '& .rbc-toolbar button:focus ': {
    backgroundColor: theme.palette.mode === 'dark' ? '#505050' : '#eaf6ff',
    color: theme.palette.mode === 'dark' ? 'white' : 'black',
  },
}));

export const IconButtonComp = styled(IconButton)(({ theme }) => ({
  color: theme.palette.icon.default,
}));

export const CenterTimer = styled('div')(() => ({
  width: '100%',
  height: '100%',
  position: 'absolute',
  top: '0',
  left: '0',
  backgroundColor: 'rgba(0, 0, 0, 0.6)',
  zIndex: 1201,
  display: 'flex',
  justifyContent: 'center',
  alignItems: 'center',
}));

export const HelpIcon = styled(HelpOutlineOutlined)(({ theme }) => ({
  width: '15px',
  height: '18px',
  marginBottom: theme.spacing(1),
  marginLeft: theme.spacing(0.3),
}));

export const RadioButton = styled(Radio)(({ theme }) => ({
  radio: {
    '&.Mui-checked': {
      color: theme.palette.text.brand,
    },
  },
}));

export const ExpansionPanelComponent = styled(Accordion)(() => ({
  boxShadow: 'none',
  border: '1px solid rgb(196,196,196)',
  background: 'none',
}));

export const FormContainer = styled(FormControl)(({ theme }) => ({
  margin: theme.spacing(1),
}));
