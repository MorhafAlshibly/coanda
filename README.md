**steps:**

-  create sps for each purpose

`az ad sp create-for-rbac --name "infra-sp" --role 'Owner' --scopes /subscriptions/...`

`az role assignment create --assignee "${{AZURE_SERVICE_PRINCIPAL_APPID}}" --role "AcrPush" --subscription ...`

-  login to sp

`az login --service-principal -u ${{AZURE_SERVICE_PRINCIPAL_APPID}} -p ${{AZURE_SERVICE_PRINCIPAL_PASSWORD}} --tenant ${{AZURE_SERVICE_PRINCIPAL_TENANT}}`

`az acr login --name <your-acr-name>`

**todo:**

-  put api behind api key
-  team add id auto increment and unique to the current primary key?
-  update record its own and check for record better or not? or just have create record update a record
-  check indexes are good
-  use null structs in model folders (done for tournaments, rest not done)
-  use context for limiting request lifetime and for api key
-  tournies still using sq not gq
-  seperate sql tests into dynamic and non dynamic tests
-  remove team owner - fix team member number wont work if players leave

**some things u need:** (outdated)

-  jq
-  dos2unix
-  task
-  terraform latest for azure
-  terraform 1.5.7 for oci

**docs:**
https://morhafalshibly.github.io/coanda/
