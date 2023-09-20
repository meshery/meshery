import React, { useState, useEffect } from 'react';
import { NoSsr } from '@material-ui/core';

let ReactCountdownClock;
if (typeof window !== 'undefined') {
  ReactCountdownClock = require('react-countdown-clock');
}

const LoadTestTimerDialog = ({ countDownComplete, t, open }) => {
  const [dur, setDur] = useState(0);

  useEffect(() => {
    let tNum = 0;
    try {
      tNum = parseInt(t.substring(0, t.length - 1));
    } catch (ex) {
      console.error('Unexpected Error');
    }

    switch (t.substring(t.length - 1, t.length).toLowerCase()) {
      case 'h':
        setDur(tNum * 60 * 60);
        break;
      case 'm':
        setDur(tNum * 60);
        break;
      default:
        setDur(tNum);
    }
  }, [t]);

  if (!open) {
    return null;
  }

  return (
    <NoSsr>
      <div style={{
        marginLeft : 'auto',
        marginRight : 'auto',
        width : '400px',
        position : "relative",
        zIndex : "0"
      }}>
        <ReactCountdownClock
          seconds={dur}
          color="#667C89"
          alpha={0.9}
          size={400}
          onComplete={countDownComplete}
        />
      </div>
    </NoSsr>
  );
}

export default LoadTestTimerDialog;
