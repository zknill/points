App Engine
==========

The google app engine allows the points slackbot to be easily hosted and deployed.

## Development
Install the google app engine development plugin and then from inside this dir:

`goapp serve --host=0.0.0.0`

## Tests
The tests use the app engine testing package and can be run with all the normal go test flags:

`goapp test`

## Deployment
To deploy the app

`goapp deploy -application {app_name} app.yaml`

## Endpoints
There is a single endpoint `POST /command` which takes x-www-form-urlendoded data.

The `text` field holds the commands to be parsed by the app. e.g. `list` or `add slackbot`