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
    Box: ({ children, ...props }: any) => <div {...props}>{children}</div>,
  };
});

import * as styles from './styles';

describe('workspaces/styles', () => {
  it('exports all expected styled components', () => {
    expect(styles.TableIconsContainer).toBeDefined();
    expect(styles.IconWrapper).toBeDefined();
    expect(styles.CreateButtonWrapper).toBeDefined();
    expect(styles.BulkActionWrapper).toBeDefined();
    expect(styles.UserCommonBox).toBeDefined();
  });

  it('renders styled components without error', () => {
    const { container } = render(
      <styles.TableIconsContainer>
        <styles.IconWrapper>content</styles.IconWrapper>
      </styles.TableIconsContainer>,
    );
    expect(container).toBeInTheDocument();
  });
});
