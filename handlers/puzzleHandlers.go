package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	agentUtils "github.com/lambda-platform/agent/utils"
	"github.com/lambda-platform/dataform"
	"github.com/lambda-platform/datagrid"
	"github.com/lambda-platform/lambda/DB"

	"github.com/lambda-platform/lambda/models"
	"github.com/lambda-platform/lambda/DBSchema"
	"github.com/lambda-platform/lambda/config"
	"github.com/labstack/echo/v4"
	"os"
	"github.com/lambda-platform/datasource"
	"github.com/lambda-platform/lambda/utils"
	"net/http"
	"regexp"
	"strconv"
)

type vb_schema struct {
	ID         int        `gorm:"column:id;primary_key" json:"id"`
	Name   string `json:"name"`
	Schema string `json:"schema"`
}

func Index(c echo.Context) error {
	dbSchema := models.DBSCHEMA{}

	if(config.LambdaConfig.SchemaLoadMode == "auto"){
		dbSchema = DBSchema.GetDBSchema()
	} else {
		schemaFile, err := os.Open("app/models/db_schema.json")
		defer schemaFile.Close()
		if err != nil{
			fmt.Println("schema FILE NOT FOUND")
		}
		dbSchema = models.DBSCHEMA{}
		jsonParser := json.NewDecoder(schemaFile)
		jsonParser.Decode(&dbSchema)
	}


	gridList := []models.VBSchemaList{}
	userRoles := []models.UserRoles{}

	DB.DB.Where("type = ?", "grid").Find(&gridList)
	DB.DB.Find(&userRoles)

	//gridList, err := models.VBSchemas(qm.Where("type = ?", "grid")).All(context.Background(), DB)
	//dieIF(err)

	User := agentUtils.AuthUserObject(c)

	//csrfToken := c.Get(middleware.DefaultCSRFConfig.ContextKey).(string)
	csrfToken := ""
	return c.Render(http.StatusOK, "puzzle.html", map[string]interface{}{
		"lambda_config": config.LambdaConfig,
		"title":                     config.LambdaConfig.Title,
		"favicon":                     config.LambdaConfig.Favicon,
		"app_logo":                     config.LambdaConfig.Logo,
		"app_text":                     "СИСТЕМИЙН УДИРДЛАГА",
		"dbSchema":                  dbSchema,
		"gridList":                  gridList,
		"User":                      User,
		"user_fields":               config.LambdaConfig.UserDataFields,
		"user_roles":               userRoles,
		"data_form_custom_elements": config.LambdaConfig.DataFormCustomElements,
		"mix":                       utils.Mix,
		"csrfToken":                       csrfToken,
	})

}


func GetTableSchema(c echo.Context) error {
	table := c.Param("table")
	tableMetas := DBSchema.TableMetas(table)
	return c.JSON(http.StatusOK, tableMetas)

}

