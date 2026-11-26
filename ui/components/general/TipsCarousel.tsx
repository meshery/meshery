import React from 'react';
import PropTypes from 'prop-types';
import { Box, Typography, useTheme } from '@sistent/sistent';
import TipsIcon from '../../assets/icons/Tipsicon'; // add this icon
import { useInterval } from '@/utils/hooks';

const TipsCarousel = ({ tips }) => {
  const theme = useTheme();
  const [activeStep, setActiveStep] = React.useState(0);

  const handleClick = (index) => {
    setActiveStep(index);
  };

  useInterval(() => {
    setActiveStep((prevActiveStep) =>
      prevActiveStep === tips.length - 1 ? 0 : prevActiveStep + 1,
    );
  }, 3000);

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
      <Typography style={{ marginTop: '1rem' }}>{tips[String(activeStep)]}</Typography>
      <div style={{ textAlign: 'center' }}>
        {tips.map((tip, index) => (
          <button
            key={index}
            type="button"
            onClick={() => handleClick(index)}
            aria-label={`Show tip ${index + 1} of ${tips.length}`}
            aria-current={index === activeStep ? 'true' : undefined}
            style={{
              cursor: 'pointer',
              background: 'none',
              border: 'none',
              padding: 0,
              fontSize: '3rem',
              lineHeight: 1,
              marginRight: '0.25rem',
              color:
                index === activeStep
                  ? theme.palette.background.brand?.default
                  : theme.palette.background.constant.white,
            }}
          >
            &bull;
          </button>
        ))}
      </div>
    </Box>
  );
};

export default TipsCarousel;

TipsCarousel.propTypes = {
  tips: PropTypes.array.isRequired,
};
