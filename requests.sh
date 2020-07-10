
for j in {01..30}; do
  for i in {1..5}; do
    (echo "[$j,$i] $(date) $(curl -sS http://35.190.83.212/hi)" >> requests.out)&
  done
  sleep 1
done
