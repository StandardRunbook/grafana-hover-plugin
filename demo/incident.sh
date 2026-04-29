#!/usr/bin/env bash
# Inject a 30s incident into /tmp/varied_logs.log. Distribution shifts:
#   • cpu_usage drops to 5–15% (vs 60–85% in baseline)
#   • ERROR rate ~5x baseline
#   • A NEW template "CRITICAL: db pool exhausted (waiters=N) on shard-S"
#     appears that did not exist in the baseline at all
#
# Records the absolute incident window to /tmp/incident_window.txt so
# the dashboard URL can be aimed at the exact slice.
LOG=/tmp/varied_logs.log
i=${SEED:-50000}

incident_start_iso=$(date -u +%Y-%m-%dT%H:%M:%SZ)
echo "[incident] start: $incident_start_iso"

end=$(($(date +%s) + 30))
while [ $(date +%s) -lt $end ]; do
  i=$((i+1))
  case $((RANDOM % 100)) in
    0|1)
      v=$((5 + RANDOM % 10)).$((RANDOM % 100))
      echo "cpu_usage: ${v}% - hover-$i" ;;
    2|3)
      v=$((6 + RANDOM % 6)).$((RANDOM % 100))
      echo "memory_usage: ${v}GB - hover-$i" ;;
    4|5|6|7|8|9)
      echo "User hover-$i logged in from 10.0.$((RANDOM % 60)).$((RANDOM % 200))" ;;
    10|11|12|13|14|15|16|17|18|19|20|21|22|23|24|25|26|27|28|29|30|31|32|33|34|35|36|37|38|39|40|41|42|43|44|45|46|47|48|49)
      echo "ERROR: connection refused to host-$((RANDOM % 50)) while calling auth (hover-$i)" ;;
    *)
      echo "CRITICAL: db pool exhausted (waiters=$((RANDOM % 200))) on shard-$((RANDOM % 8))" ;;
  esac >> "$LOG"
  sleep 0.05
done

incident_end_iso=$(date -u +%Y-%m-%dT%H:%M:%SZ)
echo "[incident] end: $incident_end_iso"
echo "${incident_start_iso} ${incident_end_iso}" > /tmp/incident_window.txt
