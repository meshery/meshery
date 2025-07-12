import { renderHook, act } from '@testing-library/react';
import { SnackbarProvider } from 'notistack';
import useConnectionLoader from '../useConnectionLoader';

const wrapper = ({ children }) => (
  <SnackbarProvider>
    {children}
  </SnackbarProvider>
);

describe('useConnectionLoader', () => {
  beforeEach(() => {
    // Mock the event bus
    window.mesheryEventBus = {
      subscribe: jest.fn(),
      unsubscribe: jest.fn(),
    };
  });

  afterEach(() => {
    delete window.mesheryEventBus;
  });

  it('initializes with default state', () => {
    const { result } = renderHook(() => useConnectionLoader(), { wrapper });
    
    expect(result.current.isVerifying).toBe(false);
    expect(result.current.verificationMessage).toBe('');
  });

  it('starts verification with custom message', () => {
    const { result } = renderHook(() => useConnectionLoader(), { wrapper });
    
    act(() => {
      result.current.startVerification('Custom verification message');
    });
    
    expect(result.current.isVerifying).toBe(true);
    expect(result.current.verificationMessage).toBe('Custom verification message');
  });

  it('stops verification', () => {
    const { result } = renderHook(() => useConnectionLoader(), { wrapper });
    
    // Start verification first
    act(() => {
      result.current.startVerification('Testing...');
    });
    
    expect(result.current.isVerifying).toBe(true);
    
    // Stop verification
    act(() => {
      result.current.stopVerification();
    });
    
    expect(result.current.isVerifying).toBe(false);
    expect(result.current.verificationMessage).toBe('');
  });

  it('subscribes to connection events', () => {
    renderHook(() => useConnectionLoader(), { wrapper });
    
    expect(window.mesheryEventBus.subscribe).toHaveBeenCalledWith('connection', expect.any(Function));
  });

  it('unsubscribes from events on cleanup', () => {
    const { unmount } = renderHook(() => useConnectionLoader(), { wrapper });
    
    unmount();
    
    expect(window.mesheryEventBus.unsubscribe).toHaveBeenCalledWith('connection', expect.any(Function));
  });
}); 