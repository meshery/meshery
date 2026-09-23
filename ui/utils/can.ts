import { createCanShow } from '@sistent/sistent';
import _ from 'lodash';
import { ProviderUiAccessControl } from './disabledComponents';
import { store } from '../store';
import { mesheryEventBus } from './eventBus';

type AbilityRule = { action: string; subject: string };

class SimpleAbility {
  private rules: AbilityRule[] = [];

  update(rules: AbilityRule[]) {
    this.rules = rules;
  }

  can(action: string, subject: string) {
    return this.rules.some((rule) => rule.action === action && rule.subject === subject);
  }
}

export const ability = new SimpleAbility();

export default function CAN(action, subject) {
  return ability.can(action, _.lowerCase(subject));
}

const getProviderUiAccessControl = () =>
  new ProviderUiAccessControl(store.getState().providerCapabilities);

export const CanShow = createCanShow(getProviderUiAccessControl, CAN, () => mesheryEventBus);
