import { makeStyles, Typography } from '@material-ui/core';
import clsx from "clsx";
import React from 'react';
import AnimatedMeshSync from './AnimatedMeshSync';
import AnimatedMeshPattern from "./AnimatedMeshPattern";
import AnimatedFilter from "./AnimatedFilter";

const useStyles = makeStyles(() => ({
  loadingWrapper : {
    textAlign : "center",
    marginTop : "calc(50vh - 141px)",
    transform : "translateY(-50%)",
  }
}));

const animatedIconList = {
  AnimatedMeshPattern : <AnimatedMeshPattern style={{ height : "100px",margin : "4px 0px 8px" }} />,
  AnimatedMeshSync : <>
    <AnimatedMeshSync style={{ height : "75px" }} />
    <img
      src="/static/img/meshery-logo/meshery-black.svg"
      alt="mehsery-logo"
      width="125px"
      style={{ margin : "4px 0px 8px" }}
    />
  </>,
  AnimatedFilter : <AnimatedFilter style={{ height : "75px",margin : "4px 0px 8px" }}/>
}

function LoadingScreen({ message, className,animatedIcon, ...other }) {
  const classes = useStyles();

  return (
    <div className={clsx(classes.loadingWrapper, className)} {...other}>
      {animatedIconList[animatedIcon]}
      <Typography variant="caption" component="div">
        {message}
      </Typography>
    </div>
  );
}

export default LoadingScreen;