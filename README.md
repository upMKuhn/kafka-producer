
# Kafka Stocks Producer
A investigation project sending dummy data to kafka.  
### The Producer Workflow:  
  * Load schema
  * Push updated protobuf definitions to Confluent Schema
  * Create Dummy Data and push to Kafka Topic

### Learnings
1. Protobuf is a second class citizen in the Confluent Eco system. Better stick to Avro.
2. Data serialized for Kafka Connect for schema registry are incompatible with a regular serializer
   * The starting byte of messages specifies the schema version. [(The magic byte)](https://docs.confluent.io/home/overview.html#wire-format)       
![kafka connect](assets/schema-registry-and-kafka.png)


