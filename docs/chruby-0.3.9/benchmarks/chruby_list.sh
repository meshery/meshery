root="${0%/*}/.."
n=1000

. "$root/share/chruby/chruby.sh"

RUBIES=(
	/opt/rubies/ruby-1.9.3-p448
	/opt/rubies/ruby-2.0.0-p247
	/opt/rubies/jruby-1.7.4
	/opt/rubies/ruby-1.8.7-p374
	/opt/rubies/rubinius-2.0.0-rc1
)
RUBY_ROOT="/opt/rubies/rubinius-2.0.0-rc1"

for i in {1..3}; do
	echo "Listing rubies $n times ..."

	time (
		for ((i=0; i<$n; i++)); do
			chruby >/dev/null
		done
	)
done
