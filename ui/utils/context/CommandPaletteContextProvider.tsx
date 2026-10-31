import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from 'react';
import CommandPalette from '@/components/command-palette/CommandPalette';

export type CommandPaletteContextValue = {
  open: boolean;
  openPalette: () => void;
  closePalette: () => void;
  togglePalette: () => void;
};

export const CommandPaletteContext = createContext<CommandPaletteContextValue>({
  open: false,
  openPalette: () => {},
  closePalette: () => {},
  togglePalette: () => {},
});

export const CommandPaletteContextProvider = ({ children }: { children: ReactNode }) => {
  const [open, setOpen] = useState(false);

  const openPalette = useCallback(() => {
    setOpen(true);
  }, []);

  const closePalette = useCallback(() => {
    setOpen(false);
  }, []);

  const togglePalette = useCallback(() => {
    setOpen((prev) => !prev);
  }, []);

  const value = useMemo(
    () => ({
      open,
      openPalette,
      closePalette,
      togglePalette,
    }),
    [open, openPalette, closePalette, togglePalette],
  );

  return <CommandPaletteContext.Provider value={value}>{children}</CommandPaletteContext.Provider>;
};

export const CommandPaletteHost = () => {
  const { togglePalette } = useContext(CommandPaletteContext);

  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.defaultPrevented || event.isComposing || event.repeat) {
        return;
      }

      if (
        (event.ctrlKey || event.metaKey) &&
        !event.altKey &&
        !event.shiftKey &&
        event.key?.toLowerCase() === 'k'
      ) {
        event.preventDefault();
        togglePalette();
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => {
      window.removeEventListener('keydown', handleKeyDown);
    };
  }, [togglePalette]);

  // Always rendered so Dialog exit transitions play
  return <CommandPalette />;
};

export const useCommandPalette = () => {
  return useContext(CommandPaletteContext);
};

export default CommandPaletteContextProvider;
