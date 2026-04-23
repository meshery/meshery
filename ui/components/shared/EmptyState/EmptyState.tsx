import React from 'react';
import { Box, Typography, Button, styled } from '@sistent/sistent';
import CurvedArrowIcon from '../../../assets/icons/CurvedArrowIcon';

const StyledEmptyStateRoot = styled(Box)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  justifyContent: 'center',
  padding: theme.spacing(4, 2),
  textAlign: 'center',
  minHeight: '400px',
}));

export interface EmptyStateProps {
  message: string;
  icon?: React.ReactNode;
  pointerLabel?: string;
  action?: {
    label: string;
    onClick: () => void;
  };
}

/**
 * Standardized EmptyState component for Meshery UI.
 * Supports optional icon, message, and pointer label with arrow icon.
 */
const EmptyState: React.FC<EmptyStateProps> = ({ message, icon, pointerLabel, action }) => {
  return (
    <StyledEmptyStateRoot>
      {pointerLabel && (
        <Box
          sx={{
            display: 'flex',
            width: '100%',
            padding: '0 40px',
            mb: 5,
            justifyContent: 'flex-start',
          }}
        >
          <CurvedArrowIcon />
          <Typography
            sx={{
              fontSize: 24,
              color: '#808080',
              px: 2,
              display: 'flex',
              alignItems: 'flex-end',
              mb: -3,
              textAlign: 'left',
              lineHeight: 1.5,
            }}
          >
            {pointerLabel}
          </Typography>
        </Box>
      )}
      <Box sx={{ mt: pointerLabel ? 10 : 0 }}>
        {icon && (
          <Box mb={2} sx={{ opacity: 0.8 }}>
            {icon}
          </Box>
        )}
        <Typography variant="h5" sx={{ color: '#808080', fontWeight: 400 }}>
          {message}
        </Typography>
        {action && (
          <Button
            variant="contained"
            color="primary"
            onClick={action.onClick}
            sx={{ mt: 3, borderRadius: '8px' }}
          >
            {action.label}
          </Button>
        )}
      </Box>
    </StyledEmptyStateRoot>
  );
};

export default EmptyState;
