import React from "react";
import BaseInput from "./CustomBaseInput";
import { localToUTC,utcToLocal } from "@rjsf/utils"
const CustomDateTimeWidget = (props) => {
  let inputType="datetime-local"
  const value=utcToLocal(props.value);
  const onChange = (value) => {
    props.onChange(localToUTC(value));
  }
  return (
    <BaseInput {...props} options={{ ...props.options, inputType ,focused : true }} onChange={onChange} value={value} />
  )
}
export default CustomDateTimeWidget;