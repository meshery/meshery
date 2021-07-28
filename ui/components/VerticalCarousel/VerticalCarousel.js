/* eslint-disable no-unused-vars */
import React, { createRef, useCallback, useEffect } from 'react';
import Slider from "react-slick";
import { Grid } from '@material-ui/core';
// const VerticalCarousel = ({
//   item,
//   setActiveIndex
// }) => {
// 
//   const handleClick = (direction) => {
//     if(direction === "prev") setActiveIndex(prev => prev -1)
//     if(direction === "next") setActiveIndex(prev => prev + 1)
//     return
//   }
// 
//   return (
//     <div style={{position: "relative"}}>
//       <IconButton
//         onClick={() => handleClick('prev')}
//         style={{position: "absolute", top:"-3rem" , zIndex: 3, marginLeft: "auto", marginRight: "auto", left: 0, right: 0, textAlign:"center"}}
//       >
//         <ArrowUpwardIcon />
//       </IconButton>
//       <IconButton
//         onClick={() => handleClick('next')}
//         style={{position: "absolute", bottom:"2rem" , zIndex: 3, marginLeft: "auto", marginRight: "auto", left: 0, right: 0, textAlign:"center"}}
//       >
//         <ArrowDownwardIcon />
//       </IconButton>
//       {item} 
//     </div>
//   );
// }
// 
// VerticalCarousel.propTypes = {
//   data: PropTypes.array.isRequired,
//   leadingText: PropTypes.string.isRequired,
// };

    

const VerticalCarousel = ({handleAfterSlideChange, slides, sliderRef}) => {

  const scroll = useCallback(
    y => {
      if (y > 0) {
        return sliderRef?.current?.slickNext(); 
      } else {
        return sliderRef?.current?.slickPrev();
      }
    },
    [sliderRef]
  );

  function disableScrolling(){
    // var x=window.scrollX;
    // var y=window.scrollY;
    // window.onscroll=function(){window.scrollTo(x, y);};
    document.body.style.overflow = 'hidden';
  }

  function enableScrolling(){
    // window.onscroll=function(){};
    document.body.style.overflow = 'scroll';
  }

  useEffect(() => {
    document.getElementById("carousel-div").addEventListener("wheel", e => {
      scroll(e.deltaY);
    });
    document.getElementById("carousel-div").addEventListener("mouseenter", e => {
      disableScrolling()
    });
    document.getElementById("carousel-div").addEventListener("mouseleave", e => {
      enableScrolling()
    });
    
  }, [scroll]);

  const settings = {
    dots: false,
    infinite: false,
    slidesToShow: 1,
    arrows: true,
    slidesToScroll: 1,
    vertical: true,
    verticalSwiping: true,
    swipeToSlide: true,
    afterChange: handleAfterSlideChange
  };

  return (
    <>
      <Grid item id="carousel-div" xs={8} >
        <Slider {...settings} ref={sliderRef} >{slides} </Slider>
      </Grid>
    </>
  )


} 


export default VerticalCarousel;
