import { PureAbility } from '@casl/ability';
import _ from 'lodash';

export const ability = new PureAbility([]);

export default function CAN(action, subject) {
  return ability.can(action, _.lowerCase(subject));
}
