import { Card, FormControl, Paper, styled, Select, FormGroup, Grid } from '@layer5/sistent';

export const StyledSelect = styled(Select)(({ theme }) => ({
  '&.Mui-focused .MuiOutlinedInput-notchedOutline': {
    borderColor: 'transparent !important',
  },
  width: '85%',
  backgroundColor: theme.palette.text.secondary,
  color: theme.palette.text.primary,
  '@media (max-width: 368px)': {
    width: '180px',
  },
  position: 'relative',
  overflow: 'hidden',
  '& svg': {
    backgroundColor: theme.palette.icon.secondary,
    fill: theme.palette.icon.secondary,
    float: 'right',
    top: '0',
    right: '0',
    height: '100%',
  },
}));

export const SelectItem = styled('div')({
  width: '100%',
  '& svg': {
    backgroundColor: 'inherit',
  },
});

export const OrgText = styled('span')({
  marginLeft: '10px',
  overflow: 'hidden',
  whiteSpace: 'nowrap',
  textOverflow: 'ellipsis',
});

export const OrgIconContainer = styled('div')({
  width: '32px',
  height: '100%',
  position: 'relative',
});

export const FormGroupWrapper = styled(FormControl)(() => ({
  padding: 20,
  border: '1.5px solid #969696',
  display: 'flex',
  width: '70%',
}));

export const FormContainerWrapper = styled('div')(() => ({
  display: 'flex',
  flexWrap: 'wrap',
  justifyContent: 'space-evenly',
  padding: 50,
}));

export const StatsWrapper = styled(Paper)({
  maxWidth: '100%',
  height: 'auto',
  borderTopLeftRadius: 0,
  borderTopRightRadius: 0,
  borderBottomLeftRadius: 3,
  borderBottomRightRadius: 3,
});

export const PaperRoot = styled(Paper)({
  flexGrow: 1,
  maxWidth: '100%',
  marginLeft: 0,
  borderTopLeftRadius: 3,
  borderTopRightRadius: 3,
});

export const RootContainer = styled('div')(({ theme }) => ({
  width: '100%',
  paddingLeft: theme.spacing(15),
  paddingRight: theme.spacing(15),
  paddingBottom: theme.spacing(10),
  paddingTop: theme.spacing(5),
}));

export const ProviderCard = styled(Card)(({ theme }) => ({
  border: '1px solid rgba(0,179,159,0.3)',
  margin: '20px 0px',
  backgroundColor: theme.palette.mode === 'dark' ? '#293B43' : '#C9DBE3',
}));

export const StyledPaper = styled(Paper)(({ theme }) => ({
  '& .MuiTab-root': {
    [theme.breakpoints.up('sm')]: {
      fontSize: '1em',
    },
    [theme.breakpoints.between('xs', 'sm')]: {
      fontSize: '0.8em',
    },
  },
  '& .hideScrollbar': {
    overflowX: 'auto',
    '&::-webkit-scrollbar': {
      display: 'none',
    },
    '&::-moz-scrollbar': {
      display: 'none',
    },
  },
}));

export const BoxWrapper = styled('div')({
  width: '100%',
  display: 'flex',
  justifyContent: 'space-between',
  whiteSpace: 'nowrap',
  paddingRight: '10px',
});

export const GridCapabilityHeader = styled(Grid)(({ theme }) => ({
  backgroundColor: theme.palette.mode === 'dark' ? '#293B43' : '#7493A1',
}));

export const GridExtensionHeader = styled(Grid)(({ theme }) => ({
  backgroundColor: theme.palette.mode === 'dark' ? '#293B43' : '#C9DBE3',
}));

export const GridExtensionItem = styled(Grid)(({ theme }) => ({
  backgroundColor: theme.palette.mode === 'dark' ? '#3D4F57' : '#E7EFF3',
}));

export const Divider = styled('hr')({
  border: '1px solid rgba(116,147,161, 0.3)',
  width: '100%',
  margin: '30px 0',
});

export const TabLabel = styled('span')(({ theme }) => ({
  [theme.breakpoints.up('sm')]: {
    fontSize: '1em',
  },
  [theme.breakpoints.between('xs', 'sm')]: {
    fontSize: '0.8em',
  },
}));

export const StyledFormGroup = styled(FormGroup)({
  display: 'flex',
  flexWrap: 'wrap',
  justifyContent: 'space-evenly',
  padding: 50,
});

export const FormLegend = styled('legend')({
  fontSize: '1rem',
});

export const FormLegendSmall = styled('legend')({
  fontSize: 16,
});

export const IconStyled = styled('span')(({ theme }) => ({
  display: 'inline',
  verticalAlign: 'text-top',
  marginLeft: theme.spacing(0.5),
}));

export const IconText = styled('span')({
  display: 'inline',
  verticalAlign: 'middle',
});

export const BackToPlay = styled('div')(({ theme }) => ({
  margin: theme.spacing(2),
}));

export const LinkStyled = styled('div')({
  cursor: 'pointer',
});

export const HideScrollbar = styled('div')({
  overflowX: 'auto',
  '&::-webkit-scrollbar': {
    display: 'none',
  },
  '&::-moz-scrollbar': {
    display: 'none',
  },
});

export const FormContainer = styled('div')({
  padding: '1rem',
  marginBottom: '2rem',
});

export const PreferenceLabel = styled('label')({
  display: 'block',
  marginBottom: '0.5rem',
  fontWeight: 500,
});

export const PreferenceGroup = styled('div')({
  marginBottom: '1.5rem',
});

export const PreferenceValue = styled('div')(({ theme }) => ({
  color: theme.palette.text.secondary,
  marginTop: '0.25rem',
}));

export const PreferenceSection = styled('section')({
  marginBottom: '2rem',
  padding: '1rem',
});

export const ListItemStyled = styled('div')(({ theme }) => ({
  padding: theme.spacing(1),
  '&:hover': {
    backgroundColor: theme.palette.action.hover,
  },
}));

export const ContentWrapper = styled('div')({
  flex: 1,
  padding: '1rem',
});
