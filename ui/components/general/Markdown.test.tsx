import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

// Stub the CSS side-effect imports.
vi.mock('@uiw/react-markdown-preview/markdown.css', () => ({}));
vi.mock('@uiw/react-md-editor/markdown-editor.css', () => ({}));

// next/dynamic returns a thin component that just renders its mocked target.
vi.mock('next/dynamic', () => ({
  default: (factory: any) => {
    const MockedDynamicComponent = (props: any) => (
      <div data-testid="md-editor" {...props}>
        {props.value}
      </div>
    );
    MockedDynamicComponent.displayName = 'MockedDynamicComponent';
    // Capture the factory to allow assertion that next/dynamic was used.
    (MockedDynamicComponent as any).__factory = factory;
    return MockedDynamicComponent;
  },
}));

vi.mock('@uiw/react-md-editor', () => ({
  default: (props: any) => <div data-testid="real-md-editor" {...props} />,
}));

import { MDEditor } from './Markdown';

describe('Markdown (MDEditor)', () => {
  it('exports an MDEditor component', () => {
    expect(MDEditor).toBeDefined();
  });

  it('renders the editor', () => {
    render(<MDEditor value="hello" />);
    expect(screen.getByTestId('md-editor')).toBeInTheDocument();
  });

  it('passes props through to the editor', () => {
    render(<MDEditor value="some markdown content" />);
    expect(screen.getByTestId('md-editor')).toHaveTextContent('some markdown content');
  });
});
