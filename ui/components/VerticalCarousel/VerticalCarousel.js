import React from 'react';
import PropTypes from 'prop-types';
import {Grid} from "@material-ui/core"

const VerticalCarousel = ({
  item,
  setActiveIndex
}) => {

  const handleClick = (direction) => {
    if(direction === "prev") setActiveIndex(prev => prev -1)
    if(direction === "next") setActiveIndex(prev => prev + 1)
    return
  }

  return (
    <Grid direction="column" container style={{position: "relative"}}>
	    <div>
        <button
          type="button"
          onClick={() => handleClick('prev')}
          style={{position: "absolute", top:"-2rem" , zIndex: 3}}
        >
          <div> Prev </div>
        </button>
        <button
          type="button"
          onClick={() => handleClick('next')}
          style={{position: "absolute", bottom:"2rem" , zIndex: 3}}
        >
          <div> Next </div>
        </button>
      </div>
      <Grid item> {item} </Grid>  
    </Grid>
  );
}

VerticalCarousel.propTypes = {
  data: PropTypes.array.isRequired,
  leadingText: PropTypes.string.isRequired,
};

export default VerticalCarousel;
