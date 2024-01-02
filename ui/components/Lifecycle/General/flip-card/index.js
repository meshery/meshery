import { makeStyles } from '@material-ui/core';
import React, { useState } from 'react';

/**
 * Wrapper component for flip cards.
 *
 * @param {Object} props - The component props.
 * @param {Object} props.frontComponents - The components of the card front.
 * @param {string} props.backComponents - The components of the card back.
 *
 */

const useStyles = makeStyles(() => ({
  flipCardWrapper: {
    background: 'transparent',
    perspective: '1000px',
  },

  flipCardInner: {
    display: 'flex',
    flexDirection: 'column',
    width: '100%',
    height: '100%',
    textAlign: 'center',
    transition: 'transform 0.6s',
    transformStyle: 'preserve-3d',
    boxShadow: '0 4px 8px 0 rgba(0,0,0,0.2)',
  },

  flipCardBack: {
    flex: 1,
    display: 'flex',
    width: '100%',
    height: 'fit-content',
    webkitBackfaceVisibility: 'hidden',
    backfaceVisibility: 'hidden',
    transform: 'rotateY(180deg)',
  },

  flipCardFront: {
    flex: 1,
    display: 'flex',
    position: 'absolute',
    top: 0,
    width: '100%',
    height: '100%',
    webkitBackfaceVisibility: 'hidden',
    backfaceVisibility: 'hidden',
  },
}));

const FlipCard = ({ frontComponents, backComponents, disableFlip }) => {
  const classes = useStyles();

  const [isFlipped, setIsFlipped] = useState(false);

  const handleFlip = () => {
    if (!disableFlip) {
      setIsFlipped(!isFlipped);
    }
  };

  return (
    <>
      <div className={classes.flipCardWrapper}>
        <div
          className={classes.flipCardInner}
          onClick={handleFlip}
          style={{ transform: isFlipped ? 'rotateY(180deg)' : 'rotateY(0deg)' }}
        >
          <div className={classes.flipCardFront}>{frontComponents}</div>
          <div className={classes.flipCardBack}>{backComponents}</div>
        </div>
      </div>
    </>
  );
};

export default FlipCard;
