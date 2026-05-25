import React from 'react';
import { render } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

// Capture all enqueueSnackbar / closeSnackbar calls so we can assert on them.
const enqueueSnackbar = vi.fn(() => 'snack-key-1');
const closeSnackbar = vi.fn();
let mockShowProgress = false;

vi.mock('notistack', () => ({
  useSnackbar: () => ({ enqueueSnackbar, closeSnackbar }),
}));

vi.mock('react-redux', () => ({
  useSelector: (selector: any) => selector({ ui: { showProgress: mockShowProgress } }),
}));

vi.mock('@sistent/sistent', () => ({
  LinearProgress: () => <div data-testid="linear-progress" />,
}));

// Import after mocks are set up.
import MesheryProgressBar from './MesheryProgressBar';

describe('MesheryProgressBar', () => {
  beforeEach(() => {
    enqueueSnackbar.mockClear();
    closeSnackbar.mockClear();
    mockShowProgress = false;
  });

  it('renders null', () => {
    const { container } = render(<MesheryProgressBar />);
    expect(container.firstChild).toBeNull();
  });

  it('enqueues a snackbar when showProgress is true', () => {
    mockShowProgress = true;
    render(<MesheryProgressBar />);
    expect(enqueueSnackbar).toHaveBeenCalledTimes(1);
    const [, options] = enqueueSnackbar.mock.calls[0];
    expect(options).toEqual({ variant: 'default', persist: true });
  });

  it('does not enqueue when showProgress is false', () => {
    mockShowProgress = false;
    render(<MesheryProgressBar />);
    expect(enqueueSnackbar).not.toHaveBeenCalled();
    expect(closeSnackbar).not.toHaveBeenCalled();
  });

  it('closes the snackbar when showProgress transitions to false', () => {
    mockShowProgress = true;
    const { rerender } = render(<MesheryProgressBar />);
    expect(enqueueSnackbar).toHaveBeenCalledTimes(1);

    mockShowProgress = false;
    rerender(<MesheryProgressBar />);
    expect(closeSnackbar).toHaveBeenCalledWith('snack-key-1');
  });
});
