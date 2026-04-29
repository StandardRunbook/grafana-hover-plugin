#!/usr/bin/env bash
# Background streamer: appends mixed normal-traffic log lines to
# /tmp/varied_logs.log at ~25 lines/sec. Runs until killed. Used by
# demo/run.sh to populate the analyzer's trailing baseline.
#
# Five log shapes, roughly balanced — no incident templates.
LOG=/tmp/varied_logs.log
i=${SEED:-30000}
while true; do
  i=$((i+1))
  case $((RANDOM % 100)) in
    0|1|2|3|4|5|6|7|8|9|10|11|12|13|14|15|16|17|18|19)
      v=$((60 + RANDOM % 25)).$((RANDOM % 100))
      echo "cpu_usage: ${v}% - hover-$i" ;;
    20|21|22|23|24|25|26|27|28|29|30|31|32|33|34|35|36|37|38|39)
      v=$((6 + RANDOM % 6)).$((RANDOM % 100))
      echo "memory_usage: ${v}GB - hover-$i" ;;
    40|41|42|43|44|45|46|47|48|49|50|51|52|53|54|55|56|57|58|59|60|61|62|63|64)
      echo "User hover-$i logged in from 10.0.$((RANDOM % 60)).$((RANDOM % 200))" ;;
    65|66|67|68|69|70|71|72|73|74|75|76|77|78|79|80|81|82|83|84|85|86|87|88|89)
      v=$((50 + RANDOM % 250))
      echo "disk_io: ${v}MB/s - hover-$i" ;;
    *)
      echo "ERROR: connection refused to host-$((RANDOM % 50)) while calling auth (hover-$i)" ;;
  esac >> "$LOG"
  sleep 0.04
done
