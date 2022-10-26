# test

## apis

#### create api

parameters:

- name: api name
- method: http method
- sqlType: chain or template
- sqlTemplate: sql template (golang template)
- sqlTemplateParameters: parameters for template
- sqlChain: gorm chain for sql 

example:

1. Post http://localhost:8888/api/create

```json
{
    "name": "updateSystem",
    "method": "Put",
    "sqlType": "chain",
    "sqlChain": {
        "table": "settings",
        "where": [
            {"key": "id", "value": "id = ?","valid":"number,min=0"},
        ],
        "update": [
            {"key": "content", "valid":"str,len>0"},
            {"key": "desc", "valid":"str,len>0"}
        ]
    }
}
```

2. Post http://localhost:8888/api/create

```json
{
    "db": "db1",
    "id": 0,
    "name": "listApi2",
    "method": "POST",
    "desc": "test",
    "sqlType": "template",
    "sqlTemplate": "select api_name,method,desc,sql_type from apis where api_name like @name {{ if gt .limit 0 }} limit {{ .limit }} {{ end }}",
    "SqlTemplateParameters": {
        "limit": "number,gt=0",
        "name": "gt=0"
    },
    "SqlTemplateResult": {}
}
```

### request api

- name: api name
- method: http method (Get Post Delete)
- sqlType: chain or template
- sqlTemplateParameters: parameters for template
- sqlChainParameters: gorm chain for sql 

example:

- Get http://localhost:8888/api/<name>/<sqlType>?<sqlTemplateParameters>
- Delete http://localhost:8888/api/<name>/<sqlType>?<sqlTemplateParameters>
- Post http://localhost:8888/api/<name>/<sqlType>   with json body: {}


1. Put http://localhost:8888/api/updateSystem/chain

```json
{
  "chainParameters": {
    "where": {
        "id": 18
    },
    "update":{
        "content": "xxxx",
        "desc": "xxx"
    }
  }
}
```

2. Put http://localhost:8888/api//template

```json
{
  "templateParameters": {
    "id": 18,
    "desc": "xxx",
    "key": "xxx",
    "eid": ""
  }
}
```
