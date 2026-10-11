import React from 'react';
import { render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { alpha, darkModePalette, lightModePalette } from '@sistent/sistent';
import ReactSelectWrapper from './ReactSelectWrapper';
import MultiSelectWrapper, { getMultiSelectOptionHighlight } from './multi-select-wrapper';

let capturedSelectProps: any[] = [];
let themeMode: 'light' | 'dark' = 'light';

const darkCardToken = darkModePalette.background.card;

vi.mock('@sistent/sistent', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@sistent/sistent')>();
  return {
    ...actual,
    useTheme: () => {
      const isDark = themeMode === 'dark';
      return {
        palette: {
          mode: themeMode,
          background: {
            card: darkCardToken,
            paper: '#fff',
            elevatedComponents: darkModePalette.background.elevatedComponents,
            secondary: lightModePalette.background.secondary,
          },
          common: { white: '#ffffff' },
          primary: { main: '#00B39F' },
          text: {
            default: isDark ? darkModePalette.text.default : lightModePalette.text.default,
            primary: isDark ? darkModePalette.text.default : lightModePalette.text.default,
            disabled: '#999',
          },
          error: { main: '#B32700' },
          border: { strong: '#666' },
          action: { selected: 'rgba(0,0,0,0.08)', hover: 'rgba(0,0,0,0.04)', active: '#999' },
          getContrastText: (bg: string) => actual.createTheme({}).palette.getContrastText(bg),
        },
        typography: { body2: { fontSize: '0.875rem', lineHeight: 1.43 } },
      };
    },
  };
});

vi.mock('react-select/creatable', () => ({
  default: (props) => {
    capturedSelectProps.push(props);
    return <div data-testid="creatable-select" />;
  },
}));

vi.mock('react-select', () => ({
  components: {
    Input: ({ children, ...props }) => <div {...props}>{children}</div>,
  },
}));

const lastStyles = () => capturedSelectProps.at(-1).styles;

describe('select wrappers', () => {
  beforeEach(() => {
    capturedSelectProps = [];
    themeMode = 'light';
  });

  it('passes a custom noOptionsMessage to CreatableSelect', () => {
    const customMessage = () => 'Custom empty state';
    render(
      <MultiSelectWrapper
        onChange={vi.fn()}
        options={[]}
        value={[]}
        noOptionsMessage={customMessage}
      />,
    );
    expect(capturedSelectProps.at(-1).noOptionsMessage).toBe(customMessage);
  });

  it('uses the theme card background for the menu in dark mode', () => {
    themeMode = 'dark';
    render(<MultiSelectWrapper onChange={vi.fn()} options={[]} value={[]} />);
    expect(lastStyles().menu({}).backgroundColor).toBe(darkCardToken);
  });

  it('keeps the closed control transparent so it stays attached to the row', () => {
    render(<MultiSelectWrapper onChange={vi.fn()} options={[]} value={[]} />);
    expect(lastStyles().control({}).backgroundColor).toBe('transparent');
  });

  it('hides the placeholder when focused or when input has a value', () => {
    render(<MultiSelectWrapper onChange={vi.fn()} options={[]} value={[]} />);
    const { placeholder } = lastStyles();
    const base = { display: 'block' };
    expect(placeholder(base, { isFocused: true, selectProps: { inputValue: '' } }).display).toBe(
      'none',
    );
    expect(
      placeholder(base, { isFocused: false, selectProps: { inputValue: 'abc' } }).display,
    ).toBe('none');
    expect(placeholder(base, { isFocused: false, selectProps: { inputValue: '' } }).display).toBe(
      'block',
    );
  });

  it('uses a primary-tinted option highlight (not action.hover)', () => {
    expect(
      getMultiSelectOptionHighlight({ palette: { mode: 'dark', primary: { main: '#00B39F' } } }),
    ).toBe(alpha('#00B39F', 0.28));
    expect(
      getMultiSelectOptionHighlight({ palette: { mode: 'light', primary: { main: '#00B39F' } } }),
    ).toBe(alpha('#00B39F', 0.12));
  });

  it('paints chip label and remove with the same opaque fill and contrast text', () => {
    const check = (mode: 'light' | 'dark') => {
      themeMode = mode;
      capturedSelectProps = [];
      render(<MultiSelectWrapper onChange={vi.fn()} options={[]} value={[]} />);
      const label = lastStyles().multiValueLabel({});
      const remove = lastStyles().multiValueRemove({});
      const expectedBg =
        mode === 'dark'
          ? darkModePalette.background.elevatedComponents
          : lightModePalette.background.secondary;
      expect(label.backgroundColor).toBe(expectedBg);
      expect(remove.backgroundColor).toBe(expectedBg);
      expect(label.color).toBe(remove.color);
      expect(label.color).not.toBe(expectedBg);
    };
    check('dark');
    check('light');
  });

  it('marks Mui-selected from isSelected only, not isFocused alone', () => {
    render(
      <MultiSelectWrapper
        onChange={vi.fn()}
        options={[
          { label: 'All', value: '*' },
          { label: 'Production', value: 'prod' },
        ]}
        value={[]}
      />,
    );
    const Option = capturedSelectProps.at(-1).components.Option;
    const option = (isFocused: boolean, isSelected: boolean) => (
      <Option
        innerRef={vi.fn()}
        innerProps={{ onClick: vi.fn(), role: 'option' }}
        isFocused={isFocused}
        isSelected={isSelected}
        label="Production"
        value="prod"
      />
    );

    const { container, rerender } = render(option(true, true));
    expect(container.querySelector('[data-testid="multi-select-option"]')).toHaveClass(
      'Mui-selected',
    );
    rerender(option(false, true));
    expect(container.querySelector('[data-testid="multi-select-option"]')).toHaveClass(
      'Mui-selected',
    );
    rerender(option(true, false));
    expect(container.querySelector('[data-testid="multi-select-option"]')).not.toHaveClass(
      'Mui-selected',
    );
  });

  it('renders ReactSelectWrapper option content without a menu context error', () => {
    render(
      <ReactSelectWrapper
        label="Environment"
        placeholder="Select environment"
        onChange={vi.fn()}
        options={[{ label: 'Production', value: 'prod' }]}
        value={null}
      />,
    );
    const Option = capturedSelectProps.at(-1).components.Option;
    render(
      <Option
        innerRef={vi.fn()}
        innerProps={{ onClick: vi.fn(), role: 'option' }}
        isFocused={false}
        isSelected={false}
      >
        Production
      </Option>,
    );
    expect(screen.getByText('Production')).toBeInTheDocument();
  });
});
