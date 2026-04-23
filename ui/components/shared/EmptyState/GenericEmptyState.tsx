import React from 'react';
import { Box, Typography, Button, styled } from '@sistent/sistent';

const StyledEmptyStateRoot = styled(Box)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  justifyContent: 'center',
  padding: theme.spacing(8, 2),
  textAlign: 'center',
  backgroundColor: 'transparent',
  borderRadius: '12px',
  border: `1px dashed ${theme.palette.divider}`,
  minHeight: '300px',
}));

export interface GenericEmptyStateProps {
  title: string;
  subtitle?: string;
  icon?: React.ReactNode;
  action?: {
    label: string;
    onClick: () => void;
  };
}

/**
 * Shared EmptyState component to be used across the application.
 */
const GenericEmptyState: React.FC<GenericEmptyStateProps> = ({ title, subtitle, icon, action }) => {
  return (
    <StyledEmptyStateRoot>
      {icon && (
        <Box
          mb={2}
          sx={{ fontSize: '4rem', color: (theme) => theme.palette.text.secondary, opacity: 0.6 }}
        >
          {icon}
        </Box>
      )}
      <Typography variant="h5" component="h2" gutterBottom>
        {title}
      </Typography>
      {subtitle && (
        <Typography variant="body1" color="textSecondary" sx={{ mb: 3, maxWidth: '500px' }}>
          {subtitle}
        </Typography>
      )}
      {action && (
        <Button variant="contained" color="primary" onClick={action.onClick}>
          {action.label}
        </Button>
      )}
    </StyledEmptyStateRoot>
  );
};

export default GenericEmptyState;
