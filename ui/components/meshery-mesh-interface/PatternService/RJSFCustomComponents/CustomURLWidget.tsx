import BaseInput from './CustomBaseInput';

const CustomURLWidget = (props: { options?: object } & Record<string, unknown>) => {
  return <BaseInput {...props} options={{ ...props.options }} />;
};

export default CustomURLWidget;
