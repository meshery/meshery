import React, { createContext, useCallback, useContext, useEffect, useMemo, useState } from 'react';
import {
  ArticleIcon,
  Box,
  CloseIcon,
  IconButton,
  Tab,
  TerminalIcon,
  Tooltip,
} from '@sistent/sistent';
import type { ConsoleKind, ConsoleTarget } from 'lib/console/protocol';
import ConsolePanel from './ConsolePanel';
import ConsoleShell, { type ConsolePanelMode } from './ConsoleShell';
import { bindDock, publishDockState } from './dock-store';
import { ConsoleTabs, TabLabel, TabTitle } from './console-styles';

export interface OpenConsoleRequest {
  connectionId: string;
  target: ConsoleTarget;
  kind: ConsoleKind;
}

interface OpenConsole extends OpenConsoleRequest {
  id: string;
  label: string;
}

export interface ConsoleContextValue {
  /** Opens a console in the shared drawer, or focuses it if already open. */
  openConsole: (request: OpenConsoleRequest) => void;
  /** Closes a console, tearing its socket down. */
  closeConsole: (request: OpenConsoleRequest) => void;
  /**
   * Whether a console for this exact target is already open. Callers use it to
   * render an indicator on a resource that already has a shell attached.
   */
  hasConsole: (request: OpenConsoleRequest) => boolean;
  /** Every console currently open in the drawer. */
  getOpenConsoles: () => OpenConsoleRequest[];
  /**
   * Observes the set of open consoles. The listener fires immediately with the
   * current set and again on every change; the returned function unsubscribes.
   *
   * This exists so a non-React consumer — Kanvas mirrors it into its own store to
   * badge graph nodes that already have a console — can stay in sync with the
   * drawer as the source of truth, rather than tracking what it opened and
   * silently drifting when the user closes a console from the drawer itself.
   */
  subscribe: (listener: (open: OpenConsoleRequest[]) => void) => () => void;
}

const ConsolesContext = createContext<ConsoleContextValue | null>(null);

const NO_CONSOLES: ConsoleContextValue = {
  openConsole: () => {},
  closeConsole: () => {},
  hasConsole: () => false,
  getOpenConsoles: () => [],
  subscribe: () => () => {},
};

/**
 * The imperative handle on the shared console panel.
 *
 * This is the API surface remote components (the Kanvas operator view) use to
 * launch a shell or a log tail from a graph node, so it is exported through
 * `remote-component.config.js` alongside the components themselves.
 *
 * Outside a provider it degrades to no-ops rather than throwing, so a component
 * that merely *offers* a console action can be dropped anywhere without the
 * tree it lands in having to know about consoles.
 */
export const useConsole = (): ConsoleContextValue => useContext(ConsolesContext) ?? NO_CONSOLES;

/**
 * The mounted provider's handle, for callers that are not React components.
 *
 * Kanvas dispatches its context-menu actions from plain functions, not from a
 * component, so it cannot read the context. This module-scope handle is bound
 * while a provider is mounted and reset when it unmounts, so a stale reference
 * can never resurrect a dead drawer.
 */
let mountedConsoles: ConsoleContextValue | null = null;

/** The open-console set, mirrored outside React so non-React callers can read it. */
let openConsoleRequests: OpenConsoleRequest[] = [];
const consoleListeners = new Set<(open: OpenConsoleRequest[]) => void>();

const publishOpenConsoles = (open: OpenConsoleRequest[]) => {
  openConsoleRequests = open;
  consoleListeners.forEach((listener) => listener(open));
};

/**
 * Subscribes to the open-console set. Defined at module scope rather than as a
 * closure over the provider, so both the imperative handle and the context value
 * hand out the same stable function — and so it does not collide with the
 * provider's own `consoles` state variable.
 */
const subscribeToConsoles = (listener: (open: OpenConsoleRequest[]) => void) => {
  consoleListeners.add(listener);
  // Fire immediately so a subscriber that attaches after consoles were opened
  // does not have to wait for the next change to learn about them.
  listener(openConsoleRequests);
  return () => {
    consoleListeners.delete(listener);
  };
};

/**
 * Imperative console control for non-React callers. Falls back to no-ops when no
 * provider is mounted, so an extension loading before the app tree is ready
 * degrades rather than throwing.
 */
export const consoles: ConsoleContextValue = {
  openConsole: (request) => (mountedConsoles ?? NO_CONSOLES).openConsole(request),
  closeConsole: (request) => (mountedConsoles ?? NO_CONSOLES).closeConsole(request),
  hasConsole: (request) => (mountedConsoles ?? NO_CONSOLES).hasConsole(request),
  getOpenConsoles: () => openConsoleRequests,
  subscribe: subscribeToConsoles,
};

/** Identifies a console, so re-requesting the same one focuses it rather than duplicating it. */
const consoleId = ({ connectionId, target, kind }: OpenConsoleRequest): string =>
  [kind, connectionId, target.namespace ?? '', target.name, target.container ?? ''].join('/');

/**
 * The tab's label. Just the resource — the kind is carried by the tab's icon,
 * because a "Shell: " / "Logs: " prefix on every tab spent a third of the strip
 * repeating what the icon already says and pushed the strip into scroll.
 */
const consoleLabel = ({ target }: OpenConsoleRequest): string =>
  target.container ? `${target.name}/${target.container}` : target.name;

