-  full test suite (mid way)
-  toTeams toRecords are all same, combine em

**steps:**

-  create sps for each purpose

`az ad sp create-for-rbac --name "terraform-sp" --role 'Contributor' --scopes /subscriptions/`

`az ad sp create-for-rbac --name "docker-sp" --role 'AcrPush' --scopes /subscriptions/`

-  login to terraform with terraform sp (github ci cd)
-  login to docker with docker sp

`az login --service-principal -u ${{AZURE_SERVICE_PRINCIPAL_APPID}} -p ${{AZURE_SERVICE_PRINCIPAL_PASSWORD}} --tenant ${{AZURE_SERVICE_PRINCIPAL_TENANT}}`

`az acr login --name <your-acr-name>`
