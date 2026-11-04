import { useEffect, useRef, useState } from 'react';
import { EVENT_TYPES } from 'lib/event-types';
import { useNotification } from '@/utils/hooks/useNotification';
import type { ControllersConfigDoc, VersionedControllersConfigDoc } from './ControllersConfigForm';

/** Drops the server-stamped schemaVersion, which is not part of the editable document. */
export const stripSchemaVersion = (
  doc?: VersionedControllersConfigDoc | null,
): ControllersConfigDoc => {
  if (!doc) return {};
  const rest: VersionedControllersConfigDoc = { ...doc };
  delete rest.schemaVersion;
  return rest;
};

type ControllersConfigDraftMessages = {
  /** Notified once per failed load of the persisted document. */
  loadError: string;
  /** Notified when the save request fails; the server error is attached as details. */
  saveError: string;
  /** Notified after the document is persisted. */
  saveSuccess: string;
};

export type UseControllersConfigDraftOptions = {
  /** True once the query has resolved; the draft is seeded from `source` while clean. */
  isLoaded: boolean;
  /** The persisted document the draft is seeded from. */
  source: VersionedControllersConfigDoc | null | undefined;
  /** The query error, if the document failed to load. */
  loadError?: unknown;
  /** Persists the draft. Must reject on failure (i.e. pass the unwrapped mutation). */
  save: (body: ControllersConfigDoc) => Promise<unknown>;
  messages: ControllersConfigDraftMessages;
  /** Invoked after a successful save, e.g. to close the editor. */
  onSaved?: () => void;
};

export type ControllersConfigDraft = {
  /** The document currently being edited. */
  draft: ControllersConfigDoc;
  /** True once the draft has diverged from the persisted document. */
  dirty: boolean;
  /** Accepts an edited document from ControllersConfigForm. */
  onChange: (next: ControllersConfigDoc) => void;
  /** Resets the draft back to the persisted document. */
  discard: () => void;
  /** Persists the draft and notifies the outcome. Never rejects. */
  save: () => Promise<void>;
};

/**
 * The editing shell shared by the two controllers-configuration editors (the
 * Settings tab's server-wide defaults and the per-connection override modal):
 * seed a local draft from the persisted document, track dirtiness, save,
 * notify, discard. Both editors go through this hook so the two cannot drift
 * apart in their save, discard or error-reporting behaviour.
 */
export const useControllersConfigDraft = ({
  isLoaded,
  source,
  loadError,
  save,
  messages,
  onSaved,
}: UseControllersConfigDraftOptions): ControllersConfigDraft => {
  const { notify } = useNotification();
  const [draft, setDraft] = useState<ControllersConfigDoc>({});
  const [dirty, setDirty] = useState(false);

  // Seed the draft from the persisted document when that document arrives or
  // changes - never merely because `dirty` went false.
  //
  // Depending on `dirty` re-seeded the draft the moment a save completed, while
  // the query cache could still be holding the pre-save response until its
  // invalidation and refetch landed. The editor visibly snapped back to stale
  // values for that window, immediately after telling the user the save had
  // succeeded. Keying on the source reference alone means a save updates the
  // draft only when the refreshed document actually arrives.
  const seededFrom = useRef<VersionedControllersConfigDoc | null | undefined>(undefined);
  useEffect(() => {
    if (!isLoaded) return;
    if (seededFrom.current === source) return;
    seededFrom.current = source;
    setDraft(stripSchemaVersion(source));
    setDirty(false);
  }, [isLoaded, source]);

  // Notify once per load failure rather than on every render.
  useEffect(() => {
    if (loadError) {
      notify({ message: messages.loadError, event_type: EVENT_TYPES.ERROR });
    }
  }, [loadError]);

  return {
    draft,
    dirty,
    onChange: (next) => {
      setDraft(next);
      setDirty(true);
    },
    discard: () => {
      setDraft(stripSchemaVersion(source));
      setDirty(false);
    },
    save: async () => {
      try {
        await save(draft);
        setDirty(false);
        notify({ message: messages.saveSuccess, event_type: EVENT_TYPES.SUCCESS });
        onSaved?.();
      } catch (err) {
        notify({
          message: messages.saveError,
          event_type: EVENT_TYPES.ERROR,
          details: String((err as { data?: unknown })?.data ?? err),
        });
      }
    },
  };
};
