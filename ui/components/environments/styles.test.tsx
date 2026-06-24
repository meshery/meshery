import React from 'react';
import { render } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => {
  const styled = (_Component: any) => (_factory?: any) => {
    const Styled = ({ children, ...props }: any) => <div {...props}>{children}</div>;
    Styled.displayName = 'StyledMock';
    return Styled;
  };
  return {
    styled,
    Checkbox: ({ children, ...props }: any) => <input type="checkbox" {...props} />,
    Typography: ({ children }: any) => <span>{children}</span>,
    Button: ({ children, ...props }: any) => <button {...props}>{children}</button>,
    useTheme: () => ({
      palette: {
        background: {
          card: 'card',
          brand: { default: 'brand' },
          constant: { table: 'table', white: 'white' },
          default: 'default',
        },
        text: { default: 'text-default', inverse: 'inverse' },
        mode: 'light',
      },
    }),
    Card: ({ children, ...props }: any) => <div {...props}>{children}</div>,
    Box: ({ children, ...props }: any) => <div {...props}>{children}</div>,
  };
});

import * as styles from './styles';

describe('environments/styles', () => {
  it('exports all expected styled components', () => {
    expect(styles.BulkActionWrapper).toBeDefined();
    expect(styles.CardWrapper).toBeDefined();
    expect(styles.Statistic).toBeDefined();
    expect(styles.StatisticName).toBeDefined();
    expect(styles.TabCardContent).toBeDefined();
    expect(styles.TabIconBox).toBeDefined();
    expect(styles.TabNameBox).toBeDefined();
    expect(styles.TabTitle).toBeDefined();
    expect(styles.TabCount).toBeDefined();
    expect(styles.AllocationButton).toBeDefined();
    expect(styles.AllocationWorkspace).toBeDefined();
    expect(styles.PopupButton).toBeDefined();
    expect(styles.Record).toBeDefined();
    expect(styles.BulkSelectCheckbox).toBeDefined();
    expect(styles.CardTitle).toBeDefined();
    expect(styles.DateLabel).toBeDefined();
    expect(styles.EmptyDescription).toBeDefined();
    expect(styles.DescriptionLabel).toBeDefined();
    expect(styles.Name).toBeDefined();
    expect(styles.Status).toBeDefined();
    expect(styles.StyledChip).toBeDefined();
    expect(styles.CreateButtonWrapper).toBeDefined();
    expect(styles.EditButton).toBeDefined();
    expect(styles.TextButton).toBeDefined();
    expect(styles.IconButton).toBeDefined();
  });

  it('renders styled components without error', () => {
    const { container } = render(
      <styles.BulkActionWrapper>
        <styles.Name>hello</styles.Name>
      </styles.BulkActionWrapper>,
    );
    expect(container).toBeInTheDocument();
  });
});
