import React from 'react';
import BaseInput from './CustomBaseInput';

const CustomColorWidget = (props) => {
  const { options } = props;

  return <BaseInput {...props} options={{ ...options, inputType: 'color' }} />;
};
export default CustomColorWidget;
