import React from 'react';
import PropTypes from 'prop-types';
import { Box, Typography } from '@material-ui/core';
import TipsIcon from '../../assets/icons/Tipsicon'; // add this icon
import { useTheme } from '@material-ui/core';

const TipsCarousel = ({ tips }) => {
  const theme = useTheme();
  const [activeStep, setActiveStep] = React.useState(0);

  const handleClick = (index) => {
    setActiveStep(index);
  };

  // auto increment after 3 seconds
  React.useEffect(() => {
    const interval = setInterval(() => {
      setActiveStep((prevActiveStep) =>
        prevActiveStep === tips.length - 1 ? 0 : prevActiveStep + 1,
      );
    }, 3000);
    return () => clearInterval(interval);
  }, []);

  return (
    <Box
      style={{
        width: '40%',
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
          <span
            key={index}
            onClick={() => handleClick(index)}
            style={{
              cursor: 'pointer',
              fontSize: '3rem',
              marginRight: '0.25rem',
              color: index === activeStep ? theme.palette.keppelGreen : theme.palette.white,
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

TipsCarousel.propTypes = {
  tips: PropTypes.array.isRequired,
};
