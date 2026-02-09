HOST="localhost:8080"
TOTAL_ORDERS=1000
SYMBOLS=("RELIANCE" "TCS" "INFY" "HDFC" "ICICI")
SIDES=("BUY" "SELL")

echo "🚀 Starting load test - sending $TOTAL_ORDERS orders"
echo "Target: http://$HOST/order"
echo ""

send_order() {
    local symbol=${SYMBOLS[$RANDOM % ${#SYMBOLS[@]}]}
    local side=${SIDES[$RANDOM % ${#SIDES[@]}]}
    local price=$((2000 + RANDOM % 1000))
    local qty=$((1 + RANDOM % 50))
    
    curl -s -X POST "http://$HOST/order" \
        -H "Content-Type: application/json" \
        -d "{
            \"symbol\": \"$symbol\",
            \"side\": \"$side\",
            \"price\": $price,
            \"qty\": $qty,
            \"user_id\": \"load-test\"
        }" > /dev/null
}

start_time=$(date +%s)
for i in $(seq 1 $TOTAL_ORDERS); do
    send_order &
    
    if [ $((i % 100)) -eq 0 ]; then
        echo "Sent $i orders..."
    fi
done

wait

end_time=$(date +%s)
duration=$((end_time - start_time))
throughput=$((TOTAL_ORDERS / duration))

echo ""
echo "Load test complete!"
echo "Total orders: $TOTAL_ORDERS"
echo "Duration: ${duration}s"
echo "Throughput: ~${throughput} orders/sec"
echo ""
echo "Check system stats:"
echo "curl http://$HOST/stats | jq"