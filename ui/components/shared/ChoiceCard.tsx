import { forwardRef, type ReactNode } from 'react';
import { Card, Radio, Typography } from '@sistent/sistent';
import { alpha } from '@/theme';

type ChoiceCardProps = {
  selected: boolean;
  disabled?: boolean;
  label: ReactNode;
  ariaLabel?: string;
  onSelect: () => void;
  testId?: string;
};

/** Selectable radio card. Shared by Settings and the connection wizard. */
const ChoiceCard = forwardRef<HTMLDivElement, ChoiceCardProps>(function ChoiceCard(
  { selected, disabled, label, ariaLabel, onSelect, testId, ...props },
  ref,
) {
  const select = () => {
    if (!disabled) onSelect();
  };

  return (
    <Card
      ref={ref}
      variant="outlined"
      role="radio"
      aria-checked={selected}
      aria-disabled={disabled}
      aria-label={ariaLabel ?? (typeof label === 'string' ? label : undefined)}
      tabIndex={disabled ? -1 : 0}
      data-testid={testId}
      onClick={select}
      onKeyDown={(event) => {
        if (event.key === 'Enter' || event.key === ' ') {
          event.preventDefault();
          select();
        }
      }}
      {...props}
      sx={(theme) => {
        const brand = theme.palette.background.brand?.default || theme.palette.primary.main;
        return {
          display: 'flex',
          alignItems: 'center',
          gap: 1.25,
          px: 1.5,
          py: 1.5,
          cursor: disabled ? 'not-allowed' : 'pointer',
          opacity: disabled ? 0.6 : 1,
          borderColor: selected ? brand : 'divider',
          bgcolor: selected ? alpha(brand, 0.08) : 'background.card',
          boxShadow: 'none',
          transition: 'border-color 0.15s ease, background 0.15s ease',
          '&:hover': disabled ? undefined : { borderColor: brand },
        };
      }}
    >
      <Radio checked={selected} disabled={disabled} tabIndex={-1} onChange={select} sx={{ p: 0 }} />
      <Typography variant="body2" sx={{ fontWeight: 600, minWidth: 0 }}>
        {label}
      </Typography>
    </Card>
  );
});

export default ChoiceCard;
