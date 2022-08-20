# chruby

[![Build Status](https://travis-ci.org/postmodern/chruby.svg?branch=master)](https://travis-ci.org/postmodern/chruby)

Changes the current Ruby.

## Features

* Updates `$PATH`.
  * Also adds RubyGems `bin/` directories to `$PATH`.
* Correctly sets `$GEM_HOME` and `$GEM_PATH`.
  * Users: gems are installed into `~/.gem/$ruby/$version`.
  * Root: gems are installed directly into `/path/to/$ruby/$gemdir`.
* Additionally sets `$RUBY_ROOT`, `$RUBY_ENGINE`, `$RUBY_VERSION` and
  `$GEM_ROOT`.
* Optionally sets `$RUBYOPT` if second argument is given.
* Calls `hash -r` to clear the command-lookup hash-table.
* Fuzzy matching of Rubies by name.
* Defaults to the system Ruby.
* Optionally supports auto-switching and the `.ruby-version` file.
* Supports [bash] and [zsh].
* Small (~90 LOC).
* Has tests.

## Anti-Features

* Does not hook `cd`.
* Does not install executable shims.
* Does not require Rubies be installed into your home directory.
* Does not automatically switch Rubies by default.
* Does not require write-access to the Ruby directory in order to install gems.

## Requirements

* [bash] >= 3 or [zsh]

## Install

    wget -O chruby-0.3.9.tar.gz https://github.com/postmodern/chruby/archive/v0.3.9.tar.gz
    tar -xzvf chruby-0.3.9.tar.gz
    cd chruby-0.3.9/
    sudo make install

### PGP

All releases are [PGP] signed for security. Instructions on how to import my
PGP key can be found on my [blog][1]. To verify that a release was not tampered 
with:

    wget https://raw.github.com/postmodern/chruby/master/pkg/chruby-0.3.9.tar.gz.asc
    gpg --verify chruby-0.3.9.tar.gz.asc chruby-0.3.9.tar.gz

### setup.sh

chruby also includes a `setup.sh` script, which installs chruby and the latest
releases of [Ruby], [JRuby] and [Rubinius]. Simply run the script as root or 
via `sudo`:

    sudo ./scripts/setup.sh

### Homebrew

chruby can also be installed with [homebrew]:

    brew install chruby

Or the absolute latest chruby can be installed from source:

    brew install chruby --HEAD

### Arch Linux

chruby is already included in the [AUR]:

    yaourt -S chruby
    
### FreeBSD

chruby is included in the official [FreeBSD ports collection]:

    cd /usr/ports/devel/chruby/ && make install clean

### Rubies

#### Manually

Chruby provides detailed instructions for installing additional Rubies:

* [Ruby](https://github.com/postmodern/chruby/wiki/Ruby)
* [JRuby](https://github.com/postmodern/chruby/wiki/JRuby)
* [Rubinius](https://github.com/postmodern/chruby/wiki/Rubinius)
* [MagLev](https://github.com/postmodern/chruby/wiki/MagLev)

#### ruby-install

You can also use [ruby-install] to install additional Rubies:

Installing to `/opt/rubies` or `~/.rubies`:

    ruby-install ruby
    ruby-install jruby
    ruby-install rubinius
    ruby-install maglev

#### ruby-build

You can also use [ruby-build] to install additional Rubies:

Installing to `/opt/rubies`:

    ruby-build 1.9.3-p392 /opt/rubies/ruby-1.9.3-p392
    ruby-build jruby-1.7.3 /opt/rubies/jruby-1.7.3
    ruby-build rbx-2.0.0-rc1 /opt/rubies/rubinius-2.0.0-rc1
    ruby-build maglev-1.0.0 /opt/rubies/maglev-1.0.0

## Configuration

Add the following to the `~/.bashrc` or `~/.zshrc` file:

``` bash
source /usr/local/share/chruby/chruby.sh
```

### System Wide

If you wish to enable chruby system-wide, add the following to
`/etc/profile.d/chruby.sh`:

``` bash
if [ -n "$BASH_VERSION" ] || [ -n "$ZSH_VERSION" ]; then
  source /usr/local/share/chruby/chruby.sh
  ...
fi
```

This will prevent chruby from accidentally being loaded by `/bin/sh`, which
is not always the same as `/bin/bash`.

### Rubies

When chruby is first loaded by the shell, it will auto-detect Rubies installed
in `/opt/rubies/` and `~/.rubies/`. After installing new Rubies, you _must_
restart the shell before chruby can recognize them.

For Rubies installed in non-standard locations, simply append their paths to
the `RUBIES` variable:

``` bash
source /usr/local/share/chruby/chruby.sh

RUBIES=(
  /opt/jruby-1.7.0
  "$HOME/src/rubinius"
)
```

### Migrating

If you are migrating from another Ruby manager, set `RUBIES` accordingly:

#### RVM

``` bash
RUBIES+=(~/.rvm/rubies/*)
```

#### rbenv

``` bash
RUBIES+=(~/.rbenv/versions/*)
```

#### rbfu

``` bash
RUBIES+=(~/.rbfu/rubies/*)
```

### Auto-Switching

If you want chruby to auto-switch the current version of Ruby when you `cd`
between your different projects, simply load `auto.sh` in `~/.bashrc` or
`~/.zshrc`:

``` bash
source /usr/local/share/chruby/chruby.sh
source /usr/local/share/chruby/auto.sh
```

chruby will check the current and parent directories for a [.ruby-version]
file. Other Ruby switchers also understand this file:
https://gist.github.com/1912050

If you want to automatically run the version of a gem executable specified in 
your project's Gemfile, try 
[rubygems-bundler](https://github.com/mpapis/rubygems-bundler).

### Default Ruby

If you wish to set a default Ruby, simply call `chruby` in `~/.bash_profile` or
`~/.zprofile`:

    chruby ruby-1.9

If you have enabled auto-switching, simply create a `.ruby-version` file:

    echo "ruby-1.9" > ~/.ruby-version

### RubyGems

Gems installed as a non-root user via `gem install` will be installed into
`~/.gem/$ruby/X.Y.Z`.  By default, RubyGems will use the absolute path to the
currently selected ruby for the shebang of any binstubs it generates.  In some
cases, this path may contain extra version information (e.g.
`ruby-2.0.0-p451`).  To mitigate potential problems when removing rubies, you
can force RubyGems to generate binstubs with shebangs that will search for
ruby in your `$PATH` by using `gem install --env-shebang` (or the equivalent
short option `-E`).  This parameter can also be added to your gemrc file.

### Integration

For instructions on using chruby with other tools, please see the [wiki]:

* [Capistrano](https://github.com/capistrano/chruby#readme)
* [Chef](https://supermarket.getchef.com/cookbooks/chruby_install)
* [Cron](https://github.com/postmodern/chruby/wiki/Cron)
* [Emacs](https://github.com/arnebrasseur/chruby.el#readme)
* [Pow](https://github.com/postmodern/chruby/wiki/Pow)
* [Puppet](https://github.com/dgoodlad/puppet-chruby#readme)
* [Sudo](https://github.com/postmodern/chruby/wiki/Sudo)
* [Vim](https://github.com/postmodern/chruby/wiki/Vim)
* [Fish](https://github.com/JeanMertz/chruby-fish#readme)

## Examples

List available Rubies:

    $ chruby
       ruby-1.9.3-p392
       jruby-1.7.0
       rubinius-2.0.0-rc1

Select a Ruby:

    $ chruby 1.9.3
    $ chruby
     * ruby-1.9.3-p392
       jruby-1.7.0
       rubinius-2.0.0-rc1
    $ echo $PATH
    /home/hal/.gem/ruby/1.9.3/bin:/opt/rubies/ruby-1.9.3-p392/lib/ruby/gems/1.9.1/bin:/opt/rubies/ruby-1.9.3-p392/bin:/usr/lib64/qt-3.3/bin:/usr/local/bin:/usr/bin:/bin:/usr/local/sbin:/usr/sbin:/home/hal/bin:/home/hal/bin
    $ gem env
    RubyGems Environment:
      - RUBYGEMS VERSION: 1.8.23
      - RUBY VERSION: 1.9.3 (2013-02-22 patchlevel 392) [x86_64-linux]
      - INSTALLATION DIRECTORY: /home/hal/.gem/ruby/1.9.3
      - RUBY EXECUTABLE: /opt/rubies/ruby-1.9.3-p392/bin/ruby
      - EXECUTABLE DIRECTORY: /home/hal/.gem/ruby/1.9.3/bin
      - RUBYGEMS PLATFORMS:
        - ruby
        - x86_64-linux
      - GEM PATHS:
         - /home/hal/.gem/ruby/1.9.3
         - /opt/rubies/ruby-1.9.3-p392/lib/ruby/gems/1.9.1
      - GEM CONFIGURATION:
         - :update_sources => true
         - :verbose => true
         - :benchmark => false
         - :backtrace => false
         - :bulk_threshold => 1000
         - "gem" => "--no-rdoc"
      - REMOTE SOURCES:
         - http://rubygems.org/

Switch to JRuby in 1.9 mode:

    $ chruby jruby --1.9
    $ ruby -v
    jruby 1.7.0 (1.9.3p203) 2012-10-22 ff1ebbe on OpenJDK 64-Bit Server VM 1.7.0_09-icedtea-mockbuild_2012_10_17_15_53-b00 [linux-amd64]

Switch back to system Ruby:

    $ chruby system
    $ echo $PATH
    /usr/local/bin:/usr/bin:/bin:/usr/local/sbin:/usr/sbin:/sbin:/home/hal/bin

Run a command under a Ruby with `chruby-exec`:

    $ chruby-exec jruby -- gem update

Switch to an arbitrary Ruby on the fly:

    $ chruby_use /path/to/ruby

## Uninstall

After removing the chruby configuration:

    $ sudo make uninstall

## Alternatives

* [RVM]
* [rbenv]
* [rbfu]*
* [ry]
* [ruby-version]*

\* *Deprecated in favor of chruby.*

## Endorsements

> yeah `chruby` is nice, does the limited thing of switching really good,
> the only hope it never grows 

-- [Michal Papis](https://twitter.com/mpapis/status/258049391791841280) of [RVM]

> I just looooove [chruby](#readme) For the first time I'm in total control of
> all aspects of my Ruby installation. 

-- [Marius Mathiesen](https://twitter.com/zmalltalker/status/271192206268829696)

> Written by Postmodern, it's basically the simplest possible thing that can
> work.

-- [Steve Klabnik](http://blog.steveklabnik.com/posts/2012-12-13-getting-started-with-chruby)

> So far, I'm a huge fan. The tool does what it advertises exactly and simply.
> The small feature-set is also exactly and only the features I need.

-- [Patrick Brisbin](http://pbrisbin.com/posts/chruby)

> I wrote ruby-version; however, chruby is already what ruby-version wanted to
> be. I've deprecated ruby-version in favor of chruby.

-- [Wil Moore III](https://github.com/wilmoore)

## Credits

* [mpapis](https://github.com/mpapis) for reviewing the code.
* [havenwood](https://github.com/havenwood) for handling the homebrew formula.
* [zendeavor](https://github.com/zendeavor) for style fixes.
* `#bash`, `#zsh`, `#machomebrew` for answering all my questions.

[wiki]: https://github.com/postmodern/chruby/wiki

[bash]: http://www.gnu.org/software/bash/
[zsh]: http://www.zsh.org/
[PGP]: http://en.wikipedia.org/wiki/Pretty_Good_Privacy
[homebrew]: http://brew.sh/
[AUR]: https://aur.archlinux.org/packages/chruby/
[FreeBSD ports collection]: https://www.freshports.org/devel/chruby/
[ruby-install]: https://github.com/postmodern/ruby-install#readme
[ruby-build]: https://github.com/sstephenson/ruby-build#readme
[.ruby-version]: https://gist.github.com/1912050

[RVM]: https://rvm.io/
[rbenv]: https://github.com/sstephenson/rbenv#readme
[rbfu]: https://github.com/hmans/rbfu#readme
[ry]: https://github.com/jayferd/ry#readme
[ruby-version]: https://github.com/wilmoore/ruby-version#readme

[Ruby]: http://www.ruby-lang.org/en/
[JRuby]: http://jruby.org/
[Rubinius]: http://rubini.us/

[1]: http://postmodern.github.com/contact.html#pgp
