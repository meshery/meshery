import React from "react";
import BaseInput from "./CustomBaseInput";

const CustomDateTimeWidget = (props) => {
  let inputType="datetime-local"
  return (
    <BaseInput {...props} options={{ ...props.options, inputType ,focused : true }} />
  )
}
export default CustomDateTimeWidget;