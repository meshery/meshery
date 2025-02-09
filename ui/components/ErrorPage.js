import React, { useState, useEffect } from 'react';
import { NoSsr } from '@mui/material';
import { Typography, Link, Box, styled } from '@layer5/sistent';
import { UsesSistent } from './SistentWrapper';

const StyledBox = styled(Box)(({ theme }) => ({
  backgroundColor: theme.palette.background?.default,
  padding: '2rem',
  textAlign: 'center',
  borderRadius: 4,
  height: '100%',
}));
const StyledTypographyH1 = styled(Typography)(() => ({
  fontSize: '3rem',
  lineHeight: '2rem',
  marginBottom: '2rem',
}));
const StyledTypographyH5 = styled(Typography)(({ theme }) => ({
  fontWeight: 400,
  fontSize: '1.5rem',
  color: theme.palette.text?.tertiary,
  fontStyle: 'italic',
  marginTop: '2.5rem',
}));
const StyledTypographyBody1 = styled(Typography)(({ theme }) => ({
  marginTop: '5rem',
  color: theme.palette.text?.default,
}));

const StyledLink = styled(Link)(({ theme }) => ({
  textDecoration: 'none',
  color: theme.palette.text?.brand,
}));
/**
 * CustomErrorMessage component is used to display a custom error message when a page is not found.
 * @returns {JSX.Element} JSX.Element
 */
const CustomErrorMessage = ({ message, showImage = true }) => {
  const [customMessage, setCustomMessage] = useState('Oh, no. Please pardon our meshy app.');

  useEffect(() => {
    const customMessages = [
      'Oh, no. Please pardon our meshy app.',
      'Oops. Please excuse the mesh.',
      'Things tend to get a bit meshy around here.',
      'Please pardon our mesh.',
      'How did this mesh happen?',
      "Well, isn't this a mesh?",
      'Yikes. Things are a mesh here.',
    ];
    setCustomMessage(customMessages[Math.floor(Math.random() * customMessages.length)]);
  }, []);

  return (
    <UsesSistent>
      <NoSsr>
        <StyledBox>
          <Box>
            <StyledTypographyH1 variant="h1">{customMessage}</StyledTypographyH1>
            <StyledTypographyH5 variant="h5">
              {message || 'Page does not exist.'}
            </StyledTypographyH5>
          </Box>
          {showImage && (
            <Box
              component="img"
              src="/static/img/service-mesh.svg"
              alt="service meshed"
              sx={{
                display: 'block',
                margin: 'auto',
                mt: 3.125,
                maxWidth: '50%',
                height: '45%',
              }}
            />
          )}
          <StyledTypographyBody1 variant="body1">
            Start a conversation at Layer5 community{' '}
            <StyledLink
              href="https://meshery.io/community#community-forums/c/meshery/5"
              target="_blank"
            >
              discussion forum
            </StyledLink>
            .
          </StyledTypographyBody1>
        </StyledBox>
      </NoSsr>
    </UsesSistent>
  );
};

export default CustomErrorMessage;
