import { workerfyActor } from '@layer5/sistent';
import { schemaValidatorMachine } from './schemaValidator';

console.log('Workerfying schemaValidatorMachine');

workerfyActor(schemaValidatorMachine);
