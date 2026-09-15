import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';

const h = vi.hoisted(() => ({
  openCreateConnectionMock: vi.fn(),
}));

vi.mock('@/utils/context/ConnectionWizardContextProvider', () => ({
  useConnectionWizardModal: () => ({
    openCreateConnection: h.openCreateConnectionMock,
    closeCreateConnection: vi.fn(),
    open: false,
    presetKind: null,
    skipKindSelection: false,
  }),
}));
vi.mock('@/utils/can', () => ({ default: () => true }));
vi.mock('@/utils/hooks/useTestIDs', () => ({ default: () => () => 'test-id' }));
vi.mock('../assets/icons/AddIconCircleBorder', () => ({ default: () => null }));
vi.mock('@sistent/sistent', () => ({
  Button: ({ children, onClick, disabled }: any) => (
    <button onClick={onClick} disabled={disabled}>
      {children}
    </button>
  ),
  Typography: ({ children, ...rest }: any) => <span {...rest}>{children}</span>,
}));

import MesherySettingsEnvButtons from './MesherySettingsEnvButtons';

describe('MesherySettingsEnvButtons – Add Cluster', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('opens the Create Connection wizard in place with Kubernetes pre-selected', () => {
    render(<MesherySettingsEnvButtons />);

    fireEvent.click(screen.getByText('Add Cluster'));

    expect(h.openCreateConnectionMock).toHaveBeenCalledWith({
      kind: 'kubernetes',
      skipKindSelection: true,
    });
  });
});
