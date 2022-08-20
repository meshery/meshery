. ./test/helper.sh

function test_chruby_exec_no_arguments()
{
	chruby-exec 2>/dev/null

	assertEquals "did not exit with 1" 1 $?
}

function test_chruby_exec_no_command()
{
	chruby-exec "$test_ruby_version" 2>/dev/null

	assertEquals "did not exit with 1" 1 $?
}

function test_chruby_exec()
{
	local ruby_version=$(chruby-exec "$test_ruby_version" -- ruby -e "print RUBY_VERSION")

	assertEquals "did change the ruby" "$test_ruby_version" "$ruby_version"
}

SHUNIT_PARENT=$0 . $SHUNIT2
