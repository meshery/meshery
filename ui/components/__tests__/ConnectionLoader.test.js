import React from 'react';
import { render, screen } from '@testing-library/react';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import ConnectionLoader from '../ConnectionLoader';

const theme = createTheme();

const renderWithTheme = (component) => {
  return render(
    <ThemeProvider theme={theme}>
      {component}
    </ThemeProvider>
  );
};

describe('ConnectionLoader', () => {
  it('renders with default message', () => {
    renderWithTheme(<ConnectionLoader />);
    
    expect(screen.getByText('Verifying connection...')).toBeInTheDocument();
    expect(screen.getByText('Please wait while we verify your connection...')).toBeInTheDocument();
  });

  it('renders with custom message', () => {
    const customMessage = 'Testing connection...';
    renderWithTheme(<ConnectionLoader message={customMessage} />);
    
    expect(screen.getByText(customMessage)).toBeInTheDocument();
    expect(screen.getByText('Please wait while we verify your connection...')).toBeInTheDocument();
  });

  it('displays circular progress indicator', () => {
    renderWithTheme(<ConnectionLoader />);
    
    // Check if the circular progress is rendered
    const progressElement = document.querySelector('.MuiCircularProgress-root');
    expect(progressElement).toBeInTheDocument();
  });

  it('has proper styling and layout', () => {
    renderWithTheme(<ConnectionLoader />);
    
    // Check if the container has proper styling
    const container = document.querySelector('.MuiPaper-root');
    expect(container).toBeInTheDocument();
  });
}); 