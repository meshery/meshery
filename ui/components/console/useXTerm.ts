import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { alpha, useTheme } from '@sistent/sistent';
import type { Terminal } from '@xterm/xterm';
import type { FitAddon } from '@xterm/addon-fit';
import type { SearchAddon } from '@xterm/addon-search';
import { resolveTerminalFontFamily } from './terminal-font';

export interface XTermGeometry {
  cols: number;
  rows: number;
}

export interface UseXTermOptions {
  /** A read-only terminal ignores keystrokes; used for log panes. */
  readOnly?: boolean;
  /** Called with typed input. Never called when `readOnly`. */
  onInput?: (data: string) => void;
  /** Called after a resize settles, with the new geometry. */
  onResize?: (geometry: XTermGeometry) => void;
  /** Scrollback lines to retain. Logs want far more than a shell does. */
  scrollback?: number;
}

export interface UseXTerm {
  /** Attach to the element that should host the terminal. */
  containerRef: (node: HTMLDivElement | null) => void;
  /** Writes bytes verbatim; xterm decodes UTF-8, including split sequences. */
  write: (chunk: Uint8Array) => void;
  /** Writes a line of text, for local notices rather than remote output. */
  writeln: (text: string) => void;
  clear: () => void;
  /** Re-measures and reflows. Call when the pane's size changes. */
  fit: () => void;
  search: (query: string, forward?: boolean) => void;
  /** The terminal's current geometry, or null before it exists. */
  getGeometry: () => XTermGeometry | null;
  /** True once the terminal has been created and attached. */
  ready: boolean;
}

/**
 * Creates and owns an xterm terminal bound to a container element.
 *
 * xterm touches `window` at import time, so it is loaded lazily inside an
 * effect rather than at module scope: Next.js renders this tree on the server
 * during the initial pass, and a static import would crash the build.
 */
export const useXTerm = ({
  readOnly = false,
  onInput,
  onResize,
  scrollback = 5000,
}: UseXTermOptions = {}): UseXTerm => {
  const theme = useTheme();
  const [ready, setReady] = useState(false);

  const terminalRef = useRef<Terminal | null>(null);
  const fitAddonRef = useRef<FitAddon | null>(null);
  const searchAddonRef = useRef<SearchAddon | null>(null);
  const containerNodeRef = useRef<HTMLDivElement | null>(null);
  const resizeObserverRef = useRef<ResizeObserver | null>(null);

  // Callbacks live in refs: xterm's listeners are registered once, at creation,
  // and must see the latest handler without the terminal being rebuilt.
  const onInputRef = useRef(onInput);
  onInputRef.current = onInput;
  const onResizeRef = useRef(onResize);
  onResizeRef.current = onResize;

  // xterm renders to a canvas, so it needs resolved colour strings rather than
  // CSS. They are derived from the palette so the terminal tracks the app's
  // theme instead of pinning its own.
  const xtermTheme = useMemo(
    () => ({
      background: theme.palette.background.default,
      foreground: theme.palette.text.primary,
      // A read-only log pane has no caret to place, and a blinking one there
      // reads as an editable field.
      cursor: readOnly ? 'transparent' : theme.palette.primary.main,
      cursorAccent: theme.palette.background.default,
      selectionBackground: alpha(theme.palette.primary.main, 0.35),
    }),
    [
      theme.palette.background.default,
      theme.palette.text.primary,
      theme.palette.primary.main,
      readOnly,
    ],
  );

  const containerRef = useCallback((node: HTMLDivElement | null) => {
    containerNodeRef.current = node;
  }, []);

  useEffect(() => {
    let disposed = false;

    const create = async () => {
      const [{ Terminal: XTerminal }, { FitAddon: XFitAddon }, { SearchAddon: XSearchAddon }] =
        await Promise.all([
          import('@xterm/xterm'),
          import('@xterm/addon-fit'),
          import('@xterm/addon-search'),
        ]);

      const container = containerNodeRef.current;
      // The effect may resolve after the pane unmounted, or before the ref was
      // attached; either way there is nothing to open onto.
      if (disposed || !container) return;

      const terminal = new XTerminal({
        cursorBlink: !readOnly,
        disableStdin: readOnly,
        convertEol: readOnly,
        // Probed at runtime, not hardcoded: a stack whose head family is
        // absent resolves to a proportional substitute on Linux. See terminal-font.ts.
        fontFamily: resolveTerminalFontFamily(),
        fontSize: 12,
        // xterm lays out on a fixed character grid it measures from the font. Any
        // extra tracking the page's CSS would apply desynchronises the painted
        // glyphs from that grid, so pin both.
        letterSpacing: 0,
        lineHeight: 1.2,
        scrollback,
        theme: xtermTheme,
      });

      const fitAddon = new XFitAddon();
      const searchAddon = new XSearchAddon();
      terminal.loadAddon(fitAddon);
      terminal.loadAddon(searchAddon);
      terminal.open(container);
      fitAddon.fit();

      if (!readOnly) {
        terminal.onData((data) => onInputRef.current?.(data));
      }
      terminal.onResize(({ cols, rows }) => onResizeRef.current?.({ cols, rows }));

      // The pane is resizable and lives inside a drawer and a tab panel, so a
      // window listener would miss most of the resizes that matter.
      const observer = new ResizeObserver(() => {
        try {
          fitAddon.fit();
        } catch {
          // fit() throws if the container has collapsed to zero size, which is
          // what a hidden tab looks like. The next observation will refit.
        }
      });
      observer.observe(container);

      terminalRef.current = terminal;
      fitAddonRef.current = fitAddon;
      searchAddonRef.current = searchAddon;
      resizeObserverRef.current = observer;
      setReady(true);
    };

    void create();

    return () => {
      disposed = true;
      resizeObserverRef.current?.disconnect();
      resizeObserverRef.current = null;
      terminalRef.current?.dispose();
      terminalRef.current = null;
      fitAddonRef.current = null;
      searchAddonRef.current = null;
      setReady(false);
    };
    // Recreating on theme change would drop scrollback; the theme is applied
    // live by the effect below instead.
  }, [readOnly, scrollback]);

  useEffect(() => {
    if (terminalRef.current) {
      terminalRef.current.options.theme = xtermTheme;
    }
  }, [xtermTheme]);

  const write = useCallback((chunk: Uint8Array) => terminalRef.current?.write(chunk), []);
  const writeln = useCallback((text: string) => terminalRef.current?.writeln(text), []);
  const clear = useCallback(() => terminalRef.current?.clear(), []);
  const fit = useCallback(() => {
    try {
      fitAddonRef.current?.fit();
    } catch {
      // See above: a collapsed container cannot be measured.
    }
  }, []);
  const search = useCallback((query: string, forward = true) => {
    if (!query) return;
    if (forward) searchAddonRef.current?.findNext(query);
    else searchAddonRef.current?.findPrevious(query);
  }, []);
  const getGeometry = useCallback((): XTermGeometry | null => {
    const terminal = terminalRef.current;
    return terminal ? { cols: terminal.cols, rows: terminal.rows } : null;
  }, []);

  return { containerRef, write, writeln, clear, fit, search, getGeometry, ready };
};
