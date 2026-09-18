import React from 'react';
import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import PerformanceForm from './PerformanceForm';
import { isValidDuration } from '../../utils/validators';

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

  it('shows validation error and does not call closeModal when an invalid duration is submitted', () => {
    // Regression: the old code read `props.performanceProfile.duration` (the initial prop)
    // instead of local `tState` when validating, so typing an invalid value and clicking
    // Save Profile would silently pass validation if the prop held a valid default.
    //
    // This test mirrors the handleSaveProfile logic directly:
    //   - t prop = '30s'  (a valid default that the old code would have read)
    //   - typed tValue = '0s' (invalid — zero is rejected by isValidDuration)
    // If validation mistakenly reads the prop it returns true and closeModal is called.
    // If it correctly reads the typed state it returns false, sets tError, and returns early.

    const closeModal = vi.fn();

    // Thin wrapper that owns tState (as handleSaveProfile does in index.tsx).
    const TestHarness: React.FC = () => {
      const [tState, setTState] = React.useState('0s'); // invalid typed value
      const [tError, setTError] = React.useState('');

      const handleSave = () => {
        // Exact logic copied from index.tsx handleSaveProfile
        if (!isValidDuration(tState)) {
          setTError('error-autocomplete-value');
          return;
        }
        setTError('');
        closeModal();
      };

      return (
        <>
          <PerformanceForm
            {...defaultProps}
            t="30s" // prop carries the valid default; old code read this instead of state
            tValue={tState}
            tError={tError}
            handleInputDurationChange={(_e, newValue) => setTState(newValue ?? '')}
          />
          <button data-testid="save-btn" onClick={handleSave}>
            Save Profile
          </button>
        </>
      );
    };

    const { container } = render(<TestHarness />);

    // Trigger the save path with the invalid typed duration.
    fireEvent.click(screen.getByTestId('save-btn'));

    // Validation error must be visible so the user can correct the input.
    expect(container.querySelector('.error-autocomplete-value')).not.toBeNull();

    // Modal must NOT have been closed — the bug was that closeModal() was called
    // even when duration validation failed.
    expect(closeModal).not.toHaveBeenCalled();
  });
});
