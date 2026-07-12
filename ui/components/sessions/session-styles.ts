import { styled, Tabs } from '@sistent/sistent';

/** Every edge the floating panel can be resized from. Docked, only the top edge. */
export type ResizeDirection = 'n' | 's' | 'e' | 'w' | 'ne' | 'nw' | 'se' | 'sw';

/**
 * The sessions panel's host, in all three of its states.
 *
 * The host is a child of the app's content area (`StyledAppContent`, a relative
 * flex column), which is what lets the *docked* state be a real dock: it takes
 * its height out of the flow, so the page shrinks above it instead of being
 * covered by it. Floating and minimized instead take the host out of the flow —
 * an absolute, click-through overlay of the same content area, so a floating
 * panel is draggable anywhere over the page but never over the Navigator.
 *
 * Minimized reuses the *floating* host geometry on purpose, even for a docked
 * panel. The frame inside stays bottom-anchored at the same height either way, so
 * hiding a docked panel does not change the frame's box — and a terminal that is
 * never resized is a terminal whose buffer is never reflowed behind the user's
 * back. Hiding is `visibility`, not unmounting, because every session must stay
 * mounted while it is minimized or its shell dies with it.
 */
export const SessionsHost = styled('div', {
  shouldForwardProp: (prop) =>
    prop !== '$open' && prop !== '$floating' && prop !== '$minimized' && prop !== '$dockHeight',
})<{ $open: boolean; $floating: boolean; $minimized: boolean; $dockHeight: number }>(
  ({ $open, $floating, $minimized, $dockHeight }) => ({
    // Nothing open reserves no space at all: `visibility` would leave the dock's
    // height carved out of a page that has no sessions to show in it.
    ...(!$open && { display: 'none' }),
    ...($floating || $minimized
      ? {
          position: 'absolute',
          inset: 0,
          // Click-through, so the overlay never swallows clicks meant for the
          // page beneath it. The frame re-enables them for itself.
          pointerEvents: 'none',
          zIndex: 1200,
          ...($minimized && { visibility: 'hidden' }),
        }
      : {
          position: 'relative',
          flex: `0 0 ${$dockHeight}px`,
          height: `${$dockHeight}px`,
        }),
  }),
);

/**
 * The panel itself.
 *
 * Positioned with `left`/`top` rather than a `transform`, which would be the
 * cheaper way to drag it: a transformed ancestor becomes the containing block for
 * `position: fixed` descendants, and a session's fullscreen state is exactly such
 * a descendant. Dragging the panel would silently trap fullscreen inside it.
 */
export const SessionsFrame = styled('div', {
  shouldForwardProp: (prop) => prop !== '$floating' && prop !== '$dockHeight',
})<{ $floating: boolean; $dockHeight: number }>(({ theme, $floating, $dockHeight }) => ({
  position: 'absolute',
  display: 'flex',
  flexDirection: 'column',
  overflow: 'hidden',
  pointerEvents: 'auto',
  backgroundColor: theme.palette.background.paper,
  ...($floating
    ? {
        borderRadius: '6px',
        border: `1px solid ${theme.palette.divider}`,
        boxShadow: theme.shadows[8],
      }
    : {
        left: 0,
        right: 0,
        bottom: 0,
        height: `${$dockHeight}px`,
        borderTop: `1px solid ${theme.palette.divider}`,
      }),
}));

/** The panel's title bar; also the drag handle, once the panel is floating. */
export const SessionsHeaderBar = styled('div', {
  shouldForwardProp: (prop) => prop !== '$floating',
})<{ $floating: boolean }>(({ theme, $floating }) => ({
  display: 'flex',
  alignItems: 'center',
  gap: '0.25rem',
  // Tall enough for the session controls the focused session portals in.
  minHeight: '2.5rem',
  padding: '0.25rem 0.25rem 0.25rem 0.75rem',
  borderBottom: `1px solid ${theme.palette.divider}`,
  backgroundColor: theme.palette.background.default,
  cursor: $floating ? 'move' : 'default',
  // A drag that selects the title text as it goes reads as a broken drag.
  userSelect: 'none',
}));

/**
 * The header's controls slot: the focused session's container select and search
 * box are portaled in here. It takes the slack in the bar, right-aligned, so the
 * panel's own buttons stay pinned to the corner where the user expects them.
 */
export const SessionsHeaderActions = styled('div')({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'flex-end',
  gap: '0.4rem',
  flex: 1,
  minWidth: 0,
  paddingRight: '0.25rem',
});

export const SessionsTitle = styled('span')(({ theme }) => ({
  flex: '0 1 auto',
  minWidth: 0,
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap',
  fontSize: '0.8rem',
  fontWeight: 500,
  color: theme.palette.text.secondary,
}));

/** The panel's body: the tab strip and the session panes below it. */
export const SessionsBody = styled('div')({
  flex: 1,
  minHeight: 0,
  display: 'flex',
  flexDirection: 'column',
  padding: '0.5rem',
});

const EDGE = '6px';
const CORNER = '12px';

const HANDLES: Record<ResizeDirection, object> = {
  n: { top: 0, left: 0, right: 0, height: EDGE, cursor: 'ns-resize' },
  s: { bottom: 0, left: 0, right: 0, height: EDGE, cursor: 'ns-resize' },
  e: { top: 0, bottom: 0, right: 0, width: EDGE, cursor: 'ew-resize' },
  w: { top: 0, bottom: 0, left: 0, width: EDGE, cursor: 'ew-resize' },
  ne: { top: 0, right: 0, width: CORNER, height: CORNER, cursor: 'nesw-resize' },
  nw: { top: 0, left: 0, width: CORNER, height: CORNER, cursor: 'nwse-resize' },
  se: { bottom: 0, right: 0, width: CORNER, height: CORNER, cursor: 'nwse-resize' },
  sw: { bottom: 0, left: 0, width: CORNER, height: CORNER, cursor: 'nesw-resize' },
};

/**
 * A resize grip. Grips sit *inside* the frame, which clips its overflow, and
 * above the body — a terminal fills its pane edge to edge, and a grip outside the
 * frame would be clipped away while a grip below it would never see the pointer.
 */
export const ResizeHandle = styled('div', {
  shouldForwardProp: (prop) => prop !== '$dir',
})<{ $dir: ResizeDirection }>(({ $dir }) => ({
  position: 'absolute',
  zIndex: 2,
  touchAction: 'none',
  ...HANDLES[$dir],
}));

/**
 * The restore affordance at the foot of the Navigator, shown while the panel is
 * minimized. It lives in the Navigator rather than floating over the page so that
 * it does not compete with the app's other floating buttons.
 */
export const NavSessionsButton = styled('button', {
  shouldForwardProp: (prop) => prop !== '$collapsed',
})<{ $collapsed: boolean }>(({ theme, $collapsed }) => ({
  display: 'flex',
  alignItems: 'center',
  justifyContent: $collapsed ? 'center' : 'flex-start',
  gap: '0.75rem',
  width: '100%',
  padding: $collapsed ? '0.5rem 0' : '0.5rem 1.25rem',
  border: 0,
  background: 'transparent',
  cursor: 'pointer',
  fontSize: '0.875rem',
  color: theme.palette.background.constant?.white,
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
