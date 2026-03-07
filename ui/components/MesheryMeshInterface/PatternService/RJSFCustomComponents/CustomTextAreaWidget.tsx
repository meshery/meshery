// @ts-nocheck
import BaseInput from './CustomBaseInput';

const CustomTextAreaWidget = (props) => {
  const { options } = props;
  const rows = options.rows || 3;

  return <BaseInput {...props} options={{ ...options, rows }} multiline={true} />;
};
export default CustomTextAreaWidget;
