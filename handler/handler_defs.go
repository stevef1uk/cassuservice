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
    "crypto/tls"
	"crypto/x509"
    "io/ioutil"
    "path/filepath"
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

	IMPORTFORPOST = `
    "{{.PATH}}/restapi/operations/{{.TABLENAME}}"
     `

	IMPORTFORPOST2 = `
    "fmt"
    "strconv"
     `
	IMPORTFORPOST2A = `
    "fmt"
     `

	PARSETIME = `
func parseTime ( input string) time.Time {
    var ret time.Time
    if input == "" {
        ret = time.Now()
    } else {
        ret, _ = time.Parse( time.RFC3339, input )
    }
    return ret;
}
`
	PARSERTIME_FUNC_NAME = "parseTime"

	HEADER =`
var ` + SESSION_VAR + ` *gocql.Session

func SetUp() {
  var err error
  log.Println("Trying to connect to Cassandra database using ", os.Getenv("CASSANDRA_SERVICE_HOST"))
  cluster := gocql.NewCluster(os.Getenv("CASSANDRA_SERVICE_HOST"))
  cluster.Keyspace = "{{.KEYSPACE}}"
  cluster.Consistency = gocql.One
  username := os.Getenv("CASSANDRA_USERNAME")
  password := os.Getenv("CASSANDRA_PASSWORD")
  if username != "" {
     log.Println("Using credentials, username = ", username)
          cluster.Authenticator = gocql.PasswordAuthenticator{
                Username: username,
                Password: password,
    }
  } else {
     log.Println("Are you sure you don't need to set $CASSANDRA_USERNAME and $CASSANDRA_PASSWORD")
  }` + `
  astra := os.Getenv("ASTRA_SECURE_CONNECT_PATH")
  if ( astra != "" ) {
    if os.Getenv("ASTRA_PORT") != "" {
		astra = astra + string(os.PathSeparator)
		cluster.Hosts = []string{os.Getenv("CASSANDRA_SERVICE_HOST") + ":" + os.Getenv("ASTRA_PORT")}
		certPath, _ := filepath.Abs(astra + "cert")
		keyPath, _ := filepath.Abs(astra + "key")
		caPath, _ := filepath.Abs(astra + "ca.crt")
		cert, _ := tls.LoadX509KeyPair(certPath, keyPath)
		caCert, _ := ioutil.ReadFile(caPath)
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
		}
		cluster.SslOpts = &gocql.SslOptions{
			Config: tlsConfig,
			EnableHostVerification: false,
		}
	} else {
		log.Fatal("With Datastax Astra you need to set ASTRA_PORT environment variable (in secure connect download file cqlshrc")
	}
  }
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

	GS_PARAMS = "Params"
	GS_BODY = "Body"
)

