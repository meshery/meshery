import { workerfyActor } from '@layer5/sistent';
import { schemaValidatorMachine } from './schemaValidator';

workerfyActor(schemaValidatorMachine);
