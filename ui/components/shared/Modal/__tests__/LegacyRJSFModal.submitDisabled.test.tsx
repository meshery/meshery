// The Submit button in LegacyRJSFModal lives outside the `loadingSchema` branch,
// so it renders while the modal is still pending. Its click path is guarded by
// `formRef.current`, which is null until RJSFWrapper mounts - so a click during
// that window does nothing at all: no submit, no close, no feedback.
//
// These tests pin the two inputs that close that gap (issue #21227): the button
// is disabled for the whole pending window, and a consumer can disable it for
// its own reasons through `submitDisabled`. The title-derived `canNotSubmit`
// heuristic is asserted alongside them so the OR is not silently narrowed later.

import React from 'react';
import { render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const notify = vi.fn();
vi.mock('@/utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify }),
}));

// RJSFWrapper pulls the whole RJSF stack in; the form itself is not under test
// here, only whether the footer button is disabled.
vi.mock('../../../meshery-mesh-interface/PatternService/RJSF_wrapper', () => ({
  default: () => <div data-testid="rjsf-wrapper" />,
}));

vi.mock('../../../meshery-mesh-interface/PatternService/helper', () => ({
  getSchema: () => ({}),
}));

import { SistentThemeProvider } from '@/theme';
import Modal from '../LegacyRJSFModal';

const schema = { type: 'object', properties: {} };

const renderModal = (props: Record<string, unknown> = {}) =>
  render(
    <SistentThemeProvider>
      <Modal
        open
        title="Help & Support"
        handleClose={vi.fn()}
        handleSubmit={vi.fn()}
        schema_array={[]}
        type="helpAndSupport"
        schemaChangeHandler={vi.fn()}
        {...props}
      />
    </SistentThemeProvider>,
  );

const submitButton = () => screen.getByRole('button', { name: /submit/i });

describe('LegacyRJSFModal submit button', () => {
  beforeEach(() => {
    notify.mockReset();
  });

  it('is disabled while the modal is still pending, so the dead click cannot happen', () => {
    // No schema yet - the modal renders its spinner and RJSFWrapper is unmounted,
    // which is exactly when handleFormSubmit's formRef guard makes a click inert.
    renderModal();

    expect(screen.queryByTestId('rjsf-wrapper')).not.toBeInTheDocument();
    expect(submitButton()).toBeDisabled();
  });

  it('is enabled once the schema resolves and the form is mounted', () => {
    renderModal({ schema });

    expect(screen.getByTestId('rjsf-wrapper')).toBeInTheDocument();
    expect(submitButton()).toBeEnabled();
  });

  it('stays disabled when the consumer sets submitDisabled, even with a schema', () => {
    renderModal({ schema, submitDisabled: true });

    expect(screen.getByTestId('rjsf-wrapper')).toBeInTheDocument();
    expect(submitButton()).toBeDisabled();
  });

  it('defaults submitDisabled to false, so existing consumers are unaffected', () => {
    renderModal({ schema, submitDisabled: undefined });

    expect(submitButton()).toBeEnabled();
  });

  it('still honours the title-derived canNotSubmit heuristic', () => {
    renderModal({ schema, title: 'Untitled Design' });

    expect(notify).toHaveBeenCalled();
    expect(submitButton()).toBeDisabled();
  });
});
