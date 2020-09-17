### Windows

#### Installing the `mesheryctl` binary

Download and unzip `mesheryctl` from the [Meshery releases](https://github.com/layer5io/meshery/releases/latest) page. Add `mesheryctl` to your PATH for ease of use. Then, execute:

```bash
./mesheryctl system start
```

#### Scoop

Use [Scoop](https://scoop.sh) to install Meshery on your Windows machine.

**Installing with Scoop**

Add the Meshery Scoop Bucket and install:

```bash
scoop bucket add mesheryctl https://github.com/layer5io/scoop-bucket.git
scoop install mesheryctl
```

**Upgrading with Scoop**

To upgrade `mesheryctl`, execute the following command:

```bash
scoop update mesheryctl
```

# Advanced Installation

Users can control the specific container image and tag (version) of Meshery that they would like to run by editing their local `~/.meshery/meshery.yaml` (a docker compose file).
Aligned with the Meshery container image, instead of leaving the implicit `:stable-latest` tag behind image: layer5/meshery, users will instead identify a specific image tag like so:

```bash
version: '3'
services:
  meshery:
    image: layer5/meshery:v0.5.0
    labels:
      - "com.centurylinklabs.watchtower.enable=true"
```

---

When Meshery is up and running, instructions to access Meshery will be printed on the screen and your default browser should be directed to the Meshery login screen.