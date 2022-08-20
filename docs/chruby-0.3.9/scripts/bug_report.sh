#
# chruby script to collect environment information for bug reports.
#

[[ -z "$PS1" ]] && exec "$SHELL" -i -l "$0"

function print_section()
{
	echo
	echo "## $1"
	echo
}

function indent()
{
	echo "$1" | sed 's/^/    /'
}

function print_variable()
{
	if [[ -n "$2" ]]; then echo "    $1=$2"
	else                   echo "    $1=$(eval "echo \$$1")"
	fi
}

function print_version()
{
	local full_path="$(command -v "$1")"

	if [[ -n "$full_path" ]]; then
		local version="$(("$1" --version || "$1" -V) 2>/dev/null)"

		indent "$(echo "$version" | head -n 1) [$full_path]"
	fi
}


print_section "System"

indent "$(uname -a)"
print_version "bash"
print_version "tmux"
print_version "zsh"
print_version "ruby"
print_version "bundle"
print_version "chruby-exec"

print_section "Environment"

print_variable "CHRUBY_VERSION"
print_variable "SHELL"
print_variable "PATH"
print_variable "HOME"

print_variable "RUBIES" "(${RUBIES[*]})"
print_variable "RUBY_ROOT"
print_variable "RUBY_VERSION"
print_variable "RUBY_ENGINE"
print_variable "RUBY_AUTO_VERSION"
print_variable "RUBYLIB"
print_variable "RUBYOPT"
print_variable "RUBYPATH"
print_variable "RUBYSHELL"
print_variable "GEM_ROOT"
print_variable "GEM_HOME"
print_variable "GEM_PATH"

if [[ -n "$ZSH_VERSION" ]]; then
	print_section "Hooks"
	print_variable "preexec_functions" "(${preexec_functions[*]})"
	print_variable "precmd_functions" "(${precmd_functions[*]})"
elif [[ -n "$BASH_VERSION" ]]; then
	print_section "Hooks"
	indent "$(trap -p)"
fi

if [[ -f .ruby-version ]]; then
	print_section ".ruby-version"
	echo "    $(< .ruby-version)"
fi

print_section "Aliases"

indent "$(alias)"
