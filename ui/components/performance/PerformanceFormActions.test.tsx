import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  Button: ({ children, disabled, onClick }: any) => (
    <button disabled={disabled} onClick={onClick}>
      {children}
    </button>
  ),
  CircularProgress: () => <span />,
}));

vi.mock('@/assets/icons', () => ({ SaveOutlined: () => <svg /> }));

vi.mock('@meshery/schemas/permissions', () => ({
  Keys: new Proxy({}, { get: () => ({}) }),
}));

import PerformanceFormActions from './PerformanceFormActions';

describe('PerformanceFormActions', () => {
  const baseProps = {
    disableTest: false,
    hasJsonError: false,
    blockRunTest: false,
    hasTestResult: true,
    onAbort: vi.fn(),
    onShowResults: vi.fn(),
    onSaveProfile: vi.fn(),
    onRunTest: vi.fn(),
  };

  const button = (name: string) => screen.getByRole('button', { name });

  it('enables every action for a valid form', () => {
    render(<PerformanceFormActions {...baseProps} />);

    ['Clear', 'Results', 'Save Profile', 'Run Test'].forEach((name) => {
      expect(button(name)).toBeEnabled();
    });
  });

  it('blocks saving and running while additional options are not valid JSON', () => {
    render(<PerformanceFormActions {...baseProps} hasJsonError />);

    expect(button('Save Profile')).toBeDisabled();
    expect(button('Run Test')).toBeDisabled();
  });

  it('keeps Clear available so a JSON error can be reset', () => {
    render(<PerformanceFormActions {...baseProps} hasJsonError />);

    expect(button('Clear')).toBeEnabled();
    expect(button('Results')).toBeEnabled();
  });
});
