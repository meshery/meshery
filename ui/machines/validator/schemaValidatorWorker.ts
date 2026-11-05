import { workerfyActor } from '@sistent/sistent';
import { schemaValidatorMachine } from './schemaValidator';

workerfyActor(schemaValidatorMachine);
