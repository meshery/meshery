import React, { forwardRef, useEffect, useState, type ReactNode } from 'react';
import {
  Autocomplete,
  Box,
  Card,
  Chip,
  CustomTooltip,
  FormControlLabel,
  InfoTooltip,
  MenuItem,
  Radio,
  RadioGroup,
  TextField,
  Typography,
} from '@sistent/sistent';
import { alpha } from '@/theme';
import { getPath, setPath, type FieldPath } from './fieldPath';
import {
  INHERIT,
  fitNumberWidth,
  fitWidth,
  isInertIn,
  type DeploymentModeGovernance,
} from './deploymentMode';
import type { ControllersConfigDoc } from './ControllersConfigForm';

type ChoiceCardProps = {
  value: string;
  selected: boolean;
  disabled?: boolean;
  label: ReactNode;
  testId?: string;
};

/** RadioGroup row plus the card chrome RadioGroup does not draw. */
export const ChoiceCard = forwardRef<HTMLDivElement, ChoiceCardProps>(function ChoiceCard(
  { value, selected, disabled, label, testId, ...props },
  ref,
) {
  return (
    <Card
      ref={ref}
      variant="outlined"
      {...props}
      sx={(theme) => {
        const brand = theme.palette.background.brand?.default || theme.palette.primary.main;
        return {
          boxShadow: 'none',
          opacity: disabled ? 0.6 : 1,
          borderColor: selected ? brand : 'divider',
          bgcolor: selected ? alpha(brand, 0.08) : 'background.card',
          transition: 'border-color 0.15s ease, background 0.15s ease',
          '&:hover': disabled ? undefined : { borderColor: brand },
        };
      }}
    >
      <FormControlLabel
        value={value}
        disabled={disabled}
        data-testid={testId}
        control={<Radio sx={{ p: 0 }} />}
        label={label}
        sx={{
          m: 0,
          width: '100%',
          alignItems: 'center',
          gap: 1.25,
          px: 1.5,
          py: 1.5,
          '& .MuiFormControlLabel-label': { fontWeight: 600, fontSize: '0.875rem', minWidth: 0 },
        }}
      />
    </Card>
  );
});

const CHART_PACKAGE = 'https://artifacthub.io/api/v1/packages/helm/meshery/meshery-operator';

export function OperatorVersionField({
  value,
  placeholder,
  disabled,
  onChange,
}: {
  value: string | undefined;
  placeholder: string;
  disabled: boolean;
  onChange: (next: string | undefined) => void;
}) {
  const [versions, setVersions] = useState<string[]>([]);
  useEffect(() => {
    if (typeof fetch !== 'function') return;
    const ac = new AbortController();
    fetch(CHART_PACKAGE, { signal: ac.signal })
      .then((res) => (res.ok ? res.json() : Promise.reject()))
      .then((data) => {
        const tags = (data?.available_versions ?? [])
          .map((entry: { version?: string }) => entry.version)
          .filter((version: string | undefined): version is string => Boolean(version));
        setVersions(tags);
      })
      .catch(() => undefined);
    return () => ac.abort();
  }, []);

  return (
    <Autocomplete
      freeSolo
      disabled={disabled}
      options={versions}
      value={value ?? ''}
      onChange={(_, next) => {
        const normalized = next == null ? '' : String(next).trim();
        onChange(normalized === '' ? undefined : normalized);
      }}
      onInputChange={(_, next, reason) => {
        if (reason === 'reset') return;
        const normalized = next.trim();
        onChange(normalized === '' ? undefined : normalized);
      }}
      renderInput={(params) => (
        <TextField
          {...params}
          size="small"
          placeholder={placeholder}
          aria-label="Operator version"
          slotProps={{
            ...params.slotProps,
            htmlInput: {
              ...params.inputProps,
              ...(typeof params.slotProps?.htmlInput === 'object'
                ? params.slotProps.htmlInput
                : {}),
              'aria-label': 'Operator version',
            },
          }}
        />
      )}
      sx={{
        width: fitWidth(value, placeholder, 'Inherit (server release)'),
        maxWidth: '100%',
        '& input::placeholder': {
          whiteSpace: 'nowrap',
          overflow: 'hidden',
          textOverflow: 'ellipsis',
        },
      }}
    />
  );
}

