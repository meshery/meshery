import React, {useState} from 'react'

function GetChild(children, key) {
    // if (children.length != 2) throw Error("FlipCard requires exactly two child components");
  
    return children[key];
  }

function FlipCard(children, duration = 500,) {

    const [flipped, setFlipped] = useState(false);
    const [activeBack, setActiveBack] = useState(false);

    const Front = GetChild(children, 0);
    const Back = GetChild(children, 1);
  
  return (
        <div
      onClick={() => {
        setFlipped((flipped) => !flipped);
        onClick();
      }}
    >
              <div
        style={{ transform : flipped
          ? "scale(-1,1)"
          : undefined,
        transition : `transform ${duration}ms`,
        transformOrigin : "50% 50% 10%" }}
      >
                 {!activeBack
          ? (
            <div >
              {React.isValidElement(Front)
                ? Front
                : null}
            </div>
          )
          : (
            <div>{React.isValidElement(Back)
              ? Back
              : null}</div>
          )}
      </div>

        </div>
  )
}

export default FlipCard