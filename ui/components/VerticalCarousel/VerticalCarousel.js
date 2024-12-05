/* eslint-disable no-unused-vars */
import React, { createRef, useCallback, useEffect } from 'react';
import Slider from 'react-slick';
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

const VerticalCarousel = ({ handleAfterSlideChange, slides, sliderRef }) => {
  // const scroll = useCallback(
  //   y => {
  //     if (y > 0) {
  //       return sliderRef?.current?.slickNext();
  //     } else {
  //       return sliderRef?.current?.slickPrev();
  //     }
  //   },
  //   [sliderRef]
  // );

  // useEffect(() => {
  //   document.getElementById("carousel-div").addEventListener("wheel", e => {
  //     e.preventDefault()
  //     scroll(e.deltaY);
  //   });
  // }, [scroll]);

  const settings = {
    // dots : true,
    infinite: false,
    slidesToShow: 1,
    arrows: true,
    // slidesToScroll : 1,
    // adaptiveHeight : true,
    // centerMode : true,
    // centerPadding : "0px",
    vertical: true,
    // verticalSwiping : true,
    // swipeToSlide : true,
    afterChange: handleAfterSlideChange,
  };

  return (
    <>
      <Grid item id="carousel-div" xs={8}>
        <Slider {...settings} ref={sliderRef}>
          {slides}{' '}
        </Slider>
      </Grid>
    </>
  );
};

export default VerticalCarousel;
