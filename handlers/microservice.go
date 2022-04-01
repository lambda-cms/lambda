package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/lambda-platform/datasource"
	"github.com/lambda-platform/lambda/DBSchema"
	"github.com/lambda-platform/lambda/config"
	"github.com/lambda-platform/lambda/models"
	krudModels "github.com/lambda-platform/krud/models"
	agentModels "github.com/lambda-platform/agent/models"
	"strings"
	"io"
	"errors"
	"archive/zip"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"fmt"
	"path/filepath"
	"strconv"
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/utils"
)

func UploadSCHEMA(c echo.Context) error {

	UploadDBSCHEMA()


	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": true,
	})
}
func UploadDBSCHEMA()  {

	DBSchema.GenerateSchemaForCloud()



	url := config.LambdaConfig.LambdaMainServicePath+"/console/upload/"+config.LambdaConfig.ProjectKey
	//url := "http://localhost/console/upload/"+config.LambdaConfig.ProjectKey
	method := "POST"


	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, errFile1 := os.Open("app/models/db_schema.json")
	defer file.Close()
	part1,
	errFile1 := writer.CreateFormFile("file",filepath.Base("app/models/db_schema.json"))
	_, errFile1 = io.Copy(part1, file)
	if errFile1 != nil {
		fmt.Println(errFile1)

	}


	file2, errFile2 := os.Open("lambda.json")
	defer file2.Close()
	part2,
	errFile2 := writer.CreateFormFile("lambda_config",filepath.Base("lambda.json"))
	_, errFile2 = io.Copy(part2, file2)
	if errFile2 != nil {
		fmt.Println(errFile2)

	}
	err := writer.Close()
	if err != nil {
		fmt.Println(err)

	}


	client := &http.Client {
	}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)

	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)

	}
	defer res.Body.Close()

	_, err2 := ioutil.ReadAll(res.Body)
	if err2 != nil {
		fmt.Println(err)


	}




}
func ASyncFromCloud()  {


	userWithUUID := "false"

	if config.Config.SysAdmin.UUID{
		userWithUUID = "true"
	}

	url := config.LambdaConfig.LambdaMainServicePath+"/console/project-data/"+config.LambdaConfig.ProjectKey+"/"+config.LambdaConfig.ModuleName+"/"+userWithUUID

	method := "GET"

	client := &http.Client {
	}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)

	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)

	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)

	}


	data := CloudData{}
	json.Unmarshal(body, &data)

	//DSVBS := []models.VBSchema{}
	FormVbs := []models.VBSchema{}
	GridVbs := []models.VBSchema{}
	MenuVbs := []models.VBSchema{}
	ChartVbs := []models.VBSchema{}
	MoqupVbs := []models.VBSchema{}
	DatasourceVbs := []models.VBSchema{}
	cruds := []krudModels.Krud{}
	FormSchemasJSON, _ := json.Marshal(data.FormSchemas)
	GridSchemasJSON, _ := json.Marshal(data.GridSchemas)
	MenuSchemasJSON, _ := json.Marshal(data.MenuSchemas)
	ChartSchemasJSON, _ := json.Marshal(data.ChartSchemas)
	MoqupSchemasJSON, _ := json.Marshal(data.MoqupSchemas)
	DatasourceSchemasJSON, _ := json.Marshal(data.DatasourceSchemas)
	KrudJSON, _ := json.Marshal(data.Cruds)
	json.Unmarshal([]byte(FormSchemasJSON), &FormVbs)
	json.Unmarshal([]byte(GridSchemasJSON), &GridVbs)
	json.Unmarshal([]byte(MenuSchemasJSON), &MenuVbs)
	json.Unmarshal([]byte(ChartSchemasJSON), &ChartVbs)
	json.Unmarshal([]byte(MoqupSchemasJSON), &MoqupVbs)
	json.Unmarshal([]byte(DatasourceSchemasJSON), &DatasourceVbs)
	json.Unmarshal([]byte(KrudJSON), &cruds)

	//DB.DB.Where("type = ?", "datasource").Find(&DSVBS)

	//DB.DB.Exec("TRUNCATE krud")
	//DB.DB.Exec("TRUNCATE vb_schemas")
	//
	//for _, vb := range FormVbs {
	//	DB.DB.Create(&vb)
	//}
	//for _, vb := range GridVbs {
	//	DB.DB.Create(&vb)
	//}
	//for _, vb := range MenuVbs {
	//	DB.DB.Create(&vb)
	//}
	//for _, vb := range ChartVbs {
	//	DB.DB.Create(&vb)
	//}
	//for _, vb := range MoqupVbs {
	//	DB.DB.Create(&vb)
	//}

	for _, ds := range DatasourceVbs {
		datasource.DeleteView(ds.Name)
		errSave:= datasource.CreateView(ds.Name, ds.Schema)
		if errSave != nil {
			fmt.Println(errSave.Error())
		}
	}

	var downloadError error = DownloadGeneratedCodes()

	if(downloadError != nil){
		fmt.Println(downloadError)
	} else {

		var unzip error = UnZipLambdaCodes()


		if(unzip != nil){
			fmt.Println(unzip)
		}

	}



}
func GetRolesData(c echo.Context) error {

	GetRoleData()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": true,
	})
}
func GetRoleData() error{
	url := config.LambdaConfig.LambdaMainServicePath+"/role-data"
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("file not found error")
	}

	data := map[string]interface{}{}


	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	json.Unmarshal(bodyBytes, &data)

	roleData := map[int]map[string]interface{}{}
	roleDataPre, _ := json.Marshal(data["roleData"])
	json.Unmarshal(roleDataPre, &roleData)

	Roles := []agentModels.Role{}
	RolesPre, _ := json.Marshal(data["roles"])
	json.Unmarshal(RolesPre, &Roles)

	for k, data := range roleData {


		bolB, _ := json.Marshal(data)
		_ = ioutil.WriteFile("lambda/role_"+strconv.Itoa(k)+".json", bolB, 0777)
	}

	DB.DB.Exec("TRUNCATE roles")
	for _, Role := range Roles {
		DB.DB.Create(&Role)
	}



	return nil
}
func DownloadGeneratedCodes() error{
	url := config.LambdaConfig.LambdaMainServicePath+"/console/get-codes/"+config.LambdaConfig.ProjectKey

	resp, err := http.Get(url)
	if err != nil {
		return err
	}


	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("file not found error")
	}

	// Create the file
	out, err := os.Create("lambda.zip")
	if err != nil {

		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)


	return err
}
func UnZipLambdaCodes() error{
	var dest string = "lambda"
	var src string = "lambda.zip"
	if(!utils.FileExists(src)){
		return errors.New("Lambda file Not found")
	} else {
		formPatch :="lambda/models/form/"
		gridPatch := "lambda/models/grid/"
		if _, err := os.Stat(formPatch); os.IsNotExist(err) {
			os.MkdirAll("lambda/models/form", 0755)
			os.MkdirAll(formPatch, 0755)

			os.MkdirAll("lambda/models/form/caller/", 0755)
		} else {
			os.MkdirAll("lambda/models/form", 0755)
			os.RemoveAll(formPatch)
			os.MkdirAll(formPatch, 0755)

			os.MkdirAll("lambda/models/form/caller/", 0755)
		}
		if _, err := os.Stat(gridPatch); os.IsNotExist(err) {
			os.MkdirAll("lambda/models/grid", 0755)
			os.MkdirAll(gridPatch, 0755)
			os.MkdirAll("lambda/models/grid/caller", 0755)
		} else {
			os.MkdirAll("lambda/models/grid", 0755)
			os.RemoveAll(gridPatch)
			os.MkdirAll(gridPatch, 0755)
			os.MkdirAll("lambda/models/grid/caller", 0755)
		}

		graphqlPatch :=  "lambda/graph"
		graphqlGeneratedPatch :=  "lambda/graph/generated"
		modelsPatch :=  "lambda/graph/models"
		schemaPatch :=  "lambda/graph/schemas"
		resolversPatch :=  "lambda/graph/resolvers"
		schemaCommonPatch :=  "lambda/graph/schemas-common"
		if _, err := os.Stat(modelsPatch); os.IsNotExist(err) {

			os.MkdirAll(graphqlPatch, 0755)
			os.MkdirAll(graphqlGeneratedPatch, 0755)
			os.MkdirAll(modelsPatch, 0755)
			os.MkdirAll(schemaPatch, 0755)
			os.MkdirAll(resolversPatch, 0755)
			os.MkdirAll(schemaCommonPatch, 0755)

		} else {

			os.RemoveAll(graphqlPatch)
			os.RemoveAll(graphqlGeneratedPatch)
			os.RemoveAll(modelsPatch)
			os.RemoveAll(schemaPatch)
			os.RemoveAll(resolversPatch)
			os.RemoveAll(schemaCommonPatch)
			os.MkdirAll(graphqlPatch, 0755)
			os.MkdirAll(graphqlGeneratedPatch, 0755)
			os.MkdirAll(modelsPatch, 0755)
			os.MkdirAll(schemaPatch, 0755)
			os.MkdirAll(resolversPatch, 0755)
			os.MkdirAll(schemaCommonPatch, 0755)
		}

		var filenames []string

		r, err := zip.OpenReader(src)
		if err != nil {
			return  err
		}
		defer r.Close()

		for _, f := range r.File {

			// Store filename/path for returning and using later on
			fpath := filepath.Join(dest, f.Name)

			// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
			if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
				return fmt.Errorf("%s: illegal file path", fpath)
			}

			filenames = append(filenames, fpath)

			if f.FileInfo().IsDir() {
				// Make Folder
				os.MkdirAll(fpath, os.ModePerm)
				continue
			}

			// Make File
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}

			rc, err := f.Open()
			if err != nil {
				return  err
			}

			_, err = io.Copy(outFile, rc)

			// Close the file without defer to close before next iteration of loop
			outFile.Close()
			rc.Close()

			if err != nil {
				return  err
			}
		}
		e := os.Remove(src)
		if e != nil {
			return e
		}
		return nil
	}







}
type CloudData struct {
	Cruds []struct {
		Form       int    `json:"form"`
		Grid       int    `json:"grid"`
		ID         int    `json:"id"`
		ProjectsID int    `json:"projects_id"`
		Template   string `json:"template"`
		Title      string `json:"title"`
	} `json:"cruds"`
	GridSchemas []struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		ProjectsID int    `json:"projects_id"`
		Schema     string `json:"schema"`
		Type       string `json:"type"`
	} `json:"form-schemas"`
	FormSchemas []struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		ProjectsID int    `json:"projects_id"`
		Schema     string `json:"schema"`
		Type       string `json:"type"`
	} `json:"grid-schemas"`
	MenuSchemas []struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		ProjectsID int    `json:"projects_id"`
		Schema     string `json:"schema"`
		Type       string `json:"type"`
	} `json:"menu-schemas"`
	ChartSchemas []struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		ProjectsID int    `json:"projects_id"`
		Schema     string `json:"schema"`
		Type       string `json:"type"`
	} `json:"chart-schemas"`
	MoqupSchemas []struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		ProjectsID int    `json:"projects_id"`
		Schema     string `json:"schema"`
		Type       string `json:"type"`
	} `json:"moqup-schemas"`
	DatasourceSchemas []struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		ProjectsID int    `json:"projects_id"`
		Schema     string `json:"schema"`
		Type       string `json:"type"`
	} `json:"datasource-schemas"`
}
