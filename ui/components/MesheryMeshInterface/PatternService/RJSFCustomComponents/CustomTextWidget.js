import React from 'react';
import BaseInput from './CustomBaseInput';

const CustomTextWidget = (props) => {
  let inputType = 'string';
  if (props.schema?.type === 'number' || props.schema?.type === 'integer') {
    inputType = 'number';
  }
  return <BaseInput {...props} options={{ ...props.options, inputType }} />;
};

export default CustomTextWidget;
