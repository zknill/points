package board

func empty() string {
	return `+------+-------+
| NAME | SCORE |
+------+-------+
+------+-------+
`
}

func aliceZero() string {
	return `+-------+-------+
| NAME  | SCORE |
+-------+-------+
| alice |     0 |
+-------+-------+
`
}

func aliceBobOneWanted() string {
	return `+-------+-------+
| NAME  | SCORE |
+-------+-------+
| alice |     1 |
| bob   |     1 |
+-------+-------+
`
}

func aliceOneBobTwoWanted() string {
	return `+-------+-------+
| NAME  | SCORE |
+-------+-------+
| bob   |     2 |
| alice |     1 |
+-------+-------+
`
}
