import React, { forwardRef, Ref, ComponentProps } from 'react';
import { PromptComponent } from '@sistent/sistent';

type PromptComponentProps = ComponentProps<typeof PromptComponent>;

export interface PromptRef {
  show: (_options: {
    title: string;
    subtitle?: string;
    options?: string[];
    showInfoIcon?: string;
    variant?: 'default' | 'danger' | 'warning' | 'info' | 'success';
  }) => Promise<string>;
}

const _PromptComponent = forwardRef<PromptRef, PromptComponentProps>((props, ref) => {
  return <PromptComponent {...props} ref={ref as Ref<unknown>} />;
});

_PromptComponent.displayName = '_PromptComponent';

export default _PromptComponent;
