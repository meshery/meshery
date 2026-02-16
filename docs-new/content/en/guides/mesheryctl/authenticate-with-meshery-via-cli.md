---
title: Authenticating Meshery via CLI
categories: [mesheryctl]
description: Get your authentication token from Meshery CLI.
---

To authenticate Meshery through `mesheryctl`, the Meshery CLI, you will use the `mesheryctl system login` command. On executing this command, you will be provided with a list of providers to choose from. You can then select a Provider of your choice to complete the authentication and authorization process.

As of this writing, you will be presented with two providers, _Layer5_ and _None_. 
```bash
Use the arrow keys to navigate: ↓ ↑ → ← 
? Select a Provider: 
  ▸ Layer5
    None
```

- Selecting _Layer5_ will open a browser to complete the login and authentication process with Layer5 cloud. On successful authentication, you can close the window and return to the command prompt. 

  Verify that an `auth.json` file was created in the `.meshery` folder in your home directory.

  ```bash
  ls -l $HOME/.meshery/

  total 12
  -rw-rw-r-- 1 ubuntu ubuntu  39 Dec 21 06:13 auth.json
  -rw-rw-r-- 1 ubuntu ubuntu 260 Dec 21 06:04 config.yaml
  -rw-rw-r-- 1 ubuntu ubuntu 988 Dec 21 06:04 meshery.yaml
  ```

  **_The need for authentication to `Layer5` [provider](/extensibility/providers) is to save your environment setup while also having persistent/steady sessions and to be able to retrieve performance test results._**

- Selecting _None_ will create an empty `auth.json` file. All your work remains local and ephemeral. 

If `mesheryctl` is running in a system that does not have a browser, you can download an auth token file from your Layer5 cloud account and copy it into the `.meshery` folder in your home directory. The following steps show how you can generate and download a token:

1. Navigate to [https://cloud.layer5.io/security/tokens](https://cloud.layer5.io/security/tokens) and sign-in.
Ensure you are in the right organization and click **Create**.

<a href="./images/create-token.png"><img alt="Create Token" src="./images/create-token.png" /></a>

2. Provide a token name and purpose. Click **Create** to generate.

    <a href="./images/generate-token.png"><img alt="Generate Token" src="./images/generate-token.png" /></a>x

3. Click the **Download** icon to download the `auth.json` file.

    <a href="./images/download-token.png"><img alt="Download Token" src="./images/download-token.png" /></a>

Then run `mesheryctl system check` to ensure you do not see an authentication error.   

For an exhaustive list of `mesheryctl` commands and syntax, visit [`mesheryctl` Command Reference](/reference/mesheryctl).

{{< related-discussions tag="mesheryctl" >}}

