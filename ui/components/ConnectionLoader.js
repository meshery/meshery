import React from 'react';
import { CircularProgress, Typography, Box, Paper } from '@mui/material';
import { styled } from '@mui/material/styles';

const LoaderContainer = styled(Paper)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  justifyContent: 'center',
  padding: theme.spacing(4),
  margin: theme.spacing(2),
  minHeight: '200px',
  backgroundColor: theme.palette.background.paper,
  borderRadius: theme.spacing(1),
  boxShadow: theme.shadows[2],
}));

const LoaderText = styled(Typography)(({ theme }) => ({
  marginTop: theme.spacing(2),
  textAlign: 'center',
  color: theme.palette.text.secondary,
}));

const ConnectionLoader = ({ message = "Verifying connection..." }) => {
  return (
    <LoaderContainer>
      <CircularProgress 
        size={60} 
        thickness={4}
        sx={{ color: 'primary.main' }}
      />
      <LoaderText variant="body1">
        {message}
      </LoaderText>
      <LoaderText variant="body2" sx={{ mt: 1, opacity: 0.7 }}>
        Please wait while we verify your connection...
      </LoaderText>
    </LoaderContainer>
  );
};

export default ConnectionLoader; 