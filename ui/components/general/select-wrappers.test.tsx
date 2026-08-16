import React from 'react';
import { render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { alpha, darkModePalette, lightModePalette } from '@sistent/sistent';
import ReactSelectWrapper from './ReactSelectWrapper';
import MultiSelectWrapper, { getMultiSelectOptionHighlight } from './multi-select-wrapper';

let capturedSelectProps: any[] = [];
let themeMode: 'light' | 'dark' = 'light';

const darkCardToken = darkModePalette.background.card; // #212121

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
            secondary: isDark ? darkModePalette.text.secondary : lightModePalette.text.secondary,
            disabled: '#999',
          },
          error: { main: '#B32700' },
          // Resting multi-select outline uses Sistent border.strong when present.
          border: { strong: isDark ? '#666' : '#BDBDBD' },
          action: {
            selected: 'rgba(0,0,0,0.08)',
            hover: 'rgba(0,0,0,0.04)',
            active: 'rgba(0,0,0,0.54)',
          },
          getContrastText: (bg: string) => actual.createTheme({}).palette.getContrastText(bg),
        },
        typography: {
          body2: { fontSize: '0.875rem', lineHeight: 1.43 },
        },
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

/** Relative luminance + WCAG contrast for solid hex colors. */
function hexToRgb(hex: string) {
  const h = hex.replace('#', '');
  const full =
    h.length === 3
      ? h
          .split('')
          .map((c) => c + c)
          .join('')
      : h;
  return [0, 2, 4].map((i) => parseInt(full.slice(i, i + 2), 16));
}

function relativeLuminance(hex: string) {
  const [r, g, b] = hexToRgb(hex).map((v) => {
    const s = v / 255;
    return s <= 0.03928 ? s / 12.92 : ((s + 0.055) / 1.055) ** 2.4;
  });
  return 0.2126 * r + 0.7152 * g + 0.0722 * b;
}

function contrastRatio(a: string, b: string) {
  const L1 = relativeLuminance(a);
  const L2 = relativeLuminance(b);
  const [hi, lo] = L1 > L2 ? [L1, L2] : [L2, L1];
  return (hi + 0.05) / (lo + 0.05);
}

