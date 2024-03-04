# Meshery Models

## A unit of packaging for Meshery’s logical object representation

# Prologue and Challenge

As a cloud-native management plane, Meshery should be able to:

- Empower its users with a wide range of tools that provide support for the majority of the systems in the cloud and cloud native ecosystems.
- Abstract away the system specific requirements and help the users focus on getting things done.

Meshery currently has an object model that is internal to Meshery that can capture all these use cases and serve as a way for consistently creating well-defined boundaries around entities that Meshery manages.

## Guiding Principles

Adhere to the following design principles:

1. **Establish a set of constructs that will be used across Meshery**
2. **Make sure the boundaries are clearly defined with vigil eye to extensibility**

   Clearly define what the responsibilities of a construct are and are not.

   Allow users to extend the constructs and create new constructs.

3. **Reuse where possible. Differentiate where feasible.**

   Build upon other works while uplifting and sprinkling in unique value. If it makes sense, make use of the things we already have in our model.

4. **Be end user-centric; Enable the business function.**

   Simplify. Networking and distributed systems are hard; Connect the dots between dynamic infrastructure and the business.

## Design Goals

The goal is to have a model for Meshery that:

1. Sets a foundation for the project to build and expand on (are versioned).
2. Facilitates a consistent way of communicating between multiple components.
3. Provides the ability to understand constructs of different systems in the cloud native ecosystem and beyond.

#### Open Questions:

1. Where’s the database? - Meshery Cloud / Git
2. Where is the pattern engine doc?

# MeshModel Model Constructs

# Model Packaging

Each model package can be imported and exported from the system as OCI-compatible images, making them portable (a design goal), abstracting their complexity, and encapsulating potential intellectual property that users might have invested into their models. Model packages are versioned and bundle any number of components, relationships, policies, connections and credentials. For example:

### Model Schema

```
apiVersion: core.meshery.io/v1alpha1
kind: ModelDefinition
metadata:
  name:
  version: 1.19
spec:
  schematic:
		cue: |
		designations: vs : kubernetes
	    	outputs : vs :  {
	        	apiVersion: networking.istio.io/v1alpha3
	        	kind:   	VirtualService
	        	spec: {
			# rest of the configuration for VS
	        	}
	    	}

	    	parameters: {
			param1: string
			param2: int
			# user inputs
	    	}
```

# Glossary: Core Constructs

Every construct will be represented in three forms (two static and one dynamic).

1. **Schema** (static) - the skeletal structure representing a logical view of the size, shape, characteristics of a construct.
   1. _Example: meshery/schemas_
2. **Definition** (static) - An implementation of the Schema containing specific configuration for the construct at-hand. 2. _Example: Kubernetes Pod_
3. **Declaration** (static) - A defined construct; A specific deof the Definition. 3. _Example: NGINX container as a Kubernetes Pod_
4. **Instance** (dynamic) - A realized construct (deployed/discovered); An instantiation of the Declaration. 4. _Example: NGINX-as234z2 pod running in cluster_

MeshModel supports a common substrate of operands to support granularly, expressively selectors for matching (or not matching) model constructs. The set of supported operands in MeshModel is a direct representation of those supported by Cuelang and defined by RFC xxx / regex xxx. Sample operands:

1. - - wildcard
2. ? - ….
3. ! - ! means NOT
4. &lt;, >, =, !=, …

The absence of a value for a given attribute in a Definition means: — something that has to be written per Construct.

## Component

Components are designed to be portable, accessible, versioned and easily manageable.

Components are a way of exposing the capabilities of the underlying platform. In the case of application development, Components are entities that are registered by platform engineers and used by operators.

`ComponentDefinition`, as the name implies, is the definition of a `Component` which will be stored in the `Registry`. `Components` are instances of `ComponentDefinitions`.

A `Component` on evaluation will output a set of manifests which will be sent to the appropriate `Registrant` (host) for further handling.

Any user should be able to easily create his own `ComponentDefinition` and publish it so that others can make use of it.

### Component Schema

