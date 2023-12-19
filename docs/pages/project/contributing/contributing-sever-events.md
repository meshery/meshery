---
layout: default
title: Contributing to Meshery Server Events
abstract: Guide is to help backend contributors send server events using Golang.
permalink: project/contributing/contributing-server-events
type: project
category: contributing
language: en
list: include
---

Meshery incorporates an internal events publication mechanism that provides users with real-time updates on the processes occurring within the Meshery server when interacting with its endpoints. It ensures that users are kept in the loop regarding the ongoing activities within the API, and guides users towards future steps to resolve issues. This guide will provide step-by-step instructions on sending events from the server, including when to trigger events and what information to include.


First, let's take a look at how the event object is constructed,

<pre class="codeblock-pre">
 <div class="codeblock">
 eventBuilder := events.NewEvent()
                  .ActedUpon(UUID)
                  .FromUser(UUID)
                  .FromSystem(UUID)
                  .WithCategory(string)
                  .WithAction(string)
                  .WithSeverity(EventSeverity)
                  .WithDescription(string)
                  .WithMetadata(map[string]interface{})
                  .Build()
  </div>
 </pre>
&nbsp;&nbsp;

*Note: `events` is a package [github.com/meshery/meshkit/models/events](https://github.com/meshery/meshkit) from MeshKit containing source code of event related functionality*

The event mechanism utilizes [builder pattern](https://en.wikipedia.org/wiki/Builder_pattern) to construct event objects, `events.NewEvent()` creates an instance of [EventBuilder](https://github.com/meshery/meshkit/blob/ea3c60907a1cd1902902a4113206579992772083/models/events/build.go#L9) type, which functions as a builder class for constructing [Event](https://github.com/meshery/meshkit/blob/ea3c60907a1cd1902902a4113206579992772083/models/events/events.go#L37) objects. 

- `ActedUpon(UUID)` : it takes [UUID](https://pkg.go.dev/github.com/gofrs/uuid#UUID) of the resource being provisioned on the server, and the events are associated with that resource.
- `FromUser()` : it takes [UUID](https://pkg.go.dev/github.com/gofrs/uuid#UUID) of the user initiating the HTTP request to the server. This enables the server to publish events intended for that specific user in response.
- `FromSystem(UUID)` : it takes [UUID](https://pkg.go.dev/github.com/gofrs/uuid#UUID) of the running meshery instance (Handler instance contains it [`h.SystemID`](https://github.com/meshery/meshery/blob/c73d3687da3eb3b2aaaabbd27e91f99906fc5838/server/handlers/handler_instance.go#L29)).
- `WithCategory` : The string argument to specify under which category this event falls, For example, when a user is provisioning a Kubernetes connection, the specified category for the event is "connection.
- `WithAction` : The string argument to specify the type of action the server is performing, for instance, when registering a Kubernetes connection, the action would be "register".

*Note: categories and actions must be in snake case,*
      example categories: **"signup_request", "github_artifacts"**.
      example actions: **"remove_from_organization", "add_to_organization"**.

- `WithSeverity`: it takes [EventSeverity](https://github.com/meshery/meshkit/blob/ea3c60907a1cd1902902a4113206579992772083/models/events/events.go#L75) (underlying type is string), to specify severity of the event to bring user's attention to the event accordingly.

*Note: In certain conditions you must add some fields with specific  keys:*
- If the severity is "Error", include a field in Metadata using the key `err` with the corresponding error value.
- When you want to add a link to meshery docs, include a field in Metadata using the key `doc` with the corresponding link value.

- `WithDescription` : The string argument provides a detailed, descriptive message that depicts the specific operation or action being executed on the server.
- `WithMetadata`: it takes a Map `map[string]interface{}` data structure containing any supplementary information that the developer/contributor deems essential for the user to be informed about.
- `Build` : returns the [Event](https://github.com/meshery/meshkit/blob/ea3c60907a1cd1902902a4113206579992772083/models/events/events.go#L37) instance constructed upon the previously described functions and prepares it for publication through the [Broadcast](https://github.com/meshery/meshery/blob/1b5d78ed34648e0a91df8c2273026b930f748fbc/server/models/event_broadcast.go#L14), which is responsible for disseminating events.



## An Example in Code

<pre class="codeblock-pre">
 <div class="codeblock">
func (h *Handler) KubernetesPingHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	userID := uuid.FromStringOrNil(user.ID)
	token, ok := req.Context().Value(models.TokenCtxKey).(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "failed to get the token for the user")
		return
	}

	connectionID := req.URL.Query().Get("connection_id")
	eventBuilder := events.NewEvent().FromUser(userID).FromSystem(*h.SystemID).WithCategory("connection").WithAction("ping")
	if connectionID != "" {
		// Get the context associated with this ID
		k8sContext, err := provider.GetK8sContext(token, connectionID)
		eventBuilder.ActedUpon(uuid.FromStringOrNil(connectionID))
		if err != nil {
			eventBuilder.WithSeverity(events.Error).WithDescription(fmt.Sprintf("Encountered an issue while retrieving kubernetes context for connection ID `%s`", connectionID)).
			WithMetadata(map[string]interface{}{
				"error": err,
			})
			event := eventBuilder.Build()
		    _ = provider.PersistEvent(event)
		    go h.config.EventBroadcaster.Publish(userID, event)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "failed to get kubernetes context for the given ID")
			return
		}
		eventBuilder.WithSeverity(events.Success).WithDescription("Retrieved kubernetes context successfully").
		WithMetadata(map[string]interface{}{
			"Context Name":   k8sContext.Name,
			"Server Address": k8sContext.Server,
		})
		event := eventBuilder.Build()
		_ = provider.PersistEvent(event)
		go h.config.EventBroadcaster.Publish(userID, event)

		// Remaining code ...
    }
}
 </div>
</pre>
&nbsp;&nbsp;


In the given code block, the "resource" refers to the connection on which the server is performing the "ping" action. The "ActedUpon" field is supplied with the "connectionID" value obtained from the query parameter. `provider.GetK8sContext` returns the [K8sContext](https://github.com/meshery/meshery/blob/995ab671a12013088f874430cfa2c0f025b073d2/server/models/k8s_context.go#L26) object which serves as a representation of a Kubernetes context within the Meshery. if `provider.GetK8sContext` returns an error, the event is constructed with severity [Error](https://github.com/meshery/meshkit/blob/ea3c60907a1cd1902902a4113206579992772083/models/events/events.go#L18), description and enriched with metadata fields to provide the user with a complete context about the event.

If `provider.GetK8sContext` succeeds, an event with a [Success](https://github.com/meshery/meshkit/blob/ea3c60907a1cd1902902a4113206579992772083/models/events/events.go#L21) severity level is constructed. The event includes a success message and metadata containing the context name and the address of the Kubernetes API server, obtained from the operation. This information aims to provide users with relevant details about the successful operation.

after constructing the Event object using the Build function, event is stored into database for records and subsequently published in a new goroutine through [EventBroadcaster](https://github.com/meshery/meshery/blob/995ab671a12013088f874430cfa2c0f025b073d2/server/models/handlers.go#L241) member variable of [HandlerConfig](https://github.com/meshery/meshery/blob/995ab671a12013088f874430cfa2c0f025b073d2/server/models/handlers.go#L206), which is part of [Handler](https://github.com/meshery/meshery/blob/995ab671a12013088f874430cfa2c0f025b073d2/server/handlers/handler_instance.go#L18) instance.


## Suggested Reading

To gain a deeper understanding and examine events more closely, you can explore the source code files [k8config_handler.go](https://github.com/meshery/meshery/blob/master/server/handlers/k8sconfig_handler.go) and [design_engine_handler.go](https://github.com/meshery/meshery/blob/master/server/handlers/design_engine_handler.go) within the Meshery project's codebase.