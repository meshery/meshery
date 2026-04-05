// @ts-nocheck
import BaseInput from './CustomBaseInput';

const CustomURLWidget = (props) => {
  return <BaseInput {...props} options={{ ...props.options }} />;
};

export default CustomURLWidget;
