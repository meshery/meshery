import React, { forwardRef } from 'react';
import { PromptComponent } from '@layer5/sistent';
import { UsesSistent } from './SistentWrapper';

const _PromptComponent = forwardRef((props, ref) => {
  return (
    <UsesSistent>
      <PromptComponent {...props} ref={ref} />
    </UsesSistent>
  );
});

_PromptComponent.displayName = '_PromptComponent';

export default _PromptComponent;
