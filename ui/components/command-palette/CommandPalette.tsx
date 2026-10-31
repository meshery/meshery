import {
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
  type ComponentType,
  type HTMLAttributes,
  type Key,
  type ReactNode,
} from 'react';
import { useRouter } from 'next/router';
import { useSelector } from 'react-redux';
import {
  Autocomplete,
  Box,
  Dialog,
  InputAdornment,
  SearchIcon,
  TextField,
  Typography,
  useTheme,
  type AutocompleteProps,
} from '@sistent/sistent';
import type { RootState } from '@/store';
import { ProviderUiAccessControl } from '@/utils/disabledComponents';
import { resolveNavigatorComponents } from '@/components/layout/Navigator/Navigator';
import { useCommandPalette } from '@/utils/context/CommandPaletteContextProvider';
import { buildCommands, type Command } from './commands';
import { filterCommands } from './filter';

const InlinePopper = ({ children, className }: { children?: ReactNode; className?: string }) => {
  return <div className={className}>{children}</div>;
};

const TypedAutocomplete = Autocomplete as unknown as ComponentType<
  AutocompleteProps<Command, false, false, false>
>;

type RenderOptionProps = HTMLAttributes<HTMLLIElement> & { key?: Key };

const visuallyHiddenSx = {
  position: 'absolute',
  width: '1px',
  height: '1px',
  padding: 0,
  margin: '-1px',
  overflow: 'hidden',
  clip: 'rect(0 0 0 0)',
  whiteSpace: 'nowrap',
  border: 0,
} as const;

