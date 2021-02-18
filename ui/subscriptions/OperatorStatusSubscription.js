import {
    graphql,
    requestSubscription
  } from 'react-relay'
import environment from '../lib/relayEnvironment'
  
const newVoteSubscription = graphql`
    subscription OperatorStatusSubscription {
        listenToOperatorEvents {
        core
        controllers {
            name
            status
        }
        }
    }
    `

  export default () => {
  
    const subscriptionConfig = {
      subscription: newVoteSubscription,
      variables: {},
      updater: proxyStore => {
        const createOperatorStatusField = proxyStore.getRootField('listenToOperatorEvents')
        const core = createOperatorStatusField.getLinkedRecord('core')
        const controllers = createOperatorStatusField.getLinkedRecord('controllers')
        
    },
      onError: error => console.log(`An error occured:`, error)
    }
  
    requestSubscription(
      environment,
      subscriptionConfig
    )
  
  }