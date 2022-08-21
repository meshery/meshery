import React from "react";
import { NoSsr } from "@mui/material";

let ReactCountdownClock;
if (typeof window !== "undefined") {
  ReactCountdownClock = require("react-countdown-clock");
}

const LoadTestTimerDialog = (props) => {
  const { countDownComplete, t, open } = props;
  if (!open) {
    return "";
  }
  let tNum = 0;
  let dur;
  try {
    tNum = parseInt(t.substring(0, t.length - 1));
  } catch (ex) {
    console.error("Unexpected Error");
  }
  switch (t.substring(t.length - 1, t.length).toLowerCase()) {
    case "h":
      dur = tNum * 60 * 60;
      break;
    case "m":
      dur = tNum * 60;
      break;
    default:
      dur = tNum;
  }
  return (
    <NoSsr>
      <div
        style={{
          marginLeft: "auto",
          marginRight: "auto",
          width: "400px",
          position: "relative",
          zIndex: "0",
          // height: '400',
        }}
      >
        <ReactCountdownClock seconds={dur} color="#667C89" alpha={0.9} size={400} onComplete={countDownComplete} />
      </div>
    </NoSsr>
  );
};

export default LoadTestTimerDialog;