func GetVB(c echo.Context) error {

	type_ := c.Param("type")
	id := c.Param("id")
	condition := c.Param("condition")

	if id != "" {

		match, _ := regexp.MatchString("_", id)

		if(match){
			VBSchema := models.VBSchemaAdmin{}

			DB.DB.Where("id = ?", id).First(&VBSchema)

			return c.JSON(http.StatusOK, map[string]interface{}{
				"status": true,
				"data":   VBSchema,
			})
		} else {

			VBSchema := models.VBSchema{}

			DB.DB.Where("id = ?", id).First(&VBSchema)

			if type_ == "form"{

				if condition != ""{
					if condition != "builder"{
						return dataform.SetCondition(condition, c, VBSchema)
					}
				}
			}

			return c.JSON(http.StatusOK, map[string]interface{}{
				"status": true,
				"data":   VBSchema,
			})
		}





	} else {

		VBSchemas := []models.VBSchemaList{}

		DB.DB.Select("id, name, type, created_at, updated_at").Where("type = ?", type_).Find(&VBSchemas)

		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": true,
			"data":   VBSchemas,
		})
	}

	return c.JSON(http.StatusBadRequest, map[string]interface{}{
		"status": false,
	})

}
func SaveVB(modelName string) echo.HandlerFunc {
	return func(c echo.Context) error {
		type_ := c.Param("type")
		id := c.Param("id")
		//condition := c.Param("condition")

		vbs := new(vb_schema)
		if err := c.Bind(vbs); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"status": false,
				"error": err.Error(),
			})
		}

		if id != "" {
			id_, _ := strconv.ParseUint(id, 0, 64)

			vb := models.VBSchema{}

			DB.DB.Where("id = ?", id_).First(&vb)

			vb.Name = vbs.Name
			vb.Schema = vbs.Schema
			//_, err := vb.Update(context.Background(), DB, boil.Infer())

			BeforeSave(id_, type_)

			err := DB.DB.Save(&vb).Error

			if type_ == "form" {
				//WriteModelData(id_)
				//WriteModelData(modelName)
				//WriteModelDataById(modelName, vb.ID)
			} else if type_ == "grid" {
				//WriteGridModel(modelName)
				//WriteGridModelById(modelName, vb.ID)
			}

			if err != nil {

				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"status": false,
					"error": err.Error(),
				})
			} else {

				error := AfterSave(vb, type_)

				if(error != nil){
					return c.JSON(http.StatusOK, map[string]interface{}{
						"status": false,
						"error":error.Error(),
					})
				} else {
					return c.JSON(http.StatusOK, map[string]interface{}{
						"status": true,
					})
				}
			}

		} else {
			vb := models.VBSchema{
				Name:   vbs.Name,
				Schema: vbs.Schema,
				Type:   type_,
				ID:0,
			}

			//err := vb.Insert(context.Background(), DB, boil.Infer())

			DB.DB.NewRecord(vb) // => returns `true` as primary key is blank

			err := DB.DB.Create(&vb).Error

			if type_ == "form" {
				//WriteModelData(vb.ID)
				//WriteModelData(modelName)
				//WriteModelDataById(modelName, vb.ID)
			} else if type_ == "grid" {
				//WriteGridModelById(modelName, vb.ID)
				//WriteGridModel(modelName)
			}



			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{
					"status": "false",
				})
			} else {
				error := AfterSave(vb, type_)

				if(error != nil){
					return c.JSON(http.StatusOK, map[string]interface{}{
						"status": false,
						"error":error.Error(),
					})
				} else {
					return c.JSON(http.StatusOK, map[string]interface{}{
						"status": true,
					})
				}

			}

		}

		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status": false,
		})
	}
}
func DeleteVB(c echo.Context) error {

	type_ := c.Param("type")
	id := c.Param("id")
	//condition := c.Param("condition")

	vbs := new(vb_schema)
	id_, _ := strconv.ParseUint(id, 0, 64)

	BeforeDelete(id_, type_)

	err := DB.DB.Where("id = ?", id).Where("type = ?", type_).Delete(&vbs).Error

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "false",
		})
	} else {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "true",
		})
	}



}