export const CommandPalette = () => {
  const { open, closePalette } = useCommandPalette();
  const theme = useTheme();
  const router = useRouter();

  const providerCapabilities = useSelector((state: RootState) => state.ui?.providerCapabilities);
  const meshAdapters = useSelector((state: RootState) => state.adapter?.meshAdapters);
  const catalogVisibility = useSelector((state: RootState) => state.ui?.catalogVisibility);

  const [inputValue, setInputValue] = useState('');
  const pendingCommandRef = useRef<Command | null>(null);
  const inputRef = useRef<HTMLInputElement | null>(null);

  // Reset search and pending commands every time the palette opens
  useEffect(() => {
    if (open) {
      setInputValue('');
      pendingCommandRef.current = null;
    }
  }, [open]);

  const providerUiAccessControl = useMemo(() => {
    if (!providerCapabilities) {
      return null;
    }
    return new ProviderUiAccessControl(providerCapabilities);
  }, [providerCapabilities]);

  const nav = useMemo(() => {
    if (!providerCapabilities || !providerUiAccessControl) {
      return [];
    }
    const currentPath = (router.asPath || '').split('?')[0];
    return resolveNavigatorComponents({
      providerUiAccessControl,
      theme,
      meshAdapters,
      catalogVisibility,
      currentPath,
    });
  }, [
    providerCapabilities,
    providerUiAccessControl,
    theme,
    meshAdapters,
    catalogVisibility,
    router.asPath,
  ]);

  // Re-evaluated every time palette opens because `ability` is a mutable singleton
  const commands = useMemo(() => {
    if (!providerCapabilities || !providerUiAccessControl) {
      return [];
    }
    return buildCommands(nav, providerUiAccessControl);
  }, [open, nav, providerCapabilities, providerUiAccessControl]);

  const filtered = useMemo(() => {
    return filterCommands(commands, inputValue);
  }, [commands, inputValue]);

  // Contiguous grouping by section order: Actions, Navigate, Settings
  // Preserves rank order within each section so MUI groupBy renders each header at most once
  const contiguousOptions = useMemo(() => {
    const actions = filtered.filter((c) => c.section === 'Actions');
    const navigate = filtered.filter((c) => c.section === 'Navigate');
    const settings = filtered.filter((c) => c.section === 'Settings');
    return [...actions, ...navigate, ...settings];
  }, [filtered]);

  const handleSelect = useCallback(
    (_event: unknown, option: Command | null) => {
      if (!option || option.disabled) {
        return;
      }
      pendingCommandRef.current = option;
      closePalette();
    },
    [closePalette],
  );

  const handleExited = useCallback(() => {
    const cmd = pendingCommandRef.current;
    pendingCommandRef.current = null;
    if (!cmd) {
      return;
    }

    if (
      cmd.href.startsWith('/configuration/designs?action=import') &&
      router.pathname === '/configuration/designs'
    ) {
      router.push({
        pathname: '/configuration/designs',
        query: {
          ...router.query,
          action: 'import',
        },
      });
    } else {
      router.push(cmd.href);
    }
  }, [router]);

  const handleClose = useCallback(
    (_event: unknown, reason: 'backdropClick' | 'escapeKeyDown') => {
      if (reason === 'escapeKeyDown' || reason === 'backdropClick') {
        pendingCommandRef.current = null;
        closePalette();
      }
    },
    [closePalette],
  );

  const statusText = useMemo(() => {
    if (commands.length === 0) {
      return 'No commands available';
    }
    if (filtered.length === 0) {
      return 'No commands found';
    }
    return `${filtered.length} result${filtered.length === 1 ? '' : 's'}`;
  }, [commands.length, filtered.length]);

  const noOptionsText = commands.length === 0 ? 'No commands available' : 'No commands found';

  return (
    <Dialog
      open={open}
      onClose={handleClose}
      data-testid="command-palette"
      aria-labelledby="command-palette-title"
      fullWidth
      maxWidth="sm"
      slotProps={{
        paper: {
          'aria-label': 'Command palette',
        },
        transition: {
          onEntered: () => {
            if (open) {
              inputRef.current?.focus();
            }
          },
          onExited: handleExited,
        },
      }}
      sx={{
        '& .MuiDialog-container': {
          alignItems: 'flex-start',
          pt: { xs: 4, sm: 8 },
        },
        '& .MuiDialog-paper': {
          backgroundColor: theme.palette.background.card,
          backgroundImage: 'none',
          borderRadius: 2,
          overflow: 'hidden',
          border: `1px solid ${theme.palette.border.default}`,
          boxShadow: theme.shadows?.[8],
        },
      }}
    >
      <Box sx={{ position: 'relative' }}>
        <Typography id="command-palette-title" sx={visuallyHiddenSx}>
          Command palette
        </Typography>

        <Box role="status" aria-live="polite" sx={visuallyHiddenSx}>
          {statusText}
        </Box>

        <TypedAutocomplete
          open
          freeSolo={false}
          autoHighlight
          disabledItemsFocusable
          clearOnEscape={false}
          forcePopupIcon={false}
          filterOptions={(x) => x}
          options={contiguousOptions}
          value={null}
          inputValue={inputValue}
          onInputChange={(_e, v, reason) => {
            if (reason !== 'reset') {
              setInputValue(v);
            }
          }}
          onClose={(_e, reason) => {
            if (reason === 'escape') {
              pendingCommandRef.current = null;
              closePalette();
            }
          }}
          onChange={handleSelect}
          groupBy={(option) => option.section}
          getOptionLabel={(option) => option.title}
          getOptionDisabled={(option) => option.disabled}
          noOptionsText={noOptionsText}
          slots={{
            popper: InlinePopper,
          }}
          slotProps={{
            paper: {
              elevation: 0,
              sx: {
                backgroundColor: 'transparent',
                boxShadow: 'none',
                border: 'none',
              },
            },
            listbox: {
              sx: {
                maxHeight: '380px',
                p: 0,
                backgroundColor: theme.palette.background.card,
                '& .MuiAutocomplete-option.Mui-focused, & .MuiAutocomplete-option:hover, & .MuiAutocomplete-option[aria-selected="true"]':
                  {
                    backgroundColor: theme.palette.background.hover,
                    boxShadow: `inset 3px 0 0 ${theme.palette.border.brand}`,
                  },
                '& .MuiAutocomplete-option[aria-disabled="true"]': {
                  opacity: 1,
                  cursor: 'not-allowed',
                },
              },
            },
          }}
          renderInput={(params) => {
            const { slotProps: paramsSlotProps, ...otherParams } = params;
            return (
              <TextField
                {...otherParams}
                autoFocus
                variant="outlined"
                placeholder="Search commands and pages..."
                slotProps={{
                  ...paramsSlotProps,
                  htmlInput: {
                    ...paramsSlotProps?.htmlInput,
                    ref: (node: HTMLInputElement | null) => {
                      inputRef.current = node;
                      const origRef = paramsSlotProps?.htmlInput?.ref;
                      if (typeof origRef === 'function') {
                        origRef(node);
                      } else if (origRef && typeof origRef === 'object' && 'current' in origRef) {
                        (origRef as { current: HTMLInputElement | null }).current = node;
                      }
                    },
                    'aria-label': 'Search commands',
                    'data-testid': 'command-palette-input',
                  },
                  input: {
                    ...paramsSlotProps?.input,
                    startAdornment: (
                      <>
                        <InputAdornment position="start">
                          <SearchIcon
                            sx={{ color: theme.palette.icon.default, fontSize: '1.25rem' }}
                          />
                        </InputAdornment>
                        {paramsSlotProps?.input?.startAdornment}
                      </>
                    ),
                    sx: {
                      borderRadius: 0,
                      borderBottom: `1px solid ${theme.palette.border.default}`,
                      '& fieldset': { border: 'none' },
                      backgroundColor: theme.palette.background.card,
                      color: theme.palette.text.primary,
                    },
                  },
                }}
              />
            );
          }}
          renderGroup={(params) => {
            const headingId = `command-palette-group-${params.group.toLowerCase().replace(/\s+/g, '-')}`;
            return (
              <li key={params.key || params.group} role="presentation">
                <Typography
                  component="div"
                  role="heading"
                  aria-level={3}
                  id={headingId}
                  variant="overline"
                  sx={{
                    display: 'block',
                    px: 2,
                    py: 0.75,
                    fontWeight: 700,
                    letterSpacing: '0.08em',
                    color: theme.palette.text.secondary,
                    backgroundColor: theme.palette.background.paper,
                  }}
                >
                  {params.group}
                </Typography>
                <Box
                  component="ul"
                  role="group"
                  aria-labelledby={headingId}
                  sx={{ p: 0, m: 0, listStyle: 'none' }}
                >
                  {params.children}
                </Box>
              </li>
            );
          }}
          renderOption={(props: RenderOptionProps, option: Command) => {
            const { key, ...otherProps } = props;
            return (
              <li
                key={key ?? option.id}
                {...otherProps}
                data-testid={`command-palette-item-${option.id}`}
              >
                <Box
                  sx={{
                    display: 'flex',
                    alignItems: 'center',
                    width: '100%',
                    gap: 1.5,
                    py: 0.5,
                    px: 1,
                  }}
                >
                  <Box
                    aria-hidden="true"
                    sx={{
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                      width: 20,
                      height: 20,
                      flexShrink: 0,
                      color: option.disabled
                        ? theme.palette.text.secondary
                        : theme.palette.icon.default,
                      '& svg': { width: 20, height: 20 },
                    }}
                  >
                    {option.icon}
                  </Box>
                  <Box sx={{ flexGrow: 1, minWidth: 0 }}>
                    <Typography
                      variant="body2"
                      sx={{
                        color: option.disabled
                          ? theme.palette.text.secondary
                          : theme.palette.text.primary,
                        fontWeight: 500,
                      }}
                    >
                      {option.title}
                    </Typography>
                    {option.subtitle && (
                      <Typography
                        variant="caption"
                        sx={{
                          display: 'block',
                          color: theme.palette.text.secondary,
                        }}
                      >
                        {option.subtitle}
                      </Typography>
                    )}
                    {option.disabled && option.requires && (
                      <Typography
                        variant="caption"
                        sx={{
                          display: 'block',
                          color: theme.palette.text.secondary,
                          fontStyle: 'italic',
                        }}
                      >
                        Requires: {option.requires}
                      </Typography>
                    )}
                  </Box>
                </Box>
              </li>
            );
          }}
        />
      </Box>
    </Dialog>
  );
};

export default CommandPalette;
