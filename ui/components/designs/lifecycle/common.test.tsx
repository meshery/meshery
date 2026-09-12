import React from 'react';
import { fireEvent, render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const normalizeStaticImagePath = vi.fn((src?: string) => src);
const processDesignMock = vi.fn();

vi.mock('@/utils/fallback', () => ({
  normalizeStaticImagePath: (...args: any[]) => normalizeStaticImagePath(...args),
}));

vi.mock('@/utils/utils', () => ({
  processDesign: (...args: any[]) => processDesignMock(...args),
}));

vi.mock('@/assets/icons/Pattern', () => ({
  default: () => <svg data-testid="pattern-icon" />,
}));

vi.mock('@sistent/sistent', () => {
  const styled = (_tag: any) => () => {
    const Styled = ({ children, ...rest }: any) => <div {...rest}>{children}</div>;
    Styled.displayName = 'StyledMock';
    return Styled;
  };
  return {
    styled,
    Box: ({ children, 'data-testid': testId, ...rest }: any) => (
      <div data-testid={testId} {...rest}>
        {children}
      </div>
    ),
    Checkbox: ({ value, onChange, disabled, slotProps }: any) => (
      <input
        type="checkbox"
        checked={!!value}
        onChange={onChange}
        disabled={disabled}
        data-testid={slotProps?.input?.['data-testid']}
        readOnly
      />
    ),
    CircularProgress: ({ 'data-testid': testId }: any) => (
      <span data-testid={testId || 'spinner'} />
    ),
    Stack: ({ children, 'data-testid': testId, ...rest }: any) => (
      <div data-testid={testId} {...rest}>
        {children}
      </div>
    ),
    Typography: ({ children, 'data-testid': testId }: any) => (
      <span data-testid={testId}>{children}</span>
    ),
    CustomTooltip: ({ children, title }: any) => (
      <div data-testid="tooltip" data-title={String(title)}>
        {children}
      </div>
    ),
    InfoCircleIcon: () => <svg data-testid="info-icon" />,
    useTheme: () => ({
      palette: {
        text: {
          disabled: '#aaa',
          neutral: { default: '#000' },
        },
      },
    }),
    IconButton: ({ children, 'data-testid': testId, ...rest }: any) => (
      <button type="button" data-testid={testId} {...rest}>
        {children}
      </button>
    ),
  };
});

import {
  CheckBoxField,
  ComponentIcon,
  DEPLOYMENT_TYPE,
  Loading,
  StepHeading,
  getSvgWhiteForComponent,
} from './common';

describe('DEPLOYMENT_TYPE', () => {
  it('exposes canonical deploy/undeploy strings', () => {
    expect(DEPLOYMENT_TYPE).toEqual({ DEPLOY: 'deploy', UNDEPLOY: 'undeploy' });
  });
});

describe('ComponentIcon', () => {
  beforeEach(() => normalizeStaticImagePath.mockClear());

  it('returns a default PatternIcon when no iconSrc is supplied', () => {
    normalizeStaticImagePath.mockReturnValueOnce(null);
    render(<ComponentIcon iconSrc={null} label="lbl" />);
    expect(screen.getByTestId('pattern-icon')).toBeInTheDocument();
  });

  it('renders an image when icon source is valid', () => {
    normalizeStaticImagePath.mockReturnValueOnce('/img/foo.svg');
    const { container } = render(<ComponentIcon iconSrc="/raw/foo.svg" label="lbl" />);
    const img = container.querySelector('img');
    expect(img).not.toBeNull();
    expect(img).toHaveAttribute('src', '/img/foo.svg');
  });

  it('falls back to label text when image fails to load', () => {
    normalizeStaticImagePath.mockReturnValueOnce('/img/foo.svg');
    const { container } = render(<ComponentIcon iconSrc="/raw/foo.svg" label="ModelName" />);
    const img = container.querySelector('img') as HTMLImageElement;
    fireEvent.error(img);
    expect(screen.getByTitle('ModelName')).toBeInTheDocument();
  });
});

describe('Loading', () => {
  it('renders a spinner and a message', () => {
    render(<Loading message="working..." />);
    expect(screen.getByTestId('loading')).toBeInTheDocument();
    expect(screen.getByTestId('loading-spinner')).toBeInTheDocument();
    expect(screen.getByTestId('loading-message')).toHaveTextContent('working...');
  });

  it('honors a custom testId', () => {
    render(<Loading message="x" data-testid="custom" />);
    expect(screen.getByTestId('custom')).toBeInTheDocument();
    expect(screen.getByTestId('custom-spinner')).toBeInTheDocument();
    expect(screen.getByTestId('custom-message')).toHaveTextContent('x');
  });
});

describe('CheckBoxField', () => {
  it('renders label and checkbox, propagates onChange', () => {
    const onChange = vi.fn();
    render(<CheckBoxField label="agree" checked={true} onChange={onChange} />);
    expect(screen.getByText('agree')).toBeInTheDocument();
    const checkbox = screen.getByTestId('checkbox-field-checkbox') as HTMLInputElement;
    fireEvent.click(checkbox);
    expect(onChange).toHaveBeenCalled();
  });

  it('renders helpText tooltip when supplied', () => {
    render(<CheckBoxField label="opt" checked={false} onChange={vi.fn()} helpText="more info" />);
    expect(screen.getByTestId('checkbox-field-help-icon')).toBeInTheDocument();
  });

  it('disables the checkbox when disabled is true', () => {
    render(<CheckBoxField label="opt" checked={false} onChange={vi.fn()} disabled />);
    expect(screen.getByTestId('checkbox-field-checkbox')).toBeDisabled();
  });
});

describe('StepHeading', () => {
  it('renders its children inside a Typography', () => {
    render(<StepHeading>title</StepHeading>);
    expect(screen.getByText(/title/)).toBeInTheDocument();
  });
});

describe('getSvgWhiteForComponent', () => {
  it('normalizes the svgWhite field of a component style', () => {
    normalizeStaticImagePath.mockReturnValueOnce('NORMALIZED');
    const result = getSvgWhiteForComponent({ styles: { svgWhite: 'raw.svg' } });
    expect(normalizeStaticImagePath).toHaveBeenCalledWith('raw.svg');
    expect(result).toBe('NORMALIZED');
  });

  it('handles missing components gracefully', () => {
    normalizeStaticImagePath.mockReturnValueOnce(undefined);
    expect(getSvgWhiteForComponent(undefined)).toBeUndefined();
  });
});
