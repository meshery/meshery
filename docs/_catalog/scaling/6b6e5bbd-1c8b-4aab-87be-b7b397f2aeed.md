---
layout: item
name: 'K8''s-Cluster-overprovisioner '
publishedVersion: 0.0.5
userId: 9222bde1-64c6-4fb2-971a-3402d5ae2fd9
userName: Deepak Reddy
userAvatarURL: https://lh3.googleusercontent.com/a/ACg8ocIGbiDtE0q65qVvAUdzHw8Qky81rM0kSAknIqbgysfDCw=s96-c
type: scaling
compatibility: 
    - kubernetes
patternId: 6b6e5bbd-1c8b-4aab-87be-b7b397f2aeed
image: https://raw.githubusercontent.com/layer5labs/meshery-extensions-packages/master/action-assets/design-assets/6b6e5bbd-1c8b-4aab-87be-b7b397f2aeed-light.png,https://raw.githubusercontent.com/layer5labs/meshery-extensions-packages/master/action-assets/design-assets/6b6e5bbd-1c8b-4aab-87be-b7b397f2aeed-dark.png
patternInfo: |
  This%20design%20provide%20a%20buffer%20for%20cluster%20autoscaling%20to%20allow%20overprovisioning%20of%20cluster%20nodes.%20This%20is%20desired%20when%20you%20have%20work%20loads%20that%20need%20to%20scale%20up%20quickly%20without%20waiting%20for%20the%20new%20cluster%20nodes%20to%20be%20created%20and%20join%20the%20cluster.%0A%0AIt%20works%20by%20creating%20a%20deployment%20that%20creates%20pods%20of%20a%20lower%20than%20default%20PriorityClass.%20These%20pods%20request%20resources%20from%20the%20cluster%20but%20don't%20actually%20consume%20any%20resources.%20These%20pods%20are%20then%20evicted%20allowing%20other%20normal%20pods%20to%20be%20created%20while%20also%20triggering%20a%20scale-up%20by%20the%20.
patternCaveats: |
  get%20info%20from%20this%20https%3A%2F%2Fgithub.com%2Fkubernetes%2Fautoscaler%2Fblob%2Fmaster%2Fcluster-autoscaler%2FFAQ.md%23how-can-i-configure-overprovisioning-with-cluster-autoscaler
permalink: catalog/scaling/k8-s-cluster-overprovisioner-6b6e5bbd-1c8b-4aab-87be-b7b397f2aeed.html
URL: 'https://raw.githubusercontent.com/meshery/meshery.io/master/catalog/6b6e5bbd-1c8b-4aab-87be-b7b397f2aeed/0.0.5/design.yml'
downloadLink: 6b6e5bbd-1c8b-4aab-87be-b7b397f2aeed/design.yml
---