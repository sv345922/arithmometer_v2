package dataBase

import (
	"context"
	"database/sql"
	"github.com/sv345922/arithmometer_v2/internal/entities"
	"reflect"
	"testing"
)

func TestCreateEmptyDb(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *sql.DB
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateEmptyDb(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateEmptyDb() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateEmptyDb() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataBase_Save(t *testing.T) {
	type fields struct {
		Tasks       []*entities.Task
		Expressions []*entities.Expression
		AllNodes    []*entities.Node
		Timings     *entities.Timings
		Users       []*entities.User
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &DataBase{
				Tasks:       tt.fields.Tasks,
				Expressions: tt.fields.Expressions,
				AllNodes:    tt.fields.AllNodes,
				Timings:     tt.fields.Timings,
				Users:       tt.fields.Users,
			}
			if err := db.Save(); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewDB(t *testing.T) {
	tests := []struct {
		name string
		want *DataBase
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDB(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDB() = %v, want %v", got, tt.want)
			}
		})
	}
}
