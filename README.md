**steps:**

-  create sps for each purpose

`az ad sp create-for-rbac --name "infra-sp" --role 'Owner' --scopes /subscriptions/...`

`az role assignment create --assignee "${{AZURE_SERVICE_PRINCIPAL_APPID}}" --role "AcrPush" --subscription ...`

-  login to sp

`az login --service-principal -u ${{AZURE_SERVICE_PRINCIPAL_APPID}} -p ${{AZURE_SERVICE_PRINCIPAL_PASSWORD}} --tenant ${{AZURE_SERVICE_PRINCIPAL_TENANT}}`

`az acr login --name <your-acr-name>`

**todo:**

-  cloudflare auth, stop ddos and block out non authenticated members
-  check indexes are good
-  use context for limiting request lifetime and for api key
-  matchmaking tests and cron job tests
-  disable graphql schema checking?
-  check if parent and child pagination is working as intended for events and mb matchmaiking
-  use null structs in model folders (done for tournaments, rest not done)
-  seperate sql tests into dynamic and non dynamic tests

**some things u need:** (outdated)

-  jq
-  dos2unix
-  task
-  terraform latest for azure
-  terraform 1.5.7 for oci

**docs:**
https://morhafalshibly.github.io/coanda/
