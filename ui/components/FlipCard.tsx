import React, { useState, useRef, useEffect } from 'react';
import { FlipCardWrapper, InnerCard } from './MesheryPatterns/style';

function GetChild(children, key) {
  if (children.length != 2) throw Error('FlipCard requires exactly two child components');

  return children[key];
}

function FlipCard({ duration = 500, onClick, onShow, children }) {
  const [flipped, setFlipped] = useState(false);
  const [activeBack, setActiveBack] = useState(false);

  const timeout = useRef(null);

  const Front = GetChild(children, 0);
  const Back = GetChild(children, 1);

  useEffect(() => {
    // This function makes sure that the inner content of the card disappears roughly
    // after 30 deg rotation has already occured. It will ensure that the user doesn't gets
    // a "blank" card while the card is rotating
    //
    // This guarantee can be offered because of two main reasons:
    // 1. In sufficiently modern browsers JS and CSS are handled in different threads
    // hence ones execution doesn't blocks another.
    // 2. setTimeout will put its callback at the end of current context's end hence ensuring
    // this callback doesn't gets blocked by another JS process.
    if (timeout.current) clearTimeout(timeout.current);

    timeout.current = setTimeout(() => {
      setActiveBack(flipped);
    }, duration / 6);
  }, [flipped]);

  return (
    <>
      <FlipCardWrapper
        onClick={() => {
          setFlipped((flipped) => !flipped);
          onClick && onClick();
          onShow && onShow();
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
            <div
              style={{
                backfaceVisibility: 'hidden',
              }}
            >
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
    </>
  );
}

export default FlipCard;
