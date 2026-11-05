import BaseInput from './CustomBaseInput';

type TextAreaWidgetProps = { options?: { rows?: number } } & Record<string, unknown>;

const CustomTextAreaWidget = (props: TextAreaWidgetProps) => {
  const { options } = props;
  const rows = options?.rows || 3;

  return <BaseInput {...props} options={{ ...options, rows }} multiline={true} />;
};
export default CustomTextAreaWidget;
