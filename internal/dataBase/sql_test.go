package dataBase

import (
	"context"
	"fmt"
	"github.com/sv345922/arithmometer_v2/internal/entities"
	"log"
	"os"
	"testing"
)

const dbfile = `./db/test_db.db`

func TestCreateDb(t *testing.T) {
	ctx := context.Background()
	dataBase, err := CreateEmptyDb(ctx)
	if err != nil {
		log.Printf("%v", err)
	}
	defer dataBase.Close()
	_ = dataBase
}
func TestInsertExpression(t *testing.T) {
	expr := entities.Expression{
		Id:         0,
		UserId:     123,
		UserTask:   "2 + 2",
		ResultExpr: 0,
		Status:     "status",
		RootId:     222,
	}
	ctx := context.Background()

	dataBase, err := CreateDb(ctx, "./db/test_db.db")
	if err != nil {
		log.Printf("%v", err)
	}
	defer dataBase.Close()
	id, err := InsertExpression(ctx, dataBase, &expr)
	if err != nil {
		t.Errorf("%v", err)
	}
	expr.Id = id
	log.Println(expr)

}
func TestIsFileExist(t *testing.T) {
	fileName := "/test.file"
	wd, _ := os.Getwd()
	file := wd + fileName
	log.Printf("временный файл %s", file)
	f, err := os.Create(file)
	if err != nil {
		log.Printf("%v\n", err)
	}
	f.Close()
	res, err := IsFileExist(file)
	if err != nil && res == false {
		log.Printf("%v", err)
	}
	os.Remove(file)
}
func TestInsertNode(t *testing.T) {
	node := entities.Node{
		ExpressionId: 1,
		Op:           "+",
		X:            11,
		Y:            22,
		Val:          3,
		Sheet:        true,
		Calculated:   false,
		Parent:       33,
	}
	ctx := context.Background()
	dataBase, err := CreateDb(ctx, dbfile)
	if err != nil {
		log.Printf("%v", err)
	}
	id, err := InsertNode(ctx, dataBase, &node)
	if err != nil {
		t.Errorf("%v", err)
	}
	node.Id = id
	nodes2, err := GetNodes(ctx, dataBase)
	if *nodes2[0] != node {
		t.Errorf("не совпадают вставленный %v в БД и полученный из неё объект %v", node, *nodes2[0])
	}
	dataBase.Close()
	os.Remove(dbfile)
}
func TestGetTimings(t *testing.T) {
	ctx := context.Background()
	dataBase, err := CreateDb(ctx, dbfile)
	if err != nil {
		log.Printf("%v", err)
	}
	timings, err := GetTimings(ctx, dataBase)
	if err != nil {
		log.Printf("%v", err)
	}
	fmt.Println(timings)
	dataBase.Close()
	os.Remove(dbfile)
}
