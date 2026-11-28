// The inline-error half of the load-parameter contract, asserted against the
// real Sistent fields: a validation message for Concurrent requests or Queries
// per second must reach the user attached to the field it belongs to, so each
// assertion follows the input's own `aria-describedby` to the rendered text
// instead of searching the document for a loose string.
import React, { useState } from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { SistentThemeProvider } from '@sistent/sistent';
import { describe, expect, it, vi } from 'vitest';
import PerformanceForm from './PerformanceForm';
import PerformanceFormActions from './PerformanceFormActions';
import { createPerformanceFormChangeHandler } from './performance-handlers';
import { validateLoadParams } from './helper';

const renderForm = (overrides: Record<string, unknown> = {}) =>
  render(
    <SistentThemeProvider initialMode="light">
      <PerformanceForm
        profileName=""
        meshName=""
        selectedMesh=""
        meshModels={[]}
        url="http://example.com"
        urlError={false}
        c="1"
        cError=""
        qps="0"
        qpsError=""
        t="30s"
        tValue="30s"
        tError=""
        headers=""
        cookies=""
        contentType=""
        reqBody=""
        additionalOptions=""
        jsonError={false}
        caCertificate={{}}
        metadata={{ ca_certificate: {} }}
        loadGenerator="fortio"
        handleChange={() => vi.fn()}
        handleDurationChange={vi.fn()}
        handleInputDurationChange={vi.fn()}
        handleCertificateUpload={vi.fn()}
        {...overrides}
      />
    </SistentThemeProvider>,
  );

const field = (name: RegExp) => screen.getByRole('spinbutton', { name });

/** The message the field itself points at, or null when it points at nothing. */
const messageFor = (input: HTMLElement) => {
  const describedBy = input.getAttribute('aria-describedby');
  if (!describedBy) return null;
  const text = describedBy
    .split(' ')
    .map((id) => document.getElementById(id)?.textContent?.trim())
    .filter(Boolean)
    .join(' ');
  return text === '' ? null : text;
};

describe('PerformanceForm load parameter errors', () => {
  it('shows the concurrency message on the Concurrent requests field only', () => {
    renderForm({ cError: 'Enter a whole number of 1 or more', c: '0' });

    const concurrency = field(/Concurrent requests/i);
    expect(concurrency).toHaveAttribute('aria-invalid', 'true');
    expect(messageFor(concurrency)).toContain('Enter a whole number of 1 or more');

    const qps = field(/Queries per second/i);
    expect(qps).toHaveAttribute('aria-invalid', 'false');
    expect(messageFor(qps)).toBeNull();
  });

  it('shows the QPS message on the Queries per second field only', () => {
    renderForm({ qpsError: 'Enter a whole number of 0 or more (0 = max QPS)', qps: '-5' });

    const qps = field(/Queries per second/i);
    expect(qps).toHaveAttribute('aria-invalid', 'true');
    expect(messageFor(qps)).toContain('Enter a whole number of 0 or more (0 = max QPS)');

    const concurrency = field(/Concurrent requests/i);
    expect(concurrency).toHaveAttribute('aria-invalid', 'false');
    expect(messageFor(concurrency)).toBeNull();
  });

  it('shows no load parameter message while both values are valid', () => {
    renderForm();

    expect(field(/Concurrent requests/i)).toHaveAttribute('aria-invalid', 'false');
    expect(field(/Queries per second/i)).toHaveAttribute('aria-invalid', 'false');
    expect(screen.queryByText(/Enter a whole number/i)).not.toBeInTheDocument();
  });
});

// The same wiring `performance/index.tsx` uses, with nothing mocked between the
// keystroke and the buttons: the real change handler stores what was typed,
// `validateLoadParams` judges it, and the result reaches both the field and the
// footer. The tests above hand each component its props directly, so only this
// one shows that typing a bad value actually locks Save Profile and Run Test.
const LoadParamsHarness = () => {
  const [c, setC] = useState<string | number>('1');
  const [qps, setQps] = useState<string | number>('0');
  const handleChange = createPerformanceFormChangeHandler({
    setProfileName: vi.fn(),
    setMeshName: vi.fn(),
    setC,
    setQps,
    setHeaders: vi.fn(),
    setCookies: vi.fn(),
    setContentType: vi.fn(),
    setReqBody: vi.fn(),
    setLoadGenerator: vi.fn(),
    setUrl: vi.fn(),
    setCaCertificate: vi.fn(),
    setAdditionalOptions: vi.fn(),
    setJsonError: vi.fn(),
    setDisableTest: vi.fn(),
    setUrlError: vi.fn(),
  });
  const { cError, qpsError } = validateLoadParams(c, qps);

  return (
    <SistentThemeProvider initialMode="light">
      <PerformanceForm
        profileName=""
        meshName=""
        selectedMesh=""
        meshModels={[]}
        url="http://example.com"
        urlError={false}
        c={c}
        cError={cError}
        qps={qps}
        qpsError={qpsError}
        t="30s"
        tValue="30s"
        tError=""
        headers=""
        cookies=""
        contentType=""
        reqBody=""
        additionalOptions=""
        jsonError={false}
        caCertificate={{}}
        metadata={{ ca_certificate: {} }}
        loadGenerator="fortio"
        handleChange={handleChange}
        handleDurationChange={vi.fn()}
        handleInputDurationChange={vi.fn()}
        handleCertificateUpload={vi.fn()}
      />
      <PerformanceFormActions
        disableTest={false}
        invalidLoadParams={!!(cError || qpsError)}
        blockRunTest={false}
        hasTestResult={false}
        onAbort={vi.fn()}
        onShowResults={vi.fn()}
        onSaveProfile={vi.fn()}
        onRunTest={vi.fn()}
      />
    </SistentThemeProvider>
  );
};

describe('PerformanceForm load parameters, typed into the real form', () => {
  it('locks Save Profile and Run Test while Concurrent requests is invalid, and unlocks them once fixed', async () => {
    const user = userEvent.setup();
    render(<LoadParamsHarness />);

    const concurrency = field(/Concurrent requests/i);
    const saveProfile = screen.getByRole('button', { name: /^Save Profile$/ });
    const runTest = screen.getByTestId('run-performance-test');
    expect(saveProfile).toBeEnabled();
    expect(runTest).toBeEnabled();

    await user.clear(concurrency);
    await user.type(concurrency, '0');

    expect(concurrency).toHaveAttribute('aria-invalid', 'true');
    expect(messageFor(concurrency)).toContain('Enter a whole number of 1 or more');
    expect(saveProfile).toBeDisabled();
    expect(runTest).toBeDisabled();

    await user.clear(concurrency);
    await user.type(concurrency, '5');

    expect(concurrency).toHaveAttribute('aria-invalid', 'false');
    expect(messageFor(concurrency)).toBeNull();
    expect(saveProfile).toBeEnabled();
    expect(runTest).toBeEnabled();
  });
});
