import React from 'react';
import { Box, Button, Grid2, Typography } from '@sistent/sistent';
import { alpha, useTheme } from '@/theme';

interface NotificationProps {
  type: 'success' | 'error' | string;
  message: string;
  retry?: boolean;
  onRetry?: () => void;
}

const Notification: React.FC<NotificationProps> = ({ type, message, retry, onRetry }) => {
  const theme = useTheme();
  const isSuccess = type === 'success';

  return (
    <Box
      sx={{
        backgroundColor: isSuccess
          ? alpha(theme.palette.success.main, 0.15)
          : alpha(theme.palette.error.main, 0.1),
        borderRadius: '0.5rem',
        padding: '0.5rem',
        display: 'flex',
        marginBottom: '1rem',
      }}
    >
      <Grid2 sx={{ width: '80%' }}>
        <Typography variant="body2">
          <b>{isSuccess ? 'Success' : 'Verification Failed'}</b>
        </Typography>
        <Typography variant="body2" color="text.secondary">
          {message}
        </Typography>
      </Grid2>
      <Grid2
        sx={{
          width: '20%',
          display: 'flex',
          justifyContent: 'flex-end',
          alignItems: 'center',
        }}
      >
        {retry && (
          <Button
            variant="contained"
            color="error"
            size="small"
            onClick={onRetry}
            sx={{ height: '2rem' }}
          >
            <Typography variant="body2">Retry</Typography>
          </Button>
        )}
      </Grid2>
    </Box>
  );
};

export default Notification;
