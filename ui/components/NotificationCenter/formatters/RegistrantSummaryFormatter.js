import React from 'react';
import { TitleLink } from './common';
import { Box, Typography } from '@sistent/sistent';

export const RegistrantSummaryFormatter = ({ event }) => {
  const docLink = 'https://docs.meshery.io/concepts/logical#logical-concepts';

  return (
    <Box>
      <Typography>{event.description}</Typography>
      <TitleLink href={docLink} target="_blank" rel="noopener noreferrer">
        Understanding Models, Components, and Relationships
      </TitleLink>
    </Box>
  );
};
