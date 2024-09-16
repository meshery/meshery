// Notification.js
import React from 'react';
import { Grid, Typography, Button } from '@material-ui/core';

const Notification = ({ type, message, retry, onRetry }) => {
  return (
    <div
      style={{
        background: type === 'success' ? '#00B39F40' : '#ff000010',
        borderRadius: '0.5rem',
        padding: '0.5rem',
        display: 'flex',
        marginBottom: '1rem',
      }}
    >
      <Grid style={{ width: '80%' }}>
        <Typography variant="body2">
          <b>{type === 'success' ? 'Success' : 'Verification Failed'}</b>
        </Typography>
        <Typography variant="body2" sx={{ color: '#00000020' }}>
          {message}
        </Typography>
      </Grid>
      <Grid
        style={{
          width: '20%',
          display: 'flex',
          justifyContent: 'flex-end',
          alignItems: 'center',
        }}
      >
        {retry && (
          <Button
            style={{
              backgroundColor: '#ff0000',
              padding: '0.5rem',
              borderRadius: '0.5rem',
              border: '0',
              color: '#fff',
              height: '2rem',
            }}
            onClick={onRetry}
          >
            <Typography variant="body2">Retry</Typography>
          </Button>
        )}
      </Grid>
    </div>
  );
};

export default Notification;
