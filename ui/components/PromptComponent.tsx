import React, { Ref, ComponentProps } from 'react';
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

const _PromptComponent = ({ ref, ...props }: PromptComponentProps & { ref?: Ref<PromptRef> }) => (
  <PromptComponent {...props} ref={ref as Ref<unknown>} />
);

export default _PromptComponent;
