package query

var InsertExpression = `INSERT INTO expressions (
                         userId, 
                         userTask, 
                         resultExpr, 
                         status, 
                         rootId) 
VALUES ($1, $2, $3, $4, $5);`

var InsertNode = `INSERT INTO allNodes (
                         ExpressionId, 
                         Op, 
                         X, 
                         Y, 
                         Val, 
                         Sheet,
                     	 Calculated,
                      	 Parent) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`

var InsertTask = `INSERT INTO tasks (
                         nodeId, 
                         x, 
                         xReady, 
                         y, 
                         yReady, 
                         calcID,
                      	 deadline,
                   		 duration) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`
