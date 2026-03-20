---
title: "Extensibility: Build-time"
description: System integrators can supply build-time configuration, data, packages, and other extensions to Meshery's container image.
---

Meshery supports build-time extensibility, allowing system integrators to customize the Meshery container image during the Docker build process. This extension point enables organizations to inject custom configurations, data, provider extensions, and other resources directly into the Meshery container image at build-time.

## Build-time Extension Point

The build-time extension point allows you to supply custom content by placing a `.meshery` directory at the root of the Meshery repository before building the Docker image. The contents of this directory will be copied into the container image and made available to Meshery at runtime.

### How It Works

1. **Placement**: Create a `.meshery` directory at the root of the Meshery repository (same level as the `Dockerfile`).
2. **Content**: Add any build-time configuration, data, packages, or other resources you want to include in the container image to this directory.
3. **Build**: Run `docker build` as usual. The Docker build process will automatically copy the `.meshery` directory and its contents into the container image at `/home/appuser/.meshery/`.
4. **Runtime Access**: The injected content will be available to Meshery at runtime under `/home/appuser/.meshery/`.

### Directory Structure

```
meshery/                     # Repository root
├── .meshery/               # Build-time extension directory (optional)
│   ├── config/            # Custom configuration files
│   ├── provider/          # Provider-specific extensions
│   ├── content/           # Custom patterns, filters, applications
│   └── ...                # Any other custom content
├── install/
│   └── docker/
│       └── Dockerfile     # Main Dockerfile
└── ...
```

### Example Use Cases

#### Custom Configuration Files

Place custom configuration files that should be bundled with your Meshery deployment:

```bash
mkdir -p .meshery/config
echo "custom_setting=value" > .meshery/config/custom.conf
```

#### Provider-specific Extensions

Include proprietary or organization-specific provider extensions:

```bash
mkdir -p .meshery/provider/MyProvider/v1.0.0
cp /path/to/provider-extension.tar.gz .meshery/provider/MyProvider/v1.0.0/
```

#### Custom Content

Add organization-specific patterns, filters, or applications:

```bash
mkdir -p .meshery/content/patterns
cp /path/to/custom-patterns/* .meshery/content/patterns/
```

#### Pre-installed Models

Include pre-downloaded models or manifests that might not be available at runtime:

```bash
mkdir -p .meshery/models
cp -r /path/to/custom-models/* .meshery/models/
```

### Building with Build-time Extensions

To build a Meshery container image with your custom build-time extensions:

```bash
# Create your .meshery directory with custom content
mkdir -p .meshery/config
echo "my_config=value" > .meshery/config/custom.conf

# Build the Docker image (from the repository root)
docker build -f install/docker/Dockerfile -t meshery:custom .
```

The build process will:
1. Copy all files from the build context (including `.meshery` directory)
2. Extract the `.meshery` directory in a build stage
3. Copy the `.meshery` directory contents to `/home/appuser/.meshery/` in the final image

### Runtime Behavior

At runtime, Meshery will have access to:
- **Seed Content**: `/home/appuser/.meshery/seed_content/` (standard Meshery seed content)
- **Config**: `/home/appuser/.meshery/config/` (Meshery configuration files)
- **Custom Content**: Any additional content you placed in `.meshery/` during build

Your custom content in `.meshery/` will be merged with the standard Meshery directories. The contents are owned by the `appuser` user in the container.

### Important Notes

- **Optional**: The `.meshery` directory is optional. If not present during build, the Docker build will continue without errors.
- **Merge Behavior**: If you create subdirectories with the same names as existing Meshery directories (e.g., `config`), the contents will be merged.
- **Permissions**: All files copied from the build-time `.meshery` directory will be owned by `appuser:appuser` in the container.
- **Not in Version Control**: The `.meshery` directory at the repository root is typically not committed to version control. Add it to `.gitignore` if you're creating temporary build-time extensions.

### Security Considerations

When using build-time extensions:

- **Review Content**: Always review the contents of the `.meshery` directory before building to ensure no sensitive information or malicious content is included.
- **Secrets Management**: Do not include secrets, passwords, or API keys in the build-time `.meshery` directory. Use environment variables or runtime configuration for sensitive data.
- **Trust Sources**: Only include content from trusted sources to avoid introducing security vulnerabilities.

### Comparison with Runtime Extension Points

Meshery supports multiple [extension points](/extensibility) at different stages:

| Extension Point | When | Use Case |
|----------------|------|----------|
| **[Build-time](/extensibility/build-time)** | During `docker build` | Pre-package custom configurations, offline packages, proprietary plugins |
| **[Boot-time](/extensibility)** | Container startup | Dynamic configuration that can change between deployments |
| **[Runtime](/extensibility)** | While Meshery is running | User-installed plugins, dynamically loaded extensions |

Choose the build-time extension point when you need to:
- Bundle organization-specific customizations in a custom container image
- Include content that must be available before Meshery starts
- Create a self-contained deployment with all dependencies pre-packaged
- Deploy in air-gapped or restricted network environments

### Troubleshooting

**My `.meshery` directory isn't being copied**
- Ensure the directory is at the repository root (same level as `install/docker/Dockerfile`)
- Check that the directory isn't excluded in `.dockerignore`
- Verify the build context includes the repository root

**Permission denied errors at runtime**
- The build process should automatically set ownership to `appuser:appuser`
- If issues persist, check file permissions in your source `.meshery` directory

**Build fails with "no such file or directory"**
- This should not happen as the `.meshery` directory is optional
- If you see this error, check your Docker version and build logs
