import React, { useState, useRef, useEffect, type ReactNode } from 'react';
import { FlipCardWrapper, InnerCard } from './MesheryPatterns/style';

interface FlipCardProps {
  duration?: number;
  onClick?: () => void;
  onShow?: () => void;
  children: [ReactNode, ReactNode];
}

function FlipCard({ duration = 500, onClick, onShow, children }: FlipCardProps) {
  const [flipped, setFlipped] = useState(false);
  const [activeBack, setActiveBack] = useState(false);
  const timeout = useRef<ReturnType<typeof setTimeout> | null>(null);

  const Front = children[0];
  const Back = children[1];

  useEffect(() => {
    // Delay the back-face swap until ~30° of rotation has passed so the user
    // never sees a blank card mid-flip. JS and CSS run on separate threads in
    // modern browsers, and setTimeout defers past the current call stack.
    if (timeout.current) clearTimeout(timeout.current);
    timeout.current = setTimeout(() => {
      setActiveBack(flipped);
    }, duration / 6);
  }, [flipped, duration]);

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
