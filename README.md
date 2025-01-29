**steps:**

-  create sps for each purpose

`az ad sp create-for-rbac --name "infra-sp" --role 'Owner' --scopes /subscriptions/...`

`az role assignment create --assignee "${{AZURE_SERVICE_PRINCIPAL_APPID}}" --role "AcrPush" --subscription ...`

-  login to sp

`az login --service-principal -u ${{AZURE_SERVICE_PRINCIPAL_APPID}} -p ${{AZURE_SERVICE_PRINCIPAL_PASSWORD}} --tenant ${{AZURE_SERVICE_PRINCIPAL_TENANT}}`

`az acr login --name <your-acr-name>`

**todo:**

-  server error in one microservice errors entire request, find solution to fix.
-  delete old matches, tickets and users
-  check indexes are good
-  matchmaking tests, cron job tests, seperate sql tests into dynamic and non dynamic tests, renaming tests to common format
-  either mm arena or arena, consistent naming
-  consisteny sq gq use and using booleanexpressions in filters
-  set consistent default value when converting enum e.g mtErr

**some things u need:** (outdated)

-  jq
-  dos2unix
-  task
-  terraform latest for azure
-  terraform 1.5.7 for oci

**docs:**
https://morhafalshibly.github.io/coanda/
