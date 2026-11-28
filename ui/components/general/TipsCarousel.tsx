import React from 'react';
import PropTypes from 'prop-types';
import { Box, Typography, useTheme } from '@sistent/sistent';
import { styled } from '@/theme';
import TipsIcon from '../../assets/icons/Tipsicon';
import { useInterval } from '@/utils/hooks';

const CarouselContainer = styled(Box)(({ theme }) => ({
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
  [theme.breakpoints.down(780)]: {
    display: 'none',
  },
}));

const IndicatorButton = styled('button')(({ theme }) => ({
  appearance: 'none',
  background: 'none',
  border: 0,
  padding: 0,
  font: 'inherit',
  lineHeight: 1,
  cursor: 'pointer',
  fontSize: '3rem',
  marginRight: '0.25rem',
  '&:focus-visible': {
    outline: `2px solid ${theme.palette.background.brand?.default ?? '#00B39F'}`,
    outlineOffset: '2px',
  },
}));

type TipsCarouselProps = {
  tips: string[];
};

const TipsCarousel = ({ tips }: TipsCarouselProps) => {
  const theme = useTheme();
  const [activeStep, setActiveStep] = React.useState(0);

  const handleClick = (index: number) => {
    setActiveStep(index);
  };

  useInterval(() => {
    setActiveStep((prevActiveStep) =>
      prevActiveStep === tips.length - 1 ? 0 : prevActiveStep + 1,
    );
  }, 3000);

  return (
    <CarouselContainer>
      <TipsIcon />

      <Typography style={{ marginTop: '1rem' }}>{tips[activeStep]}</Typography>

      <div style={{ textAlign: 'center' }}>
        {tips.map((_, index) => (
          <IndicatorButton
            key={index}
            type="button"
            aria-label={`Show tip ${index + 1} of ${tips.length}`}
            aria-current={index === activeStep ? 'step' : undefined}
            onClick={() => handleClick(index)}
            style={{
              color:
                index === activeStep
                  ? (theme.palette.background.brand?.default ?? '#00B39F')
                  : (theme.palette.background.constant?.white ?? '#fff'),
            }}
          >
            &bull;
          </IndicatorButton>
        ))}
      </div>
    </CarouselContainer>
  );
};

export default TipsCarousel;

TipsCarousel.propTypes = {
  tips: PropTypes.array.isRequired,
};
