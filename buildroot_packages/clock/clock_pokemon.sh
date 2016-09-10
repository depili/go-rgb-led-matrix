#!/bin/sh

start() {
	echo -"Starting Clock"
	cd /root/rpi-matrix
	while true
	do
		echo -e "\033[9;0]"
		/root/rpi-matrix/clock_cmd.sh
		echo "CRASHED!"
		sleep 2
	done

}

stop() {
	true
}

restart() {
	stop
	start
}

case "$1" in
	start)
		start
		;;
	stop)
		stop
	;;
	restart|reload)
		restart
		;;
	*)
		echo "Usage: $0 {start|stop|restart}"
		exit 1
esac

exit $?
