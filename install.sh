#!/usr/bin/env bash
cp etc/lvodQuery.monit /etc/monit.d
monit reload
sleep 5
monit start lvodQuery

echo "install done!"
