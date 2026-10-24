import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  ErrorBoundary: ({ children }: any) => <div data-testid="error-boundary">{children}</div>,
}));

vi.mock('@/components/workspaces/SpacesSwitcher/WorkspaceContent', () => ({
  default: ({ workspace }: any) => (
    <div data-testid="workspace-content">
      workspace:{(workspace as { id?: string } | null | undefined)?.id ?? 'none'}
    </div>
  ),
}));

import WorkSpaceContentDataTable from './WorkSpaceContentDataTable';

describe('WorkSpaceContentDataTable', () => {
  it('wraps WorkspaceContent in an ErrorBoundary and forwards the workspace prop', () => {
    render(<WorkSpaceContentDataTable workspace={{ id: 'ws-1' }} />);
    expect(screen.getByTestId('error-boundary')).toBeInTheDocument();
    expect(screen.getByTestId('workspace-content')).toHaveTextContent('workspace:ws-1');
  });

  it('handles when workspace is not provided', () => {
    render(<WorkSpaceContentDataTable workspace={undefined as unknown as object} />);
    expect(screen.getByTestId('workspace-content')).toHaveTextContent('workspace:none');
  });
});
