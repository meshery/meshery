import React from 'react';
import { render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import ReactSelectWrapper from './ReactSelectWrapper';
import MultiSelectWrapper from './multi-select-wrapper';

let capturedSelectProps: any[] = [];

vi.mock('react-select/creatable', () => ({
  default: (props) => {
    capturedSelectProps.push(props);
    return <div data-testid="creatable-select" />;
  },
}));

vi.mock('react-select', () => ({
  components: {
    Input: ({ children, ...props }) => <div {...props}>{children}</div>,
  },
}));

describe('select wrappers', () => {
  beforeEach(() => {
    capturedSelectProps = [];
  });

  it('renders ReactSelectWrapper option content without a menu context error', () => {
    render(
      <ReactSelectWrapper
        label="Environment"
        placeholder="Select environment"
        onChange={vi.fn()}
        options={[{ label: 'Production', value: 'prod' }]}
        value={null}
      />,
    );

    const Option = capturedSelectProps.at(-1).components.Option;

    render(
      <Option
        innerRef={vi.fn()}
        innerProps={{ onClick: vi.fn(), role: 'option' }}
        isFocused={false}
        isSelected={false}
      >
        Production
      </Option>,
    );

    expect(screen.getByText('Production')).toBeInTheDocument();
  });

  it('renders MultiSelectWrapper option content without a menu context error', () => {
    render(
      <MultiSelectWrapper
        onChange={vi.fn()}
        options={[
          { label: 'All', value: '*' },
          { label: 'Production', value: 'prod' },
        ]}
        value={[]}
      />,
    );

    const Option = capturedSelectProps.at(-1).components.Option;

    render(
      <Option
        innerRef={vi.fn()}
        innerProps={{ onClick: vi.fn(), role: 'option' }}
        isFocused={true}
        isSelected={true}
        label="Production"
        value="prod"
      />,
    );

    expect(screen.getByText('Production')).toBeInTheDocument();
  });
});
