import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const notify = vi.fn();
const dispatch = vi.fn();
const updateLoadTestPrefs = vi.fn().mockImplementation(() => ({
  unwrap: () => Promise.resolve(),
}));
let isSaving = false;

let selectedStore: any = {
  prefTest: {
    loadTestPref: { qps: 10, c: 5, t: '30s', gen: 'fortio' },
  },
  ui: { selectedK8sContexts: [] },
};

vi.mock('../../lib/prePopulatedOptions', () => ({
  durationOptions: ['10s', '30s', '1m', '5m'],
}));

vi.mock('../../utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify }),
}));

vi.mock('../../lib/event-types', () => ({
  EVENT_TYPES: { ERROR: 'error', SUCCESS: 'success' },
}));

vi.mock('@sistent/sistent', () => ({
  FormControl: ({ children, ...rest }: any) => <div {...rest}>{children}</div>,
  CircularProgress: ({ size }: any) => <div data-testid="loading" data-size={size} />,
  RadioGroup: ({ value, onChange, children, ...rest }: any) => (
    <div data-testid="radio-group" data-value={value} {...rest}>
      {React.Children.map(children, (child) =>
        React.cloneElement(child, {
          onClick: () => onChange?.({ target: { value: child.props.value } }),
        }),
      )}
    </div>
  ),
  FormControlLabel: ({ label, value, control, onClick, ...rest }: any) => (
    <label data-radio-value={value} onClick={onClick}>
      {control}
      <span>{label}</span>
    </label>
  ),
  TextField: ({ id, label, value, onChange, type }: any) => (
    <label>
      {label}
      <input
        data-testid={`textfield-${id}`}
        id={id}
        type={type}
        value={value || ''}
        onChange={onChange}
      />
    </label>
  ),
  Grid2: ({ children }: any) => <div>{children}</div>,
  Button: ({ children, onClick, disabled, ...rest }: any) => (
    <button onClick={onClick} disabled={disabled} {...rest}>
      {children}
    </button>
  ),
  CustomTooltip: ({ children }: any) => <>{children}</>,
  useTheme: () => ({
    palette: { mode: 'light', primary: { main: '#06f' } },
    spacing: (n: number) => n,
  }),
  styled: (Component: any) => () => {
    const StyledComponent = ({ children, ...props }: any) =>
      React.createElement(Component, props, children);
    return StyledComponent;
  },
  Autocomplete: ({ value, inputValue, onChange, onInputChange, renderInput }: any) => (
    <div data-testid="autocomplete">
      <input
        data-testid="duration-input"
        value={inputValue || ''}
        onChange={(e) => onInputChange?.(e, e.target.value)}
      />
      <button onClick={() => onChange?.({}, value)}>autocomplete-change</button>
      {renderInput && renderInput({})}
    </div>
  ),
  NoSsr: ({ children }: any) => <>{children}</>,
  Radio: () => <input type="radio" />,
  Box: ({ children, sx, ...rest }: any) => (
    <div data-sx={JSON.stringify(sx || {})} {...rest}>
      {children}
    </div>
  ),
  SaveIcon: () => <svg data-testid="save-icon" />,
}));

vi.mock('@/rtk-query/user', () => ({
  useGetLoadTestPrefsQuery: () => ({ data: undefined }),
  useUpdateLoadTestPrefsMutation: () => [
    (...args: any[]) => updateLoadTestPrefs(...args),
    { isLoading: isSaving },
  ],
}));

vi.mock('react-redux', () => ({
  useSelector: (fn: any) => fn(selectedStore),
  useDispatch: () => dispatch,
}));

vi.mock('@/store/slices/mesheryUi', () => ({
  updateProgress: vi.fn(),
}));

vi.mock('@/store/slices/prefTest', () => ({
  updateLoadTestPref: vi.fn((arg: any) => arg),
}));

vi.mock('../../lib/load-test-prefs', () => ({
  normalizeLoadTestPrefs: (p: any) => p,
}));

import MesherySettingsPerformanceComponent from './MesherySettingsPerformanceComponent';

describe('MesherySettingsPerformanceComponent', () => {
  beforeEach(() => {
    notify.mockReset();
    dispatch.mockReset();
    updateLoadTestPrefs.mockClear();
    isSaving = false;
    selectedStore = {
      prefTest: {
        loadTestPref: { qps: 10, c: 5, t: '30s', gen: 'fortio' },
      },
      ui: { selectedK8sContexts: [] },
    };
  });

  it('renders the form with initial values from redux state', () => {
    render(<MesherySettingsPerformanceComponent />);

    expect(screen.getByTestId('textfield-c')).toHaveValue(5);
    expect(screen.getByTestId('textfield-qps')).toHaveValue(10);
  });

  it('updates c (concurrent requests) when input changes', () => {
    render(<MesherySettingsPerformanceComponent />);
    const cInput = screen.getByTestId('textfield-c') as HTMLInputElement;
    fireEvent.change(cInput, { target: { value: '20' } });
    expect(cInput).toHaveValue(20);
  });

  it('updates qps when input changes', () => {
    render(<MesherySettingsPerformanceComponent />);
    const qpsInput = screen.getByTestId('textfield-qps') as HTMLInputElement;
    fireEvent.change(qpsInput, { target: { value: '50' } });
    expect(qpsInput).toHaveValue(50);
  });

  it('renders only fortio as a load generator radio option', () => {
    render(<MesherySettingsPerformanceComponent />);
    expect(screen.getByText('fortio')).toBeInTheDocument();
    expect(screen.queryByText('wrk2')).not.toBeInTheDocument();
  });

  it('calls updateLoadTestPrefs with the form values when Save is clicked', async () => {
    const user = userEvent.setup();
    render(<MesherySettingsPerformanceComponent />);

    const saveButton = screen.getByRole('button', { name: /Save/ });
    await user.click(saveButton);

    expect(updateLoadTestPrefs).toHaveBeenCalledWith(
      expect.objectContaining({
        loadTestPrefs: { qps: 10, c: 5, t: '30s', gen: 'fortio' },
      }),
    );
  });

  it('does not submit when t is invalid (no unit suffix)', async () => {
    render(<MesherySettingsPerformanceComponent />);
    const durationInput = screen.getByTestId('duration-input') as HTMLInputElement;

    fireEvent.change(durationInput, { target: { value: 'invalid' } });
    fireEvent.click(screen.getByRole('button', { name: /Save/ }));

    expect(updateLoadTestPrefs).not.toHaveBeenCalled();
  });
});