func GetProjectVB(c echo.Context) error {

	pid := c.Param("pid")
	type_ := c.Param("type")
	id := c.Param("id")
	condition := c.Param("condition")

	if id != "" {





			VBSchema := models.VBSchema{}

			DB.DB.Table("project_schemas").Where("id = ? AND projects_id = ?", id, pid).First(&VBSchema)

			if type_ == "form"{

				if condition != ""{
					if condition != "builder"{
						return dataform.SetCondition(condition, c, VBSchema)
					}
				}
			}

			return c.JSON(http.StatusOK, map[string]interface{}{
				"status": true,
				"data":   VBSchema,
			})






	} else {
		VBSchemas := []models.VBSchemaList{}

		DB.DB.Table("project_schemas").Select("id, name, type, created_at, updated_at").Where("type = ? AND projects_id = ?", type_, pid).Find(&VBSchemas)

		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": true,
			"data":   VBSchemas,
		})
	}

	return c.JSON(http.StatusBadRequest, map[string]interface{}{
		"status": false,
	})

}
func SaveProjectVB(modelName string) echo.HandlerFunc {
	return func(c echo.Context) error {
		pid := c.Param("pid")
		ProjectID, _ := strconv.Atoi(pid)
		type_ := c.Param("type")
		id := c.Param("id")
		//condition := c.Param("condition")

		vbs := new(vb_schema)
		if err := c.Bind(vbs); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"status": false,
				"error": err.Error(),
			})
		}

		if id != "" {
			id_, _ := strconv.ParseUint(id, 0, 64)

			vb := models.ProjectVBSchema{}

			DB.DB.Where("id = ?", id_).First(&vb)

			vb.Name = vbs.Name
			vb.ProjectID = ProjectID
			vb.Schema = vbs.Schema
			//_, err := vb.Update(context.Background(), DB, boil.Infer())

			BeforeSave(id_, type_)

			err := DB.DB.Save(&vb).Error

			if type_ == "form" {
				//WriteModelData(id_)
				//WriteModelData(modelName)
				//WriteModelDataById(modelName, vb.ID)
			} else if type_ == "grid" {
				//WriteGridModel(modelName)
				//WriteGridModelById(modelName, vb.ID)
			}

			if err != nil {

				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"status": false,
					"error": err.Error(),
				})
			} else {

				//error := AfterSave(vb, type_)
				//
				//if(error != nil){
				//	return c.JSON(http.StatusOK, map[string]interface{}{
				//		"status": false,
				//		"error":error.Error(),
				//	})
				//} else {
					return c.JSON(http.StatusOK, map[string]interface{}{
						"status": true,
					})
				//}
			}

		} else {
			vb := models.ProjectVBSchema{
				Name:   vbs.Name,
				Schema: vbs.Schema,
				Type:   type_,
				ProjectID: ProjectID,
				ID:0,
			}

			//err := vb.Insert(context.Background(), DB, boil.Infer())

			DB.DB.NewRecord(vb) // => returns `true` as primary key is blank

			err := DB.DB.Create(&vb).Error

			if type_ == "form" {
				//WriteModelData(vb.ID)
				//WriteModelData(modelName)
				//WriteModelDataById(modelName, vb.ID)
			} else if type_ == "grid" {
				//WriteGridModelById(modelName, vb.ID)
				//WriteGridModel(modelName)
			}



			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{
					"status": "false",
				})
			} else {
				//error := AfterSave(vb, type_)
				//
				//if(error != nil){
				//	return c.JSON(http.StatusOK, map[string]interface{}{
				//		"status": false,
				//		"error":error.Error(),
				//	})
				//} else {
					return c.JSON(http.StatusOK, map[string]interface{}{
						"status": true,
					})
				//}

			}

		}

		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status": false,
		})
	}
}
func DeleteProjectVB(c echo.Context) error {

	pid := c.Param("pid")
	type_ := c.Param("type")
	id := c.Param("id")
	//condition := c.Param("condition")

	vbs := new(vb_schema)
	//id_, _ := strconv.ParseUint(id, 0, 64)
	//
	//BeforeDelete(id_, type_)

	err := DB.DB.Table("project_schemas").Where("id = ? AND projects_id = ? AND type = ?", id, pid, type_).Delete(&vbs).Error

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "false",
		})
	} else {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "true",
		})
	}



}

func BeforeDelete(id uint64, type_ string){

	if type_ == "datasource"{
		vb := models.VBSchema{}

		DB.DB.Where("id = ?", id).First(&vb)

		datasource.DeleteView("ds_"+vb.Name)
	}

}
func BeforeSave(id uint64, type_ string){

	if type_ == "datasource"{
		vb := models.VBSchema{}

		DB.DB.Where("id = ?", id).First(&vb)

		datasource.DeleteView("ds_"+vb.Name)
	}

}
func AfterSave(vb models.VBSchema, type_ string) error{

	if type_ == "datasource"{
		return datasource.CreateView(vb.Name, vb.Schema)
	}

	return nil

}


/*GRID*/

func GridVB(GetGridMODEL func(schema_id string) (interface{}, interface{}, string, string, interface{}, string)) echo.HandlerFunc {
	return func(c echo.Context) error {
		schemaId := c.Param("schemaId")
		action := c.Param("action")
		id := c.Param("id")

		return datagrid.Exec(c, schemaId, action, id, GetGridMODEL)
	}
}


/*FROM*/

func GetOptions(c echo.Context) error {

	r := new(dataform.Relations)
	if err := c.Bind(r); err != nil {

		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status": false,
			"error": err.Error(),
		})
	}
	optionsData := map[string][]map[string]interface{}{}

	var DB_ *sql.DB
	DB_ = DB.DB.DB()
	for table, relation := range r.Relations {
		data := dataform.OptionsData(DB_, relation, c)
		optionsData[table] = data

	}
	return c.JSON(http.StatusOK, optionsData)

}