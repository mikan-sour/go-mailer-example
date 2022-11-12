# A mailer app in go

## Setup
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


