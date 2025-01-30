import { styled } from '@layer5/sistent';
import React, { useState } from 'react';

const FlipCardWrapper = styled('div')({
  background: 'transparent',
  perspective: '1000px',
});

const FlipCardInner = styled('div')({
  display: 'flex',
  flexDirection: 'column',
  width: '100%',
  height: '100%',
  textAlign: 'center',
  transition: 'transform 0.6s',
  transformStyle: 'preserve-3d',
  boxShadow: '0 4px 8px 0 rgba(0,0,0,0.2)',
});

const FlipCardBack = styled('div')({
  flex: 1,
  display: 'flex',
  width: '100%',
  height: 'fit-content',
  WebkitBackfaceVisibility: 'hidden', // Updated to camelCase
  backfaceVisibility: 'hidden',
  transform: 'rotateY(180deg)',
});

const FlipCardFront = styled('div')({
  flex: 1,
  display: 'flex',
  position: 'absolute',
  top: 0,
  width: '100%',
  height: '100%',
  WebkitBackfaceVisibility: 'hidden', // Updated to camelCase
  backfaceVisibility: 'hidden',
});

const FlipCard = ({ frontComponents, backComponents, disableFlip }) => {
  const [isFlipped, setIsFlipped] = useState(false);

  const handleFlip = () => {
    if (!disableFlip) {
      setIsFlipped(!isFlipped);
    }
  };

  return (
    <FlipCardWrapper>
      <FlipCardInner
        onClick={handleFlip}
        sx={{ transform: isFlipped ? 'rotateY(180deg)' : 'rotateY(0deg)' }}
      >
        <FlipCardFront>{frontComponents}</FlipCardFront>
        <FlipCardBack>{backComponents}</FlipCardBack>
      </FlipCardInner>
    </FlipCardWrapper>
  );
};

export default FlipCard;
