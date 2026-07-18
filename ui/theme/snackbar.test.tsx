import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  BasicMarkdown: ({ content }: { content: string }) => <span>{content}</span>,
  Box: ({ children }: { children?: React.ReactNode }) => <div>{children}</div>,
  CheckCircleIcon: () => null,
  CircularProgress: () => null,
  ErrorIcon: () => null,
  InfoIcon: () => null,
  WarningIcon: () => null,
  lighten: (color: string) => color,
  styled: (Component: any) => () => {
    const Styled = ({ children, variant: _variant, ...props }: any) =>
      typeof Component === 'string' ? (
        React.createElement(Component, props, children)
      ) : (
        <div {...props}>{children}</div>
      );
    Styled.displayName = 'StyledSnackbarMock';
    return Styled;
  },
}));

vi.mock('notistack', () => ({
  SnackbarContent: ({ children }: { children?: React.ReactNode }) => <div>{children}</div>,
}));

import { ThemeResponsiveSnackbar } from './snackbar';

describe('ThemeResponsiveSnackbar', () => {
  it('invokes the action callback with the notistack snackbar id, not the reserved key prop', () => {
    const action = vi.fn((snackbarId) => <button data-testid="close">{String(snackbarId)}</button>);

    render(
      <ThemeResponsiveSnackbar id="snack-42" variant="info" message="hello" action={action} />,
    );

    expect(action).toHaveBeenCalledWith('snack-42');
    expect(screen.getByTestId('close')).toHaveTextContent('snack-42');
  });

  it('renders the message inside the variant-tagged content wrapper', () => {
    render(<ThemeResponsiveSnackbar id={1} variant="success" message="saved" />);

    expect(screen.getByTestId('SnackbarContent-success')).toHaveTextContent('saved');
  });
});
