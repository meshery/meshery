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
    Box: ({ children }: any) => <div>{children}</div>,
    Button: ({ children }: any) => <button>{children}</button>,
    DialogTitle: ({ children }: any) => <div>{children}</div>,
    Typography: ({ children }: any) => <span>{children}</span>,
  };
});

import * as styles from './Filters.styled';

describe('Filters.styled', () => {
  it('exports all expected styled components', () => {
    expect(styles.CreateButton).toBeDefined();
    expect(styles.ViewSwitchButton).toBeDefined();
    expect(styles.YmlDialogTitle).toBeDefined();
    expect(styles.YmlDialogTitleText).toBeDefined();
    expect(styles.BtnText).toBeDefined();
    expect(styles.ActionsBox).toBeDefined();
  });

  it('renders the styled components without error', () => {
    const { container } = render(
      <styles.CreateButton>
        <styles.BtnText>Hi</styles.BtnText>
      </styles.CreateButton>,
    );
    expect(container).toBeInTheDocument();
  });
});
