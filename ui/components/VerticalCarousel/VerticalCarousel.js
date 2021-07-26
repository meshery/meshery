/* eslint-disable no-unused-vars */
import React from 'react';
import ArrowDownwardIcon from '@material-ui/icons/ArrowDownward';
import ArrowUpwardIcon from '@material-ui/icons/ArrowUpward';
import PropTypes from 'prop-types';
import {Button, Grid, IconButton} from "@material-ui/core"

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
    <Grid direction="column" container style={{position: "relative"}} alignItems="center" item>
      <IconButton
        onClick={() => handleClick('prev')}
        style={{position: "absolute", top:"-3rem" , zIndex: 3, marginLeft: "auto", marginRight: "auto", left: 0, right: 0, textAlign:"center"}}
      >
        <ArrowUpwardIcon />
      </IconButton>
      <IconButton
        onClick={() => handleClick('next')}
        style={{position: "absolute", bottom:"2rem" , zIndex: 3, marginLeft: "auto", marginRight: "auto", left: 0, right: 0, textAlign:"center"}}
      >
        <ArrowDownwardIcon />
      </IconButton>
      <Grid item> {item} </Grid>  
    </Grid>
  );
}

VerticalCarousel.propTypes = {
  data: PropTypes.array.isRequired,
  leadingText: PropTypes.string.isRequired,
};

export default VerticalCarousel;
