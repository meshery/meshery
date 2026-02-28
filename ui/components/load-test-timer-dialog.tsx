import React, { useEffect, useState } from 'react';
import { Box, CircularProgress, NoSsr, Typography, IconButton } from '@sistent/sistent';
import CloseIcon from '@mui/icons-material/Close';

const parseDuration = (value) => {
  if (!value) {
    return 0;
  }

  const tNum = parseInt(value.substring(0, value.length - 1), 10);

  if (Number.isNaN(tNum)) {
    console.error(`Failed to parse duration value: ${value}`);
    return 0;
  }

  switch (value.substring(value.length - 1, value.length).toLowerCase()) {
    case 'h':
      return tNum * 60 * 60;
    case 'm':
      return tNum * 60;
    default:
      return tNum;
  }
};

const LoadTestTimerDialog = ({ countDownComplete, t, open, onClose }) => {
  const [dur, setDur] = useState(() => parseDuration(t));
  const [timeLeft, setTimeLeft] = useState(() => parseDuration(t));

  useEffect(() => {
    setDur(parseDuration(t));
  }, [t]);

  useEffect(() => {
    if (!open) {
      setTimeLeft(0);
      return;
    }

    setTimeLeft(dur);
  }, [dur, open]);

  useEffect(() => {
    if (!open || dur <= 0) {
      return undefined;
    }

    const interval = setInterval(() => {
      setTimeLeft((prev) => {
        if (prev <= 1) {
          clearInterval(interval);
          if (countDownComplete) {
            countDownComplete();
          }
          return 0;
        }

        return prev - 1;
      });
    }, 1000);

    return () => clearInterval(interval);
  }, [countDownComplete, dur, open]);

  const formatTime = (seconds) => {
    const totalSeconds = Math.max(0, seconds);
    const hours = Math.floor(totalSeconds / 3600);
    const minutes = Math.floor((totalSeconds % 3600) / 60);
    const secs = Math.floor(totalSeconds % 60);

    if (hours > 0) {
      return `${String(hours).padStart(2, '0')}:${String(minutes).padStart(2, '0')}:${String(
        secs,
      ).padStart(2, '0')}`;
    }

    if (minutes > 0) {
      return `${String(minutes).padStart(2, '0')}:${String(secs).padStart(2, '0')}`;
    }

    return `${secs}`;
  };

  const progress = dur > 0 ? (timeLeft / dur) * 100 : 0;
  const formattedTime = formatTime(timeLeft);

  if (!open) {
    return null;
  }

  return (
    <NoSsr>
      <Box
        sx={{
          marginLeft: 'auto',
          marginRight: 'auto',
          width: 400,
          position: 'relative',
        }}
        aria-label={`Performance test countdown: ${formattedTime} remaining`}
      >
        {onClose && (
          <IconButton
            onClick={onClose}
            sx={{
              position: 'absolute',
              top: 8,
              right: 8,
              color: 'white',
              zIndex: 1,
            }}
            aria-label="Close countdown timer"
          >
            <CloseIcon />
          </IconButton>
        )}
        <Box sx={{ position: 'relative', display: 'inline-flex', width: 400, height: 400 }}>
          <CircularProgress
            variant="determinate"
            value={progress}
            size={400}
            sx={{ color: 'white' }}
          />
          <Box
            sx={{
              top: 0,
              left: 0,
              bottom: 0,
              right: 0,
              position: 'absolute',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
            }}
          >
            <Typography variant="h4" component="div" color="white">
              {formattedTime}
            </Typography>
          </Box>
        </Box>
      </Box>
    </NoSsr>
  );
};

export default LoadTestTimerDialog;
