Points
======

Points is a small command line app written in go that lets you easily keep a score or leaderboard.

##Overview

`points` allows you to add names and scores and keeps a tally of who's in the lead.

`points` outputs the leaderboard as an ascii table that can easily be shared in your team's chat.

```
+-------+--------+
| NAME  | POINTS |
+-------+--------+
| Alice |     10 |
| Bob   |      5 |
+-------+--------+
```

The table is ordered by who has the most points in the points column, and then alphabetically if rows have the same number of points.

##Commands

####`init`

`init` allows your to create a leaderboard. If no column names are specified then the defaults of `name` and `points` are used.

Calling `init` with the names of the columns for the leaderboard, the `init` command creates a file `points.json` in the current directory.

The recommended usage is `init name points` as it is assumed that the second column will be the number of points. However other meta data can be passed to init.

`init name points meta` could be used to set up a table that looked like this:
```
+-------+--------+------+
| NAME  | POINTS | META |
+-------+--------+------+
| Alice |     10 | yes  |
| Bob   |      5 | no   |
+-------+--------+------+
```

####`add`

`add` allows you to add points to a row or create a new row in the table.

For the table:
```
+-------+--------+
| NAME  | POINTS |
+-------+--------+
| Alice |     10 |
+-------+--------+
```
`add bob` would create an entry for bob with zero points
`add bob 5` would create an entry for bob with 5 points

If you set up your leaderboard to have meta data after the points column, pass that after the points.

`add points bob 5 yes`

```
+-------+--------+------+
| NAME  | POINTS | META |
+-------+--------+------+
| Bob   |      5 | yes  |
+-------+--------+------+
```

####`list`

`list` prints to stdout the ascii table.

####`slack`

`slack` prints to stdout the ascii table, but surrounded by back ticks to make it easy to paste into a slack channel or message.

####`reset`

`reset` clears all the scores from the leaderboard but keeps the names in the table.

a single entry in the table can be reset using the `--entry` flag to the reset command.

`reset --entry bob`

####`history`

`history` shows the times and commands run.

the table:
```
+-------+--------+
| NAME  | POINTS |
+-------+--------+
| Alice |      9 |
| Bob   |      5 |
+-------+--------+
```
may have the history:
```
2001-02-03 16:50:01: init name points
2001-02-03 16:50:02: add alice 2
2001-02-03 16:50:03: add bob 5
2001-02-03 16:50:04: add alice 7
```