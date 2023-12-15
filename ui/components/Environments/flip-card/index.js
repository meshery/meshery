import React, { useState } from 'react';
import { FlipCardWrapper, FlipCardBack, FlipCardFront, FlipCardInner } from './style';

/**
 * Wrapper component for flip cards.
 *
 * @param {Object} props - The component props.
 * @param {Object} props.frontComponents - The components of the card front.
 * @param {string} props.backComponents - The components of the card back.
 *
 */

const FlipCard = ({ frontComponents, backComponents, disableFlip }) => {
  const [isFlipped, setIsFlipped] = useState(false);

  const handleFlip = () => {
    if (!disableFlip) {
      setIsFlipped(!isFlipped);
    }
  };

  return (
    <>
      <FlipCardWrapper>
        <FlipCardInner
          onClick={handleFlip}
          style={{ transform: isFlipped ? 'rotateY(180deg)' : 'rotateY(0deg)' }}
        >
          <FlipCardFront>{frontComponents}</FlipCardFront>
          <FlipCardBack>{backComponents}</FlipCardBack>
        </FlipCardInner>
      </FlipCardWrapper>
    </>
  );
};

export default FlipCard;
