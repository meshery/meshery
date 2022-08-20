### 0.3.9 / 2014-11-23

#### chruby.sh

* Stop searching `RUBIES` if an exact match is found. (@havenwood)

#### auto.sh

* Fixed a bug where `/.ruby-version` was being ignored. (@havenwood)

#### chruby-exec

* Ensure that all parameters are properly shell-escaped. (@havenwood)

#### scripts/setup.sh

* No longer install [ruby-install](https://github.com/postmodern/ruby-install#readme).

#### scripts/bug_report.sh

* Print `RUBYLIB`, `RUBYOPT`, `RUBYPATH` and `RUBYSHELL` env variables.
* Include `chruby-exec` in the versions section.
* Fall back to `-V` if `--version` did not work.

### 0.3.8 / 2013-12-04

#### chruby.sh

* Remove trailing slashes from ruby directories when iterating over `RUBIES`.
  (@halostatue)
* Ensure all temporary variables are local or unset.

#### auto.sh

* Ensure that `chruby_auto` can read `.ruby-version` files that do not end with
  a new-line. (@hosiawak)

#### scripts/setup.sh

* Install ruby-install 0.3.3.

#### scripts/bug_report.sh

* Print `$HOME`, `$RUBY_AUTO_VERSION`.
* Print `trap -p`, `$preexec_functions` and `$precmd_functions`.
* Print env variables even when they are empty.

### 0.3.7 / 2013-08-18

* Multiple style changes and optimizations. (@zendeavor)
* Safely glob the contents of `/opt/rubies` and `~/.rubies`.
  This prevents nullglob errors under zsh and `.rbx` directories from being
  added to `RUBIES`.
* Unset `GEM_PATH` in `chruby_reset` if it has become empty.
  Allows the RubyGems to use the default `GEM_PATH`.
* Safely quote `RUBIES[@]` to prevent implicit word-splitting when listing
 `RUBIES`.
* Map `-V` to `--version` in `chruby`. (@havenwood)
* Added benchmarks.

#### auto.sh

* Unset `RUBY_AUTO_VERSION` when loaded. Forces sub-shells to re-detect any
  `.ruby-version` file. (@KevinSjoberg)
* No longer export `RUBY_AUTO_VERSION`. Allows new windows in tmux to detect
  the `.ruby-version` file.
* Set `RUBY_AUTO_VERSION` even if `.ruby-version` contains an unknown Ruby.
  Prevents `chruby` from printing errors after every command.
* Fixed a typo where `RUBY_VERSION_FILE` was still being used. (@KevinSjoberg)

#### chruby-exec

* If stdin is a TTY, then spawn an interactive shell.

### 0.3.6 / 2013-06-23

* `chruby_use` no longer echos the selected Ruby.

#### chruby-exec

* Now runs under bash.
* Load `chruby.sh` for `CHRUBY_VERSION`.

#### auto.sh

* Record and compare the contents of `.ruby-version` files in order to detect
  modifications, such as when switching between branches.

### 0.3.5 / 2013-05-28

* Added a RPM spec.
* Respect `PREFIX` when auto-detecting `/opt/rubies/*`.
* Do not set `GEM_ROOT` if rubygems is not installed (Charlie Somerville).
* `chruby_use` now echos the select ruby and the version, only if the shell is
  in interactive mode (Brian D. Burns).
* `chruby_reset` no longer accidentally removes `/bin` if `GEM_HOME` or
  `GEM_ROOT` are empty (David Grayson).
* `chruby` now selects the last lexical match for the given ruby.

#### auto.sh

* Ensure that auto-switching works in non-interactive mode:
  * zsh: use `preexec_functions` which runs in both interactive and
    non-interactive sessions.
  * bash: use `trap DEBUG` which runs before every command, in both interactive
    and non-interactive mode. `PROMPT_COMMAND` only runs in interactive mode.
* Fixed a serious design flaw, where `chruby_auto` passed the contents of
  `.ruby-version` as multiple arguments to the `chruby` function. Originally,
  this allowed for `.ruby-version` files to specify additional `RUBYOPT` options
  (ex: `jruby --1.8`). However, an attacker could craft a malicious
  `.ruby-version` file that would require arbitrary code
  (ex: `1.9.3 -r./evil.rb`). The `./evil.rb` file would then be required when
  `ruby` is invoked by `chruby_use` in order to determine `RUBY_ENGINE`,
  `RUBY_VERSION`, `GEM_ROOT`.

  In order to prevent the abuse of this feature, `chruby_auto` now passes the
  entire contents of `.ruby-version` as a first and only argument to the
  `chruby` function.

  If you have `auto.sh` enabled, it is recommended that you upgrade.
  If you cannot upgrade, consider disabling `auto.sh`.
  If you want to scan your entire system for malicious `.ruby-version` files:

        find / -name .ruby-version 2>/dev/null | xargs -i{} grep -H " " {}

  Thanks to David Grayson for reporting this flaw.

#### scripts/setup.sh

* Do not assume bash is installed at `/bin/bash` (Shannon Skipper).

### 0.3.4 / 2013-02-28

* Prepend the new gem paths to `GEM_PATH` in `chruby_use`, instead of
  overriding the variable. This allows users to add common gem paths to
  `GEM_PATH` in `~/.bashrc`.
* Only remove the gem paths used by the Ruby in `chruby_reset`.

#### auto.sh

* Detect when `PROMPT_COMMAND=" "` before checking if `PROMPT_COMMAND` is an
  empty String. This appears to only happen on OSX Mountain Lion.

#### scripts/bug_report.sh

* Include `CHRUBY_VERSION` in the output.

### 0.3.3 / 2013-02-18

* Added `-v` `--version` options to `chruby` and `chruby-exec`.
* Added `scripts/bug_report.sh` for collecting environment information
  for bug reports.
* Initialize `RUBIES` to `()` to avoid double-loading `chruby.sh`.
* Invoke `ruby` using the absolute path to avoid shell aliases.
  This fixes a bug caused by [ohmyzsh] aliases.

#### auto.sh

* Unset `RUBY_VERSION_FILE` on initial load for [tmux].
* Remove trailing `;` and whitespace from `PROMPT_COMMAND` before
  appending `; chruby_auto`.

#### scripts/setup.sh

* Bump MRI version to 1.9.3-p385.
* Use `\x1b` instead of `\e` for OSX.

### 0.3.2 / 2013-01-15

* Prevent `auto.sh` from being loaded more than once.
* Recommend using `~/.bash_profile` and `~/.zprofile`.
* Use `cp` and `mkdir` instead of `install` in the `Makefile`.

#### chruby-exec

* Run under [bash], to avoid the [dash] shell.
* Invoke `$SHELL` with the `-i` option, so [zsh] will load shell configuration.

#### scripts/setup.sh

* Install [JRuby] 1.7.2.
* Use special `./configure` options for [homebrew].
* Also install openssl and readline via homebrew.

### 0.3.1 / 2012-12-29

* Fixed the auto-detection of `~/.rubies/*`.
* Check if `bin/ruby` exists and is executable before switching to a Ruby.
* Prevent `export=""` from accidentally being set under [zsh].
* Prevent `script/setup.sh` from exiting if a `brew install` fails because all
  packages are already installed.
* Updated the example `/etc/profile.d/chruby.sh` to only load under [bash]
  and [zsh].

### 0.3.0 / 2012-12-20

* Added the `chruby-exec` utility for use in `crontab` or with Continuous
  Integration (CI).
* Added support for auto-detecting Rubies installed into `/opt/rubies/` or
  `~/.rubies/`.
* Added `share/chruby/auto.sh`, which provides support for auto-switching
  to Rubies specified in the [.ruby-version](https://gist.github.com/1912050)
  file.
* Removed the "short circuit" check in `chruby_use`, to allow forcibly
  switching to the current Ruby, in case `PATH` or `GEM_PATH` become corrupted.

### 0.2.6 / 2012-12-18

* Forcibly switch to system Ruby when loading `share/chruby/chruby.sh`.
  This fixes switching issues for [tmux] users.

### 0.2.5 / 2012-12-15

* Renamed the `RUBY` environment variable to `RUBY_ROOT` to avoid breaking
  the `FileUtils#ruby` method in [rake](http://rake.rubyforge.org/).
* Do not unset `GEM_HOME`, `GEM_PATH`, `GEM_ROOT` if running under root.

### 0.2.4 / 2012-12-13

* Added a `Vagrantfile` for testing chruby in various environments.
* Changed all code and examples to reference `/opt/rubies/`.
* Ensure all error messages are printed to stderr.
* Refactored `scripts/setup.sh` to manually install all Rubies and install any
  dependencies via the System Package Manager.
* PGP signatures are now stored in `pkg/`.

#### Makefile

* Updated the `Makefile` to be compatible with BSD automake.
* Do not override `PREFIX`.
* Added a `test` task.

#### Homebrew

* Use `HOMEBREW_PREFIX`.
* Use `sha1` instead of `md5` (deprecated).
* No longer dynamically generate the example configuration.

### 0.2.3 / 2012-11-19

* Updated the `Makefile` to be compatible with the [dash] shell.
* Use inline substring substitutions instead of `sed`.

### 0.2.2 / 2012-11-17

* Use `typeset` to declare `RUBIES` as an indexed Array.
* Use the correct globbed Array syntax for both [zsh] and [bash].
* Improved the post-installation message in the [homebrew] recipe to auto-detect
  [RVM], [rbenv] and [rbfu].

### 0.2.1 / 2012-10-23

* Fixed `make install` to work on OS X.
* Added a [homebrew] recipe.

### 0.2.0 / 2012-10-16

* Install `chruby.sh` into `$PREFIX/share/chruby/`.

### 0.1.2 / 2012-08-29

* Check if `$RUBY` _and_ `$RUBYOPT` are different from the arguments passed to
  `chruby_use`.
* Fixed a spelling error in the README (thanks Ian Barnett).

### 0.1.1 / 2012-08-24

* Added unit-tests using [shunit2](http://code.google.com/p/shunit2/)
* Improved sanitation of `$PATH` in `chruby_reset`. (thanks mpapis)
* If the desired Ruby is already in use, immediately return from `chruby_use`.
* Export `$RUBY_ENGINE`, `$RUBY_VERSION`, `$GEM_ROOT` in `chruby_use`.

### 0.1.0 / 2012-08-18

* Added support for [zsh].
* Renamed the `$RUBY_PATH` variable to `$RUBY`.
* Set the `$RUBY_ENGINE` variable.
* Set the `$GEM_ROOT` variable to `Gem.default_dir`.
  This supports the custom RubyGems directory used by [Rubinius].
* Only initialize the `$RUBIES` variable if it does not have a value.

### 0.0.2 / 2012-08-14

* Added a `LICENSE.txt`.
* Added a `ChangeLog.md`.
* Updated the `Makefile` to generate proper tar archives.

### 0.0.1 / 2012-08-01

* Initial release.

[dash]: http://gondor.apana.org.au/~herbert/dash/
[bash]: http://www.gnu.org/software/bash/
[zsh]: http://www.zsh.org/
[tmux]: http://tmux.sourceforge.net/
[ohmyzsh]: https://github.com/robbyrussell/oh-my-zsh#readme

[Rubinius]: http://rubini.us/
[homebrew]: http://mxcl.github.com/homebrew/

[RVM]: https://rvm.io/
[rbenv]: https://github.com/sstephenson/rbenv#readme
[rbfu]: https://github.com/hmans/rbfu#readme

[MRI]: http://www.ruby-lang.org/en/
[JRuby]: http://jruby.org/
[Rubinius]: http://rubini.us/
