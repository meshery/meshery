import { PureAbility } from '@casl/ability';
import { createCanShow } from '@sistent/sistent';
import _ from 'lodash';
import { ProviderUiAccessControl } from './disabledComponents';
import { store } from '../store';
import { mesheryEventBus } from './eventBus';

export const ability = new PureAbility([]);

export default function CAN(action, subject) {
  return ability.can(action, _.lowerCase(subject));
}

const getProviderUiAccessControl = () =>
  new ProviderUiAccessControl(store.getState().providerCapabilities);

export const CanShow = createCanShow(getProviderUiAccessControl, CAN, () => mesheryEventBus);
