/**
 * Imperative open/close hook for the shared `Modal` primitive.
 *
 * Provides a tiny state shell so callers don't have to repeat
 * `const [open, setOpen] = useState(false)` for every dialog.
 *
 *   const modal = useModal();
 *   <Button onClick={modal.open}>Edit</Button>
 *   <Modal isOpen={modal.isOpen} onClose={modal.close} title="Edit" />
 *
 * `toggle` is provided as a convenience for menu/disclosure patterns.
 * `setIsOpen` is exposed so consumers that already track state externally
 * (e.g. via Redux) can drive the same prop surface.
 */
import { useCallback, useState } from 'react';

export interface UseModalReturn {
  isOpen: boolean;
  open: () => void;
  close: () => void;
  toggle: () => void;
  setIsOpen: (next: boolean) => void;
}

export const useModal = (initialOpen = false): UseModalReturn => {
  const [isOpen, setIsOpen] = useState<boolean>(initialOpen);

  const open = useCallback(() => setIsOpen(true), []);
  const close = useCallback(() => setIsOpen(false), []);
  const toggle = useCallback(() => setIsOpen((prev) => !prev), []);

  return { isOpen, open, close, toggle, setIsOpen };
};

export default useModal;
