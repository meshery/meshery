import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

const formCoreMock = vi.fn();

vi.mock('./PatternServiceFormCore', () => ({
  default: (props: any) => {
    formCoreMock(props);
    const SettingsForm = () => <div data-testid="settings-form" />;
    return <div data-testid="pattern-form-core">{props.children(SettingsForm)}</div>;
  },
}));

vi.mock('@sistent/sistent', () => ({
  AppBar: ({ children }: any) => <div data-testid="appbar">{children}</div>,
  Box: ({ children }: any) => <div>{children}</div>,
  IconButton: ({ children, onClick }: any) => (
    <button data-testid="icon-btn" onClick={onClick}>
      {children}
    </button>
  ),
  Toolbar: ({ children }: any) => <div data-testid="toolbar">{children}</div>,
  Tooltip: ({ children }: any) => <div>{children}</div>,
  useTheme: () => ({ palette: { mode: 'light' } }),
}));

vi.mock('@/theme', () => ({
  darken: (c: string) => `dark(${c})`,
}));

vi.mock('@/assets/icons', () => ({
  Delete: () => <svg data-testid="delete-icon" />,
  HelpOutlined: () => <svg data-testid="help-icon" />,
  Settings: () => <svg data-testid="settings-icon" />,
}));

vi.mock('../../css/icons.styles', () => ({ iconSmall: {} }));

import PatternServiceForm from './PatternServiceForm';

describe('PatternServiceForm', () => {
  beforeEach(() => {
    formCoreMock.mockClear();
  });

  it('augments the workload schema with name/namespace/labels/annotations', () => {
    const workload: any = { properties: {} };
    render(
      <PatternServiceForm
        formData={{}}
        schemaSet={{ workload, type: 'workload' }}
        onSubmit={vi.fn()}
        onDelete={vi.fn()}
        namespace="default"
      />,
    );

    expect(workload.properties.name).toBeDefined();
    expect(workload.properties.namespace).toBeDefined();
    expect(workload.properties.labels).toBeDefined();
    expect(workload.properties.annotations).toBeDefined();
  });

  it('renders the settings toolbar with a Settings label', () => {
    const workload: any = { properties: {}, description: 'A description' };
    render(
      <PatternServiceForm
        formData={{}}
        schemaSet={{ workload, type: 'workload' }}
        onSubmit={vi.fn()}
        onDelete={vi.fn()}
        namespace="default"
      />,
    );
    expect(screen.getByText('Settings')).toBeInTheDocument();
    expect(screen.getByTestId('settings-icon')).toBeInTheDocument();
  });

  it('invokes reference.current.delete when the trash icon is clicked', () => {
    const deleteSpy = vi.fn();
    const reference = { current: { delete: deleteSpy } };
    const workload: any = { properties: {} };
    render(
      <PatternServiceForm
        formData={{}}
        schemaSet={{ workload, type: 'workload' }}
        onSubmit={vi.fn()}
        onDelete={vi.fn()}
        namespace="default"
        reference={reference}
      />,
    );
    // Find the icon button that contains the Delete icon. There are two icon buttons
    // (help + delete). Click the one that renders the delete icon.
    const allButtons = screen.getAllByTestId('icon-btn');
    const deleteBtn = allButtons.find((btn) => btn.querySelector('[data-testid="delete-icon"]'));
    expect(deleteBtn).toBeDefined();
    fireEvent.click(deleteBtn!);
    expect(deleteSpy).toHaveBeenCalled();
  });
});
