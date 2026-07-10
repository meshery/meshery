import { IconButton, styled, Tabs } from '@sistent/sistent';

/**
 * Full-viewport, click-through host for the floating sessions panel.
 *
 * Sistent's `Panel` positions itself absolutely, so it needs a positioned
 * ancestor; a fixed, inset-0 host makes it float over the whole app and lets it
 * be dragged anywhere. The host ignores pointer events so it never swallows
 * clicks meant for the page — the panel itself re-enables them.
 *
 * `$hidden` collapses the host without unmounting the panel, because `Panel`
 * renders null when closed and that would take every live shell with it.
 */
export const SessionsPanelHost = styled('div', {
  shouldForwardProp: (prop) => prop !== '$hidden',
})<{ $hidden?: boolean }>(({ $hidden }) => ({
  position: 'fixed',
  inset: 0,
  pointerEvents: 'none',
  zIndex: 1200,
  visibility: $hidden ? 'hidden' : 'visible',
}));

/** The restore affordance shown while the sessions panel is minimized. */
export const MinimizedSessionsButton = styled(IconButton)(({ theme }) => ({
  position: 'fixed',
  bottom: '1rem',
  right: '1rem',
  zIndex: 1201,
  backgroundColor: theme.palette.background.paper,
  boxShadow: theme.shadows[4],
  '&:hover': {
    backgroundColor: theme.palette.action.hover,
  },
}));

/**
 * The element a session renders into, in both normal and fullscreen states.
 *
 * Fullscreen is a style change on this one element rather than a different
 * wrapper: moving the session in the element tree would make React unmount and
 * remount it, killing the live shell or log stream underneath.
 */
export const SessionHost = styled('div', {
  shouldForwardProp: (prop) => prop !== '$fullScreen',
})<{ $fullScreen?: boolean }>(({ theme, $fullScreen }) => ({
  position: 'relative',
  display: 'flex',
  flexDirection: 'column',
  flex: 1,
  minHeight: 0,
  ...($fullScreen && {
    position: 'fixed',
    inset: 0,
    zIndex: theme.zIndex.modal,
    backgroundColor: theme.palette.background.paper,
    padding: '0.75rem',
  }),
}));

/**
 * The dark, monospaced surface a terminal or log stream renders onto.
 *
 * `position: relative` anchors the connecting overlay; `min-height: 0` lets the
 * surface shrink inside a flex column so xterm's fit addon can measure it
 * rather than the pane growing without bound.
 */
export const SessionSurface = styled('div')(({ theme }) => ({
  position: 'relative',
  flex: 1,
  minHeight: 0,
  width: '100%',
  overflow: 'hidden',
  backgroundColor: theme.palette.background.default,
  border: `1px solid ${theme.palette.divider}`,
  borderRadius: '4px',
  // Just enough to keep glyphs off the border. More than this and the terminal
  // stops looking like a terminal.
  padding: '0.25rem 0.4rem',
}));

/** The element xterm opens onto; it must fill its surface for fit() to measure. */
export const TerminalMount = styled('div')({
  width: '100%',
  height: '100%',
});

/**
 * The single control row above a session.
 *
 * One row, not two: the session state used to live in a footer of its own, which
 * cost a line of terminal height and sat on top of the panel's resize grip.
 */
export const SessionToolbar = styled('div')({
  display: 'flex',
  flexWrap: 'nowrap',
  alignItems: 'center',
  gap: '0.4rem',
  paddingBottom: '0.4rem',
  minHeight: '2.25rem',
});

/**
 * Session state, inline in the toolbar. It yields its width first — the controls
 * beside it are interactive and must never be pushed off the end of a panel the
 * user has narrowed.
 */
export const StatusText = styled('span', {
  shouldForwardProp: (prop) => prop !== '$error',
})<{ $error?: boolean }>(({ theme, $error }) => ({
  flex: '1 1 auto',
  minWidth: 0,
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap',
  paddingInline: '0.25rem',
  fontSize: '0.75rem',
  color: $error ? theme.palette.error.main : theme.palette.text.secondary,
}));

/** Right-aligned toolbar slot the host fills, e.g. with a fullscreen control. */
export const ToolbarEnd = styled('div')({
  display: 'flex',
  alignItems: 'center',
  flexShrink: 0,
  marginLeft: 'auto',
});

/** The session tab strip: dense, so more tabs fit before it has to scroll. */
export const SessionTabs = styled(Tabs)(({ theme }) => ({
  minHeight: '2rem',
  borderBottom: `1px solid ${theme.palette.divider}`,
  '& .MuiTab-root': {
    minHeight: '2rem',
    minWidth: 0,
    maxWidth: '16rem',
    padding: '0 0.5rem',
    textTransform: 'none',
    fontSize: '0.8rem',
    gap: '0.35rem',
  },
  '& .MuiTabs-scrollButtons.Mui-disabled': {
    // Reserving space for arrows that can never fire wastes width the tab
    // labels need, and reads as a broken control.
    display: 'none',
  },
}));

/** A tab's contents: title, then its close control. */
export const TabLabel = styled('span')({
  display: 'flex',
  alignItems: 'center',
  gap: '0.15rem',
  minWidth: 0,
  maxWidth: '100%',
});

/**
 * A tab's title. Truncates with an ellipsis rather than reversing text direction
 * to clip the head: pod names are full of hyphens, which are bidi-neutral and
 * reorder visibly under `direction: rtl`. Dropping the "Shell: "/"Logs: " prefix
 * bought back enough width that most names now fit whole anyway.
 */
export const TabTitle = styled('span')({
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap',
});
