interface LayoutItem {
  i: string;
  x: number;
  y: number;
  w: number;
  h: number;
  minW?: number;
  minH?: number;
  moved?: boolean;
  static?: boolean;
}

interface WidgetDefaultSizing {
  w: number;
  h: number;
}

type Layouts = Record<string, LayoutItem[]>;
type ColsConfig = Record<string, number>;
type WidgetSizingMap = Record<string, WidgetDefaultSizing>;

export const applyMinSizeConstraints = (
  layouts: Layouts,
  defaultLayouts: Layouts,
  colsConfig: ColsConfig,
  widgetSizing: WidgetSizingMap = {},
): Layouts => {
  const constrained: Layouts = {};

  for (const [bp, items] of Object.entries(layouts || {})) {
    const maxCols = colsConfig[bp] ?? 12;
    const defaults = defaultLayouts?.[bp] ?? [];

    constrained[bp] = items.map((item) => {
      const defaultItem = defaults.find((d) => d.i === item.i);
      const sizing = widgetSizing[item.i];

      const baseMinW = defaultItem?.w ?? sizing?.w ?? 1;
      const baseMinH = defaultItem?.h ?? sizing?.h ?? 1;

      const minW = Math.min(Math.max(baseMinW, 1), maxCols);
      const minH = Math.max(baseMinH, 1);

      return {
        ...item,
        w: Math.max(item.w, minW),
        h: Math.max(item.h, minH),
        minW,
        minH,
      };
    });
  }

  return constrained;
};
