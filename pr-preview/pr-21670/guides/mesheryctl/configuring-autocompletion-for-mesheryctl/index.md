# Configuring Autocompletion for `mesheryctl`

> Bash, Zsh, Oh My Zsh, and fish autocompletion for `mesheryctl` commands.

Source: /pr-preview/pr-21670/guides/mesheryctl/configuring-autocompletion-for-mesheryctl/

If you would like to have `mesheryctl` commands automatically completed for use as you use `mesheryctl`, then use the following instructions to configure automatic completion within your environment.

## Autocompletion for Bash

### bash <= 3.2

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">source /dev/stdin <<< "$(mesheryctl completion bash)"</div></div>
 </pre>

### bash >= 4.0

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">source <(mesheryctl completion bash)</div></div>
 </pre>

### bash <= 3.2 on MacOS

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">brew install bash-completion # ensure you have bash-completion 1.3+
mesheryctl completion bash > $(brew --prefix)/etc/bash_completion.d/mesheryctl</div></div>
 </pre>

### bash >= 4.0 on MacOS

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">brew install bash-completion@2
mesheryctl completion bash > $(brew --prefix)/etc/bash_completion.d/mesheryctl</div></div>
 </pre>

## Autocompletion for zsh

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">source <(mesheryctl completion zsh)</div></div>
 </pre><br>

If shell completion is not already enabled in your environment you will need to enable it. You can execute the following once:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">~/.zshrc > echo "autoload -U compinit; compinit"</div></div>
 </pre>

_Note_ : You might need to restart your shell for this setup to take effect.

#### zsh on MacOS and Oh My zsh

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">COMPLETION_DIR=$(echo $fpath | grep -o '[^ ]*completions' | grep -v cache) && mkdir -p $COMPLETION_DIR && mesheryctl completion zsh > "${COMPLETION_DIR}/_mesheryctl"</div></div>
 </pre>

### Autocompletion for fish

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">mesheryctl completion fish | source</div></div>
 </pre><br>

To load fish shell completions for each session, execute once:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">mesheryctl completion fish > ~/.config/fish/completions/mesheryctl.fish</div></div>
 </pre>

For an exhaustive list of `mesheryctl` commands and syntax:

- See [`mesheryctl` Command Reference](/pr-preview/pr-21670/reference/references/mesheryctl/).

Guides to using Meshery's various features and components.

<div class="related-discussions">
  <h3>Recent Discussions with "mesheryctl" Tag</h3><ul><li>
          <a href="https://discuss.meshery.io/t/pairing-up-with-a-meshmate-to-gacefully-start-contributing-to-meshery-and-its-project-ecosystem/6957" target="_blank" rel="noopener noreferrer">
            Pairing up with a Meshmate to gacefully start contributing to Meshery and its project ecosystem
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/development-meeting-contributing-to-meshery-cli-april-30-2025/6927" target="_blank" rel="noopener noreferrer">
            [Development Meeting] Contributing to Meshery CLI - April 30, 2025
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/meshery-development-meeting-april-30th-2025/6926" target="_blank" rel="noopener noreferrer">
            Meshery Development Meeting | April 30th, 2025
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/newcomers-meeting-end-to-end-testing-in-meshery-cli-using-bats-april-17-2025/6897" target="_blank" rel="noopener noreferrer">
            [Newcomers’ Meeting] End-to-End Testing in Meshery CLI using BATs – April 17, 2025
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/unable-to-lint-mesheryctl/6854" target="_blank" rel="noopener noreferrer">
            unable to lint mesheryctl 
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/mesheryctl-system-login-problem/6687" target="_blank" rel="noopener noreferrer">
            Mesheryctl system login problem
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/error-while-launching-the-meshery-dashboard/6600" target="_blank" rel="noopener noreferrer">
            Error while launching the Meshery Dashboard
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/cant-find-the-file-path-for-meshery-designs/6319" target="_blank" rel="noopener noreferrer">
            Can&#39;t find the file path for meshery designs
          </a>
        </li></ul><p>
    <a href="https://discuss.meshery.io/tag/mesheryctl" target="_blank" rel="noopener noreferrer">
      View all discussions tagged with <code>mesheryctl</code>
    </a>
  </p>
</div>
