import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

const promptComponentMock = vi.fn();

vi.mock('@sistent/sistent', () => ({
  PromptComponent: React.forwardRef((props: any, ref: any) => {
    promptComponentMock(props);
    return (
      <div data-testid="prompt" data-foo={props.foo}>
        <span ref={ref}>prompt-child</span>
      </div>
    );
  }),
}));

import PromptComponent from './PromptComponent';

describe('PromptComponent wrapper', () => {
  beforeEach(() => {
    promptComponentMock.mockClear();
  });

  it('forwards props to the underlying PromptComponent', () => {
    render(<PromptComponent foo="bar" />);
    expect(promptComponentMock).toHaveBeenCalledWith(expect.objectContaining({ foo: 'bar' }));
  });

  it('renders the underlying component output', () => {
    render(<PromptComponent />);
    expect(screen.getByTestId('prompt')).toBeInTheDocument();
  });

  it('forwards refs to the prompt component', () => {
    const refSpy = vi.fn();
    render(<PromptComponent ref={refSpy as any} />);
    expect(refSpy).toHaveBeenCalled();
  });
});
