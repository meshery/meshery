import React, { useCallback, useEffect, useState } from 'react';
import { Box, Typography, useTheme } from '@sistent/sistent';
import TipsIcon from '../../assets/icons/Tipsicon'; // add this icon

interface TipsCarouselProps {
  tips: string[];
}

const TipsCarousel = ({ tips }: TipsCarouselProps) => {
  const theme = useTheme();
  const [activeStep, setActiveStep] = useState(0);

  const handleClick = useCallback((index: number) => {
    setActiveStep(index);
  }, []);

  // auto increment after 3 seconds
  useEffect(() => {
    if (tips.length === 0) {
      return undefined;
    }

    const interval = setInterval(() => {
      setActiveStep((prevActiveStep) =>
        prevActiveStep === tips.length - 1 ? 0 : prevActiveStep + 1,
      );
    }, 3000);

    return () => clearInterval(interval);
  }, [tips.length]);

  return (
    <Box
      style={{
        width: '40%',
        minWidth: '320px',
        maxWidth: '320px',
        borderRadius: '0.438rem',
        background: 'linear-gradient(359deg, #00B39F -6.26%, #000000 85.81%)',
        padding: '2rem',
        color: '#fff',
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        '@media (max-width:780px)': {
          display: 'none',
        },
      }}
    >
      <TipsIcon />
      <Typography style={{ marginTop: '1rem' }}>{tips[activeStep] ?? ''}</Typography>
      <div style={{ textAlign: 'center' }}>
        {tips.map((tip, index) => (
          <span
            key={index}
            onClick={() => handleClick(index)}
            style={{
              cursor: 'pointer',
              fontSize: '3rem',
              marginRight: '0.25rem',
              color:
                index === activeStep
                  ? theme.palette.background.brand?.default
                  : theme.palette.background.constant.white,
            }}
          >
            &bull;
          </span>
        ))}
      </div>
    </Box>
  );
};

export default TipsCarousel;
