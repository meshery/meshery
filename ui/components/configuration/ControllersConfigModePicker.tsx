import type { ReactNode } from 'react';
import { Box, CustomTooltip } from '@sistent/sistent';
import ChoiceCard from '@/components/shared/ChoiceCard';
import { INHERIT } from './controllersConfigForm.shared';
import { modeConsequence, type DeploymentMode } from './deploymentMode';

type ControllersConfigModePickerProps = {
  label: ReactNode;
  selected: string;
  inheritModeLabel: string;
  inheritedMode: DeploymentMode;
  scope: 'connection' | 'serverDefault';
  disabled: boolean;
  onChange: (selected: string | undefined) => void;
};

export default function ControllersConfigModePicker({
  label,
  selected,
  inheritModeLabel,
  inheritedMode,
  scope,
  disabled,
  onChange,
}: ControllersConfigModePickerProps) {
  const options = [
    {
      value: INHERIT,
      label: `Inherit (${inheritModeLabel})`,
      description: `Use the next layer (server default or built-in). ${modeConsequence(inheritedMode, scope)}`,
      testId: 'controllers-config-mode-inherit',
    },
    {
      value: 'operator',
      label: 'Operator (in-cluster)',
      description: modeConsequence('operator', scope),
      testId: 'controllers-config-mode-operator',
    },
    {
      value: 'embedded',
      label: 'Embedded (in Meshery Server)',
      description: modeConsequence('embedded', scope),
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
      <Box
        role="radiogroup"
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
                selected={selected === option.value}
                disabled={disabled}
                label={option.label}
                testId={option.testId}
                onSelect={() => onChange(option.value === INHERIT ? undefined : option.value)}
              />
            </Box>
          </CustomTooltip>
        ))}
      </Box>
    </Box>
  );
}
