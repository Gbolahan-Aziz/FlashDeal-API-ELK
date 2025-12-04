### How To Run

make tidy
make run
# In another terminal:
curl -s -X POST localhost:8080/deals -d '{"title":"Gaming Mouse","price":25000,"stock":5,"active":true}' -H "Content-Type: application/json" | jq
curl -s localhost:8080/deals | jq
curl -s -X POST localhost:8080/orders -d '{"deal_id":"<COPY_ID_FROM_ABOVE>","qty":1}' -H "Content-Type: application/json" | jq
# Watch events (SSE) in a browser: http://localhost:8080/events

                          ┌──────────────────────┐
                          │      config.Load()   │
                          │   (ENV vars, PORT)   │
                          └──────────┬───────────┘
                                     │
┌────────────┐   ┌────────────┐   ┌──┴──────────┐   ┌───────────────┐
│ logx.New() │   │ store.New()│   │ sse.NewHub()│   │ events.NewNoop│
│ JSON logger│   │ In-memory  │   │ Broadcaster │   │ (Kafka later) │
└─────┬──────┘   └──────┬─────┘   └────┬────────┘   └──────┬────────┘
      │                 │             │                    │
      └─────────────────┴─────────────┴────────────────────┘
                              │
                              ▼
               ┌────────────────────────────┐
               │ httpapi.Router httpapi.Deps│
               │   all dependencies passed  │
               └──────────────┬─────────────┘
                              │
                              ▼
                     HTTP Server starts
                     (handlers + middleware)


It’s a tiny HTTP API with a real-time SSE stream. Typical consumers:

CLI / cURL for quick tests

Frontend (React/Vue/etc.) that:

hits /deals and /orders

opens /events to show live “order.created” & “deal.created”

Log pipeline (ELK) consuming the JSON logs for dashboards

(Optional) Kafka/RabbitMQ listening to audit events once we add the publisher

curl -v -X POST http://127.0.0.1:8080/deals \
  -H "Content-Type: application/json" \
  -d '{"title":"Gaming Mouse","price":25000,"stock":5,"active":true}'

curl -s :8080/deals | jq

curl -s -X POST :8080/orders \
  -H "Content-Type: application/json" \
  -H "X-Request-ID: test-123" \
  -d '{"deal_id":"<PASTE_ID>", "qty": 1, "user_id":"u-42"}' | jq

curl -N http://localhost:8080/events

for i in {1..100}; do
  curl -s -X POST http://localhost:8080/deals \
  -H "Content-Type: application/json" \
  -d "{\"title\":\"Item $i\", \"price\":$((RANDOM % 50000 + 1000)), \"stock\":$((RANDOM % 10 + 1)), \"active\":true}" \
  > /dev/null
  echo "Created deal $i"
done


for i in {1..200}; do
  curl -s http://localhost:8080/deals > /dev/null
done
echo "GET flood completed"


DEALS=$(curl -s http://localhost:8080/deals | jq -r '.[].id')

for i in {1..300}; do
  DEAL=$(echo "$DEALS" | shuf -n 1)
  curl -s -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d "{\"deal_id\":\"$DEAL\", \"qty\":1}" > /dev/null
  echo "Order $i for $DEAL"
done
