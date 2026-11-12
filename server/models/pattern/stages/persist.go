package stages

// func Persist(_ ServiceInfoProvider, act ServiceActionProvider, log logger.Handler) ChainStageFunction {
// 	return func(data *Data, err error, next ChainStageNextFunction) {
// 		if err != nil {
// 			act.Terminate(err)
// 			return
// 		}

// 		data.Lock.Lock()
// 		for k := range data.Other {
// 			// Find the services which where provisioned successfully
// 			if strings.HasSuffix(k, ProvisionSuffixKey) {
// 				k = strings.TrimSuffix(k, ProvisionSuffixKey)

// 				_, ok := data.Other[fmt.Sprintf("%s%s", k, UpdateSuffixKey)]
// 				if err := act.Persist(k, *data.Pattern.Services[k], ok); err != nil {
// 					// Just log the error - non critical issue
// 					err = models.ErrCreateResourceEntry(err)
// 					log.Warn(err)
// 				}
// 			}
// 		}
// 		data.Lock.Unlock()

// 		if next != nil {
// 			next(data, nil)
// 		}
// 	}
// }
