# How to write MeshKit compatible errors

> How to declare errors in Meshery components.

---

LLMS index: [llms.txt](/llms.txt)

---

Meshery pervasively uses MeshKit as a golang and infrastructure management-specific library in all of its components. MeshKit helps populate error messages with a uniform and useful set of informative attributes.

To help with creating error codes, MeshKit contains a tool that analyzes, verifies and updates error codes in Meshery source code trees. It extracts error details into a file that can be used for publishing all error code references on the Meshery [error codes reference page](/reference/error-codes/). The objective to create this was to avoid centralized handling of error codes and automating everything

In order to create a Meshery error object, you will need to create a custom wrapper object for the native golang error. This can be done from the <a href="https://github.com/meshery/meshkit/tree/master/errors">MeshKit Error</a> package.

This tool will create a couple of files, one of them is designed to be used to generate the error reference on the Meshery Documentation website. The file errorutil_analyze_summary.json contains a summary of the analysis, notably lists of duplicates etc.

## Conventions to follow while creating errors

1. Errors names and codes are namespaced to components, i.e. they need to be unique within a component, which is verified by this tool.
1. Errors are not to be reused across components and modules.
1. Error codes are not to be set as integer. CI will take care of updating Error codes from a string to an integer.
1. Running `make error` analyzes your code and returns any warnings to be aware of.
1. Capitalize the first letter of the every error description.
1. Using errors.NewDefault(...) is deprecated. This tool emits a warning if its use is detected.
1. Use errors.NewV3(...) for an error whose causes each need their own remediations (see [Errors with several causes](#errors-with-several-causes-errorsnewv3)); otherwise errors.NewV2(...) or errors.New(...).
1. Use errors.New(...) from MeshKit to create actual errors with all the details.
  This is often done in a factory function. It is important that the error code variable is used here, not a literal.
  Specify detailed descriptions, probable causes, and remedies. They need to be string literals, call expressions are ignored.
  This tool extracts this information from the code and exports it.
  For the Code argument in the errors.New use the same Error name and append a "Code" after it. e.g error name : ErrApplyManifest then the error code is ErrApplyManifestCode
1. Set the value to any string, like "replace_me" (no convention here), e.g. ErrApplyManifestCode = "replace_me".
1. By convention, error codes and the factory functions live in files called error.go. The tool checks all files, but updates only error.go files.

Use the `errors.New()` function to create a new instance of the error object and pass situation-specific attributes as function arguments. These attributes are:

- Code
- Short Description
- Long Description
- Probable Cause
- Suggested Remediation

### Syntax



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">errors.New(ErrExampleCode, errors.Alert, []string{&#34;&lt;short-description&gt;&#34;}, []string{&#34;&lt;long-description&gt;&#34;}, []string{&#34;&lt;probable-cause&gt;&#34;}, []string{&#34;&lt;suggested remediation&gt;&#34;})</code>
	</div>
</pre>

## Example

In this example we are creating an Error for being unable to marshal JSON



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">var (
	// Error code
    ErrMarshalCode= &#34;replace_me&#34;

    //Static errors (for example)
    ErrExample = errors.New(ErrExampleCode, errors.Alert, []string{&#34;&lt;short-description&gt;&#34;}, []string{&#34;&lt;long-description&gt;&#34;}, []string{&#34;&lt;probable-cause&gt;&#34;}, []string{&#34;&lt;suggested remediation&gt;&#34;})
)

// Dynamic errors
//Error Name
func ErrMarshal(err error, obj string) error {
	return errors.New(ErrMarshalCode, errors.Alert, []string{&#34;Unable to marshal the : &#34;, obj}, []string{err.Error()}, []string{}, []string{})
}</code>
	</div>
</pre>


### Replacing old Error Codes

Old



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">bd, err := json.Marshal(providers)
	if err != nil {
		http.Error(w, &#34;unable to marshal the providers&#34;, http.StatusInternalServerError)
		return
	}</code>
	</div>
</pre>


New



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">bd, err := json.Marshal(providers)
  if err != nil {
          marshalErr := ErrMarshal(err, &#34;providers&#34;)
          h.log.Error(marshalErr)
          writeMeshkitError(w, marshalErr, http.StatusInternalServerError)
          return
      }</code>
	</div>
</pre>


`http.Error` is rejected by CI in the `./server` module: it writes a plain-text
body and strips the MeshKit code, severity, and remediation that clients parse.
See [HTTP Error Response Contract](/project/contributing/error-contract/)
for the response shape, the `writeMeshkitError` / `writeJSONError` helpers, and
how to choose the status code.

## Errors with several causes: `errors.NewV3`

`errors.New` takes probable causes and suggested remediations as two flat lists, and nothing tells a reader which remediation answers which cause. When an error has causes that each need their own remediations, declare it with `errors.NewV3`. Each cause then carries the remediations that address it, and a remediation can name an action a client may run and the policy for running it ([error definition contract](https://github.com/meshery/schemas/blob/master/docs/error-definition-contract.md)).



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">// remRetryConnection answers every cause, so it is declared once.
var remRetryConnection = errors.Remediation{
	ID:          &#34;retry-connection&#34;,
	Description: &#34;Once the cause is fixed, reconnect the connection.&#34;,
	Action:      &#34;connection.reconnect&#34;,
	ActionPolicy: errors.RequiresApproval,
}

func ErrUnreachableCluster(err error, server string) error {
	e := errors.NewV3(ErrUnreachableClusterCode, errors.Alert,
		[]string{&#34;Error communicating with the Kubernetes API server.&#34;},
		[]string{&#34;Meshery Server did not get an answer from the cluster&#39;s API server.&#34;},
		[]errors.Cause{
			{
				ID:          &#34;api-server-unreachable&#34;,
				Description: &#34;The Kubernetes API server is not reachable.&#34;,
				Remediations: []errors.Remediation{
					{ID: &#34;check-network&#34;, Description: &#34;Verify network connectivity to the API server.&#34;},
					remRetryConnection,
				},
			},
			{
				ID:          &#34;credentials-rejected&#34;,
				Description: &#34;The cluster rejected the connection&#39;s credentials.&#34;,
				Remediations: []errors.Remediation{
					{ID: &#34;renew-credential&#34;, Description: &#34;Renew the connection&#39;s credential.&#34;,
						Action: &#34;connection.refresh-credentials&#34;, ActionPolicy: errors.Automatic},
					remRetryConnection,
				},
			},
		})
	// Instance text is appended after NewV3 returns.
	e.LongDescription = append(e.LongDescription, fmt.Sprintf(&#34;API server: %s&#34;, server), err.Error())
	return e
}</code>
	</div>
</pre>


`NewV3` returns a `*errors.ErrorV2` whose `Causes` hold the nesting and whose flat `ProbableCause` and `SuggestedRemediation` are derived from it. `GetCause`, `GetRemedy`, `Error()` and every reader of the flat lists keep working.

`errorutil` reads these definitions from source, so they follow strict rules:

- **Literals only.** Write every `errors.Cause` and `errors.Remediation` as a keyed composite literal with string-literal `ID`, `Description` and `Action`. Every element of the short- and long-description lists is a string literal too: build instance text, such as an address or `err.Error()`, into `LongDescription` after `NewV3` returns, as above.
- **Stable ids.** Ids are lowercase words joined by hyphens (`api-server-unreachable`) and never change once published: events name the causes that applied in `metadata.error.matchedCauses`, and each code's reference page anchors on them. Cause ids are unique within an error, and every cause has at least one remediation.
- **Shared remediations.** A remediation that answers several causes is declared once as a package-level `var` (in any file of the package) and referenced from each cause. `errorutil` resolves the variable, and the flat remediation list carries it once.
- **Action policies.** `errors.Manual` (also when empty) means a person applies the remediation. `errors.RequiresApproval` lets a client run the `Action` after a person approves it, and `errors.Automatic` lets it run the `Action` on its own and say so. Both need an `Action`, dot-separated lowercase words such as `connection.reconnect`. An unknown policy is never coerced: `errors.ValidateCauses` and `errorutil` reject it, and a reader treats it as manual.
- **Tests.** Unit-test a definition with `errors.ValidateCauses(errors.GetCauses(err))`. `TestErrUnreachableKubeAPICauses` in `server/models/event_error_test.go` is an example.

**An invalid definition fails the build.** `errorutil` lists every `NewV3` definition it rejects in `errorutil_analyze_summary.json` under `invalid_definitions`, logs each one with its problems, leaves it out of both exports (never in a reduced form), and exits non-zero from `analyze` and `update`. The `error-codes-updater` workflow then fails the pull request. `errors.New` and `errors.NewV2` definitions are unaffected.

When the server can tell which cause applies to an occurrence, record it where the underlying error is still in hand, so the error event names it. `ErrUnreachableKubeAPI` (`server/models/error.go`) classifies the Kubernetes client error, and `models.AttachEventError` puts the matched cause on the event; see [Contributing to Meshery Server Events](/project/contributing/contributing-server-events/).

## Replacing logrus

There already exists an [interface for logger](https://github.com/meshery/meshkit/blob/master/logger/logger.go) in MeshKit.<br><br>

<div class="alert alert-warning" role="alert"><div class="h4 alert-heading" role="heading">WARNING</div>


To enforce the use of meshkit errors, the meshkit logger reads its structured fields - code, severity, probable cause, suggested remediation - straight off a meshkit error. A plain Go error still logs, but every one of those fields renders as `None`, leaving the operator a bare message with no code to look up and no remediation to follow. Wrap an error in a meshkit error before logging it. See: [meshkit#119](https://github.com/meshery/meshkit/pull/119) for more insight.
</div>


#### Defining a Logger



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">type Logger struct {
	log   logger.Handler
}</code>
	</div>
</pre>


#### Debug

##### Old



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">logrus.Debugf(&#34;meshLocationURL: %s&#34;, meshLocationURL)</code>
	</div>
</pre>


##### New



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">l.log.Debug(&#34;meshLocationURL: &#34;, meshLocationURL)</code>
	</div>
</pre>


#### Error

##### Old



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">logrus.Errorf(&#34;error marshaling data: %v.&#34;, err)</code>
	</div>
</pre>


##### New



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">l.log.Error(ErrMarshal(err, obj))</code>
	</div>
</pre>


## A small program using meshkit errors and logger



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">package main

import (
	&#34;fmt&#34;
	&#34;os&#34;

	meshkitErrors &#34;github.com/meshery/meshkit/errors&#34;
	&#34;github.com/meshery/meshkit/logger&#34;
	&#34;github.com/sirupsen/logrus&#34;
	&#34;github.com/spf13/viper&#34;
)

var (
	// CI will replace test_code with new error code
	ErrOpeningFileCode = &#34;test_code&#34;
)

func main() {
	logLevel := viper.GetInt(&#34;LOG_LEVEL&#34;)
	if viper.GetBool(&#34;DEBUG&#34;) {
		logLevel = int(logrus.DebugLevel)
	}
	log, err := logger.New(&#34;test&#34;, logger.Options{
		Format:   logger.SyslogLogFormat,
		LogLevel: logLevel,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// logging meshkit error
	err = openFileWithMeshkitError(&#34;some.txt&#34;)
	if err != nil {
		log.Error(err)
	}
	// OUTPUT
	// ERRO[2021-11-10T17:31:28+05:30] open some.txt: no such file or directory
	// app=test code=1001 probable-cause=&#34;empty string passed as argument .file with this name doesn&#39;t exist&#34;
	// severity=2 short-description=&#34;unable to open file&#34; suggested-remediation=&#34;pass a non-empty string as
	// filename .create file before opening it&#34;

	// logging non meshkit error
	err = openFile(&#34;some.txt&#34;)
	if err != nil {
		log.Error(err)
	}
	// OUTPUT
	// ERRO[2024-07-01T19:09:09+05:30] open some.txt: no such file or directory
	// app=test code= probable-cause= severity=0 short-description= suggested-remediation=

}

// this returns a non meshkit error
func openFile(name string) error {
	_, err := os.Open(name)
	return err
}

// this returns a meshkit error
func openFileWithMeshkitError(name string) error {
	_, err := os.Open(name)
	return ErrOpeningFile(err)
}

func ErrOpeningFile(err error) error {
	return meshkitErrors.New(ErrOpeningFileCode, meshkitErrors.Alert, []string{&#34;unable to open file&#34;},
		[]string{err.Error()},
		[]string{&#34;empty string passed as argument &#34;, &#34;file with this name doesn&#39;t exist&#34;},
		[]string{&#34;pass a non-empty string as filename &#34;, &#34;create file before opening it&#34;})
}</code>
	</div>
</pre>



## Generating error codes in meshery/meshery

`make error` regenerates codes for the server but **skips `mesheryctl`**. The two components keep the same contract in different files.

- **`mesheryctl`**: take the next code from `mesheryctl/helpers/component_info.json` (`next_error_code`) and bump that value in the same commit.
- **Server**: the same contract lives in `server/helpers/component_info.json`. `errorutil` refuses to run at all - "next_error_code is lower than or equal to highest used code" - until `next_error_code` is bumped past every code you added, so bump it in the same commit.

`.github/workflows/error-codes-updater.yaml` re-runs `errorutil` on every pull request and fails it if the analysis reports anything.

### Naming and formatting

Name each constant `<BuilderFuncName>Code` - `errorutil` keys the export off that pairing. Adding a constant longer than the block's current widest name makes `gofmt` realign the entire `error.go` const block, so prefer a shorter name over a 300-line whitespace diff.

### Regenerating the docs reference

`errorutil` writes four files to the output directory, all gitignored: `errorutil_analyze_errors.json` and `errorutil_analyze_summary.json` (the analysis: duplicates, convention violations and `invalid_definitions`) and two exports:

- `errorutil_errors_export.json`, the legacy export. Its shape is frozen (snake_case keys, joined strings) because the docs shortcodes read it; `NewV3` entries additionally carry `probableCauses`.
- `errorutil_errors_index.json`, the component's definitions as a schemas `ErrorReferenceIndex` (`schemaVersion` 2): `code` (`meshery-server-1304`), `moniker`, `component`, `numericCode`, severity, descriptions, nested causes and the derived flat lists. `errorutil` validates every entry against the schema before writing it; the docs build assigns `referenceUrl`.

The reference data built from them is tracked:

- `docs/data/errorref/meshery-server_errors_export.json`, from the legacy export. Regenerate it with the `jq --slurpfile` wrapper the workflow uses, or the published error reference silently omits the new codes.
- `docs/data/errorindex/meshery-server_errors_index.json`, from the nested export (`jq . server/helpers/errorutil_errors_index.json`). The docs build takes a code's causes from it, else from the legacy export's `probableCauses`, and otherwise shows the flat lists.

The workflow regenerates both. The docs build publishes every component's codes as `reference/error-codes/index.json`, a schemas `ErrorReferenceIndex` (`schemaVersion` 2).

### Rendering errors to the user

Only `utils.Log.Error(err)` renders a MeshKit error's code, cause and remediation; cobra's default print shows just the message. In `mesheryctl` commands, log the structured error for the user *and* return it for the exit path.
