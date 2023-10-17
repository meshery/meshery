import { PureAbility } from '@casl/ability';

export const ability = new PureAbility([]);
export default function checkPermissions(action, subject) {
  return ability.can(action, subject);
}
