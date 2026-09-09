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

  // `BasicMarkdown` emits sibling block elements. As direct children of the
  // content row they became sibling flex items and laid out horizontally, so a
  // multi-block MeshKit error read "Unable to create the environmentmeshery-
  // server-1448". The markdown needs its own column inside the row.
  it('stacks the markdown blocks in a column instead of the content row', () => {
    render(<ThemeResponsiveSnackbar id={1} variant="error" message={'**Failed**\n\n`code-1`'} />);

    const row = screen.getByTestId('SnackbarContent-error');
    const messageColumn = screen.getByTestId('SnackbarContent-message');

    expect(row.style.flexDirection).not.toBe('column');
    expect(messageColumn.parentElement).toBe(row);
    expect(messageColumn.style.flexDirection).toBe('column');
    // Long unbreakable strings must wrap rather than overflow the toast.
    expect(messageColumn.style.minWidth).toBe('0px');
  });
});
