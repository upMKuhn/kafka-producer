package schema
// Kafka Schema Registry Utility Functions


import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io/ioutil"

	"github.com/riferrei/srclient"
	"google.golang.org/protobuf/proto"
)

var schemaMap map[string]srclient.Schema
var schemaKeyMap map[string]srclient.Schema
var schemaRegistry srclient.ISchemaRegistryClient

func init() {
	schemaMap = make(map[string]srclient.Schema)
	schemaKeyMap = make(map[string]srclient.Schema)
}

// ConnectSchemaRegistry creates registry obj
func ConnectSchemaRegistry(address string) {
	schemaRegistry = srclient.CreateSchemaRegistryClient(address)

}

// LoadSchemas for later encoding
func LoadSchemas(topics []string) {

	for _, topic := range topics {
		schema, _ := schemaRegistry.GetLatestSchema(topic, true)
		if schema != nil {
			schemaKeyMap[topic] = *schema
		}
		schema, _ = schemaRegistry.GetLatestSchema(topic, false)
		if schema != nil {
			schemaMap[topic] = *schema
		}
	}
}

// MarschalProtobuf Serializes Protobuf message to be schema registry compliant
func MarschalProtobuf(topic string, isKey bool, message proto.Message) ([]byte, error) {
	messageData, err := proto.Marshal(message)

	if err != nil {
		return nil, err
	}

	schema := getSchema(topic, isKey)
	schemaIDBytes := make([]byte, 4)
	messageIndexBytes := []byte{byte(2), byte(0)} // I guess it means (arra.len, Index of Message in Protobuf Schema)
	binary.BigEndian.PutUint32(schemaIDBytes, uint32(schema.ID()))
	// https://github.com/confluentinc/demo-scene/blob/master/getting-started-with-ccloud-golang/ClientApp.go
	// https://riferrei.com/2020/07/09/data-sharing-between-java-go-using-kafka-and-protobuf/
	var data []byte
	data = append(data, byte(0))
	data = append(data, schemaIDBytes...)
	data = append(data, messageIndexBytes...)
	data = append(data, messageData...)
	return data, nil
}

// RegisterSchemaFromFile create or update entry from file
func RegisterSchemaFromFile(topic string, isKey bool, filePath string) {
	schemaBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(fmt.Sprintf("Error loading the schema %s", err))
	}
	RegisterSchema(topic, isKey, string(schemaBytes))
}

// RegisterSchema create or update entry from source String
func RegisterSchema(topic string, isKey bool, newSchemaSrc string) {
	oldSchema, _ := schemaRegistry.GetLatestSchema(topic, isKey)

	if oldSchema == nil || areSchemaEqual(oldSchema.Schema(), newSchemaSrc) {
		schema, err := schemaRegistry.CreateSchema(topic, newSchemaSrc, srclient.Protobuf, isKey)
		if err != nil {
			panic(fmt.Sprintf("Error creating the schema %s", err))
		}
		fmt.Printf("Created schema %s", schema.Schema())
	}
}

func areSchemaEqual(schemaA string, schemaB string) bool {
	encodedA := base64.StdEncoding.EncodeToString([]byte(schemaA))
	encodedB := base64.StdEncoding.EncodeToString([]byte(schemaB))
	return encodedA == encodedB
}

func getSchema(topic string, isKey bool) srclient.Schema {
	if isKey {
		return schemaKeyMap[topic]
	}
	return schemaMap[topic]
}
