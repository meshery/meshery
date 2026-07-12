import React, { useCallback, useEffect, useRef, useState } from 'react';
import {
  ArrowDownwardIcon,
  CloseIcon,
  ExpandMoreIcon,
  IconButton,
  OpenInNewIcon,
  Tooltip,
} from '@sistent/sistent';
import {
  ResizeHandle,
  SessionsBody,
  SessionsFrame,
  SessionsHeaderActions,
  SessionsHeaderBar,
  SessionsHost,
  SessionsTitle,
  type ResizeDirection,
} from './session-styles';
import { SessionHeaderSlotContext } from './header-slot';

/** Docked to the foot of the content area, or floating free over it. */
export type SessionsPanelMode = 'docked' | 'floating';

interface FloatRect {
  x: number;
  y: number;
  width: number;
  height: number;
}

const MIN_WIDTH = 360;
const MIN_HEIGHT = 180;
const DEFAULT_DOCK_HEIGHT = 320;

/** Where a panel lands if it is somehow floated without ever having been docked. */
const DEFAULT_FLOAT: FloatRect = { x: 48, y: 48, width: 880, height: 400 };

/** Leaves enough of the page visible that a dock never swallows it whole. */
const MIN_PAGE_HEIGHT = 120;

const FLOATING_HANDLES: ResizeDirection[] = ['n', 's', 'e', 'w', 'ne', 'nw', 'se', 'sw'];
const DOCKED_HANDLES: ResizeDirection[] = ['n'];

const clamp = (value: number, min: number, max: number) =>
  Math.min(Math.max(value, min), Math.max(min, max));

export interface SessionsShellProps {
  /** Whether there is anything to show. A closed shell reserves no space. */
  open: boolean;
  minimized: boolean;
  mode: SessionsPanelMode;
  onModeChange: (mode: SessionsPanelMode) => void;
  onMinimize: () => void;
  onClose: () => void;
  title: string;
  /** The tab strip, rendered flush against the header. */
  tabs?: React.ReactNode;
  children: React.ReactNode;
}

/**
 * The chrome around the session panes: a title bar, dock/detach, minimize, close,
 * and the drag and resize behaviour behind them.
 *
 * Docked and floating are two styles of one element tree, never two trees. A
 * session is a live shell on the other end of a WebSocket, so React unmounting it
 * to move it into a different wrapper would kill it — which rules out swapping
 * between a docked box and Sistent's `Panel` (whose drag/resize is what this
 * otherwise re-implements). Everything here therefore changes CSS, not structure.
 */
