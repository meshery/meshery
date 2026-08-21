---
title: Ensuring Extension Compatibility
description: A guide for Meshery extension developers on maintaining compatibility with the Meshery platform, covering shared frameworks, dependencies, and build tooling.
categories: [extensibility]
aliases:
- /extensibility/verify-compatibility
---

The Meshery ecosystem is designed for extensibility, with the core platform (`meshery/meshery`) providing extension points for a variety of plugins, collectively managed in the `meshery-extensions` GitHub organization. To ensure a stable and seamless user experience, it is critical that extensions remain compatible with the Meshery platform, especially when shared frameworks or dependencies are updated.

This guide outlines the key areas of compatibility and provides a general process for extension developers to follow when Meshery's underlying frameworks are updated.

{{% alert color="info" title="Meshery's Bill of Materials" %}}
Meshery's version is defined by its [Bill of Materials (BOM)]({{< ref "project/contributing/_index.md#meshery-contribution-flow" >}}), which specifies versions for Go, Node.js, and other key dependencies. Extension developers should reference the BOM to ensure alignment.
{{% /alert %}}

## Key Areas for Compatibility

Compatibility between Meshery and its extensions spans several areas. When the Meshery platform is updated, one or more of the following are likely to change:

1. **Language Runtimes**: Meshery and its extensions are often built using the same programming languages. Ensuring version alignment is the first step to compatibility.
    - **Go Version**: Meshery is primarily written in Go. Extensions written in Go must be compiled with the same major Go version to ensure binary compatibility.
    - **Node.js Version**: UI extensions may depend on the Node.js version used by Meshery's UI.

2. **Shared Go Libraries**: Meshery and its extensions share a common set of libraries (`meshery/meshkit`) that provide foundational utilities, API clients, and data models.
    - **MeshKit**: Extensions must depend on the same version of `meshkit` as the Meshery server to avoid mismatches in data structures and function signatures.
    - **gRPC and Protobuf**: Meshery uses gRPC for communication. Extensions must use the same Protobuf definitions (`.proto` files) and generated Go code to ensure API compatibility.
    - **Other Dependencies**: Align versions of other critical shared dependencies listed in Meshery's `go.mod`.

3. **UI Frameworks and Component Libraries**: UI-based extensions must align with Meshery's frontend technology stack to ensure a consistent user experience and functional correctness.
    - **React and Next.js**: Align with the versions used in Meshery's UI.
    - **Material-UI (MUI)**: Use the same version of the MUI component library to maintain visual consistency.

4. **Build and CI/CD Tooling**: The build process for extensions should be compatible with the environment and tooling used by Meshery.
    - **Makefiles**: Variables and targets in the core Meshery Makefiles (e.g., `install/Makefile.core.mk`) may be updated.
    - **Docker Images**: Base Docker images used for building and running Meshery may be updated, affecting the extension's build environment.

## General Process for Verifying Compatibility

When Meshery's Bill of Materials is updated, follow this general process to verify and update your extension.

### 1. Identify the Changes in Meshery

Start by reviewing the release notes or commit history of the `meshery/meshery` repository to understand what has changed. Pay close attention to:

- Go version updates in the root `go.mod` and `install/Makefile.core.mk`.
- Version changes for `meshery/meshkit` and other key dependencies in `go.mod`.
- Updates to UI framework versions in `ui/package.json`.
- Changes to `.proto` files for gRPC services.

### 2. Update Your Extension's Dependencies

In your extension's repository, update your dependency files to match the versions used in Meshery.

- **For Go-based extensions**:
    1. Update the `go` directive in your `go.mod` file to match Meshery's.
    2. Use the [`syncmodutil` utility](https://github.com/meshery/meshkit/tree/master/cmd/syncmodutil) (located in `meshery/meshkit`'s `cmd/syncmodutil`) to automate dependency synchronization. `syncmodutil` automates the manual dependency-sync process in a way that is safe for Go plugin builds.

        **Usage:**

        ```sh
        go run github.com/meshery/meshkit/cmd/syncmodutil@vX.X.X <path/to/meshery/go.mod> <path/to/extension/go.mod>
        ```

        Replace `vX.X.X` with the MeshKit version pinned in Meshery's `go.mod`. (Note: `syncmodutil`'s idempotent replace-block emission was introduced in MeshKit v1.0.20.)

        For CI-gate usage, you can include the `--err` flag to fail (with a non-zero exit code) instead of rewriting in place.

        **What it does:**

        - Rewrites matching `require` versions to the source's.
        - Emits one canonical `replace` block pinning the full source module graph (source replaces take precedence, verbatim); leaves destination-only replaces untouched. Note it is idempotent as of v1.0.20.

        ![syncmodutil flow](../images/go-mod-sync-flow.svg)

        **Caveats:**

        - Always build-test the destination module after syncing.
        - A destination that relies on a specific dependency version's behavior beyond its declared API contract may still break.
        - The emitted `replace` block is intentionally broad (pins the entire source graph); unused pins are harmless.
        - For detailed technical reference, always consult the [`syncmodutil` README](https://github.com/meshery/meshkit/tree/master/cmd/syncmodutil).

        **Worked Example:**
        The real Kanvas integration serves as a concrete case study. The [`layer5labs/meshery-extensions`](https://github.com/layer5labs/meshery-extensions) repository uses `Makefile` targets (`graphql-sync` / `graphql-sync-err`) which wrap `syncmodutil`; Meshery also runs these in CI (see [`meshery/meshery` `.github/workflows/go-mod-sync.yaml`](https://github.com/meshery/meshery/blob/master/.github/workflows/go-mod-sync.yaml)). This gives system integrators a real, currently-in-use reference implementation to model their own extension's sync automation on.

        **Manual Fallback:**
        If you are not using `syncmodutil`, you must manually update the required version of `github.com/meshery/meshkit` and other shared modules, and then run:

        ```sh
        go get github.com/meshery/meshkit@vX.X.X
        go mod tidy
        ```

        *Note: The manual fallback does not provide the `replace`-pinning required for strict Go plugin ABI compatibility.*

- **For UI-based extensions**:
    1. Update the versions of React, Next.js, MUI, and other relevant packages in your `package.json`.
    2. Run `npm install` or `yarn install` to update your dependencies.

### 3. Rebuild and Test Your Extension

After updating the dependencies, rebuild your extension to ensure it compiles without errors.

- **For Go-based extensions**:

    ```sh
    go build -o <your-extension-binary>
    ```

- **For UI-based extensions**:

    ```sh
    npm run build
    ```

### 4. Perform Integration Testing

Run your updated extension with the latest version of Meshery to verify end-to-end compatibility.

1. **Run the latest Meshery server**. Ensure you are using the build corresponding to the dependency updates.
2. **Run your updated extension**.
3. **Monitor the logs** of both Meshery and your extension for any runtime errors, panics, or warnings related to version mismatches.
4. **Thoroughly test your extension's functionality**. Verify that all interactions with Meshery are working as expected, including API calls and UI integrations.

By following this process, you can ensure that your extension remains compatible with the ever-evolving Meshery platform, providing a reliable and consistent experience for all users.
