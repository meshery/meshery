import React from 'react';
import { render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const processDesignMock = vi.fn(() => ({
  configurableComponents: [{ id: 'c1' }, { id: 'c2' }],
  annotationComponents: [{ id: 'a1' }],
}));
const useDesignSchemaValidationResultsMock = vi.fn();
const useIsValidatingDesignSchemaMock = vi.fn(() => false);

vi.mock('@xstate/react', () => ({
  useSelector: () => false,
}));

vi.mock('../../../machines/validator/designValidator', () => ({
  designValidatorCommands: { validateDesignSchema: vi.fn() },
  designValidatorEvents: { tapOnError: vi.fn() },
  selectValidator: vi.fn(),
  useDesignSchemaValidationResults: (...args: any[]) =>
    useDesignSchemaValidationResultsMock(...args),
  useIsValidatingDesignSchema: (...args: any[]) => useIsValidatingDesignSchemaMock(...args),
}));

vi.mock('./common', () => ({
  ComponentIcon: ({ iconSrc }: any) => <span data-testid="component-icon" data-src={iconSrc} />,
  Loading: ({ message }: any) => <div data-testid="loading">{message}</div>,
  getSvgWhiteForComponent: () => 'icon.svg',
  processDesign: (...args: any[]) => processDesignMock(...args),
}));

vi.mock('@sistent/sistent', () => {
  const styled = (Component: any) => () => {
    if (typeof Component === 'string') {
      const Styled = ({ children, ...rest }: any) => <div {...rest}>{children}</div>;
      Styled.displayName = 'StyledHostMock';
      return Styled;
    }
    const Styled = ({ children, ...rest }: any) => <Component {...rest}>{children}</Component>;
    Styled.displayName = 'StyledMock';
    return Styled;
  };
  return {
    styled,
    List: ({ children, subheader }: any) => (
      <div>
        {subheader}
        <ul>{children}</ul>
      </div>
    ),
    ListItemText: ({ primary }: any) => <span>{primary}</span>,
    ListItemIcon: ({ children }: any) => <span>{children}</span>,
    Typography: ({ children, ...rest }: any) => <span {...rest}>{children}</span>,
    Collapse: ({ children, in: isOpen }: any) => (isOpen ? <div>{children}</div> : null),
    useTheme: () => ({
      palette: {
        text: { default: '#000' },
        background: { cta: { default: '#abc' } },
        warning: { main: '#ffaa00' },
        success: { main: '#00ff00' },
      },
    }),
  };
});

vi.mock('@/assets/icons', () => ({
  ExpandLess: () => <svg data-testid="expand-less" />,
  ExpandMore: () => <svg data-testid="expand-more" />,
}));

vi.mock('@/assets/icons/AlertIcon', () => ({
  default: () => <svg data-testid="alert-icon" />,
}));

vi.mock('./styles', () => ({
  ValidatedComponent: ({ children }: any) => (
    <div data-testid="validated-component">{children}</div>
  ),
  ValidationErrorListItem: ({ children, onClick }: any) => <li onClick={onClick}>{children}</li>,
  ComponentValidationListItem: ({ children, onClick }: any) => (
    <div data-testid="component-val-row" onClick={onClick}>
      {children}
    </div>
  ),
  ValidationResultsListWrapper: ({ children, subheader }: any) => (
    <div>
      {subheader}
      <div>{children}</div>
    </div>
  ),
  ValidationSubHeader: ({ children }: any) => <div>{children}</div>,
}));

import { ValidateDesign, ValidationResults } from './ValidateDesign';

beforeEach(() => {
  vi.clearAllMocks();
  useIsValidatingDesignSchemaMock.mockReturnValue(false);
});

describe('ValidateDesign', () => {
  it('shows the loading state when the validator is running', () => {
    useIsValidatingDesignSchemaMock.mockReturnValue(true);
    useDesignSchemaValidationResultsMock.mockReturnValue(null);

    render(
      <ValidateDesign
        design={{ name: 'D1' }}
        currentNodeId="n1"
        validationMachine={{ send: vi.fn() }}
      />,
    );
    expect(screen.getByTestId('loading')).toHaveTextContent('Validating Design');
  });

  it('renders results when validation completes', () => {
    useIsValidatingDesignSchemaMock.mockReturnValue(false);
    useDesignSchemaValidationResultsMock.mockReturnValue({});

    render(
      <ValidateDesign
        design={{ name: 'D1' }}
        currentNodeId="n1"
        validationMachine={{ send: vi.fn() }}
      />,
    );
    expect(screen.queryByTestId('loading')).not.toBeInTheDocument();
  });

  it('renders the typeof-string validation message as a typography', () => {
    useIsValidatingDesignSchemaMock.mockReturnValue(false);
    useDesignSchemaValidationResultsMock.mockReturnValue('All good');

    render(
      <ValidateDesign
        design={{ name: 'D1' }}
        currentNodeId={null}
        validationMachine={{ send: vi.fn() }}
      />,
    );
    expect(screen.getByText('All good')).toBeInTheDocument();
  });
});

describe('ValidationResults', () => {
  it('renders the success message when there are zero errors', () => {
    render(
      <ValidationResults
        errorCount={0}
        compCount={3}
        annotationCount={2}
        validationResults={{}}
        currentNodeId={null}
        validationMachine={{ send: vi.fn() }}
        canTapErrors
      />,
    );

    expect(screen.getByText('No validation errors.')).toBeInTheDocument();
    expect(screen.getByText('3 components and 2 annotations')).toBeInTheDocument();
  });

  it('renders an entry per component with errors', () => {
    render(
      <ValidationResults
        errorCount={1}
        compCount={1}
        annotationCount={0}
        validationResults={{
          svc1: {
            component: { id: 'c1', displayName: 'Service A' },
            errors: [{ instancePath: '/spec/name', message: 'is required' }],
          },
        }}
        currentNodeId={null}
        validationMachine={{ send: vi.fn() }}
        canTapErrors
      />,
    );

    expect(screen.getByTestId('validated-component')).toBeInTheDocument();
    expect(screen.getByText('Service A')).toBeInTheDocument();
  });
});
