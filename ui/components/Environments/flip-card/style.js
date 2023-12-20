import { styled } from '@mui/system';

export const FlipCardWrapper = styled('div')(() => ({
  background: 'transparent',
  perspective: '1000px',
}));

export const FlipCardInner = styled('div')(() => ({
  display: 'flex',
  flexDirection: 'column',
  width: '100%',
  height: '100%',
  textAlign: 'center',
  transition: 'transform 0.6s',
  transformStyle: 'preserve-3d',
  boxShadow: '0 4px 8px 0 rgba(0,0,0,0.2)',
}));

export const FlipCardBack = styled('div')(() => ({
  flex: 1,
  display: 'flex',
  width: '100%',
  height: 'fit-content',
  webkitBackfaceVisibility: 'hidden',
  backfaceVisibility: 'hidden',
  transform: 'rotateY(180deg)',
}));

export const FlipCardFront = styled('div')(() => ({
  flex: 1,
  display: 'flex',
  position: 'absolute',
  top: 0,
  width: '100%',
  height: '100%',
  webkitBackfaceVisibility: 'hidden',
  backfaceVisibility: 'hidden',
}));
