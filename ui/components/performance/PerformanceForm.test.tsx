// The inline-error half of the load-parameter contract, asserted against the
// real Sistent fields: a validation message for Concurrent requests or Queries
// per second must reach the user attached to the field it belongs to, so each
// assertion follows the input's own `aria-describedby` to the rendered text
// instead of searching the document for a loose string.
import React from 'react';
import { render, screen } from '@testing-library/react';
import { SistentThemeProvider } from '@sistent/sistent';
import { describe, expect, it, vi } from 'vitest';
import PerformanceForm from './PerformanceForm';

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
