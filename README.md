Points
======

Points is a small slack integration written in go that lets you easily keep a score or leaderboard.

##Overview

`points` allows you to add names and scores and keeps a tally of who's in the lead.

`points` outputs the leaderboard as an ascii table that's wrapped in a slack preformatted code block.

```
+-------+--------+
| NAME  | POINTS |
+-------+--------+
| Alice |     10 |
| Bob   |      5 |
+-------+--------+
```

The table is ordered by who has the most points in the points column, and then alphabetically if rows have the same number of points.

##Install and deploy

`points` heavily relies on google app engine and inside the `appengine` dir is a default app.yaml for deploying it. 

Slack includes a security token that can be added to an additional `environment.yaml` to set an env var `SLACK_TOKEN` to allow request to be verified. 

##Commands

Once hosted in app engine and once the slash command is set up in slack you can use the following commands. Assuming you set up the `/points` slash command.

####`init`

`init` allows your to create a leaderboard. Run `/points init` to start!

####`add`

`/points add [name]` allows you to add points to a row or create a new row in the leaderboard.

For the board:
```
+-------+--------+
| NAME  | POINTS |
+-------+--------+
| Alice |     10 |
+-------+--------+
```
`add bob` would create an entry for bob with 1 point. Only 1 point can be added to a user at once.

####`list`

`/ points list` prints to stdout the ascii table.

####`reset`

`/points reset` clears all the scores from the leaderboard but keeps the names in the table.
