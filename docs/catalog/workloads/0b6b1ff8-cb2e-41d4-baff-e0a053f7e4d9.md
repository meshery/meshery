---
layout: item
name: Apache Airflow
publishedVersion: 0.0.1
userId: 9222bde1-64c6-4fb2-971a-3402d5ae2fd9
userName: Deepak Reddy
userAvatarURL: https://lh3.googleusercontent.com/a/ACg8ocIGbiDtE0q65qVvAUdzHw8Qky81rM0kSAknIqbgysfDCw=s96-c
type: workloads
compatibility: 
    - kubernetes
patternId: 0b6b1ff8-cb2e-41d4-baff-e0a053f7e4d9
image: /assets/images/logos/service-mesh-pattern.svg
patternInfo: |
  Apache%20Airflow%20(or%20simply%20Airflow)%20is%20a%20platform%20to%20programmatically%20author%2C%20schedule%2C%20and%20monitor%20workflows.%0A%0AWhen%20workflows%20are%20defined%20as%20code%2C%20they%20become%20more%20maintainable%2C%20versionable%2C%20testable%2C%20and%20collaborative.%0A%0AUse%20Airflow%20to%20author%20workflows%20as%20directed%20acyclic%20graphs%20(DAGs)%20of%20tasks.%20The%20Airflow%20scheduler%20executes%20your%20tasks%20on%20an%20array%20of%20workers%20while%20following%20the%20specified%20dependencies.%20Rich%20command%20line%20utilities%20make%20performing%20complex%20surgeries%20on%20DAGs%20a%20snap.%20The%20rich%20user%20interface%20makes%20it%20easy%20to%20visualize%20pipelines%20running%20in%20production%2C%20monitor%20progress%2C%20and%20troubleshoot%20issues%20when%20needed.%0A%0AAirflow%20works%20best%20with%20workflows%20that%20are%20mostly%20static%20and%20slowly%20changing.%20When%20the%20DAG%20structure%20is%20similar%20from%20one%20run%20to%20the%20next%2C%20it%20clarifies%20the%20unit%20of%20work%20and%20continuity.%20Other%20similar%20projects%20include%20Luigi%2C%20Oozie%20and%20Azkaban.%0A%0AAirflow%20is%20commonly%20used%20to%20process%20data%2C%20but%20has%20the%20opinion%20that%20tasks%20should%20ideally%20be%20idempotent%20(i.e.%2C%20results%20of%20the%20task%20will%20be%20the%20same%2C%20and%20will%20not%20create%20duplicated%20data%20in%20a%20destination%20system)%2C%20and%20should%20not%20pass%20large%20quantities%20of%20data%20from%20one%20task%20to%20the%20next%20(though%20tasks%20can%20pass%20metadata%20using%20Airflow's%20XCom%20feature).%20For%20high-volume%2C%20data-intensive%20tasks%2C%20a%20best%20practice%20is%20to%20delegate%20to%20external%20services%20specializing%20in%20that%20type%20of%20work.%0A%0AAirflow%20is%20not%20a%20streaming%20solution%2C%20but%20it%20is%20often%20used%20to%20process%20real-time%20data%2C%20pulling%20data%20off%20streams%20in%20batches.%0A%0APrinciples%0ADynamic%3A%20Airflow%20pipelines%20are%20configuration%20as%20code%20(Python)%2C%20allowing%20for%20dynamic%20pipeline%20generation.%20This%20allows%20for%20writing%20code%20that%20instantiates%20pipelines%20dynamically.%0AExtensible%3A%20Easily%20define%20your%20own%20operators%2C%20executors%20and%20extend%20the%20library%20so%20that%20it%20fits%20the%20level%20of%20abstraction%20that%20suits%20your%20environment.%0AElegant%3A%20Airflow%20pipelines%20are%20lean%20and%20explicit.%20Parameterizing%20your%20scripts%20is%20built%20into%20the%20core%20of%20Airflow%20using%20the%20powerful%20Jinja%20templating%20engine.%0AScalable%3A%20Airflow%20has%20a%20modular%20architecture%20and%20uses%20a%20message%20queue%20to%20orchestrate%20an%20arbitrary%20number%20of%20workers.
patternCaveats: |
  Make%20sure%20to%20fill%20out%20your%20own%20postgres%20username%20%2Cpassword%2C%20host%2Cport%20etc%20to%20see%20airflow%20working%20as%20per%20your%20database%20requirements.%20pass%20them%20as%20environment%20variables%20or%20create%20secrets%20%20%20for%20password%20and%20config%20map%20for%20ports%20%2Chost%20.
permalink: catalog/workloads/apache-airflow-0b6b1ff8-cb2e-41d4-baff-e0a053f7e4d9.html
URL: 'https://raw.githubusercontent.com/meshery/meshery.io/master/catalog/0b6b1ff8-cb2e-41d4-baff-e0a053f7e4d9/0.0.1/design.yml'
downloadLink: 0b6b1ff8-cb2e-41d4-baff-e0a053f7e4d9/design.yml
---