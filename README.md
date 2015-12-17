# Hubroulette

Pull request roulette! Assign a teammate to review newly opened Github Pull
Requests while announcing the play-by-play in the team slack channel.

### Heroku installation

    $ git clone git@github.com:rjz/hubroulette

    $ heroku create
    $ git push heroku master
    $ heroku ps:scale web=1

### Configuration

Use environment variables to define global settings:

    $ heroku config:set \
      SLACK_CHANNEL='#github' \
      SLACK_TOKEN='xxxxxxxx-xxxxxx-xxxxxxx-xxxxxxxx' \
      GITHUB_ACCESS_TOKEN='foobar' \
      GITHUB_WEBHOOK_SECRET='xyz' \
      HUBROULETTERC='{"team":[{"github":"rjz","slack":"rj"}]}'

Per-repository configuration can be managed using a JSON `.hubrouletterc` file
in the top level of the repository.

```json
{
  "slackChannel": "#general",
  "team": [
    {
      "github": "<github_login>",
      "slack": "<slack_handle>"
    }
  ]
}
```

### License

MIT