/**
 * Hosts the consoles panel and hands its opener to the tree below.
 *
 * The panel docks to the foot of the content area, where it shares space with the
 * page rather than covering it, and detaches into a free-floating window that can
 * be dragged and resized anywhere over that area. Minimizing hides it and leaves
 * a badged control at the foot of the Navigator. The panel is mounted app-wide,
 * so a console survives navigation in any of those states.
 *
 * Every open console stays mounted for as long as it is open — including the
 * ones behind an inactive tab, and including while the panel is minimized. A
 * terminal is stateful on the remote end, so unmounting it to switch tabs, to
 * minimize, or to move between docked and floating would kill the shell. Hidden
 * panes are hidden with CSS, the shell changes style rather than structure, and
 * xterm refits from its ResizeObserver when a pane is shown or resized again.
 */
export const ConsoleProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [consoles, setConsoles] = useState<OpenConsole[]>([]);
  const [activeId, setActiveId] = useState<string | null>(null);
  const [minimized, setMinimized] = useState(false);
  const [mode, setMode] = useState<ConsolePanelMode>('docked');

  const openConsole = useCallback((request: OpenConsoleRequest) => {
    const id = consoleId(request);
    setConsoles((current) =>
      current.some((entry) => entry.id === id)
        ? current
        : [...current, { ...request, id, label: consoleLabel(request) }],
    );
    setActiveId(id);
    // Opening a console while minimized should show it, not silently queue it.
    setMinimized(false);
  }, []);

  const closeConsoleById = useCallback((id: string) => {
    setConsoles((current) => {
      const remaining = current.filter((entry) => entry.id !== id);
      setActiveId((active) => {
        if (active !== id) return active;
        return remaining.length ? remaining[remaining.length - 1].id : null;
      });
      return remaining;
    });
  }, []);

  const closeConsole = useCallback(
    (request: OpenConsoleRequest) => closeConsoleById(consoleId(request)),
    [closeConsoleById],
  );

  // Reads `consoles` and so must not be memoized against an empty dep list, or
  // callers would see a stale answer after a console opened.
  const hasConsole = useCallback(
    (request: OpenConsoleRequest) => {
      const id = consoleId(request);
      return consoles.some((entry) => entry.id === id);
    },
    [consoles],
  );

  const closeAll = useCallback(() => {
    setConsoles([]);
    setActiveId(null);
    setMinimized(false);
  }, []);

  const openRequests = useMemo(
    () => consoles.map(({ connectionId, target, kind }) => ({ connectionId, target, kind })),
    [consoles],
  );

  const getOpenConsoles = useCallback(() => openRequests, [openRequests]);

  const value = useMemo(
    () => ({
      openConsole,
      closeConsole,
      hasConsole,
      getOpenConsoles,
      subscribe: subscribeToConsoles,
    }),
    [openConsole, closeConsole, hasConsole, getOpenConsoles],
  );

  // Publish to the non-React mirror whenever the set changes.
  useEffect(() => {
    publishOpenConsoles(openRequests);
  }, [openRequests]);

  // Feed the Navigator's restore control, which sits outside this context.
  useEffect(() => {
    publishDockState({ count: consoles.length, minimized });
  }, [consoles.length, minimized]);

  useEffect(() => bindDock(() => setMinimized(false)), []);

  // Bind the imperative handle to *this* provider for as long as it is mounted.
  // The guard on unbind matters under StrictMode's double-mount, where the new
  // provider binds before the old one's cleanup runs.
  useEffect(() => {
    mountedConsoles = value;
    return () => {
      if (mountedConsoles === value) mountedConsoles = null;
    };
  }, [value]);
  const open = consoles.length > 0;

  return (
    <ConsolesContext.Provider value={value}>
      {children}

      <ConsoleShell
        open={open}
        minimized={minimized}
        mode={mode}
        onModeChange={setMode}
        onMinimize={() => setMinimized(true)}
        onClose={closeAll}
        title="Consoles"
        tabs={
          <ConsoleTabs
            value={activeId ?? false}
            onChange={(_, id) => setActiveId(id)}
            variant="scrollable"
            scrollButtons="auto"
            allowScrollButtonsMobile
          >
            {consoles.map((entry) => (
              <Tab
                key={entry.id}
                value={entry.id}
                disableRipple
                icon={
                  entry.kind === 'terminal' ? (
                    <TerminalIcon width={15} height={15} />
                  ) : (
                    <ArticleIcon width={15} height={15} />
                  )
                }
                iconPosition="start"
                label={
                  <TabLabel>
                    <Tooltip title={entry.label}>
                      <TabTitle>{entry.label}</TabTitle>
                    </Tooltip>
                    <IconButton
                      size="small"
                      component="span"
                      aria-label={`Close ${entry.label}`}
                      onClick={(event) => {
                        // The close control lives inside the tab, so stop the
                        // click from also selecting the tab being closed.
                        event.stopPropagation();
                        closeConsoleById(entry.id);
                      }}
                    >
                      <CloseIcon width={12} height={12} />
                    </IconButton>
                  </TabLabel>
                }
              />
            ))}
          </ConsoleTabs>
        }
      >
        {consoles.map((entry) => (
          <Box
            key={entry.id}
            sx={{
              flex: 1,
              minHeight: 0,
              // Hidden rather than unmounted: see the note on this component.
              display: entry.id === activeId ? 'flex' : 'none',
              flexDirection: 'column',
              paddingTop: '0.5rem',
            }}
          >
            <ConsolePanel
              connectionId={entry.connectionId}
              target={entry.target}
              kind={entry.kind}
              // Every console stays `active` — and so keeps its socket — while
              // only the one on show owns the header's controls.
              focused={entry.id === activeId}
            />
          </Box>
        ))}
      </ConsoleShell>
    </ConsolesContext.Provider>
  );
};

export default ConsoleProvider;
