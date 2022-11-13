# A mailer app in go
Sends mails based on request or by receiving a message via Kafka

## Setup for email
1. Ensure .env exists
> make env

2. Run mailhog so you can check the results
> make mailhog

3. Run the app
> make run

4. Send a mail with a post request. You can use this format: 
```
curl -X POST -H "Content-Type: application/json" \
    -d '{"from": "", "fromName":"","to":"","subject":"", "messageBody":{"errorMessage":"","url":""}}' \
    http://$API_HOST:$API_PORT/send
```

or run:
> make send

## Setup for Kafka
1. Run the kafka & zookeeper containers
> make kafka
2. Run the app
> make run
3. SH into the kafka container
> make produce_mode

## Cleanup
Run clean
> make clean

Sometimes you need to remove the network manually
>  docker network rm kafka-network

Also make sure the volumes or mounted directories are removed like `kafka_data`
> rm -rf kafka_data

## Resources
- [Kafka from the CLI](https://medium.com/@TimvanBaarsen/apache-kafka-cli-commands-cheat-sheet-a6f06eac01b#e260)