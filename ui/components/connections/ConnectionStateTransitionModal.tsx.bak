import React, { forwardRef, useImperativeHandle, useRef, useState } from 'react';
import { useSelector } from 'react-redux';
import {
  Box,
  CustomTooltip,
  IconButton,
  InfoOutlinedIcon,
  Modal,
  ModalBody,
  ModalButtonDanger,
  ModalButtonPrimary,
  ModalButtonSecondary,
  ModalFooter,
  Typography,
  WarningIcon,
  useTheme,
} from '@sistent/sistent';
import { iconLarge } from '../../css/icons.styles';
import {
  CONNECTION_KINDS,
  CONNECTION_STATES,
  CONNECTION_STATE_TO_TRANSITION_MAP,
} from '../../utils/Enum';
import { formatToTitleCase } from '../../utils/utils';
import {
  CONNECTION_DOCS_URL,
  getTransitionDescription,
  type ConnectionTransitionMap,
} from './ConnectionTable.constants';

export const KUBERNETES_CONNECTION_LIFECYCLE_DOCS_URL =
  'https://docs.meshery.io/guides/infrastructure-management/kubernetes-connection-lifecycle';

// Transition ramifications are NOT hardcoded here. The connection definition
// for the respective kind (models/.../connections/*.json, surfaced on
// `state.ui.connectionMetadataState[kind].transitionMap` via the connection
// definitions API) is the single source of truth for what each state
// transition means; this modal only resolves and renders it.
type ConnectionMetadataState = Record<
  string,
  { transitionMap?: ConnectionTransitionMap } | undefined
> | null;

const getDocsTooltipMarkdown = (kind: string | undefined): string =>
  kind === CONNECTION_KINDS.KUBERNETES
    ? `Every connection moves through a defined lifecycle of states. Learn more about the [Kubernetes connection lifecycle](${KUBERNETES_CONNECTION_LIFECYCLE_DOCS_URL}) and the [behavior of state transitions](${CONNECTION_DOCS_URL}) in Meshery Docs.`
    : `Every connection moves through a defined lifecycle of states. Learn more about the [lifecycle of connections and the behavior of state transitions](${CONNECTION_DOCS_URL}) in Meshery Docs.`;

/**
 * Data-quality guard for definition-authored transition copy. Registries can
 * still hold rows from older definition versions whose descriptions were
 * prompt-style ("Are you sure… transition from X to Y?") — that phrasing
 * restates the modal title instead of explaining consequences, so suppress it
 * and fall back to the generic line rather than render it twice.
 */
export const shouldShowTransitionDescription = (
  description: string | undefined,
  options: { currentStatus?: string; targetStatus: string },
): description is string => {
  if (!description?.trim()) {
    return false;
  }

  const normalized = description.trim().toLowerCase();
  // Prompt-style leftovers from older connection definition versions.
  if (normalized.startsWith('are you sure')) {
    return false;
  }

  const { currentStatus, targetStatus } = options;
  if (currentStatus) {
    const fromTo = `from ${currentStatus.toLowerCase()} to ${targetStatus.toLowerCase()}`;
    if (normalized === fromTo || normalized.includes(`transition ${fromTo}`)) {
      return false;
    }
  }

  return true;
};

export interface TransitioningConnection {
  id?: string;
  name?: string;
  /**
   * The connection's current lifecycle state, when known. Lets bulk flows
   * resolve the definition-authored description per connection; the copy is
   * shown only when every selected connection resolves to the same one.
   */
  status?: string;
}

export interface ConnectionStateTransitionShowParams {
  /** The lifecycle state being transitioned to, e.g. `deleted`. */
  targetStatus: string;
  /** The current lifecycle state, when known (single-connection flows). */
  currentStatus?: string;
  /** Connection kind, e.g. `kubernetes`. Selects the connection definition and docs links. */
  kind?: string;
  /** The connection(s) the transition applies to (one for row flows, many for bulk). */
  connections: TransitioningConnection[];
}

export interface ConnectionStateTransitionModalRef {
  /** Resolves `true` when the user confirms, `false` on cancel/dismiss. */
  show: (params: ConnectionStateTransitionShowParams) => Promise<boolean>;
}

const MAX_LISTED_NAMES = 3;

