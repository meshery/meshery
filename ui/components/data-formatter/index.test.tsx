import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  CustomTooltip: ({ children, title }: any) => (
    <div data-testid="tooltip" data-title={String(title)}>
      {children}
    </div>
  ),
  Typography: ({ children, style, variant, ...rest }: any) => (
    <span data-testid="typography" data-variant={variant} style={style} {...rest}>
      {children}
    </span>
  ),
  Box: ({ children, sx }: any) => (
    <div data-testid="box" data-sx={JSON.stringify(sx || {})}>
      {children}
    </div>
  ),
  IconButton: ({ children, onClick, ...rest }: any) => (
    <button data-testid="icon-button" onClick={onClick} {...rest}>
      {children}
    </button>
  ),
  LaunchIcon: () => <svg data-testid="launch-icon" />,
  useTheme: () => ({
    palette: {
      text: { tertiary: 'tertiary', default: 'default' },
    },
  }),
  Grid: ({ children, ...rest }: any) => (
    <div data-testid="grid" {...rest}>
      {children}
    </div>
  ),
}));

vi.mock('../../assets/icons/CopyIcon', () => ({
  default: ({ width, height }: any) => (
    <svg data-testid="copy-icon" data-width={width} data-height={height} />
  ),
}));

vi.mock('../../utils/objects', () => ({
  isEmptyAtAllDepths: (val: any) =>
    val == null || (typeof val === 'object' && Object.keys(val).length === 0),
}));

import {
  formatDate,
  formatTime,
  formatDateTime,
  FormattedDate,
  FormatId,
  createColumnUiSchema,
  Link,
  LinkFormatters,
  TextWithLinks,
  KeyValue,
  SectionHeading,
  SectionBody,
  ArrayFormatter,
  reorderObjectProperties,
  FormatStructuredData,
} from './index';

describe('formatDate/formatTime/formatDateTime', () => {
  it('formats a date in en-US Month Day, Year', () => {
    const d = new Date('2024-03-15T14:00:00Z');
    expect(formatDate(d)).toMatch(/Mar\s+\d+,\s+2024/);
  });

  it('formats a time as h:m:s', () => {
    const d = new Date('2024-03-15T14:00:00Z');
    // Locale-specific formatting — just assert it contains digits and a separator
    expect(formatTime(d)).toMatch(/\d+:\d+:\d+/);
  });

  it('combines date and time', () => {
    const d = new Date('2024-03-15T14:00:00Z');
    const both = formatDateTime(d);
    expect(both).toContain('2024');
    expect(both).toMatch(/\d+:\d+:\d+/);
  });

  it('renders a placeholder instead of "Invalid Date" for missing or unparsable input', () => {
    for (const bad of [undefined, null, '', 'not-a-date']) {
      expect(formatDate(bad)).toBe('-');
      expect(formatTime(bad)).toBe('-');
      expect(formatDateTime(bad)).toBe('-');
    }
  });
});

describe('FormattedDate', () => {
  it('renders a tooltip containing the full datetime', () => {
    const d = new Date('2024-03-15T14:00:00Z');
    render(<FormattedDate date={d} />);
    const tip = screen.getByTestId('tooltip');
    expect(tip.getAttribute('data-title')).toContain('2024');
  });
});

describe('FormatId', () => {
  let writeTextSpy: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    writeTextSpy = vi.fn().mockResolvedValue(undefined);
    Object.defineProperty(window.navigator, 'clipboard', {
      configurable: true,
      writable: true,
      value: { writeText: writeTextSpy },
    });
  });

  it('truncates long ids to ~15 chars and shows full id in tooltip', () => {
    render(<FormatId id="abcdefghijklmnopqrstuvwxyz" />);
    expect(screen.getAllByTestId('tooltip')[0].getAttribute('data-title')).toBe(
      'abcdefghijklmnopqrstuvwxyz',
    );
    expect(screen.getByText(/abcdefghij/)).toBeInTheDocument();
  });

  it('writes id to clipboard when the copy button is clicked', () => {
    render(<FormatId id="abc-xyz" />);

    const button = screen.getByTestId('icon-button');
    // fireEvent avoids userEvent which may re-pierce navigator
    fireEvent.click(button);

    expect(writeTextSpy).toHaveBeenCalledWith('abc-xyz');
  });
});

describe('createColumnUiSchema', () => {
  it('returns column widths floor(12/numCols) for each metadata key', () => {
    const schema = createColumnUiSchema({
      metadata: { a: {}, b: {} },
      numCols: { lg: 3, md: 2 },
    });

    expect(schema.a).toEqual({ lg: 4, md: 6 });
    expect(schema.b).toEqual({ lg: 4, md: 6 });
  });

  it('returns an empty mapping per key when no numCols entries', () => {
    const schema = createColumnUiSchema({ metadata: { a: {} }, numCols: {} });
    expect(schema.a).toEqual({});
  });
});

