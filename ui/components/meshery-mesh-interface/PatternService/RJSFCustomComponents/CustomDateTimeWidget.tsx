import BaseInput from './CustomBaseInput';
import { localToUTC, utcToLocal } from '@rjsf/utils';

type DateTimeWidgetProps = {
  value?: string;
  options?: object;
  onChange: (value: string) => void;
} & Record<string, unknown>;

const CustomDateTimeWidget = (props: DateTimeWidgetProps) => {
  const inputType = 'datetime-local';
  const value = utcToLocal(props.value);
  const onChange = (value: string) => {
    props.onChange(localToUTC(value));
  };
  return (
    <BaseInput
      {...props}
      options={{ ...props.options, inputType, focused: true }}
      onChange={onChange}
      value={value}
    />
  );
};
export default CustomDateTimeWidget;
