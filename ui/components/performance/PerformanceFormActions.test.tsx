// What the footer buttons are responsible for, asserted against the real
// Sistent buttons rather than against the parent's props.
//
// Two separate gates govern these four buttons, and the split is the contract:
// `disableTest` (the URL/JSON gate) governs all four, while an invalid
// concurrency/QPS pair must block only the two buttons that submit - Save
// Profile and Run Test - so Clear stays available as the form's reset escape
// hatch and Results keeps already-fetched results reachable while the user
// repairs the field.
import React from 'react';
import { render, screen } from '@testing-library/react';
import { SistentThemeProvider } from '@sistent/sistent';
import { describe, expect, it, vi } from 'vitest';
import PerformanceFormActions from './PerformanceFormActions';

const renderActions = (overrides: Record<string, unknown> = {}) =>
  render(
    <SistentThemeProvider initialMode="light">
      <PerformanceFormActions
        disableTest={false}
        invalidLoadParams={false}
        blockRunTest={false}
        hasTestResult={true}
        onAbort={vi.fn()}
        onShowResults={vi.fn()}
        onSaveProfile={vi.fn()}
        onRunTest={vi.fn()}
        {...overrides}
      />
    </SistentThemeProvider>,
  );

const button = (name: RegExp) => screen.getByRole('button', { name });

describe('PerformanceFormActions', () => {
  it('blocks only the submitting buttons while the load parameters are invalid', () => {
    renderActions({ invalidLoadParams: true });

    expect(button(/^Save Profile$/)).toBeDisabled();
    expect(button(/^Run Test$/)).toBeDisabled();
    expect(button(/^Clear$/)).toBeEnabled();
    expect(button(/^Results$/)).toBeEnabled();
  });

  it('enables every button once the load parameters are valid', () => {
    renderActions({ invalidLoadParams: false });

    expect(button(/^Save Profile$/)).toBeEnabled();
    expect(button(/^Run Test$/)).toBeEnabled();
    expect(button(/^Clear$/)).toBeEnabled();
    expect(button(/^Results$/)).toBeEnabled();
  });

  it('still disables all four buttons while the test itself is disabled', () => {
    renderActions({ disableTest: true, invalidLoadParams: false });

    expect(button(/^Save Profile$/)).toBeDisabled();
    expect(button(/^Run Test$/)).toBeDisabled();
    expect(button(/^Clear$/)).toBeDisabled();
    expect(button(/^Results$/)).toBeDisabled();
  });
});
