import BaseInput from './CustomBaseInput';

const CustomColorWidget = (props: { options?: object } & Record<string, unknown>) => {
  const { options } = props;

  return <BaseInput {...props} options={{ ...options, inputType: 'color' }} />;
};
export default CustomColorWidget;
