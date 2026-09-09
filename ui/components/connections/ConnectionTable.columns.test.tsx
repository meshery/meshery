import React from 'react';
import { render, renderHook } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

// Key-aware so a test can deny one capability at a time. A blanket boolean
// cannot tell a dropped gate from a cross-wired one.
const deniedPermissionIds = new Set<string>();

vi.mock('@sistent/sistent', () => ({
  useHasPermission: (key: { id?: string }) => !deniedPermissionIds.has(key?.id ?? ''),
  Box: ({ children }: any) => <div>{children}</div>,
  Grid2: ({ children }: any) => <div>{children}</div>,
  IconButton: ({ children }: any) => <button type="button">{children}</button>,
  TableCell: ({ children }: any) => <td>{children}</td>,
  InfoOutlinedIcon: () => <svg />,
  MoreVertIcon: () => <svg />,
  CustomTooltip: ({ children }: any) => <span>{children}</span>,
  getRelativeTime: () => 'a while ago',
  getFullFormattedTime: () => '2026-01-01',
  styled:
    () =>
    () =>
    ({ children }: any) => <div>{children}</div>,
}));

vi.mock('../../css/icons.styles', () => ({
  iconMedium: {},
  iconSmall: {},
}));

vi.mock('../general/multi-select-wrapper', () => ({
  // Surface `disabled` so permission-denied rendering is assertable.
  default: ({ disabled }: any) => (
    <div data-testid="multi-select-wrapper" data-disabled={String(Boolean(disabled))} />
  ),
}));

vi.mock('./ConnectionStatusSelect', () => ({
  ConnectionStatusSelect: ({ disabled }: any) => (
    <div data-testid="connection-status-select" data-disabled={String(Boolean(disabled))} />
  ),
}));

vi.mock('./ConnectionChip', () => ({
  TooltipWrappedConnectionChip: () => <div data-testid="connection-chip" />,
}));

vi.mock('../meshery-mesh-interface/PatternService/CustomTextTooltip', () => ({
  CustomTextTooltip: ({ children }: any) => <div>{children}</div>,
}));

// `utils/utils` re-exports through `pages/_app`, which pulls the whole app in.
// Only the row-value reader is needed here.
vi.mock('../../utils/utils', () => ({
  getColumnValue: (rowData: unknown[], name: string, columns: Array<{ name: string }>) =>
    rowData[columns.findIndex((col) => col.name === name)],
}));

import { Keys } from '@meshery/schemas/permissions';
import { useConnectionColumns } from './ConnectionTable.columns';

const renderColumns = () =>
  renderHook(() =>
    useConnectionColumns({
      url: 'https://docs.meshery.io',
      envUrl: 'https://docs.meshery.io/environments',
      environmentOptions: [{ label: 'dev', value: 'env-1' }],
      isEnvironmentsSuccess: true,
      updatingConnection: { current: false },
      handleDeleteConnection: vi.fn(),
      handleEnvironmentSelect: vi.fn(),
      handleStatusChange: vi.fn(),
      handleActionMenuOpen: vi.fn(),
      ping: vi.fn(),
      pingGrafana: vi.fn(),
      pingPrometheus: vi.fn(),
      transitionMapByKind: { kubernetes: ['connected'] },
    } as any),
  ).result.current;

/** Renders one column's body cell and reports its `disabled` state. */
const cellDisabledState = (
  columns: any[],
  columnName: string,
  value: unknown,
  testId: string,
): string | undefined => {
  const column = columns.find((col) => col.name === columnName);
  const { container, unmount } = render(
    <>{column.options.customBodyRender(value, { rowData: [] })}</>,
  );
  const state = container.querySelector<HTMLElement>(`[data-testid="${testId}"]`)?.dataset.disabled;
  unmount();
  return state;
};

const rowControlStates = () => {
  const columns = renderColumns();
  return {
    environments: cellDisabledState(columns, 'environments', [], 'multi-select-wrapper'),
    status: cellDisabledState(columns, 'status', 'connected', 'connection-status-select'),
  };
};

describe('useConnectionColumns permission gating', () => {
  beforeEach(() => {
    deniedPermissionIds.clear();
  });

  // The environment select and the status select are the table's only mutating
  // row affordances. Each is asserted in both directions and against its own
  // key, so a gate that is dropped, inverted, or cross-wired fails here.
  it('enables both row controls when their permissions are granted', () => {
    expect(rowControlStates()).toEqual({ environments: 'false', status: 'false' });
  });

  it('disables only the environment select when its permission is denied', () => {
    deniedPermissionIds.add(Keys.WorkspaceManagementAssignConnectionsToEnvironment.id);
    expect(rowControlStates()).toEqual({ environments: 'true', status: 'false' });
  });

  it('disables only the status select when its permission is denied', () => {
    deniedPermissionIds.add(Keys.LifecycleManagementChangeConnectionState.id);
    expect(rowControlStates()).toEqual({ environments: 'false', status: 'true' });
  });

  // A deleted connection has no valid transition, so it stays locked even for a
  // user who holds the state-change permission.
  it('keeps a deleted connection locked regardless of permission', () => {
    const columns = renderColumns();
    expect(cellDisabledState(columns, 'status', 'deleted', 'connection-status-select')).toBe(
      'true',
    );
  });
});
