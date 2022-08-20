#!/usr/bin/env bash
#
# Installs and configures chruby.
#

set -e

#
# Constants
#
export PREFIX="${PREFIX:-/usr/local}"

#
# Functions
#
function log() {
	if [[ -t 1 ]]; then
		printf "%b>>>%b %b%s%b\n" "\x1b[1m\x1b[32m" "\x1b[0m" \
		                          "\x1b[1m\x1b[37m" "$1" "\x1b[0m"
	else
		printf ">>> %s\n" "$1"
	fi
}

function error() {
	if [[ -t 1 ]]; then
		printf "%b!!!%b %b%s%b\n" "\x1b[1m\x1b[31m" "\x1b[0m" \
		                          "\x1b[1m\x1b[37m" "$1" "\x1b[0m" >&2
	else
		printf "!!! %s\n" "$1" >&2
	fi
}

function warning() {
	if [[ -t 1 ]]; then
		printf "%b***%b %b%s%b\n" "\x1b[1m\x1b[33m" "\x1b[0m" \
			                  "\x1b[1m\x1b[37m" "$1" "\x1b[0m" >&2
	else
		printf "*** %s\n" "$1" >&2
	fi
}

#
# Install chruby
#
log "Installing chruby ..."
make install

#
# Configuration
#
log "Configuring chruby ..."

config="if [ -n \"\$BASH_VERSION\" ] || [ -n \"\$ZSH_VERSION\" ]; then
	source $PREFIX/share/chruby/chruby.sh
	source $PREFIX/share/chruby/auto.sh
fi"

if [[ -d /etc/profile.d/ ]]; then
	# Bash/Zsh
	echo "$config" > /etc/profile.d/chruby.sh
	log "Setup complete! Please restart the shell"
else
	warning "Could not determine where to add chruby configuration."
	warning "Please add the following configuration where appropriate:"
	echo
	echo "$config"
	echo
fi
