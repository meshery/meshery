import { PureAbility } from '@casl/ability';
import * as Sistent from '@sistent/sistent';
import _ from 'lodash';
import { CapabilitiesRegistry } from './disabledComponents';
import { store } from '../store';
import { mesheryEventBus } from './eventBus';

export const ability = new PureAbility([]);

export default function CAN(action, subject) {
  return ability.can(action, _.lowerCase(subject));
}

const getCapabilitiesRegistry = () =>
  new CapabilitiesRegistry((store.getState() as any).ui?.capabilitiesRegistry);

// createCanShow exists at runtime on the Sistent package, but is not typed
const createCanShow = (Sistent as any).createCanShow;

export const CanShow = createCanShow(getCapabilitiesRegistry, CAN, () => mesheryEventBus);
