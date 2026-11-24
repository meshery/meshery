import React from 'react';
import { fireEvent, render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const processDesignMock = vi.fn(() => ({
  configurableComponents: [{ id: 'c1' }, { id: 'c2' }],
  annotationComponents: [{ id: 'a1' }],
}));
const useDesignSchemaValidationResultsMock = vi.fn();
const useIsValidatingDesignSchemaMock = vi.fn(() => false);
const tapOnErrorMock = vi.fn((data: any) => ({ type: 'TAP_ON_ERROR', data }));

vi.mock('../../../machines/validator/designValidator', () => ({
  designValidatorCommands: { validateDesignSchema: vi.fn(() => ({ type: 'VALIDATE' })) },
  designValidatorEvents: { tapOnError: (data: any) => tapOnErrorMock(data) },
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
  ValidationErrorListItem: ({ children, clickable, disablePadding: _dp, ...rest }: any) => (
    <li data-clickable={String(clickable)} {...rest}>
      {children}
    </li>
  ),
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
        canNavigateToErrors
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
        canNavigateToErrors
      />,
    );

    expect(screen.getByTestId('validated-component')).toBeInTheDocument();
    expect(screen.getByText('Service A')).toBeInTheDocument();
  });
});

const serviceA = { id: 'c1', displayName: 'Service A' };
const missingName = {
  instancePath: '/spec',
  keyword: 'required',
  message: "must have required property 'name'",
  params: { missingProperty: 'name' },
};
const missingPort = {
  instancePath: '/spec',
  keyword: 'required',
  message: "must have required property 'port'",
  params: { missingProperty: 'port' },
};
const resultsWithErrors = { c1: { component: serviceA, errors: [missingName, missingPort] } };

const tapsSent = (send: ReturnType<typeof vi.fn>) =>
  send.mock.calls.map(([event]) => event).filter((event) => event?.type === 'TAP_ON_ERROR');

describe('ValidateDesign error deep links', () => {
  beforeEach(() => {
    useDesignSchemaValidationResultsMock.mockReturnValue(resultsWithErrors);
  });

  // Regression: links were gated on the schema validator's returnAddress, which
  // ValidateDesign's own validation request clears, so Kanvas Designer lost them.
  it('links each error to its component and property by default, as Kanvas Designer renders it', () => {
    const send = vi.fn();
    render(
      <ValidateDesign design={{ name: 'D1' }} currentNodeId="c1" validationMachine={{ send }} />,
    );

    const errors = screen.getAllByTestId('validation-error');
    expect(errors).toHaveLength(2);
    errors.forEach((error) => {
      expect(error).toHaveAttribute('role', 'button');
      expect(error).toHaveAttribute('tabindex', '0');
      expect(error).toHaveAttribute('data-clickable', 'true');
    });

    fireEvent.click(errors[1]);

    expect(tapsSent(send)).toEqual([
      {
        type: 'TAP_ON_ERROR',
        data: { error: missingPort, type: 'schemaValidation', component: serviceA },
      },
    ]);
  });

  it('follows the link from the keyboard', () => {
    const send = vi.fn();
    render(
      <ValidateDesign design={{ name: 'D1' }} currentNodeId="c1" validationMachine={{ send }} />,
    );

    const [first] = screen.getAllByTestId('validation-error');
    fireEvent.keyDown(first, { key: 'Enter' });
    fireEvent.keyDown(first, { key: ' ' });
    fireEvent.keyDown(first, { key: 'a' });

    expect(tapsSent(send)).toHaveLength(2);
    expect(tapsSent(send)[0].data.error).toBe(missingName);
  });

  it('lists errors as plain text where there is no designer to navigate into', () => {
    const send = vi.fn();
    render(
      <ValidateDesign
        design={{ name: 'D1' }}
        currentNodeId="c1"
        validationMachine={{ send }}
        canNavigateToErrors={false}
      />,
    );

    const errors = screen.getAllByTestId('validation-error');
    errors.forEach((error) => {
      expect(error).not.toHaveAttribute('role');
      expect(error).not.toHaveAttribute('tabindex');
      expect(error).toHaveAttribute('data-clickable', 'false');
    });
    fireEvent.click(errors[0]);
    fireEvent.keyDown(errors[0], { key: 'Enter' });

    expect(tapsSent(send)).toEqual([]);
  });

  it('defaults ValidationResults to plain text, as the notification center renders stored events', () => {
    const send = vi.fn();
    render(
      <ValidationResults
        errorCount={2}
        compCount={1}
        annotationCount={0}
        validationResults={resultsWithErrors}
        currentNodeId="c1"
        validationMachine={{ send }}
      />,
    );

    const errors = screen.getAllByTestId('validation-error');
    expect(errors[0]).not.toHaveAttribute('role');
    fireEvent.click(errors[0]);
    expect(send).not.toHaveBeenCalled();
  });
});
