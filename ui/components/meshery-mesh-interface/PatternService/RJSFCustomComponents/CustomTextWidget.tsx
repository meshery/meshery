import BaseInput from './CustomBaseInput';

type TextWidgetProps = {
  schema?: { type?: string };
  options?: object;
} & Record<string, unknown>;

const CustomTextWidget = (props: TextWidgetProps) => {
  let inputType = 'string';
  if (props.schema?.type === 'number' || props.schema?.type === 'integer') {
    inputType = 'number';
  }
  return <BaseInput {...props} options={{ ...props.options, inputType }} />;
};

export default CustomTextWidget;
