import { PureAbility } from '@casl/ability';

export const ability = new PureAbility([]);

export default function CAN(action, subject) {
  return ability.can(action, subject);
}
