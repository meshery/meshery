---
layout: item
name: ELKK Stack with Kafka
publishedVersion: 0.0.49
userId: 985074f5-27a0-406c-b83a-7f247ba77746
userName: Kavya Katal
userAvatarURL: 
type: deployment
compatibility: 
    - kubernetes
patternId: a1eaf8c4-4da3-4cc4-ba75-edc9e808ab98
image: /assets/images/logos/service-mesh-pattern.svg
patternInfo: |
  ELKK%20(ELK%20%2B%20Kafka)%20logging%20pipeline%3A%20My%20App%20emits%20logs%20%E2%86%92%20Kafka%20Producer%20publishes%20to%20Kafka%20topic%20logs%20%E2%86%92%20Kafka%20brokers%20the%20stream%20(9092)%20%E2%86%92%20Logstash%20consumes%20from%20Kafka%20and%20indexes%20into%20Elasticsearch%20(index%20pattern%20app-logs-*%2C%20port%209200)%20%E2%86%92%20Kibana%20queries%20Elasticsearch%20for%20visualisation%20(5601).%0AThis%20Meshery%20Design%20captures%20the%20topology%20and%20intent%3B%20some%20nodes%20use%20generic%20workload%20placeholders%20where%20native%20models%20were%20not%20available.
patternCaveats: |
  Design%20is%20topology-focused%3B%20some%20components%20(Kafka%20Producer%2C%20Logstash%2C%20Kibana)%20are%20represented%20as%20generic%20workloads%2FNamespaces.%0AProvides%20a%20Logstash%20pipeline%20that%20consumes%20from%20Kafka%20and%20outputs%20to%20Elasticsearch.%0AUpdate%20hostnames%2Fports%20for%20your%20environment%20(e.g.%2C%20kafka%3A9092%2C%20elasticsearch%3A9200).%0ASecurity%2Fauth%20not%20configured%20(no%20TLS%2FSASL%20for%20Kafka%2C%20no%20X-Pack%20auth%20for%20ES%2FKibana).%0AIntended%20for%20a%20single-namespace%20demo%3B%20scaling%20and%20production%20hardening%20are%20out%20of%20scope.
permalink: catalog/deployment/elkk-stack-with-kafka-a1eaf8c4-4da3-4cc4-ba75-edc9e808ab98.html
URL: 'https://raw.githubusercontent.com/meshery/meshery.io/master/catalog/a1eaf8c4-4da3-4cc4-ba75-edc9e808ab98/0.0.49/design.yml'
downloadLink: a1eaf8c4-4da3-4cc4-ba75-edc9e808ab98/design.yml
---