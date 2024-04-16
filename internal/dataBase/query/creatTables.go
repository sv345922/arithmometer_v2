package query

const (
	Q_Tasks = `
CREATE TABLE IF NOT EXISTS tasks(
    nodeId INTEGER PRIMARY KEY NOT NULL,
    x FLOAT,
    xReady BOOLEAN DEFAULT FALSE,
    y FLOAT,
    yReady BOOLEAN DEFAULT FALSE,
    calcID INTEGER DEFAULT 0,
    deadline TIMESTAMP DEFAULT 0,
    duration INTEGER DEFAULT 0,
    FOREIGN KEY (nodeId) REFERENCES allNodes (id)
	);`
	Q_ReadyToCalc = `
CREATE TABLE IF NOT EXISTS q_ReadyToCalc(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    taskId INTEGER NOT NULL,
    FOREIGN KEY (taskId) REFERENCES tasks (id)
    );`
	Q_Working = `
CREATE TABLE IF NOT EXISTS q_Working(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    taskId INTEGER NOT NULL,
    FOREIGN KEY (taskId) REFERENCES tasks (id)
    );`
	Q_NotReady = `
CREATE TABLE IF NOT EXISTS q_NotReady(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    taskId INTEGER NOT NULL,
    FOREIGN KEY (taskId) REFERENCES tasks (id)
    );`
	Expressions = `
CREATE TABLE IF NOT EXISTS expressions(
    id INTEGER PRIMARY KEY,
    userId INTEGER,
    userTask TEXT,
    resultExpr FLOAT,
    status CHAR,
    rootId INTEGER,                                      
    FOREIGN KEY (userId) REFERENCES user (id)
    );`
	AllNodes = `
CREATE TABLE IF NOT EXISTS allNodes(
    id INTEGER PRIMARY KEY,
    expressionId INTEGER NOT NULL,
    op TEXT,
    x INTEGER DEFAULT 0,
    y INTEGER DEFAULT 0,
    val FLOAT,
    sheet BOOLEAN DEFAULT FALSE,
    calculated BOOLEAN DEFAULT FALSE,
    parent INTEGER DEFAULT 0
    -- FOREIGN KEY (user_id) REFERENCES user (id)
    );`
	Timings = `
CREATE TABLE IF NOT EXISTS timings(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    plus INTEGER,
    minus INTEGER,
    mult INTEGER,
    div INTEGER
    );`
	Users = `
CREATE TABLE IF NOT EXISTS users(
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    userName TEXT NOT NULL,
    password TEXT NOT NULL
    -- expressions TODO?
	);`
)