function blendRgbaOverHex(rgba: string, baseHex: string) {
  const m = rgba.match(/rgba?\(([^)]+)\)/);
  if (!m) throw new Error(`expected rgba color, got ${rgba}`);
  const [fr, fg, fb, fa = 1] = m[1].split(',').map((p) => parseFloat(p.trim()));
  const [br, bg, bb] = hexToRgb(baseHex);
  const r = Math.round(fr * fa + br * (1 - fa));
  const g = Math.round(fg * fa + bg * (1 - fa));
  const b = Math.round(fb * fa + bb * (1 - fa));
  return `#${[r, g, b].map((x) => x.toString(16).padStart(2, '0')).join('')}`;
}

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

    const menuStyles = capturedSelectProps.at(-1).styles.menu({});
    expect(menuStyles.backgroundColor).toBe(darkCardToken);
  });

  it('keeps the closed control transparent so it stays attached to the row', () => {
    themeMode = 'dark';
    render(<MultiSelectWrapper onChange={vi.fn()} options={[]} value={[]} />);
    expect(capturedSelectProps.at(-1).styles.control({}).backgroundColor).toBe('transparent');

    themeMode = 'light';
    capturedSelectProps = [];
    render(<MultiSelectWrapper onChange={vi.fn()} options={[]} value={[]} />);
    expect(capturedSelectProps.at(-1).styles.control({}).backgroundColor).toBe('transparent');
  });

  it('hides the placeholder when focused or when input has a value', () => {
    render(<MultiSelectWrapper onChange={vi.fn()} options={[]} value={[]} />);

    const { placeholder } = capturedSelectProps.at(-1).styles;

    expect(
      placeholder({ display: 'block' }, { isFocused: true, selectProps: { inputValue: '' } })
        .display,
    ).toBe('none');
    expect(
      placeholder({ display: 'block' }, { isFocused: false, selectProps: { inputValue: 'abc' } })
        .display,
    ).toBe('none');
    expect(
      placeholder({ display: 'block' }, { isFocused: false, selectProps: { inputValue: '' } })
        .display,
    ).toBe('block');
  });

  it('uses a primary-tinted option highlight (not action.hover)', () => {
    expect(
      getMultiSelectOptionHighlight({
        palette: { mode: 'dark', primary: { main: '#00B39F' } },
      } as any),
    ).toBe(alpha('#00B39F', 0.28));
    expect(
      getMultiSelectOptionHighlight({
        palette: { mode: 'light', primary: { main: '#00B39F' } },
      } as any),
    ).toBe(alpha('#00B39F', 0.12));
  });

  it('keeps label text WCAG AA (≥4.5:1) on the blended option highlight', () => {
    const primary = '#00B39F';

    const darkBlended = blendRgbaOverHex(alpha(primary, 0.28), darkModePalette.background.card);
    expect(contrastRatio(darkModePalette.text.default, darkBlended)).toBeGreaterThanOrEqual(4.5);

    const lightBlended = blendRgbaOverHex(alpha(primary, 0.12), '#ffffff');
    expect(contrastRatio(lightModePalette.text.default, lightBlended)).toBeGreaterThanOrEqual(4.5);

    // Checkbox (primary) on dark highlight should stay at least UI-component contrast.
    expect(contrastRatio(primary, darkBlended)).toBeGreaterThanOrEqual(3);
  });

  it('uses getContrastText for multiValue chip label and remove colors', () => {
    const assertChipContrast = (mode: 'light' | 'dark') => {
      themeMode = mode;
      capturedSelectProps = [];
      render(<MultiSelectWrapper onChange={vi.fn()} options={[]} value={[]} />);
      const { styles } = capturedSelectProps.at(-1);
      const label = styles.multiValueLabel({});
      const remove = styles.multiValueRemove({});

      const expectedBg =
        mode === 'dark'
          ? darkModePalette.background.elevatedComponents
          : lightModePalette.background.secondary;
      expect(label.backgroundColor).toBe(expectedBg);
      expect(remove.backgroundColor).toBe(expectedBg);
      expect(label.color).toBe(remove.color);
      expect(label.color).not.toBe(expectedBg);

      const fgSolid = String(label.color).startsWith('#')
        ? label.color
        : blendRgbaOverHex(label.color, expectedBg);
      const ratio = contrastRatio(fgSolid, expectedBg);
      // Text and icon share chipForeground — both must meet AA text (≥4.5).
      expect(ratio).toBeGreaterThanOrEqual(4.5);
      return { bg: expectedBg, fg: label.color, ratio };
    };

    assertChipContrast('dark');
    assertChipContrast('light');
  });

  it('marks Mui-selected from isSelected only, not isFocused alone', () => {
    themeMode = 'dark';

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

    const { container, rerender } = render(
      <Option
        innerRef={vi.fn()}
        innerProps={{ onClick: vi.fn(), role: 'option' }}
        isFocused={true}
        isSelected={true}
        label="Production"
        value="prod"
      />,
    );

    expect(container.querySelector('[data-testid="multi-select-option"]')).toHaveClass(
      'Mui-selected',
    );
    expect(screen.getByText('Production')).toBeInTheDocument();

    rerender(
      <Option
        innerRef={vi.fn()}
        innerProps={{ onClick: vi.fn(), role: 'option' }}
        isFocused={false}
        isSelected={true}
        label="Production"
        value="prod"
      />,
    );
    expect(container.querySelector('[data-testid="multi-select-option"]')).toHaveClass(
      'Mui-selected',
    );

    rerender(
      <Option
        innerRef={vi.fn()}
        innerProps={{ onClick: vi.fn(), role: 'option' }}
        isFocused={true}
        isSelected={false}
        label="Production"
        value="prod"
      />,
    );
    expect(container.querySelector('[data-testid="multi-select-option"]')).not.toHaveClass(
      'Mui-selected',
    );

    rerender(
      <Option
        innerRef={vi.fn()}
        innerProps={{ onClick: vi.fn(), role: 'option' }}
        isFocused={false}
        isSelected={false}
        label="Production"
        value="prod"
      />,
    );

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

  it('renders MultiSelectWrapper option content without a menu context error', () => {
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

    render(
      <Option
        innerRef={vi.fn()}
        innerProps={{ onClick: vi.fn(), role: 'option' }}
        isFocused={true}
        isSelected={true}
        label="Production"
        value="prod"
      />,
    );

    expect(screen.getByText('Production')).toBeInTheDocument();
  });
});
