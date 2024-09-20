package person

type Person struct {
	ID      *int    `db:"id"`
	Name    *string `db:"name"`
	Age     *int    `db:"age"`
	Address *string `db:"address"`
	Work    *string `db:"work"`
}