export function ControllersConfigModePicker({
  label,
  selected,
  inheritModeLabel,
  disabled,
  onChange,
}: {
  label: ReactNode;
  selected: string;
  inheritModeLabel: string;
  disabled: boolean;
  onChange: (selected: string | undefined) => void;
}) {
  const options = [
    {
      value: INHERIT,
      label: `Inherit (${inheritModeLabel})`,
      description:
        'Use the next layer (server default or built-in). Built-in default is Embedded. Changing the mode redeploys controllers.',
      testId: 'controllers-config-mode-inherit',
    },
    {
      value: 'operator',
      label: 'Operator (in-cluster)',
      description:
        'Installs Meshery Operator, MeshSync, and Broker into the cluster. Full controller settings apply. Changing the mode redeploys controllers.',
      testId: 'controllers-config-mode-operator',
    },
    {
      value: 'embedded',
      label: 'Embedded (in Meshery Server)',
      description:
        'Runs MeshSync inside Meshery Server. Nothing is installed on the cluster; only output filters apply from this form. Changing the mode redeploys controllers.',
      testId: 'controllers-config-mode-embedded',
    },
  ];

  return (
    <Box
      id="controllers-config-mode"
      sx={{ marginBottom: '1rem', scrollMarginTop: '1rem' }}
      data-testid="controllers-config-mode-picker"
    >
      {label}
      <RadioGroup
        value={selected}
        onChange={(_, next) => onChange(next === INHERIT ? undefined : next)}
        aria-label="Deployment mode"
        sx={{
          display: 'grid',
          gap: 1.25,
          marginTop: '0.75rem',
          gridTemplateColumns: { xs: '1fr', md: '1fr 1fr 1fr' },
        }}
      >
        {options.map((option) => (
          <CustomTooltip key={option.value} title={option.description} placement="top">
            <Box>
              <ChoiceCard
                value={option.value}
                selected={selected === option.value}
                disabled={disabled}
                label={option.label}
                testId={option.testId}
              />
            </Box>
          </CustomTooltip>
        ))}
      </RadioGroup>
    </Box>
  );
}

export const fieldRowSx = {
  flexWrap: 'wrap',
  alignItems: 'flex-start',
  '& > *': { flex: '0 1 auto', maxWidth: '100%' },
} as const;

export const controlSx = (nowrapPlaceholder = true) => ({
  width: '100%',
  ...(nowrapPlaceholder
    ? {
        '& input::placeholder, & textarea::placeholder': {
          whiteSpace: 'nowrap',
          overflow: 'hidden',
          textOverflow: 'ellipsis',
        },
      }
    : {}),
});

type SourceInfo = { label: string; overridden: boolean };

export type ControllersConfigFieldsApi = {
  value: ControllersConfigDoc;
  onChange: (next: ControllersConfigDoc) => void;
  inheritedLayers: (ControllersConfigDoc | null | undefined)[];
  inheritLabel: string;
  showSourceIndicators: boolean;
  disabled: boolean;
  liveGovernance?: DeploymentModeGovernance;
};

