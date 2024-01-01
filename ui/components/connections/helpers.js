import { CONTROLLER_STATES } from "../../utils/Enum"
import { LIGHTINDICATORS_CONNECTION_STATES_MAPPING } from "../NotificationCenter/constants"

export const getlightIncicatorForK8s = (operatorState, meshsyncState, brokerState) => { 
    if (operatorState === CONTROLLER_STATES.DEPLOYED
        && meshsyncState === CONTROLLER_STATES.CONNECTED
        && meshsyncState === CONTROLLER_STATES.CONNECTED) {
            return LIGHTINDICATORS_CONNECTION_STATES_MAPPING.GREEN
        } else if (operatorState === CONTROLLER_STATES.DEPLOYED 
        && (meshsyncState !== CONTROLLER_STATES.CONNECTED || brokerState !== CONTROLLER_STATES.CONNECTED)) {
        return LIGHTINDICATORS_CONNECTION_STATES_MAPPING.YELLOW
        } else {
            return LIGHTINDICATORS_CONNECTION_STATES_MAPPING.GREY
        }
}