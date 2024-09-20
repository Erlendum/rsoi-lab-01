package person

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"
)

const (
	defaultTimeout = 5 * time.Second
)

type repository struct {
	conn *sqlx.DB
}

func NewRepository(conn *sqlx.DB) *repository {
	return &repository{conn: conn}
}

func (r *repository) CreatePerson(ctx context.Context, person Person) (int, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Insert("persons").Columns("name", "age", "address", "work").Values(person.Name, person.Age, person.Address, person.Work)
	query, args, err := builder.Suffix("RETURNING id").ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "failed to build query")
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var id int
	err = r.conn.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "failed to execute query")
	}

	return id, nil
}

func (r *repository) createUpdateBuilderForPerson(id int, person Person) (sq.UpdateBuilder, bool) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	isEmpty := true
	updateBuilder := psql.Update("persons").Where(sq.Eq{"id": id})
	if person.Name != nil {
		updateBuilder = updateBuilder.Set("name", person.Name)
		isEmpty = false
	}
	if person.Age != nil {
		updateBuilder = updateBuilder.Set("age", person.Age)
		isEmpty = false
	}
	if person.Work != nil {
		updateBuilder = updateBuilder.Set("work", person.Work)
		isEmpty = false
	}
	if person.Address != nil {
		updateBuilder = updateBuilder.Set("address", person.Address)
		isEmpty = false
	}

	return updateBuilder, isEmpty
}

func (r *repository) UpdatePerson(ctx context.Context, id int, person *Person) error {
	builder, isEmpty := r.createUpdateBuilderForPerson(id, *person)
	if isEmpty {
		return nil
	}

	builder = builder.Suffix("RETURNING id, name, age, address, work")

	query, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build query")
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	err = r.conn.QueryRowContext(ctx, query, args...).Scan(&person.ID, &person.Name, &person.Age, &person.Address, &person.Work)
	if err != nil {
		return errors.Wrap(err, "failed to execute query")
	}

	return nil
}

func (r *repository) DeletePerson(ctx context.Context, id int) (bool, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Delete("persons").Where(sq.Eq{"id": id})
	query, args, err := builder.ToSql()
	if err != nil {
		return false, errors.Wrap(err, "failed to build query")
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	res, err := r.conn.ExecContext(ctx, query, args...)
	if res == nil || err != nil {
		return false, errors.Wrap(err, "failed to execute query")
	}

	countAffectedRows, err := res.RowsAffected()
	if err != nil {
		return false, errors.Wrap(err, "failed to get count of affected rows")
	}

	return countAffectedRows == 1, nil
}

func (r *repository) GetPersons(ctx context.Context) ([]Person, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	builder := psql.Select("id", "name", "age", "address", "work").From("persons")

	query, args, err := builder.ToSql()
	if err != nil {
		return []Person{}, errors.Wrap(err, "failed to build query")
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	res := make([]Person, 0)

	err = r.conn.SelectContext(ctx, &res, query, args...)
	if err != nil {
		return []Person{}, errors.Wrap(err, "failed to execute query")
	}

	return res, nil
}

func (r *repository) GetPerson(ctx context.Context, id int) (Person, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	builder := psql.Select("id", "name", "age", "address", "work").From("persons").Where(sq.Eq{"id": id})

	query, args, err := builder.ToSql()
	if err != nil {
		return Person{}, errors.Wrap(err, "failed to build query")
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	res := Person{}

	err = r.conn.GetContext(ctx, &res, query, args...)
	if err != nil {
		return Person{}, errors.Wrap(err, "failed to execute query")
	}

	return res, nil
}
