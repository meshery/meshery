---
layout: item
name: All relationships
publishedVersion: 0.0.310
userId: 830c0932-b05c-4c99-b12b-872f6cf12202
userName: zihan 2
userAvatarURL: https://avatars.githubusercontent.com/u/127078886?s=400&u=91e47105e7d418843a4757cd9bd3cc431e29deaf&v=4
type: deployment
compatibility: 

patternId: 2fdf5dcb-a3c0-4c2c-a64f-a3943ac87d5f
image: /assets/images/logos/service-mesh-pattern.svg
patternInfo: |
  This%20design%20incorporates%20all%20the%20key%20relationships%2C%20including%20the%20following%3A%0A%0A%0A%0A1.%C2%A0%20Hierarchical-Parent-Inventory%3A%20This%20represents%20a%20parent-child%20relationship%20where%20the%20configuration%20of%20a%20parent%20component%20influences%20the%20child%20component.%0A%0A%0A%0A2.%C2%A0%20Hierarchical-Parent-Wallet%3A%20In%20this%20relationship%2C%20one%20component%20(the%20%22wallet%22)%20serves%20as%20a%20container%20or%20host%20for%20another%2C%20similar%20to%20a%20parent-child%20structure.%0A%0A3.%C2%A0%20Hierarchical-Sibling-MatchLabels%3A%20A%20Match-Labels%20Relationship%20links%20components%20based%20on%20shared%20labels%2C%20indicating%20they%20are%20siblings%20that%20operate%20together.%0A%0A%0A%0A4.%C2%A0%20Edge-Mount%3A%20An%20Edge-Mount%20Relationship%20represents%20the%20assignment%20of%20persistent%20storage%20to%20a%20workload%20via%20a%20PersistentVolumeClaim%20(PVC).%0A%0A%0A%0A5.%C2%A0%20Edge-Permission%3A%20The%20Edge-Permission%20Relationship%20defines%20how%20components%20connect%20to%20establish%20access%20control%2C%20where%20bindings%20link%20subjects%20(like%20users%20or%20services)%20to%20roles%20with%20specific%20permissions.%0A%0A%0A%0A6.%C2%A0%20Edge-Firewall%3A%20An%20Edge-Firewall%20Relationship%20models%20a%20network%20policy%20that%20controls%20ingress%20and%20egress%20traffic%20between%20components.%0A%0A%0A%0A7.%C2%A0%20Edge-Network%3A%20An%20Edge-Network%20Relationship%20represents%20the%20networking%20configuration%20between%20components%2C%20typically%20illustrated%20by%20a%20service%20providing%20a%20stable%20endpoint%20for%20a%20deployment.%0A%0A%0A%0A8.%C2%A0%20Edge-Annotation%3A%20An%20Annotation%20Relationship%20is%20a%20visual%20indicator%20used%20to%20show%20a%20connection%20between%20two%20components%20without%20assigning%20any%20functional%2C%20semantic%20meaning%20to%20that%20relationship.%0A%0A%0A%0A9.%C2%A0%20Edge-Reference%3A%20An%20Edge-Reference%20Relationship%20represents%20a%20logical%20link%20where%20one%20component%20refers%20to%20another%20by%20its%20name%20or%20identifier.%20It%20enables%20interaction%20by%20declaring%20intent%2C%20such%20as%20a%20Pod%20referencing%20a%20ConfigMap%20for%20its%20configuration%20data.
patternCaveats: |
  For%20detailed%20considerations%20on%20each%20relationship%20type%2C%20refer%20to%20the%20corresponding%20individual%20published%20designs.%20These%20designs%20provide%20in-depth%20insights%20into%20best%20practices%2C%20configuration%20strategies%2C%20and%20potential%20impacts%20for%20each%20type%20of%20relationship.
permalink: catalog/deployment/all-relationships-2fdf5dcb-a3c0-4c2c-a64f-a3943ac87d5f.html
URL: 'https://raw.githubusercontent.com/meshery/meshery.io/master/catalog/2fdf5dcb-a3c0-4c2c-a64f-a3943ac87d5f/0.0.310/design.yml'
downloadLink: 2fdf5dcb-a3c0-4c2c-a64f-a3943ac87d5f/design.yml
---