/** Shared inherit / clamp / source-chip controls used by the form layout. */
export function createControllersConfigFields({
  value,
  onChange,
  inheritedLayers,
  inheritLabel,
  showSourceIndicators,
  disabled,
  liveGovernance,
}: ControllersConfigFieldsApi) {
  const inheritedValue = (path: FieldPath): unknown => {
    for (const layer of inheritedLayers) {
      const v = getPath(layer ?? undefined, path);
      if (v !== undefined) return v;
    }
    return undefined;
  };

  const sourceOf = (path: FieldPath): SourceInfo => {
    if (getPath(value, path) !== undefined) return { label: 'Override', overridden: true };
    if (getPath(inheritedLayers[0] ?? undefined, path) !== undefined)
      return { label: inheritLabel, overridden: false };
    return { label: 'Built-in default', overridden: false };
  };

  const sourceChip = (path: FieldPath) => {
    if (!showSourceIndicators) return null;
    const source = sourceOf(path);
    return (
      <Chip
        size="small"
        label={source.label}
        color={source.overridden ? 'primary' : 'default'}
        variant={source.overridden ? 'filled' : 'outlined'}
        sx={{ marginLeft: '0.5rem', height: '20px' }}
      />
    );
  };

  const isInert = (path: FieldPath): boolean => isInertIn(liveGovernance, path);
  const isDisabled = (path: FieldPath): boolean => disabled || isInert(path);

  const fieldLabel = (text: string, path: FieldPath, helper?: string) => (
    <Box
      sx={{
        display: 'flex',
        alignItems: 'center',
        marginBottom: '0.25rem',
        gap: '0.25rem',
        flexWrap: 'wrap',
      }}
    >
      <Typography variant="body2" sx={{ fontWeight: 500 }}>
        {text}
      </Typography>
      {helper ? <InfoTooltip helpText={helper} placement="top" /> : null}
      {sourceChip(path)}
      {isInert(path) ? (
        <Chip
          size="small"
          label="Not applied"
          variant="outlined"
          sx={{ marginLeft: '0.25rem', height: '20px' }}
        />
      ) : null}
    </Box>
  );

  const triStateBoolean = (label: string, path: FieldPath, helper?: string) => {
    const current = getPath(value, path) as boolean | undefined;
    const inherited = inheritedValue(path) as boolean | undefined;
    const inheritOption = `Inherit (${inherited === undefined ? 'unset' : inherited ? 'Enabled' : 'Disabled'})`;
    return (
      <Box sx={{ width: fitWidth(label, inheritOption, 'Enabled', 'Disabled'), maxWidth: '100%' }}>
        {fieldLabel(label, path, helper)}
        <TextField
          select
          size="small"
          disabled={isDisabled(path)}
          value={current === undefined ? INHERIT : current ? 'true' : 'false'}
          sx={controlSx(false)}
          onChange={(e) => {
            const v = e.target.value;
            onChange(setPath(value, path, v === INHERIT ? undefined : v === 'true'));
          }}
        >
          <MenuItem value={INHERIT}>{inheritOption}</MenuItem>
          <MenuItem value="true">Enabled</MenuItem>
          <MenuItem value="false">Disabled</MenuItem>
        </TextField>
      </Box>
    );
  };

  const textInput = (
    label: string,
    path: FieldPath,
    helper?: string,
    opts?: { number?: boolean; min?: number; max?: number },
  ) => {
    const current = getPath(value, path) as string | number | undefined;
    const inherited = inheritedValue(path);
    const placeholder = inherited !== undefined ? `Inherit (${inherited})` : 'Inherit';
    const width = opts?.number
      ? fitNumberWidth(label, current, placeholder, 'Inherit')
      : fitWidth(label, current, placeholder, 'Inherit');
    return (
      <Box sx={{ width, maxWidth: '100%' }}>
        {fieldLabel(label, path, helper)}
        <TextField
          size="small"
          type={opts?.number ? 'number' : 'text'}
          disabled={isDisabled(path)}
          value={current ?? ''}
          placeholder={placeholder}
          sx={controlSx()}
          slotProps={
            opts?.number ? { htmlInput: { min: opts.min, max: opts.max, step: 1 } } : undefined
          }
          onKeyDown={
            opts?.number
              ? (event) => {
                  if (['e', 'E', '+', '-', '.'].includes(event.key)) event.preventDefault();
                }
              : undefined
          }
          onChange={(e) => {
            const raw = e.target.value;
            if (raw === '') {
              onChange(setPath(value, path, undefined));
              return;
            }
            if (!opts?.number) {
              onChange(setPath(value, path, raw));
              return;
            }
            // Same integer parse as Settings → Performance; clamp to the schema
            // bounds (html min/max only constrain the stepper).
            const parsed = parseInt(raw, 10);
            if (Number.isNaN(parsed)) return;
            let next = parsed;
            if (opts.min !== undefined) next = Math.max(opts.min, next);
            if (opts.max !== undefined) next = Math.min(opts.max, next);
            onChange(setPath(value, path, next));
          }}
        />
      </Box>
    );
  };

  const listInput = (label: string, path: FieldPath, helper: string) => {
    const current = getPath(value, path) as string[] | undefined;
    const inherited = inheritedValue(path) as string[] | undefined;
    const joined = current?.join(', ') ?? '';
    const placeholder =
      inherited && inherited.length > 0 ? `Inherit (${inherited.join(', ')})` : 'Inherit (all)';
    const width = fitWidth(label, joined, placeholder, 'Inherit (all)');
    return (
      <Box sx={{ width, maxWidth: '100%' }}>
        {fieldLabel(label, path, helper)}
        <TextField
          size="small"
          disabled={isDisabled(path)}
          value={joined}
          placeholder={placeholder}
          sx={controlSx()}
          onChange={(e) => {
            const raw = e.target.value;
            if (raw.trim() === '') {
              onChange(setPath(value, path, undefined));
              return;
            }
            onChange(
              setPath(
                value,
                path,
                raw
                  .split(',')
                  .map((s) => s.trim())
                  .filter(Boolean),
              ),
            );
          }}
        />
      </Box>
    );
  };

  const enumSelect = (
    label: string,
    path: FieldPath,
    options: { value: string; label: string }[],
    helper?: string,
    postProcess?: (
      next: ControllersConfigDoc,
      selected: string | undefined,
    ) => ControllersConfigDoc,
  ) => {
    const current = getPath(value, path) as string | undefined;
    const inherited = inheritedValue(path) as string | undefined;
    const inheritOption = `Inherit (${inherited ?? 'unset'})`;
    const width = fitWidth(label, inheritOption, ...options.map((option) => option.label));
    return (
      <Box sx={{ width, maxWidth: '100%' }}>
        {fieldLabel(label, path, helper)}
        <TextField
          select
          size="small"
          disabled={isDisabled(path)}
          value={current ?? INHERIT}
          sx={controlSx(false)}
          onChange={(e) => {
            const selected = e.target.value === INHERIT ? undefined : e.target.value;
            let next = setPath(value, path, selected);
            if (postProcess) next = postProcess(next, selected);
            onChange(next);
          }}
        >
          <MenuItem value={INHERIT}>{inheritOption}</MenuItem>
          {options.map((option) => (
            <MenuItem key={option.value} value={option.value}>
              {option.label}
            </MenuItem>
          ))}
        </TextField>
      </Box>
    );
  };

  return {
    inheritedValue,
    isInert,
    isDisabled,
    fieldLabel,
    triStateBoolean,
    textInput,
    listInput,
    enumSelect,
  };
}
