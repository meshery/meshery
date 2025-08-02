# Chip Hover Styles Utility

This utility provides reusable hover styles for chips with overflow tooltip functionality across the Meshery UI.

## Features

- **Consistent hover behavior** across all chip components
- **Improved viewport boundary containment** using CSS `clamp()` function
- **Modern word breaking** using standard CSS properties
- **Responsive max-width** that adapts to viewport size
- **Smooth transitions** for better user experience

## Usage

```javascript
import { createEnhancedChipHoverStyles } from '../../utils/chipHoverStyles';

const StyledChip = styled(Chip)(({ theme }) => ({
  // Your existing styles...
  ...createEnhancedChipHoverStyles(theme),
}));
```

## API

### `createChipHoverStyles(theme)`
Basic hover styles with improved viewport boundary containment.

### `createEnhancedChipHoverStyles(theme)`
Complete chip hover styles including base chip styling and enhanced hover effects.

## Improvements over previous implementation

1. **Eliminated code duplication** - Styles are now centralized in one utility
2. **Better viewport boundary handling** - Uses `clamp()` for better positioning
3. **Modern CSS properties** - Removed deprecated `word-break: 'break-word'`
4. **Consistent behavior** - All chips now have the same hover behavior

## Files using this utility

- `ui/components/connections/styles.js` - ChipWrapper component
- `ui/components/telemetry/grafana/GrafanaDisplaySelection.js` - StyledChip component
- `ui/components/telemetry/grafana/GrafanaSelectionComponent.js` - StyledChip component
- `ui/components/telemetry/prometheus/PrometheusSelectionComponent.js` - PanelChipsContainer styles
