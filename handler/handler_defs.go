package handler

// Type used for template processing
type tableDetails struct {
	PATH   string
	KEYSPACE string
	TABLENAME string
	EXPORTPATH string
}


const (
	MAINFILE = "GeneratedHandler.go"
	TMP = "tmp"
    MODELS = "models." // this is the directory where the map types are created by go-swagger
    OPERATIONS = "operations."
    TEMP_VAR_PREFIX = "tmp_"
    RAWRESULT = TEMP_VAR_PREFIX + "_" + "resultSet"


	COMMONIMPORTS = `// GENERATED FILE so do not edit or will be overwritten upon next generate
package data

import (
    "{{.PATH}}/models"
    "{{.PATH}}/restapi/operations"
    middleware "github.com/go-openapi/runtime/middleware"
    "github.com/gocql/gocql"
    "os"
    "log"`

	IMPORTSEND = `
)
`
// strfmt "github.com/go-openapi/strfmt"
	IMPORTSTIMESTAMP = `
    "time"
     `

	// "gopkg.in/inf.v0"
	//"strconv"
	IMPORTDEC = `
	"gopkg.in/inf.v0"
	"strconv"
	`

	HEADER =`
var ` + SESSION_VAR + ` *gocql.Session

func SetUp() {
  var err error
  log.Println("Tring to connect to Cassandra database using ", os.Getenv("CASSANDRA_SERVICE_HOST"))
  cluster := gocql.NewCluster(os.Getenv("CASSANDRA_SERVICE_HOST"))
  cluster.Keyspace = "{{.KEYSPACE}}"
  cluster.Consistency = gocql.One` + `
` + "  "+ SESSION_VAR  + `, err = cluster.CreateSession()
  if ( err != nil ) {
    log.Fatal("Have you remembered to set the env var $CASSANDRA_SERVICE_HOST as connection to Cannandra failed with error = ", err)
  } else {
    log.Println("Yay! Connection to Cannandra established")
  }
}

func Stop() {
    log.Println("Shutting down the service handler")` + `
` + "  " + SESSION_VAR + `.Close()
}

func Search(params operations.Get{{.EXPORTPATH}}Params) middleware.Responder {
`

	INDENT_1 = "\n    "
	INDENT = "  "
	INDENT2 = "    "
	INDENT3 = "      "

	SELECT_OUTPUT = "codeGenRawTableResult"
	TMP_TIME_VAR_PREFIX = "cassuservice_tmp_time"
	SESSION_VAR = "cassuservice_session"
	PAYLOAD = "payLoad"
	PAYLOAD_STRUCT = "Payload"
	PARAMS_RET = "retParams"

	POST_HEADER = `
func Insert(params {{.KEYSPACE}}.Add{{.TABLENAME}}Params) middleware.Responder {
`


)

