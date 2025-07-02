# GraphQL-Markdown template

Docusaurus template for [GraphQL-Markdown](https://graphql-markdown.dev).

## Quick start

**1. Install**

```shell
npm init docusaurus my-website https://github.com/graphql-markdown/template.git
```

**2. Configure**

Update settings in `.graphqlrc` (see [documentation](https://graphql-markdown.dev/docs/configuration#graphql-config)).

```yaml
schema: 'https://api.react-finland.fi/graphql'
extensions:
  graphql-markdown:
    baseURL: '.'
    homepage: 'static/index.md'
    loaders:
      UrlLoader: '@graphql-tools/url-loader'
    docOptions:
      pagination: false
    printTypeOptions:
      deprecated: 'group'
```

**3. Generate**

```shell
npm run doc
```

**4. Start**

```shell
npm start
```
