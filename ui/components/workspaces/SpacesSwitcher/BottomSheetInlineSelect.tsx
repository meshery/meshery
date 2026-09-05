import React, { useState } from 'react';
import { Box, CheckCircleIcon, Collapse, ExpandMoreIcon, useTheme } from '@sistent/sistent';

export type BottomSheetInlineSelectOption = {
  id: string;
  label: string;
};

type BottomSheetInlineSelectProps = {
  value: string;
  options: BottomSheetInlineSelectOption[];
  onSelect: (id: string) => void;
  placeholder?: string;
  leadingIcon?: React.ReactNode;
  renderOption?: (option: BottomSheetInlineSelectOption, selected: boolean) => React.ReactNode;
  expanded?: boolean;
  onExpandedChange?: (expanded: boolean) => void;
  'data-cy'?: string;
};

/**
 * Inline dropdown for BottomSheet — tap the field to expand options below it.
 * Unlike MUI Select, the menu is not portaled and stays inside the sheet.
 */
export function BottomSheetInlineSelect({
  value,
  options,
  onSelect,
  placeholder = 'Select…',
  leadingIcon,
  renderOption,
  expanded: expandedProp,
  onExpandedChange,
  'data-cy': dataCy,
}: BottomSheetInlineSelectProps) {
  const theme = useTheme();
  const [expandedInternal, setExpandedInternal] = useState(false);
  const expanded = expandedProp ?? expandedInternal;
  const selected = options.find((option) => option.id === value);

  const setExpanded = (next: boolean) => {
    if (expandedProp === undefined) {
      setExpandedInternal(next);
    }
    onExpandedChange?.(next);
  };

  const selectedBackground = theme.palette.action.selected;

  return (
    <Box sx={{ width: '100%' }} data-cy={dataCy}>
      <Box
        component="button"
        type="button"
        aria-expanded={expanded}
        aria-haspopup="listbox"
        onClick={() => setExpanded(!expanded)}
        sx={{
          width: '100%',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          gap: 1,
          px: 2,
          py: 1.5,
          border: `1px solid ${theme.palette.divider}`,
          borderRadius: 1,
          backgroundColor: expanded ? selectedBackground : theme.palette.background.card,
          color: theme.palette.text.default,
          cursor: 'pointer',
          font: 'inherit',
          fontWeight: 500,
          fontSize: '0.9375rem',
          textAlign: 'left',
          transition: 'background-color 0.15s ease',
        }}
      >
        <Box
          sx={{
            flex: 1,
            minWidth: 0,
            display: 'flex',
            alignItems: 'center',
            gap: 1.5,
          }}
        >
          {leadingIcon ? (
            <Box sx={{ display: 'flex', flexShrink: 0, alignItems: 'center' }}>{leadingIcon}</Box>
          ) : null}
          <Box
            component="span"
            sx={{
              flex: 1,
              minWidth: 0,
              overflow: 'hidden',
              textOverflow: 'ellipsis',
              whiteSpace: 'nowrap',
            }}
          >
            {selected?.label ?? placeholder}
          </Box>
        </Box>
        <ExpandMoreIcon
          width={20}
          height={20}
          fill={theme.palette.icon.default}
          style={{
            flexShrink: 0,
            transform: expanded ? 'rotate(180deg)' : 'rotate(0deg)',
            transition: 'transform 0.2s ease',
          }}
        />
      </Box>

      <Collapse in={expanded} unmountOnExit>
        <Box
          role="listbox"
          sx={{
            mt: 0.5,
            border: `1px solid ${theme.palette.divider}`,
            borderRadius: 1,
            overflow: 'hidden',
            backgroundColor: theme.palette.background.card,
            maxHeight: 'min(32vh, 220px)',
            overflowY: 'auto',
          }}
        >
          {options.map((option, index) => {
            const isSelected = option.id === value;
            const isLast = index === options.length - 1;

            return (
              <Box
                key={option.id}
                component="button"
                type="button"
                role="option"
                aria-selected={isSelected}
                onClick={() => {
                  onSelect(option.id);
                  setExpanded(false);
                }}
                sx={{
                  width: '100%',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'space-between',
                  gap: 1.5,
                  px: 2,
                  py: 1.25,
                  border: 'none',
                  borderBottom: isLast ? 'none' : `1px solid ${theme.palette.divider}`,
                  backgroundColor: isSelected ? selectedBackground : 'transparent',
                  color: theme.palette.text.default,
                  cursor: 'pointer',
                  font: 'inherit',
                  fontSize: '0.9375rem',
                  textAlign: 'left',
                  '&:hover': {
                    backgroundColor: isSelected ? selectedBackground : theme.palette.action.hover,
                  },
                }}
              >
                <Box
                  sx={{
                    flex: 1,
                    minWidth: 0,
                    display: 'flex',
                    alignItems: 'center',
                    gap: 1.5,
                    overflow: 'hidden',
                  }}
                >
                  {renderOption ? (
                    renderOption(option, isSelected)
                  ) : (
                    <Box
                      component="span"
                      sx={{
                        overflow: 'hidden',
                        textOverflow: 'ellipsis',
                        whiteSpace: 'nowrap',
                      }}
                    >
                      {option.label}
                    </Box>
                  )}
                </Box>
                {isSelected ? (
                  <CheckCircleIcon
                    sx={{
                      flexShrink: 0,
                      fontSize: '1.125rem',
                      color: theme.palette.background.brand.default,
                    }}
                  />
                ) : null}
              </Box>
            );
          })}
        </Box>
      </Collapse>
    </Box>
  );
}
