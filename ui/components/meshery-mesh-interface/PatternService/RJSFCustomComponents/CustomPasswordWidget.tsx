import BaseInput from './CustomBaseInput';

const CustomPasswordWidget = (props: { options?: object } & Record<string, unknown>) => {
  const { options } = props;

  return <BaseInput {...props} options={{ ...options, inputType: 'password' }} />;
};

export default CustomPasswordWidget;
