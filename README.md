**steps:**

-  create sps for each purpose

`az ad sp create-for-rbac --name "infra-sp" --role 'Owner' --scopes /subscriptions/...`

`az role assignment create --assignee "${{AZURE_SERVICE_PRINCIPAL_APPID}}" --role "AcrPush" --subscription ...`

-  login to sp

`az login --service-principal -u ${{AZURE_SERVICE_PRINCIPAL_APPID}} -p ${{AZURE_SERVICE_PRINCIPAL_PASSWORD}} --tenant ${{AZURE_SERVICE_PRINCIPAL_TENANT}}`

`az acr login --name <your-acr-name>`

**todo:**

-  endpoint to set a privateServerCode for a match and if already sent, error and return current code
-  set elo across all arenas when creating user, still have setmatchmakinguser elo tho, also make users only live for ticket life, so they are temp
-  return matchmaking ticket in poll
-  api gateway for authentication (overall infra redo)
-  check indexes are good
-  matchmaking tests, cron job tests, seperate sql tests into dynamic and non dynamic tests, renaming tests to common format
-  check if parent and child pagination is working as intended for events and mb matchmaiking

**some things u need:** (outdated)

-  jq
-  dos2unix
-  task
-  terraform latest for azure
-  terraform 1.5.7 for oci

**docs:**
https://morhafalshibly.github.io/coanda/
