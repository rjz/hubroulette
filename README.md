# Github Random Assignee

Pull request roulette! Assign a teammate to review newly opened Github Pull
Requests and let a slack channel know who got lucky.

### Heroku installation

    $ git clone git@github.com:rjz/github-random-assignee

    $ heroku create
    $ git push heroku master
    $ heroku ps:scale web=1

### Configuration

    $ heroku config:set \
      SLACK_CHANNEL='#github' \
      SLACK_TOKEN='xxxxxxxx-xxxxxx-xxxxxxx-xxxxxxxx' \
      GITHUB_ACCESS_TOKEN='foobar' \
      GITHUB_WEBHOOK_SECRET='xyz' \
      ASSIGNEE_LOGINS='rjz,hubot,etc'

### License

MIT