// Sistent ModalButtonDanger sets backgroundColor on its own class, but MUI's
// `.MuiButton-contained` rule wins on specificity and paints brand teal instead
// of error red. Raise specificity here until Sistent ships the same fix.
// Theme callbacks use Sistent palette tokens (no hard-coded colors).
const dangerButtonSx = {
  '&&.MuiButton-contained': {
    backgroundColor: (theme) => theme.palette.background?.error?.default,
    color: (theme) => theme.palette.text?.constant?.white ?? theme.palette.common?.white,
    '&:hover': {
      backgroundColor: (theme) =>
        theme.palette.background?.error?.hover ?? theme.palette.background?.error?.default,
    },
  },
};

// Caution transitions (disconnect / ignore / maintenance) are reversible and
// must not look like delete (red). There is no ModalButtonWarning in Sistent.
//
// Important: ModalButtonPrimary / MUI contained force light (white) label color.
// White-on-amber fails contrast (label looks "invisible"). We raise specificity
// and pin a dark label on the warning fill for both default and hover.
// Prefer true black for WCAG contrast on amber; fall back to inverse/default.
const cautionLabelColor = (theme) =>
  theme.palette.common?.black ?? theme.palette.text?.inverse ?? theme.palette.text?.default;

const cautionFill = (theme) =>
  theme.palette.background?.warning?.default ?? theme.palette.status?.warning;

const cautionFillHover = (theme) =>
  theme.palette.background?.warning?.hover ??
  theme.palette.background?.warning?.default ??
  theme.palette.status?.warning;

const cautionButtonSx = {
  // Match primary + contained + brand class chains MUI emits on this button.
  '&&.MuiButton-contained, &&.MuiButton-containedPrimary': {
    backgroundColor: cautionFill,
    backgroundImage: 'none',
    color: cautionLabelColor,
    // Some builds paint label via -webkit-text-fill-color on the root.
    WebkitTextFillColor: cautionLabelColor,
    '&:hover': {
      backgroundColor: cautionFillHover,
      backgroundImage: 'none',
      color: cautionLabelColor,
      WebkitTextFillColor: cautionLabelColor,
    },
    '& .MuiButton-label, & .MuiButton-startIcon, & .MuiButton-endIcon': {
      color: cautionLabelColor,
      WebkitTextFillColor: cautionLabelColor,
    },
  },
};

/** Lifecycle states that pause or step back management without purging the connection. */
const CAUTION_TARGET_STATES = new Set([
  CONNECTION_STATES.DISCONNECTED,
  CONNECTION_STATES.IGNORED,
  CONNECTION_STATES.MAINTENANCE,
]);

/**
 * Action verb for titles and confirm labels. Reuses the shared transition map
 * so chip labels, dropdowns, and this modal stay consistent (Delete / Disconnect /
 * Discover / …). A few states need a fuller phrase so the title reads as English.
 */
export const getTransitionActionLabel = (targetStatus: string): string => {
  if (targetStatus === CONNECTION_STATES.NOTFOUND) {
    return 'Mark as not found';
  }
  if (targetStatus === CONNECTION_STATES.MAINTENANCE) {
    return 'Put into maintenance';
  }
  return (
    CONNECTION_STATE_TO_TRANSITION_MAP[targetStatus] || formatToTitleCase(targetStatus) || 'Confirm'
  );
};

export const buildTransitionTitle = ({
  targetStatus,
  kind,
  count,
}: {
  targetStatus: string;
  kind?: string;
  count: number;
}): string => {
  const plural = count > 1;
  const kindLabel = kind ? `${formatToTitleCase(kind)} ` : '';
  const connectionWord = `connection${plural ? 's' : ''}`;
  const countLabel = plural ? `${count} ` : '';

  if (targetStatus === CONNECTION_STATES.NOTFOUND) {
    return `Mark ${countLabel}${kindLabel}${connectionWord} as not found?`;
  }
  if (targetStatus === CONNECTION_STATES.MAINTENANCE) {
    return `Put ${countLabel}${kindLabel}${connectionWord} into maintenance?`;
  }

  const action = getTransitionActionLabel(targetStatus);
  return `${action} ${countLabel}${kindLabel}${connectionWord}?`;
};

/**
 * Lists named connections for the lead sentence. Remainder is computed against
 * the full selection size (not only named rows) so "Delete N connections" and
 * "and M more" stay consistent when some rows have empty names.
 */
