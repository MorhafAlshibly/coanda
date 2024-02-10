**steps:**

-  create sps for each purpose

`az ad sp create-for-rbac --name "infra-sp" --role 'Owner' --scopes /subscriptions/...`

`az role assignment create --assignee "${{AZURE_SERVICE_PRINCIPAL_APPID}}" --role "AcrPush" --subscription ...`

-  login to sp

`az login --service-principal -u ${{AZURE_SERVICE_PRINCIPAL_APPID}} -p ${{AZURE_SERVICE_PRINCIPAL_PASSWORD}} --tenant ${{AZURE_SERVICE_PRINCIPAL_TENANT}}`

`az acr login --name <your-acr-name>`

**todo:**

-  fix owner required in teams (parsing duplicate key errors)
-  put api behind api key
-  tournament, record, team add id auto increment and unique to the current primary key
-  update record its own and check for record better or not? or just have create record update a record
-  now with sql
   -  make apis all use sql
   -  make all apis return good data
   -  check indexes are good
   -  wrapper for sqlc that incoporates caching and error handling
   -  updates in single endpoints
   -  check for sqlinjec

**some things u need:**

-  jq
-  dos2unix
-  task
-  terraform latest for azure
-  terraform 1.5.7 for oci
