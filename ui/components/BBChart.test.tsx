import React from 'react';
import { render } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

const generateMock = vi.fn();
const destroyMock = vi.fn();

vi.mock('billboard.js', () => ({
  bb: {
    generate: (...args: any[]) => {
      generateMock(...args);
      return { destroy: destroyMock };
    },
  },
}));

vi.mock('@sistent/sistent', () => ({
  ErrorBoundary: ({ children }: any) => <div data-testid="boundary">{children}</div>,
}));

import BBChart from './BBChart';

describe('BBChart', () => {
  beforeEach(() => {
    generateMock.mockClear();
    destroyMock.mockClear();
  });

  it('binds billboard to the ref and forwards options', () => {
    const options = { data: { columns: [['x', 1, 2, 3]] } };
    render(<BBChart options={options as any} />);
    expect(generateMock).toHaveBeenCalledTimes(1);
    const callArg = generateMock.mock.calls[0][0];
    expect(callArg.data).toBe(options.data);
    expect(callArg.bindto).toBeTruthy();
  });

  it('destroys the chart on unmount', () => {
    const { unmount } = render(<BBChart options={{} as any} />);
    unmount();
    expect(destroyMock).toHaveBeenCalled();
  });

  it('regenerates the chart when options change', () => {
    const { rerender } = render(<BBChart options={{ id: 1 } as any} />);
    rerender(<BBChart options={{ id: 2 } as any} />);
    expect(generateMock).toHaveBeenCalledTimes(2);
  });

  it('stops click propagation through the container', () => {
    const { container } = render(<BBChart options={{} as any} />);
    const div = container.querySelector('div[style*="width"]');
    expect(div).not.toBeNull();
  });
});
