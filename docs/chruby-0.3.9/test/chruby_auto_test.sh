. ./share/chruby/auto.sh
. ./test/helper.sh

function setUp()
{
	chruby_reset
	unset RUBY_AUTO_VERSION
}

function test_chruby_auto_loaded_in_zsh()
{
	[[ -n "$ZSH_VERSION" ]] || return

	assertEquals "did not add chruby_auto to preexec_functions" \
		     "chruby_auto" \
		     "$preexec_functions"
}

function test_chruby_auto_loaded_in_bash()
{
	[[ -n "$BASH_VERSION" ]] || return

	local command=". $PWD/share/chruby/auto.sh && trap -p DEBUG"
	local output="$("$SHELL" -c "$command")"

	assertTrue "did not add a trap hook for chruby_auto" \
		   '[[ "$output" == *chruby_auto* ]]'
}

function test_chruby_auto_loaded_twice_in_zsh()
{
	[[ -n "$ZSH_VERSION" ]] || return

	. ./share/chruby/auto.sh

	assertNotEquals "should not add chruby_auto twice" \
		        "$preexec_functions" \
			"chruby_auto chruby_auto"
}

function test_chruby_auto_loaded_twice()
{
	RUBY_AUTO_VERSION="dirty"
	PROMPT_COMMAND="chruby_auto"

	. ./share/chruby/auto.sh

	assertNull "RUBY_AUTO_VERSION was not unset" "$RUBY_AUTO_VERSION"
}

function test_chruby_auto_enter_project_dir()
{
	cd "$test_project_dir" && chruby_auto

	assertEquals "did not switch Ruby when entering a versioned directory" \
		     "$test_ruby_root" "$RUBY_ROOT"
}

function test_chruby_auto_enter_subdir_directly()
{
	cd "$test_project_dir/sub_dir" && chruby_auto

	assertEquals "did not switch Ruby when directly entering a sub-directory of a versioned directory" \
		     "$test_ruby_root" "$RUBY_ROOT"
}

function test_chruby_auto_enter_subdir()
{
	cd "$test_project_dir" && chruby_auto
	cd sub_dir             && chruby_auto

	assertEquals "did not keep the current Ruby when entering a sub-dir" \
		     "$test_ruby_root" "$RUBY_ROOT"
}

function test_chruby_auto_enter_subdir_with_ruby_version()
{
	cd "$test_project_dir"    && chruby_auto
	cd sub_versioned/         && chruby_auto

	assertNull "did not switch the Ruby when leaving a sub-versioned directory" \
		   "$RUBY_ROOT"
}

function test_chruby_auto_modified_ruby_version()
{
	cd "$test_project_dir/modified_version" && chruby_auto
	echo "2.0.0" > .ruby-version            && chruby_auto

	assertEquals "did not detect the modified .ruby-version file" \
		     "$test_ruby_root" "$RUBY_ROOT"
}

function test_chruby_auto_overriding_ruby_version()
{
	cd "$test_project_dir" && chruby_auto
	chruby system          && chruby_auto

	assertNull "did not override the Ruby set in .ruby-version" "$RUBY_ROOT"
}

function test_chruby_auto_leave_project_dir()
{
	cd "$test_project_dir"    && chruby_auto
	cd "$test_project_dir/.." && chruby_auto

	assertNull "did not reset the Ruby when leaving a versioned directory" \
		   "$RUBY_ROOT"
}

function test_chruby_auto_invalid_ruby_version()
{
	local expected_auto_version="$(cat $test_project_dir/bad/.ruby-version)"

	cd "$test_project_dir" && chruby_auto
	cd bad/                && chruby_auto 2>/dev/null

	assertEquals "did not keep the current Ruby when loading an unknown version" \
		     "$test_ruby_root" "$RUBY_ROOT"
	assertEquals "did not set RUBY_AUTO_VERSION" \
		     "$expected_auto_version" "$RUBY_AUTO_VERSION"
}

function tearDown()
{
	cd "$PWD"
}

SHUNIT_PARENT=$0 . $SHUNIT2
