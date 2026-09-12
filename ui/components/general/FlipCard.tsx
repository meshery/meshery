import React, { useState, type ReactNode } from 'react';
import { FlipCardWrapper, InnerCard } from './FlipCard.styles';
import { useTimeout } from '@/utils/hooks';

interface FlipCardProps {
  duration?: number;
  onClick?: () => void;
  onShow?: () => void;
  children: ReactNode;
}

function FlipCard({ duration = 500, onClick, onShow, children }: FlipCardProps) {
  const [flipped, setFlipped] = useState(false);
  const [activeBack, setActiveBack] = useState(false);

  const childArray = React.Children.toArray(children);
  if (childArray.length !== 2) {
    throw new Error('FlipCard requires exactly two child components');
  }
  const [Front, Back] = childArray;

  // Delay the back-face swap until ~30° of rotation has passed so the user
  // never sees a blank card mid-flip.
  useTimeout(() => setActiveBack(flipped), duration / 6, [flipped, duration]);

  return (
    <FlipCardWrapper
      onClick={() => {
        setFlipped((f) => !f);
        onClick?.();
        onShow?.();
      }}
    >
      <InnerCard
        style={{
          transform: flipped ? 'scale(-1,1)' : undefined,
          transition: `transform ${duration}ms`,
          transformOrigin: '50% 50% 10%',
        }}
      >
        {!activeBack ? (
          <div style={{ backfaceVisibility: 'hidden' }}>
            {React.isValidElement(Front) ? Front : null}
          </div>
        ) : (
          <div
            style={{
              backfaceVisibility: 'hidden',
              transform: 'scale(-1, 1)',
              maxWidth: '50vw',
              wordBreak: 'break-word',
            }}
          >
            {React.isValidElement(Back) ? Back : null}
          </div>
        )}
      </InnerCard>
    </FlipCardWrapper>
  );
}

export default FlipCard;