```
apiVersion: core.meshery.io/v1alpha1
kind: ComponentDefinition
format: "CUE",
metadata:
  name: VirtualService
  version: 1.19
  schema:
		cue: |
		designations: vs : kubernetes
	    	outputs : vs :  {
	        	apiVersion: networking.istio.io/v1alpha3
	        	kind:   	VirtualService
	        	spec: {
			# rest of the configuration for VS
	        	}
	    	}

	    	parameters: {
			param1: string
			param2: int
			# user inputs
	    	}
```

`designations` is a map with the names of the `outputs` as keys and the `managers` as the values.

```
managers are entities that know how to manage(provision/deprovision/update) a certain manifest generated when instantiating Components.
```

`outputs` is a map with the names as keys and manifest as values.

```
parameters is a map with the names as keys and the type as keys. These parameters define the configuration parameters required while instantiating a Component.
```

Upon instantiating the `ComponentDefintion` with the necessary parameters, some manifests (like the Kubernetes `VirtualService` resource) will be sent to the `Kubernetes` host as assigned through the `designations` field. The host ( `Kubernetes` in this case) should know how to deal with this manifest.

_How will we leverage the power of cue while defining schemas?_

## Meshery’s Core Components

Currently Meshery has only 1 core component: Application. Below is the ComponentDefinition for it.

## Design

A `Design` consists of `Components` and `Patterns`. `Design` is the deployable unit in Meshery. `Designs` are how the users can describe the desired infrastructure.

There **cannot** be two components with the same name

## Pattern

A `Pattern` is an entity that augments the operational behavior of a deployed instance of a Design. A Pattern can be applied to a `Component` or a `Design`. `Patterns` define a common (and best) practice of both configuring and operating cloud-native application functionality. `Patterns` are read-only.

## Relationships

**Example**

Using the use case where you want to make it such that when an edge is drawn from a Kubernetes `Service` to a Kubernetes `Pod`, the port fields in the `pod` configuration gets automatically filled with the details on the edge. You are just interpreting the relationship that exists between the `Pod` and the `Service` components.

`Relationships` define the genealogy between one or more interconnected `Components`. Just as in familial relationships, MeshModel `Relationships` are represented in a variety of forms very much resembling familiar ancestral lineage, including hierarchical relationships such as direct parent and child relationship as well as any number of layers of indirect grandparent and grandchild relationships. Peer relationships like that of siblings take on different types such as `network` and `dependency`. This can have multiple effects depending on what context it is being used.

\*\*<code>Relationships</code> have <code>metadata</code>, <code>selectors</code> and some optional parameters. </strong>

As someone who creates `Designs` and `Components`, It would be useful If I can have some assistance in visualizing the relationships that exist between multiple `Components`. When some `Components` are related to each other in a `Design`, MeshMap should show how they are related visually to the extent possible.

Treating the Relationships the same way as that of Components is not the ideal approach since the use case for the Relationships are different. This approach will be replaced later with a better alternative like a Graph Database.

### Selectors

Selectors are structured as an array, wherein each entry comprises a `from(self)` field and a `to(other)` field (`[from: [{..}], to: [{..}]]`) (from and to will be **CUE** values), delineating the components involved in a particular relationship. Each item of the selector define the constraints to be satisfied for a particular relationship to exist, the selectors act as a scoping within a relationship enabling more flexibility and re-use when defining relationships.

The components within the same selector or each item of the selector relates to each other via 1:many kind of relation. i.e. Each object inside the `from `relates to each item inside the `two` field within a particular selector. \
 \
When defining relationships that involve a large number of combinations between `from` and `to`, having selectors as an array provides a mechanism to organize and manage these relationships. This prevents the need for crafting complex deny attributes and facilitates easier maintenance. \
 \
This arrangement enhances flexibility and reusability in the definition and configuration of relationships among components.

Eg: WASMFilter - EnvoyFilter,

ConfigMap - Pod,

ConfigMap - Deployment

- DaemonSets
- Jobs
- StatefulSets,

