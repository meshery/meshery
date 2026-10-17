import React, { useCallback, useContext, useMemo, useState } from 'react';
import ConnectionWizardModal from '@/components/connections/ConnectionWizardModal';

export type OpenCreateConnectionOptions = {
  /** Pre-select a connection kind once definitions load (e.g. "kubernetes"). */
  kind?: string | null;
  /**
   * When true with `kind`, land on the step after "Choose Connection"
   * (Import Kubeconfig for Kubernetes; details for generic kinds).
   */
  skipKindSelection?: boolean;
};

export type ConnectionWizardContextValue = {
  open: boolean;
  presetKind: string | null;
  skipKindSelection: boolean;
  openCreateConnection: (opts?: OpenCreateConnectionOptions) => void;
  closeCreateConnection: () => void;
};

export const ConnectionWizardContext = React.createContext<ConnectionWizardContextValue>({
  open: false,
  presetKind: null,
  skipKindSelection: false,
  openCreateConnection: () => {},
  closeCreateConnection: () => {},
});

/**
 * App-level Create Connection wizard — same pattern as WorkspaceModalContext /
 * RegistryModalContext. Call sites use `useConnectionWizardModal()`; the host
 * is mounted once under SnackbarProvider in `_app`.
 */
const ConnectionWizardContextProvider = ({ children }: { children: React.ReactNode }) => {
  const [open, setOpen] = useState(false);
  const [presetKind, setPresetKind] = useState<string | null>(null);
  const [skipKindSelection, setSkipKindSelection] = useState(false);

  const openCreateConnection = useCallback((opts?: OpenCreateConnectionOptions) => {
    const kind = opts?.kind?.trim() ? opts.kind.trim() : null;
    setPresetKind(kind);
    setSkipKindSelection(Boolean(kind && opts?.skipKindSelection));
    setOpen(true);
  }, []);

  const closeCreateConnection = useCallback(() => {
    setOpen(false);
    setPresetKind(null);
    setSkipKindSelection(false);
  }, []);

  const value = useMemo(
    () => ({
      open,
      presetKind,
      skipKindSelection,
      openCreateConnection,
      closeCreateConnection,
    }),
    [open, presetKind, skipKindSelection, openCreateConnection, closeCreateConnection],
  );

  return (
    <ConnectionWizardContext.Provider value={value}>{children}</ConnectionWizardContext.Provider>
  );
};

/**
 * Single mount for the Create Connection modal. Must sit under SnackbarProvider
 * (see `_app`) so wizard toasts work. Mirrors Registry: context + thin consumer.
 */
export const ConnectionWizardHost = () => {
  const { open, presetKind, skipKindSelection, closeCreateConnection } =
    useContext(ConnectionWizardContext);

  if (!open) {
    return null;
  }

  return (
    <ConnectionWizardModal
      isOpen={open}
      onClose={closeCreateConnection}
      presetKind={presetKind}
      skipKindSelection={skipKindSelection}
    />
  );
};

/** Hook for callers — same role as `useRegistryModal`. */
export const useConnectionWizardModal = () => {
  const context = useContext(ConnectionWizardContext);
  return {
    open: context.open,
    presetKind: context.presetKind,
    skipKindSelection: context.skipKindSelection,
    openCreateConnection: context.openCreateConnection,
    closeCreateConnection: context.closeCreateConnection,
  };
};

export default ConnectionWizardContextProvider;
