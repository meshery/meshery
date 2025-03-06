import React, { forwardRef } from 'react';
import { PromptComponent } from '@layer5/sistent';

const _PromptComponent = forwardRef((props, ref) => {
  return <PromptComponent {...props} ref={ref} />;
});

_PromptComponent.displayName = '_PromptComponent';

export default _PromptComponent;
