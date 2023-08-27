package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	apperr "api/src/app_error"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type SchemaField struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Kind     string `json:"kind"`
	Required bool   `json:"required"`
	Options  any    `json:"options"`
}

type Schema struct {
	Fields []*SchemaField `json:"fields"`
}

type Collection struct {
	CreatedAt time.Time `json:"createdAt,omitempty"`
	Id        int       `json:"id"`
	Kind      string    `json:"kind"`
	Name      string    `json:"name"`
	Schema    Schema    `json:"schema"              gorm:"type:json;serializer:json"`
	System    bool      `json:"system"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type Record struct {
	collection            *Collection
	exportUnknown         bool
	ignoreEmailVisibility bool
	loaded                bool
	originalData          map[string]any
}

const (
	FieldNameId      string = "id"
	FieldNameCreated string = "created_at"
	FieldNameUpdated string = "updated_at"
)

const (
	TextField     string = "text"
	NumberField   string = "number"
	BoolField     string = "bool"
	EmailField    string = "email"
	UrlField      string = "url"
	EditorField   string = "editor"
	DateField     string = "date"
	SelectField   string = "select"
	JsonField     string = "json"
	FileField     string = "file"
	RelationField string = "relation"
)

func (sf *SchemaField) ColDefinition() string {
	switch sf.Kind {
	case NumberField:
		return "NUMERIC DEFAULT 0"
	case BoolField:
		return "BOOLEAN DEFAULT FALSE"
	case JsonField:
		return "JSON DEFAULT NULL"
	default:
		return "TEXT DEFAULT ''"
	}
}

type Testie struct {
	db *gorm.DB
}

// CreateCollection godoc
// @Summary      Create new collection
// @Description  Create a new collection (database table)
// @Tags         Collection
// @Accept       json
// @Produce      json
// @Param        collection    body     Collection  true  "Create collection body"
// @Success      201  {object}  Collection
// @Failure      400  {object}  apperr.AppError
// @Failure      500  {object}  apperr.AppError
// @Router       /test [post]
func (t *Testie) CreateCollection(c echo.Context) error {
	body := new(Collection)
	if err := c.Bind(body); err != nil {
		fmt.Printf("Err parsing body: %#v\n", err)
		return echo.ErrBadRequest
	}

	fmt.Printf("Body: %#v\n", body)
	if err := t.db.Create(body).Error; err != nil {
		return echo.ErrInternalServerError
	}

	found := Collection{Id: body.Id}
	if err := t.db.Take(&found).Error; err != nil {
		return echo.ErrNotFound
	}

	cols := map[string]string{
		FieldNameId:      "TEXT PRIMARY KEY DEFAULT ('r'||lower(hex(randomblob(7)))) NOT NULL",
		FieldNameCreated: "TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL",
		FieldNameUpdated: "TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL",
	}
	for _, field := range found.Schema.Fields {
		cols[field.Name] = field.ColDefinition()
	}

	tableName := found.Name
	dbCols := make([]string, 0, len(cols))
	for fieldName, fieldType := range cols {
		col := fmt.Sprintf("%s %s", fieldName, fieldType)
		dbCols = append(dbCols, col)
	}
	dbColsString := fmt.Sprintf("CREATE TABLE %s (%s)", tableName, strings.Join(dbCols, ", "))
	if err := t.db.Exec(dbColsString).Error; err != nil {
		fmt.Printf("Create table err: %#v", err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, found)
}

// GetCollection godoc
// @Summary      Get collections
// @Description  Get a list of collections (database table)
// @Tags         Collection
// @Accept       json
// @Produce      json
// @Success      200  {object}  []Collection
// @Failure      500  {object}  apperr.AppError
// @Router       /test [get]
func (t *Testie) GetCollection(c echo.Context) error {
	collections := make([]*Collection, 50)
	if err := t.db.Find(&collections).Error; err != nil {
		fmt.Printf("Query err: %#v\n", err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, collections)
}

// GetRecord godoc
// @Summary      Get records of a collection
// @Description  Get a list of record of a collection
// @Tags         Record
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Collection ID"
// @Success      200  {object}  []Collection
// @Failure      500  {object}  apperr.AppError
// @Router       /test-record/:collection [post]
func (t *Testie) CreateRecord(c echo.Context) error {
	collectionName := c.Param("collection")
	if collectionName == "" {
		return echo.ErrBadRequest
	}

	collection := Collection{Name: "TestTable"}
	if err := t.db.Take(&collection).Error; err != nil {
		fmt.Printf("Err getting collection: %#v\n", err)
		return echo.ErrInternalServerError
	}

	tableName := collection.Name
	schema := collection.Schema
	cols := make([]string, 0)
	values := make([]string, 0)
	for _, field := range schema.Fields {
		cols = append(cols, field.Name)
		values = append(values, "test")
	}

	sqlString := fmt.Sprintf(
		"INSERT INTO %v (%v) VALUES (\"%v\")",
		tableName,
		strings.Join(cols, ", "),
		strings.Join(values, ", "),
	)
	if err := t.db.Exec(sqlString).Error; err != nil {
		fmt.Printf("Create record err: %#v\nSQL: %v\n", err, sqlString)
		return echo.ErrInternalServerError
	}

	return apperr.New("123213", 500, "Lmao", "lmao", nil)
}

type Test struct {
	LmaoTest string
}