export const ConnectionNames = ({ connections }: { connections: TransitioningConnection[] }) => {
  const names = connections.map((connection) => connection.name).filter(Boolean) as string[];
  if (names.length === 0) {
    return null;
  }
  const listed = names.slice(0, MAX_LISTED_NAMES);
  const remainder = connections.length - listed.length;
  return (
    <>
      {' '}
      (<b>{listed.join(', ')}</b>
      {remainder > 0 ? ` and ${remainder} more` : ''})
    </>
  );
};

const normalizeStatus = (status: string | undefined): string | undefined =>
  status === undefined ? undefined : status.toLowerCase();

/**
 * The single confirmation experience for every connection state transition -
 * row delete, bulk delete, status dropdown, the header's Kubernetes context
 * switcher, and the connection wizard all funnel through this modal. It
 * explains the ramifications of the selected transition (kind-aware), links
 * the relevant docs from an info tooltip, and uses the Sistent danger button
 * for destructive transitions.
 */
const ConnectionStateTransitionModal = forwardRef<ConnectionStateTransitionModalRef>(
  function ConnectionStateTransitionModal(_props, ref) {
    const theme = useTheme();
    // The respective kind's connection definition — the single source of truth
    // for transition copy — as seeded from the definitions API on app load.
    const connectionMetadataState = useSelector(
      (state: { ui?: { connectionMetadataState?: ConnectionMetadataState } }) =>
        state.ui?.connectionMetadataState,
    );
    const [params, setParams] = useState<ConnectionStateTransitionShowParams | null>(null);
    const promiseRef = useRef<{ resolve: (confirmed: boolean) => void }>({ resolve: () => {} });
    // Retains the last shown params so the modal can keep rendering its
    // content during the dialog's exit transition after params clears.
    const lastParamsRef = useRef<ConnectionStateTransitionShowParams | null>(null);
    if (params) {
      lastParamsRef.current = params;
    }

    useImperativeHandle(ref, () => ({
      show: (showParams) =>
        new Promise<boolean>((resolve) => {
          // Settle any in-flight confirmation as cancelled so a re-entrant
          // show() (e.g. a double-click) can never leave the first caller's
          // promise hanging.
          promiseRef.current.resolve(false);
          promiseRef.current = { resolve };
          // Normalize status case so isDelete / consequence maps stay reliable
          // regardless of caller (table lowercases; wizard historically did not).
          setParams({
            ...showParams,
            targetStatus: showParams.targetStatus.toLowerCase(),
            currentStatus: normalizeStatus(showParams.currentStatus),
          });
        }),
    }));

    const settle = (confirmed: boolean) => {
      setParams(null);
      const { resolve } = promiseRef.current;
      // Reset so a later settle (or re-entrant show) cannot re-invoke it.
      promiseRef.current = { resolve: () => {} };
      resolve(confirmed);
    };

    const displayParams = params ?? lastParamsRef.current;
    if (!displayParams) {
      return null;
    }

    const { targetStatus, currentStatus, kind, connections } = displayParams;
    const isDelete = targetStatus === CONNECTION_STATES.DELETED;
    const isCaution = CAUTION_TARGET_STATES.has(targetStatus);
    const count = connections.length;
    const plural = count > 1;
    const secondaryColor = theme.palette.text.secondary;
    // Resolve the transition's meaning dynamically from the respective kind's
    // connection definition. Bulk selections resolve per connection; the copy
    // is shown only when every selected connection lands on the same one
    // (mixed current states get the generic fallback line instead).
    const transitionMap = kind ? connectionMetadataState?.[kind]?.transitionMap : undefined;
    const resolvedDescriptions = new Set(
      connections.map((connection) =>
        getTransitionDescription(transitionMap, connection.status ?? currentStatus, targetStatus),
      ),
    );
    const definitionDescription =
      resolvedDescriptions.size === 1 ? [...resolvedDescriptions][0] : undefined;
    const transitionDescription = shouldShowTransitionDescription(definitionDescription, {
      currentStatus,
      targetStatus,
    })
      ? definitionDescription
      : undefined;
    const actionLabel = getTransitionActionLabel(targetStatus);
    const title = buildTransitionTitle({ targetStatus, kind, count });
    const connectionPhrase = (
      <>
        {plural ? `${count} connections` : 'the connection'}
        <ConnectionNames connections={connections} />
      </>
    );
    // Lead mirrors the title verb family. Multi-word actions need the object
    // in the middle so English stays grammatical
    // ("mark the connection as not found", not "mark as not found the connection").
    // Kind-specific ramifications (credential removal, operator undeploy, …)
    // are not restated here — they come from the connection definition below.
    const leadSentence = (
      <>
        {targetStatus === CONNECTION_STATES.NOTFOUND ? (
          <>You are about to mark {connectionPhrase} as not found</>
        ) : targetStatus === CONNECTION_STATES.MAINTENANCE ? (
          <>You are about to put {connectionPhrase} into maintenance</>
        ) : (
          <>
            You are about to {actionLabel.toLowerCase()} {connectionPhrase}
          </>
        )}
        .
      </>
    );

    // Delete = danger red. Caution (disconnect/ignore/maintenance) = warning
    // amber with dark text. Forward transitions (connect/register/discover) =
    // primary teal. Matches prior Prompt WARNING vs DANGER intent without a
    // third Sistent button primitive.
    const ConfirmButton = isDelete ? ModalButtonDanger : ModalButtonPrimary;
    const confirmButtonSx = isDelete ? dangerButtonSx : isCaution ? cautionButtonSx : undefined;

    return (
      <Modal
        open={Boolean(params)}
        closeModal={() => settle(false)}
        title={title}
        headerIcon={
          <WarningIcon
            {...iconLarge}
            // Amber via Sistent status/background warning tokens (not white on the header).
            fill={
              theme.palette.status?.warning ??
              theme.palette.background?.warning?.default ??
              theme.palette.warning?.main
            }
          />
        }
        maxWidth="sm"
        data-testid="connection-transition-modal"
      >
        <ModalBody>
          {/*
            Keep the info control inline with the lead copy so it tracks the
            text baseline (flex-row + default IconButton padding sat above/
            beside the sentence, which was most visible on short non-k8s leads
            like Prometheus delete). One layout covers all five entry paths.
          */}
          <Typography
            variant="body1"
            component="div"
            data-testid="connection-transition-lead"
            sx={{ lineHeight: 1.5 }}
          >
            {leadSentence}
            <CustomTooltip title={getDocsTooltipMarkdown(kind)} placement="top">
              <IconButton
                aria-label="Learn more about connection state transitions"
                data-testid="connection-transition-info"
                size="small"
                disableRipple
                onClick={(event) => event.stopPropagation()}
                sx={{
                  display: 'inline-flex',
                  verticalAlign: 'text-bottom',
                  // Tight padding so the 18px glyph sits on the body1 line,
                  // not in a tall hit-target that drifts off the sentence.
                  padding: '0 0 0 0.25rem',
                  marginLeft: '0.125rem',
                  color: 'inherit',
                }}
              >
                <InfoOutlinedIcon height={18} width={18} />
              </IconButton>
            </CustomTooltip>
          </Typography>

          <Box data-testid="connection-transition-ramifications" sx={{ marginTop: '0.5rem' }}>
            {transitionDescription ? (
              <Typography
                variant="body2"
                data-testid="connection-transition-description"
                sx={{ color: secondaryColor }}
              >
                {transitionDescription}
              </Typography>
            ) : (
              // No definition-authored copy resolves for this transition (kind
              // without a definition, unknown current state, or a mixed-state
              // bulk selection) — state only the mechanical effect.
              <Typography
                variant="body2"
                data-testid="connection-transition-fallback"
                sx={{ color: secondaryColor }}
              >
                This will transition {plural ? 'each selected connection' : 'the connection'} to the{' '}
                {formatToTitleCase(targetStatus)} state.
              </Typography>
            )}
            {plural && (
              <Typography
                variant="body2"
                data-testid="connection-transition-bulk-scope"
                sx={{ marginTop: '0.5rem', color: secondaryColor }}
              >
                This applies to each selected connection.
              </Typography>
            )}
          </Box>
        </ModalBody>
        <ModalFooter variant="filled">
          <Box sx={{ width: '100%', display: 'flex', justifyContent: 'flex-end', gap: '1rem' }}>
            <ModalButtonSecondary
              onClick={() => settle(false)}
              data-testid="connection-transition-cancel"
            >
              Cancel
            </ModalButtonSecondary>
            <ConfirmButton
              onClick={() => settle(true)}
              data-testid="connection-transition-confirm"
              data-severity={isDelete ? 'danger' : isCaution ? 'caution' : 'primary'}
              sx={confirmButtonSx}
            >
              {actionLabel}
            </ConfirmButton>
          </Box>
        </ModalFooter>
      </Modal>
    );
  },
);

export default ConnectionStateTransitionModal;
