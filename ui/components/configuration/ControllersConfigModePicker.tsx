import type { ReactNode } from 'react';
import { Box, CustomTooltip, Radio, Typography } from '@sistent/sistent';
import { alpha, styled } from '@/theme';
import { INHERIT } from './controllersConfigForm.shared';

const ModeCard = styled(Box, {
  shouldForwardProp: (prop) => prop !== 'selected' && prop !== 'disabled',
})<{ selected?: boolean; disabled?: boolean }>(({ theme, selected, disabled }) => {
  const brand = theme.palette.background.brand?.default || theme.palette.primary.main;
  return {
    display: 'flex',
    alignItems: 'center',
    gap: theme.spacing(1.25),
    padding: theme.spacing(1.5),
    borderRadius: theme.spacing(1),
    cursor: disabled ? 'not-allowed' : 'pointer',
    opacity: disabled ? 0.6 : 1,
    border: `1px solid ${selected ? brand : theme.palette.divider}`,
    background: selected ? alpha(brand, 0.08) : theme.palette.background.card,
    transition: 'border-color 0.15s ease, background 0.15s ease',
    '&:hover': disabled ? undefined : { borderColor: brand },
  };
});

type ControllersConfigModePickerProps = {
  label: ReactNode;
  selected: string;
  inheritModeLabel: string;
  disabled: boolean;
  onChange: (selected: string | undefined) => void;
};

export default function ControllersConfigModePicker({
  label,
  selected,
  inheritModeLabel,
  disabled,
  onChange,
}: ControllersConfigModePickerProps) {
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
      sx={{ marginBottom: '1.5rem', scrollMarginTop: '1rem' }}
      data-testid="controllers-config-mode-picker"
    >
      {label}
      <Box
        sx={{
          display: 'grid',
          gap: 1.25,
          marginTop: '0.75rem',
          gridTemplateColumns: { xs: '1fr', md: '1fr 1fr 1fr' },
        }}
        role="radiogroup"
        aria-label="Deployment mode"
      >
        {options.map((option) => {
          const isSelected = selected === option.value;
          const set = () => onChange(option.value === INHERIT ? undefined : option.value);
          return (
            <CustomTooltip key={option.value} title={option.description} placement="top">
              <ModeCard
                selected={isSelected}
                disabled={disabled}
                role="radio"
                aria-checked={isSelected}
                aria-label={option.label}
                tabIndex={disabled ? -1 : 0}
                data-testid={option.testId}
                onClick={() => {
                  if (!disabled) set();
                }}
                onKeyDown={(event) => {
                  if (disabled) return;
                  if (event.key === 'Enter' || event.key === ' ') {
                    event.preventDefault();
                    set();
                  }
                }}
              >
                <Radio
                  checked={isSelected}
                  tabIndex={-1}
                  disabled={disabled}
                  onChange={set}
                  sx={{ p: 0 }}
                />
                <Typography variant="body2" sx={{ fontWeight: 600, minWidth: 0 }}>
                  {option.label}
                </Typography>
              </ModeCard>
            </CustomTooltip>
          );
        })}
      </Box>
    </Box>
  );
}