These pairs have `hierarchical inventory relationship`, the policy backing this relationship and the visual paradigm are the same, If these pairs are defined in a single selector then the `deny` attribute needs to be carefully crafted to ensure disambiguity because relationship between ConfigMap and EnvoyFilter needs to be restricted. This can quickly become difficult to maintain due to large number of combinations possible between `from` and `to`.

### Actions

1. Patch

   The Patch represents an action which updates the configuration of the components involved in a relationship. The component from which configuration will be applied is called as `mutator` and the component to which configuration is applied is called as `mutated`.

   The jsonpath of the component property (`mutator`), from which patch will be applied is captured in the attribute `mutatorRef.`

   The jsonpath of the component property (`mutated`), to which patch will be applied is captured in the attribute `mutatedRef.`

   The schema of the `patch` action is defined at [patch schema.](https://github.com/meshery/schemas/blob/849c40f5766be18d36a7e250334d43e085b103a0/schemas/meshmodel/schemas/selectors.json#L40)

   NOTE: The `mutatorRef` sequence must correspond precisely to the `mutatedRef` sequence. For instance, if `mutatorRef` is defined as `[[config, url], [config, name]]`, and `mutatedRef `as `[[configPatch, value], [name]]` then during the patch action, the value at the path [config, url] is patched to mutated component at `[configPatch, value]`, and values at the path` [config, name]` is patched at `[name]`.

#### Operators

Operators help us in adding dynamic flavour to the static relationship definition. It helps in defining relationships which are easy to define but powerful in nature. They are similar to CEL in K8s, but instead they are evaluated in the OPA policies.

For use cases like:

1.  ConfigMap can be used as volume as well as Environment variable, the mutatorRef and mutatedRef for both of the scenario is totally different. To cater this need, we can either define another selector in the same relationship definition or we can use “OR” operator to express them concisely. The former introduces duplication, increases the length of the relationship, can cause confusion if not reviewed carefully (as only a part of it will be different in both the scenarios). The later approach, ensures concise and powerful representation.
2.  Kubernetes resources are highly configurable and forms relationships with each other in many different ways. For instance, PVC can be bound to a PV on the basis of

    1.  storageClass, matchSelelctors and labels of PV and PVC

            OR

    2.  claimRef properties of PVC.

    Suppose in a design consisting of Pod, PV and PVCs, the storageClassName listed in the PV and PVC is same but the constraints defined in 2.a is not satisfied even then as PVC has claimRef to a correct PV kuberentes will bind the PVC to that particular PV. Hence, we should be flexible and open to both type of constraints. This can be addressed by having different relationship but usage of “OR” operator seems more adequate.

    Hence for these such scenarios operators are helpful.

**Interaction during relationship creation.**

Building upon the above ConfigMap example, we still have a unresolved query..

1. Out of all available containers inside the Deployment/Pod to which specific container the configMap should be mounted?
2. Or for Network relationship, out of the all available containers which container to expose?

Such questions cannot be completely addressed by the relationship definition itself, as it involves user’s choice. Therefore, to cater these scenarios there’s a need for “Interaction during relationship creation”. (Ideally handled in [Tier-1 eval](#tier-1-client-side-23))

In some cases we can choose an answer with complete accuracy foreg: If there’s a single container or no containers then we can automatically default to 1st container or create a new one.

But in other situation it is difficult to answer this query.

To support this interaction, the policies during the Tier-2 evaluation, needs to be aware of the user’s choice.

Eg: Out of 3 available containers (1 init and 2 app container) the second container is chosen.

This information needs to be passed along when Tier-2 eval is invoked.

The functionality of [Patch](#actions-14) action and specifically, `mutatorRef` and `mutatedRef` schema needs to be updated to support variables in their property/jsonPath.

```
Eg: The path will be parameterized as follow.
mutatorRef: [[name]]
mutatedRef: [[settings, spec, containers, $container_index, envFrom, $variable_index configMapRef, Name]]
Assume values of $container_index = 2 and $variable_index = 1

The above example showcases a relationship were the 1st set of environment variables present in the 2nd container of a Pod is being replaced by the env variables specified by the configMap.

Since the parameters/variable will be part of relationship definition(mutatedRef/mutatorRef), the Tier-1 eval can reference the relationship definition and can invoke a set of interactions.

The chosen values for these parameters as part of interaction will be passed along as an input to the Tier-2 eval.

Requirement: The inputs should be key-value mapping, scoped within the context of relationship to ensure disambiguity from the parameters of other relationships.

Eg: inputs: {
	hierarchical_relationship: {
		container_index: 2,
    variable_index: 1
}
    }
// We will require to consider for which particular selector these values are for, as the variable names will not necessarily be unique in the relationship definition. Eg: The contianer_index will also be used for secret mount, PVC mount …
// The above schema is just for representation.
```

At the Tier-2: Server side eval stage-3, the inputs will be sanitized and processed.

In the example below, the `Relationship` goes from all the Kubernetes Pod Components (different versions included) to all the Components that have the `Networking` model category.

Scope: package / model level.

### Relationship Schema

```
apiVersion: core.meshery.io/v1alpha1
kind: RelationshipDefinition
metadata:
	name: Network Edge Relationship
	description: Network edge relationship
	type: edge
sub-type: network
status:
    spec:
      schematic:
    	cue: |
selectors:[
  allow:
    from: [self | other ]

```

- `kind: (assume !* when null)`

```
      model: (assume !* when null)
      		version: (assume * when null)
		patch:
		   patchStrategy: replace
		   mutatorRef: (#jsonref to value from where patch should be applied.)
    to: [self | other ]

```

- `kind: (assume * when null) (a package can reference components from another package)`

```
      model: (assume * when null)
    version: (assume * when null)
		patch:
           patchStrategy: replace
       mutatedRef: (#jsonRef to value that should be patched.)

  deny:
    to: [self | other ]

```

- `kind: (assume * when null) (a package can reference components from another package)`

```
      model: (assume * when null)
    version: (assume * when null)
    from: [self | other ]

```

- `kind: (assume * when null) (a package can reference components from another package)`

```
      model: (assume * when null)
    version: (assume * when null)
]
```

\_Figure: Example of a <code>from: self</code> reference</em>

### Things to Keep in Mind while creating RelationshipDefinitions

1. Relationships are defined only between Components.
2. The values for `Kind, Version ` and `Model` are case-sensitive
3. The convention is to use camel-casing for `Kind` and `SubType` values.
4. Absence of a field means in the selector means “\*” (or wildcard).
   1. If we have a selector with {Kind: Pod, Model: Kubernetes}, the absence of the Version field here means that all the versions of the Kubernetes Pod resource (e.g. k8s.io/v1/betav2) will match.
5. In the event of conflicting Relationship Definitions, union between them is taken. 2. If we have two Relationships, one from `(Component A)` to `(Component B and Component F)`, and another from `(Component A)` to `(Component B and Component C), `then it is similar to having a Relationship from `Component A` to `Component B, C and F`
6. In the event of an overlapping set of _complementary_ Relationship Definitions, Union.
7. In the event of an overlapping set of conflicting Relationship Definitions:

   No relationship type (**Kind**) is inherently more important than the next one, so will not be any case of _conflict_

### Order of Operations (Evaluation)

Relationship DefinitionOperations are evaluated

[draft statements] No relationship type is inherently more important than the next. With that said Hierarchical relationships are evaluated first. No, this is not universally true. What is true is that… drop targets are… upon breach of a join/leave proximity perimeter, an evaluation occurs foremost with respect to the **type** of Relationship being evaluated, which in the case of Edge vs Heirarchical/Sibling things are fairly clear. Although, at some point, all three types of rules will need to be considered simultaneously.

### Relationship Example: Edge-Network

Using the Kubernetes Service as a specific example for which the following Relationship would be registered:

```
name:
apiVersion: core.meshery.io/v1alpha1
kind: Relationship
metadata:
	Description: Network edge relationship
	name: Network Edge Relationship
	type: edge
sub-type: network
selectors:
  allow:
    from: self
    to:

```

- `kind: Pod`

  ```
  model: Kubernetes
  ```

- `kind: Deployment`

        ```
        model: Kubernetes
        ```

- `kind: Endpoints`

        ```
        model: Kubernetes

        ```

Use cases for Edges in the Visualizer:

1. Flow of traffic between services
2. Dependency representation (?)
3. Metrics
4.

Use cases for Edges in Designer:

1. Dependency representation (?)
2. Grouping Components. Eg: As selectors for services, Deployments, etc
   1. In this case,an Edge need not be a first class construct, it can just be a way of representing a configuration in the design(selectors of the service in this case)

### RelationshipDefinition Example: Namespaced Kubernetes Components

This Relationship identifies a component as being **parent-capable**. In other words, this Relationship identifies the matching component as being a compound node.

```
apiVersion: core.meshery.io/v1alpha1
kind: RelationshipDefinition
Metadata:
	name: Namespaced Kubernetes Component
    Description: A Kubernetes component that must belong to a Kubernetes Namespace in order to be instantiated.
	type: hierarchical
	sub-type: child
  selector:
    to:
      - Kubernetes.Namespace
    from:
	- self
```

### RelationshipDefinition Example: MeshMap Compound Node

This Relationship identifies a component as being **child-capable**. In other words, this Relationship identifies the matching component as being a compound node.

```
apiVersion: core.meshery.io/v1alpha1
kind: RelationshipDefinition
Metadata:
	name: MeshMap Compound Node
    Description: A component that is capable of having supporting children and any number of generations of their progeny.
	type: hierarchical
	sub-type: parent
  selector:
	allow:
        from:
          - self
        to:
    	- *
```

\
AttributeDefinition Schema

```
apiVersion: core.meshery.io/v1alpha1
kind: AttributeDefinition
metadata:
	componentType: compound-node
selectors:
nonself:
    type: namespace
```

### Scoping Relationships

1. Global Scope: The relationships can be configured to be applied to specific model, a specific model version or can be configured to be applied across model. The relationship schema has a `model` and `version` attribute which facilitates this control. Eg: If the model is specified as `aws-ec2-controller`, the relationship will work for those components which belongs to the `aws-ec2-controller` model.
2. Local Scope: It is controlled via the `selectors` attribute in the relationships.

## Policies

- category=security
- type=opa

### Policy Schema

```
apiVersion: core.meshery.io/v1alpha1
kind: PolicyDefinition
metadata:
	type: match
	sub-type: relationship

selectors:
self:
	type: node
nonself:
    type: namespace
action:
type: allow
```

# Policy Evaluation

### Tiers of Evaluation

#### Tier-1: Client side

_Input: All nodes on canvas and all registered relationships._

_Output: Assigning actions possible with each node on canvas, and animations to invoke. These are meant to enhance the UX and is in context with the current state of Canvas._

The Tier-1 eval results help answer the following queries:

Q1 Set of nodes to highlight.

Q2 Set of relationships possible with the nodes involved in the ongoing relationships (discerning inreal-time the relations possible with the node in focus).

    Q2.1 Automatically deaulting to a particular relationship if only single type of relationship is possible.

Q3 Set of all supported capabilities from a node.<sup>1</sup>

Q4 Set of all supported relationships from a node.<sup>1</sup>

[1] Done as a side-effect when component is dropped/design is loaded on the canvas, serves as a cache until the component is deleted from the canvas/canvas is reset..

The results of the above queries directly related to actions the client needs to take. Non exhaustive set of actions: \

1.  Highlighting the nodes \
2.  Discerning the relationships possible b/w the involved nodes and getting ready for the next course of actions to take. Eg: Asking follow up questions. Creating dependent node on the canvas, marking the node as annotation.

3.  **Point in time the evals are invoked: \
    **An evaluation needs to occur whenever there’s a change on the canvas. The eval in this tier is lightweight and does not affect the design. Different events invoke different evaluation queries. For eg: `DRAG_START, EDGE_CREATION_STARTED` these events are concerned with queries Q1, Q2. Where as events like `component_drop `are concerned with queries Q3, Q4.
    _ DRAG START - (Q1, Q1)
    _ EDGE CREATION STARTED (Q1, Q2
    _ COMPONENT DROP START (Q3, Q4)
    _ INITIAL LOAD (Q3, Q4)

        _COLLABORATOR UPDATED CANVAS (Automatically results in DESIGN SAVE event hence no special consideration):_

4.  **Applying the results of eval**

The results are key-value pairs, where key are the actions to take and value are the node id’s on which action is to be performed.

_Note: The Tier-2 eval results can override Tier-1 eval results._

#### Tier-2: (First Part of Tier-2 eval) Server side eval:

REST API: POST request to `/api/meshmodels/relationships/evaluate`

```
Body: {
	pattern_file: designFile
	evaluation_queries: ["<relationship_kind>_<relationship_subtype>_relationship"] (Optional)
}
```

_Queries: The queries are directly related to relationships, the queries/evauationQueries links the relationships and policies, they have the following format: &lt;relationships_kind>_&lt;relationships*subtype>*

As the potential relationships are formed based on the results of Tier-1 evals, a event (client side) fires the Tier-2 eval that is performed in the Server.

_Input: Set of queries to perform eval for, inputs for the evals and the prettified design. \
Actor: Client (ctxAwareMachine) \
\_System: Server \
\_Output: Eval results in key: value pairs, where key is the query for which eval was invoked._

1.  **De-Prettification of the Design: \
    **The policy engine expects de-prettified version of design for its evaluation. \
    While the engine can work on prettified designs and do not error out, the results will be in-accurate.
2.  **Validation of queries: \
    **The server validates the incoming queries with the relationships registered in the registry. \
    The queries that fails validation are excluded from the eval. \
    If the requests contains no valid queries the eval is aborted at this stage and an event is transmitted.
3.  **Preparing for eval:**

    1. Preparing and sanitising the eval inputs. (Preparing inputs should be done in the policy itself in some pre-processing rule, before actual evaluation rule is invoked)
    2. Updating the OPA store (if required) with correct set of data (relationship definitions).

    _If the Actor do not explicitly specify the queries to invoke evaluation for, by default all registered queries are given to OPA engine for evaluation. _

4.  **Evaluating the queries: \
    **The eval is performed on the de-prettified design, with the provided set of inputs and data (relationship defs: initialized when the engine starts). \
     \
    The eval results are stored as key-value pairs, where _key_ corresponds to the query for which eval was invoked and the _value_ corresponds to the nodes and actions to be taken by the client on the nodes. \
     \
    In case when conflict occurs during the eval, and the policy fails to handle the conflict using the[ tie-breaking strategy](https://www.openpolicyagent.org/docs/latest/faq/#conflict-resolution)/technique, the results for that particular query is empty, (only the result i.e. _value is empty._ In the final response the _key_ (eval query) is still present. Eg: example_query: [])

        _Note: The evals are not recursive._

5.  **Prettifying the response: \
    **Because eval results for some queries outputs a diff of input design (in de-prettified format) and the client only understands a prettified version, the response is again prettified, ensuring lossless evals.

Example response:

edge_mount_relationship: {

    edges_to_add: [

    	from: &lt;id>,

    	to: &lt;id>,

    binded-by: &lt;id>

    ],

    edges_to_remove: [

    	from: &lt;id>,

    	to: &lt;id>,

    binded-by: &lt;id>

    ]

}

#### \

Tier-2: (Second Part of Tier-2 eval) Client Side:

**Applying the eval results:**

_Input: Tier-2 Server side eval results \
Actor: Client side interceptor_

The evaluation results consists of actions that the client needs to perform to complete the process. \
Based on the query appropriate interceptor functions are invoked, the interceptor functions are aware of the visual representation supported by the query/relationships.

Immediate action items:

- [Send an event/inform the client that the eval for this particular query has failed].
- [Ensure the response contains standard set of supported actions and the entity on which action needs to be performed.] \
  Exhaustive set of currently supported actions: add, delete, patch
- [Notify users about the changes occurred as part of tier-2 eval]
- [Add the evaluation response schema in the relationship definition and before starting the Tier-2 Client side stage, perform validation of the response against the schema.]

Who is responsible for the evaluation of a relationship?

**Relation Evaluation can be divided into multiple phases **:

- **Phase 1 : Feasibility Evaluation** , which relations are possible for a given snapshot of design ( **_client _**)
- **Phase 2 : Relationship Initialization** , making components commit to a relation (**_client_** )
- **Phase 3 : Relationship Application** , applying the effects of what is meant to be in a particular relation .

          Policy evaluates the  patching of the mutated refs .  (**_policy_**) **


          UI applies that patch delta for the affected components ( **_ui_** )

- **Phase 4 : Validation of relations** : includes evaluations necessary for checking if current relations are still valid or if their are any implicit relations available (**_policy_**)

                   Policy returns the components to add / delete or modify (**_policy_**)  [ needs more standardization ]

                   Ui applies the changes to reach a valid state (**_UI_**)  [ need a contract with server ]

- **Phase 5 : Notifying the user : ** as the relation evaluation goes through the multiple , it is necessary to keep the user notified / insync with the on going process

        Client generated intermittent notifications for the on going process **


        Client generated persistent notifications when the application of relation is complete *

### Synchronization of the evaluation Process

The whole evaluation process is asynchronous in nature and driven by events triggered from the ui . letting the individual phase to run independently without any synchronization will result in race conditions , application of stale evaluations , and hence corruption of the design .

Here is how we synchronize the process using the State machines ( context aware machine is responsible here )

- Only single evaluation can be asked at once from the server
- Server eval is only asked after for when the initialization of the relation is done
- Application of eval happens after the eval results come from the server , ( no further evals can be asked at this phase )

**Legend :**

** \* indicates the behavior which is only proposed here and is not present in the application at the moment . **

\*\* ** indicates the behavior which is not completed yet . **

**Client’s Responsibility**

Client should perform the task necessary to initialize the relationship, which includes identification of which nodes are receivable (the reduction of opacity of those nodes that are not capable of establishing a relationship).

As the server's eval results do not mutate the design directly , ui is responsible for applying the changes requested by the server.

**Server’s Responsibility**

Contract - Where is the API and the JSON format?

Race condition

Tie-breaking

Ability to specify which policies are to be included for consideration in an evaluation moment (e.g. all known policies or a select set of policies.

Ability to specify which nodes and which edges are to be included for consideration in an evaluation moment (e.g. all or specific nodes or specific edges).

Are component IDs included? Design IDs, User IDs?

    Component IDs are included.

Can we queue evaluations? ( by leveraging the sm properly we wont need to keep a queue )

Do we have an evaluation notification formatter in Notification Center?

# Model Evaluation

Models need to be evaluated in order to ensure their desired behavior is enforced. Evaluation might need to occur server-side or client-side. Server-side evaluation currently happens during pattern engine processing (DAG creation and execution). Client-side evaluation currently happens in MeshMap. The following depicts model evaluation in context of either server-side or client-side evaluation.

## Model Evaluation Algorithm

The following pseudo-code represents the general expression to be evaluated.

```
<policy-operator> AND <policy selectors> [

FOR ALL (<relationships>

AND ALL relevant <attributes>) ],

then APPLY <policy-action>
```

## Example Model Evaluation using Kubernetes Namespace

In this example scenario a user drags a node (namespaced Kubernetes resource) over another node (Kubernetes Namespace) causing a join / leave event to be evaluated. In this case, the nonself node is a Kubernetes Namespace.

```
match.presenceOf operator AND selectors.nonself.namespace selectors [
FOR ALL relationships ( self.hierarchical.parent

AND ALL relevant attributes component.type = compound-node ) ]

then APPLY policy allow

```

“To” and “from” are edge-centric descriptors that are somewhat ill-suited for representing other relationship types. Alternatives to “to” and “from” are:

- “self” and “other”
- “self” and “nonself”
- “initiator” and “receiver”

