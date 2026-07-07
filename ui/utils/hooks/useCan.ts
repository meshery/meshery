import { useState, useEffect, useCallback } from 'react';
import _ from 'lodash';
import { ability } from '../can';

/**
 * Reactive permission hook that subscribes to {@link ability}'s `update` event.
 *
 * Unlike the imperative {@link CAN} function, `useCan` causes the consuming
 * component to **re-render** whenever the ability rules change (e.g. when the
 * user's capabilities finish loading asynchronously). This prevents the
 * "permanent DefaultError" bug described in meshery/meshery#20504.
 *
 * @param action - The action identifier (UUID from permission_constants).
 * @param subject - The subject string (will be lower-cased via lodash).
 * @returns `true` if the current ability rules allow the action on the subject.
 */
export function useCan(action: string, subject: string): boolean {
  const evaluate = useCallback(() => ability.can(action, _.lowerCase(subject)), [action, subject]);

  const [allowed, setAllowed] = useState(evaluate);

  useEffect(() => {
    // Re-evaluate immediately in case rules changed between render and effect
    setAllowed(evaluate());

    const unsubscribe = ability.on('update', () => {
      setAllowed(evaluate());
    });

    return () => {
      unsubscribe();
    };
  }, [evaluate]);

  return allowed;
}
