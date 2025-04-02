import { createContext } from 'react';
import { createDispatchHook, createSelectorHook, createStoreHook } from 'react-redux';

export const RTKContext = createContext(null);

export const useSelectorRtk = createSelectorHook(RTKContext);
export const useDispatchRtk = createDispatchHook(RTKContext);
export const useStoreRtk = createStoreHook(RTKContext);
