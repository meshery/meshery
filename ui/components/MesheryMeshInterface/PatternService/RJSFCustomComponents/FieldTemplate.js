import React from 'react';
import { getTemplate, getUiOptions } from '@rjsf/utils';

export const CustomFieldTemplate = (props) => {
  const { children, registry, uiSchema } = props;
  const uiOptions = getUiOptions(uiSchema);
  const WrapIfAdditionalTemplate = getTemplate('WrapIfAdditionalTemplate', registry, uiOptions);

  return <WrapIfAdditionalTemplate {...props}>{children}</WrapIfAdditionalTemplate>;
};
