import { Ability } from '@casl/ability';
import { createCanShow } from '@sistent/sistent';
import _ from 'lodash';
import { ProviderUiAccessControl } from './disabledComponents';
import { store } from '../store';
import { mesheryEventBus } from './eventBus';

export const ability = new Ability([]);

export default function CAN(action?: unknown, subject?: unknown): boolean {
  if (
    typeof action !== 'string' ||
    !action.trim() ||
    typeof subject !== 'string' ||
    !subject.trim()
  ) {
    return false;
  }
  return ability.can(action.trim(), _.lowerCase(subject.trim()));
}

const getProviderUiAccessControl = () =>
  new ProviderUiAccessControl(store.getState().providerCapabilities);

export const CanShow = createCanShow(getProviderUiAccessControl, CAN, () => mesheryEventBus);
