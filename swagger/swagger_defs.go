package swagger

const (
	HEADER = `
swagger: '2.0'
info:
  version: 1.0.0
  title: Simple API
  description: A generated file representing a Cassandra Table definition
schemes:
  - http
host: localhost
basePath: /v1
paths:`

	TIMESTAMP = "TIMESTAMP"
	DATE = "DATE"
	TIMEUUID = "TIMEUUID"
)

/*
paths:
  /accounts4:
    get:
      summary: Gets some accounts4
      description: Returns a list containing all accounts4.
      parameters:
        - name: id
          in: query
          description: PK
          required: true
          type: integer
          format: int32
        - name: name
          in: query
          description: PK
          required: true
          type: string
          format: string
        - name: time1
          in: query
          description: PK
          required: true
          type: string
          format: date-time
      responses:
        200:
          description: A list of records
          schema:
            type: array
            items:
              required:
                - id
                - name
                - ascii1
                - bint1
                - blob1
                - bool1
                - counter1
                - dec1
                - double1
                - flt1
                - inet1
                - int1
                - text1
                - time1
                - time2
                - uuid1
                - varchar1
                - events
                - mylist
                - myset
                - mymap
              properties:
                 id:
                   type: integer
                 name:
                   type: string
                 ascii1:
                   type: string
                 bint1:
                   type: integer
                 blob1:
                   type: string
                 bool1:
                   type: boolean
                 counter1:
                   type: integer
                 dec1:
                   type: number
                 double1:
                   type: number
                 flt1:
                   type: number
                 inet1:
                   type: string
                 int1:
                   type: integer
                 text1:
                   type: string
                 time1:
                   type: string
                   format: date-time
                 time2:
                   type: string
                   format: date-time
                 uuid1:
                   type: string
                 varchar1:
                   type: string
                 events:
                   type: array
                   items:
                     type: integer
                 mylist:
                   type: array
                   items:
                     type: number
                 myset:
                   type: array
                   items:
                     type: string
                 mymap:
                   $ref: "#/definitions/mymap"
        400 :
          description: Record not found
        default:
          description: Sorry something went wrong
definitions:
  mymap:
    additionalProperties:
      type: string
 */