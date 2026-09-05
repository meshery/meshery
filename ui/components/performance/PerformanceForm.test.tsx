import React from 'react';
import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import PerformanceForm from './PerformanceForm';

vi.mock('@/assets/icons', () => ({
  ExpandMore: () => <svg data-testid="expand-more-icon" />,
  HelpOutlineOutlined: () => <svg data-testid="help-icon" />,
}));

vi.mock('../meshery-mesh-interface/PatternService/CustomTextTooltip', () => ({
  CustomTextTooltip: ({ children }: any) => <>{children}</>,
}));

vi.mock('./style', () => ({
  ExpansionPanelComponent: ({ children }: any) => <div>{children}</div>,
  FormContainer: ({ children }: any) => <div>{children}</div>,
  HelpIcon: ({ children }: any) => <span>{children}</span>,
  RadioButton: ({ children, ...rest }: any) => <input type="radio" {...rest} />,
}));

const defaultProps = {
  profileName: 'Test Profile',
  meshName: 'istio',
  selectedMesh: 'istio',
  meshModels: ['Istio', 'Linkerd'],
  url: 'https://meshery.io',
  urlError: false,
  c: 1,
  qps: 10,
  t: '30s',
  tValue: '30s',
  tError: '',
  headers: '',
  cookies: '',
  contentType: '',
  reqBody: '',
  additionalOptions: '',
  jsonError: false,
  caCertificate: {},
  metadata: {},
  loadGenerator: 'fortio',
  handleChange: vi.fn(() => vi.fn()),
  handleDurationChange: vi.fn(),
  handleInputDurationChange: vi.fn(),
  handleCertificateUpload: vi.fn(),
};

describe('PerformanceForm', () => {
  it('renders duration input with initial value', () => {
    render(<PerformanceForm {...defaultProps} />);

    const durationInput = screen.getByRole('combobox', { name: /duration/i });
    expect(durationInput).toBeDefined();
    expect((durationInput as HTMLInputElement).value).toBe('30s');
  });

  it('applies error-autocomplete-value class when tError is set', () => {
    const { container } = render(
      <PerformanceForm {...defaultProps} tError="error-autocomplete-value" />,
    );

    const errorElement = container.querySelector('.error-autocomplete-value');
    expect(errorElement).not.toBeNull();
  });

  it('does not have error class when tError is empty', () => {
    const { container } = render(<PerformanceForm {...defaultProps} tError="" />);

    const errorElement = container.querySelector('.error-autocomplete-value');
    expect(errorElement).toBeNull();
  });

  it('invokes handleInputDurationChange when user types in duration input', () => {
    const handleInputDurationChange = vi.fn();
    render(
      <PerformanceForm {...defaultProps} handleInputDurationChange={handleInputDurationChange} />,
    );

    const durationInput = screen.getByRole('combobox', { name: /duration/i });
    fireEvent.change(durationInput, { target: { value: '45s' } });

    expect(handleInputDurationChange).toHaveBeenCalled();
  });

  it('shows error styling and does not fire handleDurationChange when tError is set for an invalid duration', () => {
    // Regression for: duration validation was evaluated against stale props instead of
    // local form state, so invalid values (e.g. "0s", "abc") silently passed validation.
    // handleSubmit / handleSaveProfile now sets tError="error-autocomplete-value" when
    // isValidDuration(tState) returns false, keeping the modal open for correction.
    const handleDurationChange = vi.fn();
    const { container } = render(
      <PerformanceForm
        {...defaultProps}
        t="abc"
        tValue="abc"
        tError="error-autocomplete-value"
        handleDurationChange={handleDurationChange}
      />,
    );

    // Error class must be present on the duration field so the user sees the validation error.
    const errorElement = container.querySelector('.error-autocomplete-value');
    expect(errorElement).not.toBeNull();

    // The parent component (not the form) is responsible for blocking submission,
    // so handleDurationChange should not have been called automatically.
    expect(handleDurationChange).not.toHaveBeenCalled();
  });
});
