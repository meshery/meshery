import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import '@testing-library/jest-dom';
import MesheryVersionCompatibilityNotice from './MesheryVersionCompatibilityNotice';

// Mock clipboard API
Object.assign(navigator, {
  clipboard: {
    writeText: jest.fn(),
  },
});

describe('MesheryVersionCompatibilityNotice', () => {
  it('renders title and default version details', () => {
    render(<MesheryVersionCompatibilityNotice />);

    expect(screen.getByText(/Version Compatibility Notice/i)).toBeInTheDocument();
    expect(screen.getByText(/Current: v0.7.0/i)).toBeInTheDocument();
    expect(screen.getByText(/Required: v0.7.1\+/i)).toBeInTheDocument();
    expect(screen.getByText(/mesheryctl system update/i)).toBeInTheDocument();
  });

  it('renders custom prop values correctly', () => {
    render(
      <MesheryVersionCompatibilityNotice
        currentVersion="v0.6.9"
        requiredVersion="v0.7.0"
        componentName="Meshery Adapter"
        upgradeCommand="mesheryctl system restart"
      />
    );

    expect(screen.getByText(/Meshery Adapter/i)).toBeInTheDocument();
    expect(screen.getByText(/Current: v0.6.9/i)).toBeInTheDocument();
    expect(screen.getByText(/Required: v0.7.0/i)).toBeInTheDocument();
    expect(screen.getByText(/mesheryctl system restart/i)).toBeInTheDocument();
  });

  it('copies upgrade command when copy button is clicked', () => {
    render(<MesheryVersionCompatibilityNotice upgradeCommand="mesheryctl system update" />);

    const copyBtn = screen.getByRole('button', { name: /copy upgrade command/i });
    fireEvent.click(copyBtn);

    expect(navigator.clipboard.writeText).toHaveBeenCalledWith('mesheryctl system update');
  });
});
