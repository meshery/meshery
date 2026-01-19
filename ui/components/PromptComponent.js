import React, { forwardRef } from 'react';
import { PromptComponent } from '@sistent/sistent';

const _PromptComponent = forwardRef(({ options, ...props }, ref) => {
  return <PromptComponent {...props} options={options} ref={ref} />;
});

_PromptComponent.displayName = '_PromptComponent';

export default _PromptComponent;