describe('Link / LinkFormatters', () => {
  it('renders an anchor with the supplied href and title', () => {
    render(<Link href="https://example.com" title="example" />);
    const a = screen.getByRole('link', { name: /example/ });
    expect(a).toHaveAttribute('href', 'https://example.com');
    expect(a).toHaveAttribute('target', '_blank');
    expect(a).toHaveAttribute('rel', 'noopener noreferrer');
  });

  it('LinkFormatters.DOC produces a "Doc" titled link', () => {
    const { container } = render(<>{LinkFormatters.DOC.formatter('https://docs.meshery.io/x')}</>);
    expect(container.textContent).toContain('Doc');
  });

  it('LinkFormatters.DEFAULT produces a truncated link title', () => {
    const { container } = render(
      <>{LinkFormatters.DEFAULT.formatter('https://example.com/foo')}</>,
    );
    expect(container.textContent).toContain('example.com');
  });
});

describe('TextWithLinks', () => {
  it('renders plain text when there are no links', () => {
    render(<TextWithLinks text="hello world" />);
    expect(screen.getByText(/hello world/)).toBeInTheDocument();
  });

  it('converts an http URL into a Link element', () => {
    const { container } = render(<TextWithLinks text="see https://docs.meshery.io/x for more" />);
    expect(container.textContent).toContain('Doc');
  });

  it('handles undefined / non-string input without crashing', () => {
    const { container } = render(<TextWithLinks text={undefined as any} />);
    expect(container).toBeInTheDocument();
  });
});

describe('KeyValue', () => {
  it('renders the key with underscores replaced by spaces', () => {
    render(<KeyValue Key="some_key_here" Value="123" />);
    expect(screen.getByText(/some key here/)).toBeInTheDocument();
  });

  it('renders the value as-is when it is a React element', () => {
    render(<KeyValue Key="x" Value={<span data-testid="custom-value">CUSTOM</span>} />);
    expect(screen.getByTestId('custom-value')).toBeInTheDocument();
  });
});

describe('SectionHeading / SectionBody', () => {
  it('renders the heading with its children', () => {
    render(<SectionHeading>My Heading</SectionHeading>);
    expect(screen.getByText('My Heading')).toBeInTheDocument();
  });

  it('SectionBody renders its body text', () => {
    render(<SectionBody body="hello body" />);
    expect(screen.getByText(/hello body/)).toBeInTheDocument();
  });
});

describe('ArrayFormatter', () => {
  it('renders one list item per array element', () => {
    const { container } = render(<ArrayFormatter items={['one', 'two', 'three']} />);
    const lis = container.querySelectorAll('li');
    expect(lis).toHaveLength(3);
  });
});

describe('reorderObjectProperties', () => {
  it('returns an object whose top-level keys are reordered', () => {
    const result = reorderObjectProperties({ b: 2, a: 1, c: 3 }, ['a', 'c']);
    const keys = Object.keys(result);
    expect(keys.slice(0, 2)).toEqual(['a', 'c']);
    expect(keys).toContain('b');
  });

  it('returns the input unchanged for non-objects', () => {
    expect(reorderObjectProperties('hello' as any, [])).toBe('hello');
    expect(reorderObjectProperties(null as any, [])).toBeNull();
  });
});

describe('FormatStructuredData', () => {
  it('returns null for empty data', () => {
    const { container } = render(<FormatStructuredData data={{}} />);
    expect(container.firstChild).toBeNull();
  });

  it('returns null for null data', () => {
    const { container } = render(<FormatStructuredData data={null as any} />);
    expect(container.firstChild).toBeNull();
  });

  it('renders structured object data with keys as headings', () => {
    const { container } = render(
      <FormatStructuredData data={{ name: 'kubernetes', version: '1.0' }} />,
    );
    expect(container.textContent).toContain('name');
    expect(container.textContent).toContain('kubernetes');
  });

  it('uses propertyFormatters when provided', () => {
    const propertyFormatters = {
      custom_key: (v: any) => <div data-testid="custom-formatter">CUSTOM:{v}</div>,
    };

    render(
      <FormatStructuredData
        data={{ custom_key: 'X', other: 'y' }}
        propertyFormatters={propertyFormatters}
      />,
    );

    expect(screen.getByTestId('custom-formatter')).toHaveTextContent('CUSTOM:X');
  });

  it('renders an array as an ordered list', () => {
    const { container } = render(<FormatStructuredData data={['a', 'b', 'c']} />);
    expect(container.querySelectorAll('li')).toHaveLength(3);
  });
});
