# Simple Kafka Consumer in Go

This is a template repository for a simple Kafka consumer in Go. It uses the [confluent pkg](https://github.com/confluentinc/confluent-kafka-go).

## Instructions

1. Clone this repository into your `$GOPATH/src` directory.

2. Make sure your imports are correct and run `go mod tidy` to update the `go.sum` file.

3. Start up a local kafka instance using the below command. This docker compose file will automatically create a topic called `test-topic-1`.Kafka-REST is also included to make it easy to interact with the kafka instance vua HTTP.

    ```bash
    docker-compose up -d 
    ```

    You can check that the topic was created by running the following command:

    ```bash
    curl --request GET --url http://localhost:38082/topics
    ```

4. Run `go run main.go` to start the consumer.

    *Note: If you are using a different kafka instance, you will need to update the properties in the `config/config.go` file.*

5. Send a message to the topic using the following command:

    ```bash
    curl --request POST \
    --url http://localhost:38082/topics/test-topic-1 \
    --header 'content-type: application/vnd.kafka.json.v2+json' \
    --data '{
        "records": [
            {
                "value": {
                    "taskType":"taskType1",
                    "data":{
                        "id":1,
                        "name": "John Doe"
                    }
                }
            }
        ]
    }'
    ```

    You should see the message printed to the console.