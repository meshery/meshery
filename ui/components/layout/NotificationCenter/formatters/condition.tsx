import React from 'react';
import { Box, Chip, CustomTooltip, FormattedTime, Typography } from '@sistent/sistent';

// An event that reports an ongoing condition (meshery/schemas event
// correlation contract) counts its occurrences on one event instead of
// recording one each: occurrenceCount, firstSeenAt and lastSeenAt describe
// them, and resolvedAt is set once the condition cleared. Every field is
// optional; an event without them is an ordinary, single occurrence.

export type EventCondition = {
  occurrenceCount: number;
  firstSeenAt?: string;
  lastSeenAt?: string;
  resolvedAt?: string;
};

type ConditionFields = {
  correlationKey?: unknown;
  occurrenceCount?: unknown;
  firstSeenAt?: unknown;
  lastSeenAt?: unknown;
  resolvedAt?: unknown;
};

const timestamp = (value: unknown): string | undefined =>
  typeof value === 'string' && value !== '' && !Number.isNaN(Date.parse(value)) ? value : undefined;

/**
 * Reads an event's condition fields. Returns null for an event that reports
 * no condition: one with no correlation key, no count above one and no
 * resolution, which is how every event recorded before correlation reads.
 */
export const getEventCondition = (event: ConditionFields | null | undefined) => {
  if (!event) return null;
  const count =
    typeof event.occurrenceCount === 'number' && event.occurrenceCount >= 1
      ? Math.floor(event.occurrenceCount)
      : 1;
  const resolvedAt = timestamp(event.resolvedAt);
  const hasKey = typeof event.correlationKey === 'string' && event.correlationKey !== '';
  if (!hasKey && count <= 1 && !resolvedAt) return null;
  return {
    occurrenceCount: count,
    firstSeenAt: timestamp(event.firstSeenAt),
    lastSeenAt: timestamp(event.lastSeenAt),
    resolvedAt,
  } satisfies EventCondition;
};

const occurrenceLabel = (count: number) => (count === 1 ? 'Once' : `${count} times`);

/**
 * The compact badges a notification's summary row shows for a condition: the
 * occurrence count when it repeated, and whether it is resolved.
 */
export const ConditionBadges = ({ condition }: { condition: EventCondition | null }) => {
  if (!condition) return null;
  const repeated = condition.occurrenceCount > 1;
  if (!repeated && !condition.resolvedAt) return null;
  return (
    <Box sx={{ display: 'flex', gap: 0.5, alignItems: 'center' }} data-testid="condition-badges">
      {repeated && (
        <CustomTooltip
          title={`Occurred ${occurrenceLabel(condition.occurrenceCount).toLowerCase()}`}
          placement="top"
        >
          <Chip
            size="small"
            variant="outlined"
            label={`×${condition.occurrenceCount}`}
            sx={{ height: 20 }}
          />
        </CustomTooltip>
      )}
      {condition.resolvedAt && (
        <Chip size="small" color="success" label="Resolved" sx={{ height: 20 }} />
      )}
    </Box>
  );
};

const Fact = ({ label, children }: { label: string; children: React.ReactNode }) => (
  <Box sx={{ display: 'flex', gap: 0.5, alignItems: 'baseline' }}>
    <Typography variant="body2" sx={{ fontWeight: 600 }}>
      {label}
    </Typography>
    <Typography variant="body2" component="span">
      {children}
    </Typography>
  </Box>
);

/**
 * The condition line of an expanded notification: how often it occurred, when
 * it was first and last seen, and when it was resolved.
 */
export const ConditionSummary = ({ condition }: { condition: EventCondition | null }) => {
  if (!condition) return null;
  return (
    <Box
      data-testid="condition-summary"
      sx={{ display: 'flex', flexWrap: 'wrap', columnGap: 2, rowGap: 0.5, mb: 1 }}
    >
      <Fact label="Occurred">{occurrenceLabel(condition.occurrenceCount)}</Fact>
      {condition.firstSeenAt && condition.occurrenceCount > 1 && (
        <Fact label="First seen">
          <FormattedTime date={condition.firstSeenAt} />
        </Fact>
      )}
      {condition.lastSeenAt && (
        <Fact label="Last seen">
          <FormattedTime date={condition.lastSeenAt} />
        </Fact>
      )}
      {condition.resolvedAt && (
        <Fact label="Resolved">
          <FormattedTime date={condition.resolvedAt} />
        </Fact>
      )}
    </Box>
  );
};