const SessionsShell: React.FC<SessionsShellProps> = ({
  open,
  minimized,
  mode,
  onModeChange,
  onMinimize,
  onClose,
  title,
  tabs,
  children,
}) => {
  const hostRef = useRef<HTMLDivElement>(null);
  const frameRef = useRef<HTMLDivElement>(null);
  const [dockHeight, setDockHeight] = useState(DEFAULT_DOCK_HEIGHT);
  const [floatRect, setFloatRect] = useState<FloatRect | null>(null);
  // State, not a ref: the sessions below portal into this element, so they must
  // re-render once it exists.
  const [headerSlot, setHeaderSlot] = useState<HTMLDivElement | null>(null);

  const floating = mode === 'floating';

  /**
   * The content area the panel is confined to. Read from the host's offset parent
   * rather than the window, so a floating panel cannot be dragged under the
   * Navigator or above the header, where its own title bar would be unreachable.
   */
  const area = useCallback(() => {
    const parent = hostRef.current?.offsetParent as HTMLElement | null;
    return {
      width: parent?.clientWidth ?? window.innerWidth,
      height: parent?.clientHeight ?? window.innerHeight,
    };
  }, []);

  /**
   * Detaching keeps the panel where it already is — same height, inset from the
   * edges it was pinned to — so it reads as the dock lifting off the page rather
   * than a new window appearing somewhere else.
   */
  const detach = useCallback(() => {
    const frame = frameRef.current?.getBoundingClientRect();
    const parent = (hostRef.current?.offsetParent as HTMLElement | null)?.getBoundingClientRect();
    if (frame && parent) {
      const width = clamp(frame.width - 96, MIN_WIDTH, Math.max(MIN_WIDTH, parent.width - 32));
      setFloatRect({
        x: clamp(frame.left - parent.left + 48, 0, Math.max(0, parent.width - width)),
        y: Math.max(0, frame.top - parent.top - 24),
        width,
        height: Math.max(MIN_HEIGHT, frame.height),
      });
    }
    onModeChange('floating');
  }, [onModeChange]);

  /** Docking keeps the height the panel was floated at, so the swap is not a jump. */
  const dock = useCallback(() => {
    const height = floatRect?.height ?? dockHeight;
    setDockHeight(clamp(height, MIN_HEIGHT, Math.max(MIN_HEIGHT, area().height - MIN_PAGE_HEIGHT)));
    onModeChange('docked');
  }, [area, dockHeight, floatRect, onModeChange]);

  /**
   * Tracks a pointer until it is released. The listeners go on the window, not the
   * frame: a pointer moving faster than React can re-render leaves the frame
   * behind, and a drag that stops the moment the cursor outruns the panel is worse
   * than no drag at all.
   */
  const track = (onMove: (event: PointerEvent) => void) => {
    const stop = () => {
      window.removeEventListener('pointermove', onMove);
      window.removeEventListener('pointerup', stop);
      window.removeEventListener('pointercancel', stop);
    };
    window.addEventListener('pointermove', onMove);
    window.addEventListener('pointerup', stop);
    window.addEventListener('pointercancel', stop);
  };

  const startDrag = (event: React.PointerEvent) => {
    // The header is the drag handle, but it also carries the panel's buttons and
    // the focused session's container select and search box. Dragging the panel
    // out from under a control the user is trying to click is worse than a drag
    // handle that has holes in it.
    if (!floating || (event.target as HTMLElement).closest('button, [data-no-drag]')) return;

    const origin = floatRect ?? DEFAULT_FLOAT;
    const { clientX, clientY } = event;
    const bounds = area();

    track(({ clientX: x, clientY: y }) =>
      setFloatRect({
        ...origin,
        x: clamp(origin.x + x - clientX, 0, bounds.width - origin.width),
        y: clamp(origin.y + y - clientY, 0, bounds.height - origin.height),
      }),
    );
    event.preventDefault();
  };

  const startResize = (event: React.PointerEvent, dir: ResizeDirection) => {
    const origin = floatRect ?? DEFAULT_FLOAT;
    const originHeight = dockHeight;
    const { clientX, clientY } = event;
    const bounds = area();

    track(({ clientX: x, clientY: y }) => {
      const dx = x - clientX;
      const dy = y - clientY;

      // Docked, the top edge is the only edge that is free to move: the other
      // three are pinned to the content area.
      if (!floating) {
        setDockHeight(
          clamp(
            originHeight - dy,
            MIN_HEIGHT,
            Math.max(MIN_HEIGHT, bounds.height - MIN_PAGE_HEIGHT),
          ),
        );
        return;
      }

      const next = { ...origin };
      if (dir.includes('e')) next.width = Math.max(MIN_WIDTH, origin.width + dx);
      if (dir.includes('s')) next.height = Math.max(MIN_HEIGHT, origin.height + dy);
      // Dragging a top or left edge moves the panel as it resizes it. Deriving the
      // new origin from the *clamped* size keeps the opposite edge pinned once the
      // minimum is hit, instead of letting the panel creep away under the pointer.
      if (dir.includes('w')) {
        next.width = Math.max(MIN_WIDTH, origin.width - dx);
        next.x = origin.x + origin.width - next.width;
      }
      if (dir.includes('n')) {
        next.height = Math.max(MIN_HEIGHT, origin.height - dy);
        next.y = origin.y + origin.height - next.height;
      }
      setFloatRect(next);
    });
    event.preventDefault();
  };

  // A window narrowed after the panel was dragged to the far side would strand it
  // off the edge of the content area, title bar and all, with no way back.
  useEffect(() => {
    if (!floating) return undefined;
    const onResize = () =>
      setFloatRect((rect) => {
        if (!rect) return rect;
        const bounds = area();
        const width = Math.min(rect.width, Math.max(MIN_WIDTH, bounds.width));
        const height = Math.min(rect.height, Math.max(MIN_HEIGHT, bounds.height));
        return {
          width,
          height,
          x: clamp(rect.x, 0, bounds.width - width),
          y: clamp(rect.y, 0, bounds.height - height),
        };
      });
    window.addEventListener('resize', onResize);
    return () => window.removeEventListener('resize', onResize);
  }, [area, floating]);

  const rect = floatRect ?? DEFAULT_FLOAT;

  return (
    <SessionsHost
      ref={hostRef}
      $open={open}
      $floating={floating}
      $minimized={minimized}
      $dockHeight={dockHeight}
      aria-hidden={!open || minimized}
    >
      <SessionsFrame
        ref={frameRef}
        $floating={floating}
        $dockHeight={dockHeight}
        style={
          floating
            ? { left: rect.x, top: rect.y, width: rect.width, height: rect.height }
            : undefined
        }
      >
        <SessionsHeaderBar $floating={floating} onPointerDown={startDrag}>
          <SessionsTitle>{title}</SessionsTitle>

          {/* The focused session portals its container select and search in here. */}
          <SessionsHeaderActions ref={setHeaderSlot} data-no-drag />

          <Tooltip title={floating ? 'Dock to bottom' : 'Detach'}>
            <IconButton
              size="small"
              aria-label={floating ? 'Dock sessions panel' : 'Detach sessions panel'}
              onClick={floating ? dock : detach}
            >
              {floating ? (
                <ArrowDownwardIcon width={14} height={14} fill="currentColor" />
              ) : (
                <OpenInNewIcon width={14} height={14} fill="currentColor" />
              )}
            </IconButton>
          </Tooltip>

          <Tooltip title="Minimize">
            <IconButton size="small" aria-label="Minimize sessions panel" onClick={onMinimize}>
              <ExpandMoreIcon width={16} height={16} fill="currentColor" />
            </IconButton>
          </Tooltip>

          <Tooltip title="Close all sessions">
            <IconButton size="small" aria-label="Close all sessions" onClick={onClose}>
              <CloseIcon width={14} height={14} fill="currentColor" />
            </IconButton>
          </Tooltip>
        </SessionsHeaderBar>

        <SessionsBody>
          <SessionHeaderSlotContext.Provider value={headerSlot}>
            {tabs}
            {children}
          </SessionHeaderSlotContext.Provider>
        </SessionsBody>

        {(floating ? FLOATING_HANDLES : DOCKED_HANDLES).map((dir) => (
          <ResizeHandle key={dir} $dir={dir} onPointerDown={(event) => startResize(event, dir)} />
        ))}
      </SessionsFrame>
    </SessionsHost>
  );
};

export default SessionsShell;
