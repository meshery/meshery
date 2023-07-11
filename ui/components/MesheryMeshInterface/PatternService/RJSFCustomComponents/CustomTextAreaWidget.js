import React from "react";
import BaseInput from "./CustomBaseInput";

const CustomTextAreaWidget = (props) => {
  const { options } = props;
  let rows = options.rows || 3;

  return (
    <BaseInput {...props} options={{ ...options, rows }} multiline={true} />
  )
}
export default CustomTextAreaWidget;