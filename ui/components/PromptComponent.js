import React, { forwardRef } from 'react';
import { PromptComponent } from '@sistent/sistent';

const _PromptComponent = forwardRef((props, ref) => {
  return <PromptComponent {...props} ref={ref} />;
});

_PromptComponent.displayName = '_PromptComponent';

export default _PromptComponent;
