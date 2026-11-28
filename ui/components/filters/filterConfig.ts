/**
 * Returns the configuration stored with a filter, or '' when it has none.
 *
 * A filter save stores its configuration inside `filterResource`, a JSON
 * WASMFilter component definition, under `configuration.config` (see
 * generateFilterComponent in server/handlers/meshery_filter_handler.go).
 * Filters saved by releases before the v1beta1 component migration keep it
 * under `settings.config`, so that shape is still read as a fallback.
 */
export const getFilterConfig = (filter?: { filterResource?: string } | null): string => {
  const raw = filter?.filterResource;
  if (!raw) {
    return '';
  }

  try {
    const resource = JSON.parse(raw);
    const config = resource?.configuration?.config ?? resource?.settings?.config;
    return typeof config === 'string' ? config : '';
  } catch {
    return '';
  }
};
