// @ts-nocheck
import { getTemplate, getUiOptions } from '@rjsf/utils';

export const CustomFieldTemplate = (props) => {
  const { children, registry, uiSchema } = props;
  const uiOptions = getUiOptions(uiSchema);
  const WrapIfAdditionalTemplate = getTemplate('WrapIfAdditionalTemplate', registry, uiOptions);
  if (typeof WrapIfAdditionalTemplate !== 'function') {
    return <>{children}</>;
  }
  return <WrapIfAdditionalTemplate {...props}>{children}</WrapIfAdditionalTemplate>;
};
