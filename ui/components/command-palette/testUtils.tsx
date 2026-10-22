import React, { createContext, useContext, useEffect, useRef, useState } from 'react';
import { configureStore, type Reducer } from '@reduxjs/toolkit';

export interface MockThemeOptions {
  mode?: 'light' | 'dark';
}

export const createMockTheme = (mode: 'light' | 'dark' = 'light') => ({
  palette: {
    mode,
    background: {
      card: mode === 'light' ? 'white' : 'black',
      paper: mode === 'light' ? 'white' : 'black',
      hover: mode === 'light' ? 'lightgray' : 'darkgray',
      constant: {
        white: 'white',
      },
      brand: {
        default: 'teal',
      },
    },
    navigation: {
      secondary: 'gray',
      active: 'teal',
      hover: 'gray',
    },
    border: {
      default: mode === 'light' ? 'lightgray' : 'darkgray',
      brand: 'teal',
    },
    text: {
      default: mode === 'light' ? 'black' : 'white',
      primary: mode === 'light' ? 'black' : 'white',
      secondary: mode === 'light' ? 'gray' : 'lightgray',
      disabled: mode === 'light' ? 'gray' : 'darkgray',
      constant: {
        white: 'white',
      },
    },
    common: {
      white: 'white',
    },
    icon: {
      default: mode === 'light' ? 'black' : 'white',
    },
  },
  shadows: Array(25).fill('none'),
});

export const sampleCapabilities = {
  providerType: 'remote',
  providerName: 'Meshery Cloud',
  capabilities: {
    canPersistData: true,
  },
  extensions: {
    navigator: [],
  },
};

export const sampleAdapters = [
  {
    adapter_location: 'localhost:10000',
    name: 'Istio',
    version: '1.20.0',
  },
];

export const createTestStoreWithReducers = (
  reducers: { ui: Reducer; adapter: Reducer },
  initialState?: Record<string, unknown>,
) => {
  return configureStore({
    reducer: reducers,
    preloadedState: initialState,
  });
};

interface CommandOption {
  id: string;
  title: string;
  section: string;
  disabled?: boolean;
}

