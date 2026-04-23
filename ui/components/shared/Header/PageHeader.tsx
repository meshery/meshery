import React from 'react';
import { Box, Typography, styled } from '@sistent/sistent';

const StyledHeaderRoot = styled(Box)(({ theme }) => ({
  display: 'flex',
  justifyContent: 'space-between',
  alignItems: 'center',
  marginBottom: theme.spacing(3),
  padding: theme.spacing(2, 0),
}));

export interface PageHeaderProps {
  title: string;
  actions?: React.ReactNode;
}

/**
 * Shared PageHeader component for consistent page titles and primary actions.
 */
const PageHeader: React.FC<PageHeaderProps> = ({ title, actions }) => {
  return (
    <StyledHeaderRoot>
      <Typography variant="h4" component="h1" fontWeight="bold">
        {title}
      </Typography>
      {actions && (
        <Box display="flex" gap={2} alignItems="center">
          {actions}
        </Box>
      )}
    </StyledHeaderRoot>
  );
};

export default PageHeader;
