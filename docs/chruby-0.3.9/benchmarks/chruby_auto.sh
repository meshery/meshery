root="${0%/*}/.."
n=100

RUBIES=("$root/test/opt/rubies/ruby-2.0.0-p353")

. "$root/share/chruby/chruby.sh"
. "$root/share/chruby/auto.sh"

for i in {1..3}; do
	echo "Auto-switching $n times..."

	time (
		for ((i=0; i<$n; i+=1)); do
			cd "$root/test/project"
			chruby_auto
			cd ../../
			chruby_auto
		done
	)
done
