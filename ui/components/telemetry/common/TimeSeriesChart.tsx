import React, { useMemo } from 'react';
import {
  Area,
  AreaChart,
  CartesianGrid,
  Legend,
  Line,
  LineChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from 'recharts';
import { Box, Typography, useTheme } from '@sistent/sistent';
import type { ChartSeries } from './types';
import { formatValue } from './time';

interface TimeSeriesChartProps {
  series: ChartSeries[];
  unit?: string;
  // explicit pixel height; defaults to filling the parent container
  height?: number | string;
  // when true, render filled areas instead of plain lines
  filled?: boolean;
}

type Row = Record<string, number>;

const formatClockTick = (ts: number) => {
  const d = new Date(ts * 1000);
  return `${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`;
};

const formatFullTime = (ts: number) => new Date(ts * 1000).toLocaleString();

/**
 * TimeSeriesChart renders one or more metric series as a responsive,
 * theme-aware time-series line/area chart using recharts. Series are aligned on
 * a shared timestamp axis so they share one X scale.
 */
const TimeSeriesChart: React.FC<TimeSeriesChartProps> = ({
  series,
  unit,
  height = '100%',
  filled = false,
}) => {
  const theme = useTheme();

  const palette = useMemo(
    () => [
      theme.palette.primary.main,
      theme.palette.success.main,
      theme.palette.warning.main,
      theme.palette.error.main,
      theme.palette.info?.main ?? theme.palette.secondary.main,
      theme.palette.secondary.main,
    ],
    [theme],
  );

  // Align all series on the union of their timestamps so recharts can render
  // them against one shared X scale.
  const data = useMemo<Row[]>(() => {
    const byTs = new Map<number, Row>();
    series.forEach((s, i) => {
      const key = `s${i}`;
      s.points.forEach(([ts, val]) => {
        let row = byTs.get(ts);
        if (!row) {
          row = { ts };
          byTs.set(ts, row);
        }
        if (isFinite(val)) row[key] = val;
      });
    });
    return Array.from(byTs.values()).sort((a, b) => a.ts - b.ts);
  }, [series]);

  const showLegend = series.length > 0 && series.length <= 12;
  const Chart = filled ? AreaChart : LineChart;

  if (data.length === 0) {
    return (
      <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height }}>
        <Typography variant="caption" color="textSecondary">
          No data
        </Typography>
      </Box>
    );
  }

  return (
    <ResponsiveContainer width="100%" height={height}>
      <Chart data={data} margin={{ top: 8, right: 16, bottom: 0, left: 4 }}>
        <CartesianGrid
          strokeDasharray="3 3"
          stroke={theme.palette.border.default}
          vertical={false}
        />
        <XAxis
          dataKey="ts"
          type="number"
          scale="time"
          domain={['dataMin', 'dataMax']}
          tickFormatter={formatClockTick}
          tick={{ fontSize: 11, fill: theme.palette.text.secondary }}
          stroke={theme.palette.border.strong}
          minTickGap={40}
        />
        <YAxis
          width={56}
          tickFormatter={(v: number) => formatValue(v, unit)}
          tick={{ fontSize: 11, fill: theme.palette.text.secondary }}
          stroke={theme.palette.border.strong}
        />
        <Tooltip
          contentStyle={{
            background: theme.palette.background.card,
            border: `1px solid ${theme.palette.border.default}`,
            borderRadius: 6,
            fontSize: 12,
          }}
          labelFormatter={(ts) => formatFullTime(Number(ts))}
          formatter={(value: number, name: string) => [formatValue(Number(value), unit), name]}
        />
        {showLegend && <Legend wrapperStyle={{ fontSize: 11 }} />}
        {series.map((s, i) => {
          const key = `s${i}`;
          const color = palette[i % palette.length];
          return filled ? (
            <Area
              key={key}
              type="monotone"
              dataKey={key}
              name={s.name}
              stroke={color}
              fill={color}
              fillOpacity={0.15}
              strokeWidth={1.5}
              dot={false}
              isAnimationActive={false}
              connectNulls={false}
            />
          ) : (
            <Line
              key={key}
              type="monotone"
              dataKey={key}
              name={s.name}
              stroke={color}
              strokeWidth={1.5}
              dot={false}
              isAnimationActive={false}
              connectNulls={false}
            />
          );
        })}
      </Chart>
    </ResponsiveContainer>
  );
};

// Memoized: recharts is the most expensive thing on the page. Series/unit/filled
// are referentially stable (computed in memoized parents), so this skips
// re-rendering charts when the page re-renders for unrelated reasons (drawer
// open/close, toolbar changes, polling).
export default React.memo(TimeSeriesChart);
