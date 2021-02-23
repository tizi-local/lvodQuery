#!/usr/bin/env bash

set -e

NAME=lvodQuery
PIDFILE=/var/run/$NAME.pid
SRV_BIN=/root/local/lvodQuery/lvodQuery
SRV_ARGS="-config /root/local/lvodQuery/etc/config.json"
COMMAND="${SRV_BIN} ${SRV_ARGS}"

start () {
    start-stop-daemon --start --quiet --background --make-pidfile --pidfile ${PIDFILE}  --exec ${SRV_BIN} -- ${SRV_ARGS}
  echo "[start lvodQuery succeed]"
}

stop () {
  start-stop-daemon --stop --quiet --pidfile $PIDFILE
  if [ -e $PIDFILE ]
    then rm $PIDFILE
  fi
  echo "[stop lvodQuery succeed]"
}

name=`basename $0`
case $1 in
 start)
        echo "start..."
        start
        ;;
 stop)
        echo "stop ..."
        stop
        ;;
 *)
        echo "Usage: $name [start|stop|reload]"
        exit 1
        ;;
esac