import { useLegacySelector } from 'lib/store';

/**
 * Generic hook to safely access state properties from either Immutable.js or plain JS Redux store
 * @param {string} stateName - The name of the state property to retrieve
 * @returns {any} - The requested state value converted to plain JS if needed
 */
export const useStateValue = (stateName) => {
  return useLegacySelector((state) => {
    try {
      if (typeof state?.get === 'function') {
        const value = state.get(stateName);
        return value?.toJS ? value.toJS() : value;
      }

      return state?.[stateName];
    } catch {
      return undefined;
    }
  });
};

export const useGetCurrentOrganization = () => {
  return useStateValue('organization');
};
