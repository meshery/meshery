import { styled } from '@mui/system';
import Tooltip, { tooltipClasses } from '@mui/material/Tooltip';
import React from 'react';
import { Button } from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';

export const CustomTooltip = styled(({ className, ...props }) => (
  <Tooltip {...props} classes={{ popper: className }} />
))({
  [`& .${tooltipClasses.tooltip}`]: {
    backgroundColor: 'rgba(122,132,142,1)',
    color: '#F5F5F5',
    padding: '1rem',
    fontSize: '0.925rem',
    '& .tooltip-dark': {
      fontWeight: 'bold',
      fontSize: '1rem',
    },
  },
});

export const ModalWrapper = styled('div')(() => ({
  zIndex: '100',
  borderRadius: '5px',
}));

export const CreateTokenModalHeader = styled('div')(() => ({
  display: 'flex',
  justifyContent: 'space-between',
  padding: '1rem',
  // borderBottom: "1px solid #ccc",
  boxShadow: 'inset 0px -1px 3px 0px rgba(0,0,0,0.2)',
  // backgroundColor: "#396679",
  background: 'linear-gradient(90deg, rgba(57,102,121,1) 0%, rgba(84,129,148,1) 100%)',
  filter:
    "progid:DXImageTransform.Microsoft.gradient(startColorstr='#396679',endColorstr='#548194',GradientType=1)",
  color: '#fff',
}));

export const CreateTokenModalBody = styled('form')(() => ({
  padding: '2rem 1rem',
  backgroundColor: 'white',
}));

export const CancelButton = styled(Button)(() => ({
  marginRight: '1rem',
  color: '#000',
  backgroundColor: '#fff',
  '&:hover': {
    backgroundColor: '#fff',
  },
}));

export const GenerateButton = styled(Button)(() => ({
  backgroundColor: '#00B39F',
  '&:hover': {
    backgroundColor: '#00D3A9',
  },
  // "&[disabled]": {
  //   backgroundColor: "rgba(255, 255, 255, 0.4)"
  // }
}));

export const ButtonContainer = styled('div')(() => ({
  padding: '1rem 1.5rem',
  display: 'flex',
  justifyContent: 'flex-end',
  // backgroundColor: "#1E2117",
  // backgroundColor: "#B1B6B8",
  backgroundColor: 'rgba(57, 102, 121, 1)',
  boxShadow: 'inset 0px 3px 5px 0px rgba(0,0,0,0.25)',
}));

export const ModalCloseIcon = styled(CloseIcon)((props) => ({
  width: '2rem',
  height: '2rem',
  cursor: 'pointer',
  transform: props.rotate ? `rotate(-${props.rotate}deg)` : '',
  transition: 'all .3s ease-in',
  '&:hover': {
    transform: props.rotate ? `rotate(${props.rotate}deg)` : '',
  },
}));
