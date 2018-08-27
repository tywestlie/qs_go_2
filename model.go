package main

import (
    "database/sql"
    "errors"
)

type food struct {
    ID    int     `json:"id"`
    Name  string  `json:"name"`
    Calories int `json:"calories"`
}

func (f *food) getFood(db *sql.DB) error {
  return db.QueryRow("SELECT name, calories FROM foods WHERE id=$1",
       f.ID).Scan(&f.Name, &f.Calories)
}

func (f *food) updateFood(db *sql.DB) error {
  return errors.New("Not implemented")
}

func (f *food) deleteFood(db *sql.DB) error {
  return errors.New("Not implemented")
}

func (f *food) createFood(db *sql.DB) error {
  return errors.New("Not implemented")
}

func getFoods(db *sql.DB, start, count int) ([]food, error) {
  return nil, errors.New("Not implemented")
}
