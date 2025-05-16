---
layout: item
name: Istio BookInfo Application
publishedVersion: 0.0.1
userId: 173202fe-b94d-4064-8a86-1cf063732884
userName: Ashish Tiwari
userAvatarURL: https://layer5.io/static/7b1f08e10d271cbfd963c7d435cf84ac/416c3/ashish-tiwari.webp
type: deployment
compatibility: 
    - istio-base
    - istio-operator
    - kubernetes
patternId: ddf8a9f1-5fd0-4bda-b074-7dc9ceb1dd43
image: https://raw.githubusercontent.com/layer5labs/meshery-extensions-packages/master/action-assets/design-assets/ddf8a9f1-5fd0-4bda-b074-7dc9ceb1dd43-light.png,https://raw.githubusercontent.com/layer5labs/meshery-extensions-packages/master/action-assets/design-assets/ddf8a9f1-5fd0-4bda-b074-7dc9ceb1dd43-dark.png
patternInfo: |
  This%20design%20deploys%20Istio%20BookInfo%20App.%20The%20Bookinfo%20application%20is%20broken%20into%20four%20separate%20microservices%3A%20%0A%0Aproductpage.%20The%20productpage%20microservice%20calls%20the%20details%20and%20reviews%20microservices%20to%20populate%20the%20page.%20%5C%5C%5C%5Cndetails.%20The%20details%20microservice%20contains%20book%20information.%20%5C%5C%5C%5Cnreviews.%20The%20reviews%20microservice%20contains%20book%20reviews.%20It%20also%20calls%20the%20ratings%20microservice.%20%0A%0Aratings.%20The%20ratings%20microservice%20contains%20book%20ranking%20information%20that%20accompanies%20a%20book%20review.
patternCaveats: |
  %20%20Networking%20should%20be%20properly%20configured%20to%20enable%20communication%20between%20the%20different%20services%20of%20the%20app
permalink: catalog/deployment/istio-bookinfo-application-ddf8a9f1-5fd0-4bda-b074-7dc9ceb1dd43.html
URL: 'https://raw.githubusercontent.com/meshery/meshery.io/master/catalog/ddf8a9f1-5fd0-4bda-b074-7dc9ceb1dd43/0.0.1/design.yml'
downloadLink: ddf8a9f1-5fd0-4bda-b074-7dc9ceb1dd43/design.yml
---