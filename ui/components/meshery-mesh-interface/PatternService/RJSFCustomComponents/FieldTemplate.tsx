import React from 'react';
import { getTemplate, getUiOptions, FieldTemplateProps } from '@rjsf/utils';

export const CustomFieldTemplate = (props: FieldTemplateProps) => {
  const { children, registry, uiSchema } = props;
  const uiOptions = getUiOptions(uiSchema);
  const WrapIfAdditionalTemplate = getTemplate('WrapIfAdditionalTemplate', registry, uiOptions);
  if (typeof WrapIfAdditionalTemplate !== 'function') {
    return <>{children}</>;
  }
  return <WrapIfAdditionalTemplate {...props}>{children}</WrapIfAdditionalTemplate>;
};
