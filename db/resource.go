package db

func (d *Database) GetResource(id string) (*Resource, error) {
	var res Resource

	err := d.sqlDB.Get(
		&res,
		`
		SELECT
			*
		FROM
			resource
		WHERE id = $1`,
		id,
	)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