export const createSistentMock = () => {
  class MockEventBus {
    listeners: Record<string, unknown[]> = {};
    on() {
      return {
        subscribe: () => ({ unsubscribe: () => {} }),
      };
    }
    off() {}
    emit() {}
  }

  const ThemeCtx = createContext(createMockTheme('light'));

  const SistentThemeProvider = ({
    children,
    initialTheme = 'light',
    initialMode,
  }: {
    children: React.ReactNode;
    initialTheme?: 'light' | 'dark';
    initialMode?: 'light' | 'dark';
  }) => (
    <ThemeCtx.Provider value={createMockTheme(initialMode || initialTheme)}>
      {children}
    </ThemeCtx.Provider>
  );

  const useTheme = () => useContext(ThemeCtx);

  const cleanDOMProps = <T extends Record<string, unknown>>(props: T): Record<string, unknown> => {
    const {
      startAdornment: _startAdornment,
      endAdornment: _endAdornment,
      primaryFill: _primaryFill,
      secondaryFill: _secondaryFill,
      tertiaryFill: _tertiaryFill,
      fullWidth: _fullWidth,
      size: _size,
      disabled: _disabled,
      ref: _ref,
      sx: _sx,
      ...rest
    } = props;
    return rest;
  };

  const Box = ({
    component = 'div',
    children,
    sx: _sx,
    ...props
  }: {
    component?: React.ElementType;
    children?: React.ReactNode;
    sx?: unknown;
    [key: string]: unknown;
  }) => React.createElement(component, cleanDOMProps(props), children);

  const Typography = ({
    variant: _variant,
    children,
    sx: _sx,
    component = 'span',
    ...props
  }: {
    variant?: string;
    children?: React.ReactNode;
    sx?: unknown;
    component?: React.ElementType;
    [key: string]: unknown;
  }) => React.createElement(component, cleanDOMProps(props), children);

  const SearchIcon = (props: Record<string, unknown>) => (
    <svg data-testid="search-icon" {...cleanDOMProps(props)} />
  );

  const InputAdornment = ({
    children,
    ...props
  }: {
    children: React.ReactNode;
    [key: string]: unknown;
  }) => <div {...cleanDOMProps(props)}>{children}</div>;

  const Button = ({
    children,
    sx: _sx,
    ...props
  }: {
    children: React.ReactNode;
    sx?: unknown;
    [key: string]: unknown;
  }) => <button {...cleanDOMProps(props)}>{children}</button>;

  const CustomTooltip = ({ title, children }: { title: string; children: React.ReactElement }) => (
    <div data-tooltip-title={title}>{children}</div>
  );

  const alpha = (color: string, _opacity: number) => color;

  const Dialog = ({
    open,
    onClose,
    children,
    slotProps,
    'data-testid': testId,
    'aria-labelledby': ariaLabelledby,
    'aria-label': ariaLabel,
    sx: _sx,
    fullWidth: _fw,
    maxWidth: _mw,
  }: {
    open: boolean;
    onClose?: (event: unknown, reason: string) => void;
    children?: React.ReactNode;
    slotProps?: {
      paper?: { 'aria-label'?: string };
      transition?: { onExited?: () => void; onEntered?: () => void };
    };
    'data-testid'?: string;
    'aria-labelledby'?: string;
    'aria-label'?: string;
    sx?: unknown;
    fullWidth?: boolean;
    maxWidth?: string;
  }) => {
    const prevOpenRef = useRef(false);
    useEffect(() => {
      if (!prevOpenRef.current && open && slotProps?.transition?.onEntered) {
        slotProps.transition.onEntered();
      }
      if (prevOpenRef.current && !open && slotProps?.transition?.onExited) {
        slotProps.transition.onExited();
      }
      prevOpenRef.current = open;
    }, [open, slotProps]);

    if (!open) {
      return null;
    }

    const resolvedAriaLabel = slotProps?.paper?.['aria-label'] || ariaLabel;

    return (
      <div
        role="dialog"
        data-testid={testId || 'command-palette'}
        aria-labelledby={ariaLabelledby}
        aria-label={resolvedAriaLabel}
        onKeyDown={(e) => {
          if (e.key === 'Escape') {
            onClose?.(e, 'escapeKeyDown');
          }
        }}
      >
        {children}
      </div>
    );
  };

  const Autocomplete = ({
    options = [],
    inputValue = '',
    onInputChange,
    onClose,
    onChange,
    groupBy,
    getOptionDisabled,
    renderInput,
    renderGroup,
    renderOption,
    noOptionsText,
  }: {
    options?: CommandOption[];
    inputValue?: string;
    onInputChange?: (event: unknown, value: string, reason: string) => void;
    onClose?: (event: unknown, reason: string) => void;
    onChange?: (event: unknown, option: CommandOption | null) => void;
    groupBy?: (option: CommandOption) => string;
    getOptionDisabled?: (option: CommandOption) => boolean;
    renderInput?: (params: Record<string, unknown>) => React.ReactNode;
    renderGroup?: (params: {
      key: string;
      group: string;
      children: React.ReactNode;
    }) => React.ReactNode;
    renderOption?: (
      props: Record<string, unknown>,
      option: CommandOption,
      state: { selected: boolean },
    ) => React.ReactNode;
    noOptionsText?: React.ReactNode;
  }) => {
    const [activeIndex, setActiveIndex] = useState(options.length > 0 ? 0 : -1);

    useEffect(() => {
      setActiveIndex(options.length > 0 ? 0 : -1);
    }, [options]);

    const handleKeyDown = (e: React.KeyboardEvent) => {
      if (e.key === 'ArrowDown') {
        e.preventDefault();
        setActiveIndex((prev) => (options.length === 0 ? -1 : (prev + 1) % options.length));
      } else if (e.key === 'ArrowUp') {
        e.preventDefault();
        setActiveIndex((prev) =>
          options.length === 0 ? -1 : (prev - 1 + options.length) % options.length,
        );
      } else if (e.key === 'Enter') {
        const indexToSelect =
          activeIndex >= 0 && activeIndex < options.length
            ? activeIndex
            : options.length > 0
              ? 0
              : -1;
        if (indexToSelect >= 0 && indexToSelect < options.length) {
          const selectedOption = options[indexToSelect];
          onChange?.(e, selectedOption);
          onInputChange?.(e, selectedOption.title, 'reset');
        }
      } else if (e.key === 'Escape') {
        onClose?.(e, 'escape');
      }
    };

    const activeOption =
      activeIndex >= 0 && activeIndex < options.length ? options[activeIndex] : null;

    const inputParams = {
      id: 'command-palette-combobox',
      disabled: false,
      fullWidth: true,
      size: undefined,
      slotProps: {
        inputLabel: {},
        input: {
          ref: null,
          className: '',
          startAdornment: null,
          endAdornment: null,
          onMouseDown: () => {},
        },
        htmlInput: {
          role: 'combobox',
          'aria-autocomplete': 'list',
          'aria-expanded': true,
          'aria-activedescendant': activeOption ? activeOption.id : undefined,
          value: inputValue,
          onChange: (e: React.ChangeEvent<HTMLInputElement>) => {
            onInputChange?.(e, e.target.value, 'input');
          },
          onKeyDown: handleKeyDown,
        },
      },
    };

    let renderedContent: React.ReactNode = null;

    if (options.length === 0) {
      renderedContent = <div role="presentation">{noOptionsText}</div>;
    } else if (groupBy) {
      const groups: Array<{ key: string; group: string; items: CommandOption[] }> = [];
      let currentGroup: { key: string; group: string; items: CommandOption[] } | null = null;

      options.forEach((opt: CommandOption) => {
        const groupName = groupBy(opt);
        if (!currentGroup || currentGroup.group !== groupName) {
          currentGroup = { key: `group-${groupName}`, group: groupName, items: [opt] };
          groups.push(currentGroup);
        } else {
          currentGroup.items.push(opt);
        }
      });

      let runningIndex = 0;
      renderedContent = (
        <ul role="listbox">
          {groups.map((grp) => {
            const groupChildren = grp.items.map((opt) => {
              const itemIdx = runningIndex++;
              const isSelected = activeIndex === itemIdx;
              const isDisabled = getOptionDisabled ? getOptionDisabled(opt) : false;
              const optionProps = {
                key: opt.id,
                id: opt.id,
                role: 'option',
                'aria-selected': isSelected,
                'aria-disabled': isDisabled,
                onClick: (e: React.MouseEvent) => {
                  onChange?.(e, opt);
                  onInputChange?.(e, opt.title, 'reset');
                },
              };
              return renderOption ? renderOption(optionProps, opt, { selected: isSelected }) : null;
            });

            return renderGroup
              ? renderGroup({
                  key: grp.key,
                  group: grp.group,
                  children: groupChildren,
                })
              : groupChildren;
          })}
        </ul>
      );
    }

    return (
      <div>
        {renderInput?.(inputParams)}
        {renderedContent}
      </div>
    );
  };

  const TextField = ({
    id,
    slotProps,
    placeholder,
    autoFocus: _af,
    variant: _variant,
    inputProps: _legacyInputProps,
    InputProps: _legacyInputPropsCapital,
    ..._restProps
  }: {
    id?: string;
    slotProps?: {
      htmlInput?: Record<string, unknown>;
      input?: {
        startAdornment?: React.ReactNode;
        className?: string;
        onMouseDown?: React.MouseEventHandler;
        [key: string]: unknown;
      };
      [key: string]: unknown;
    };
    placeholder?: string;
    autoFocus?: boolean;
    variant?: string;
    inputProps?: Record<string, unknown>;
    InputProps?: Record<string, unknown>;
    [key: string]: unknown;
  }) => {
    const htmlInputProps = slotProps?.htmlInput || {};
    const inputSlot = slotProps?.input || {};
    const { startAdornment, className, onMouseDown } = inputSlot;

    return (
      <div
        className={className as string | undefined}
        onMouseDown={onMouseDown as React.MouseEventHandler<HTMLDivElement> | undefined}
      >
        {startAdornment}
        <input
          {...cleanDOMProps(htmlInputProps)}
          ref={(node) => {
            const r = htmlInputProps.ref;
            if (typeof r === 'function') {
              r(node);
            } else if (r && typeof r === 'object' && 'current' in r) {
              (r as React.MutableRefObject<HTMLInputElement | null>).current = node;
            }
          }}
          id={id ?? (htmlInputProps.id as string | undefined)}
          placeholder={placeholder}
          autoFocus={_af}
        />
      </div>
    );
  };

  const mockExports: Record<string, unknown> = {
    SistentThemeProvider,
    useTheme,
    Box,
    Typography,
    SearchIcon,
    InputAdornment,
    Button,
    CustomTooltip,
    alpha,
    Dialog,
    Autocomplete,
    TextField,
    EventBus: MockEventBus,
    createCanShow: () => () => true,
    styled: (Component: React.ElementType) => () => {
      const StyledComponent = ({
        children,
        ...props
      }: {
        children?: React.ReactNode;
        [key: string]: unknown;
      }) => React.createElement(Component || 'div', cleanDOMProps(props), children);
      return StyledComponent;
    },
    charcoal: {
      10: 'black',
      20: 'black',
      30: 'black',
      40: 'black',
      50: 'black',
    },
    darkTeal: {
      main: 'teal',
    },
    MESHERY_EXTENSION_EVENT: {
      DispatchToMesheryStore: 'DispatchToMesheryStore',
    },
    InfoIcon: () => null,
    dataValidatorMachine: {
      provide: () => ({}),
    },
    notificationColors: {
      info: { main: 'teal' },
      error: { main: 'red', dark: 'darkred' },
      warning: { main: 'orange' },
      success: { main: 'green' },
    },
  };

  return new Proxy(mockExports, {
    get(target, prop: string | symbol) {
      if (typeof prop === 'symbol') {
        return undefined;
      }
      if (prop === 'then') {
        return undefined;
      }
      if (prop === 'default' || prop === '__esModule') {
        return target[prop as string];
      }
      if (prop in target) {
        return target[prop as string];
      }
      const Passthrough = ({
        children,
        ...props
      }: {
        children?: React.ReactNode;
        [key: string]: unknown;
      }) => <div {...cleanDOMProps(props)}>{children}</div>;
      Passthrough.displayName = `SistentMock(${String(prop)})`;
      return Passthrough;
    },
    has(_target, prop) {
      if (prop === 'then') {
        return false;
      }
      return true;
    },
  });
};
