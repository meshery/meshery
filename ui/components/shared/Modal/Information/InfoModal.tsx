/**
 * Temporary shim.
 *
 * The legacy catalog/publish-style `InfoModal` lives at `./LegacyInfoModal`
 * and is re-exported from here so existing imports such as
 *
 *   import InfoModal from '../shared/Modal/Information/InfoModal';
 *   import { VIEW_VISIBILITY } from '../shared/Modal/Information/InfoModal';
 *
 * continue to resolve. The canonical, simpler `InfoModal` primitive is at
 * `../InfoModal`; new code should consume it via the barrel:
 *
 *   import { InfoModal } from '@/components/shared/Modal';
 *
 * This shim will be deleted once the Phase 5.b migration sub-issues
 * (#18752–#18756) complete.
 */
export * from './LegacyInfoModal';
export { default } from './LegacyInfoModal';
