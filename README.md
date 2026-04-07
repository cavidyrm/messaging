# messaging
#kafka sms test
echo '{"event_id":"550e8400-e29b-41d4-a716-446655440000","aggregate_id":"550e8400-e29b-41d4-a716-446655440001","aggregate_type":"SMS","event_type":"SMSSendRequested","version":1,"data":{"phone_number":"09147786264","text":"Test SMS"},"metadata":{},"timestamp":"2026-04-07T10:00:00Z"}' | docker exec -i kafka kafka-console-producer.sh --bootstrap-server localhost:9092 --topic messaging

#kafka email test
echo '{"event_id":"660e8400-e29b-41d4-a716-446655449876","aggregate_id":"660e8400-e29b-41d4-a716-446655442301","aggregate_type":"Email","event_type":"EmailSendRequested","version":1,"data":{"address":"javid.yarmohamadi@gmail.com","subject":"Hello from Kafka","body":"This is a test email sent via the messaging service."},"metadata":{},"timestamp":"2026-04-07T10:00:00Z"}' | docker exec -i kafka kafka-console-producer.sh --bootstrap-server localhost:9092 --topic messaging
