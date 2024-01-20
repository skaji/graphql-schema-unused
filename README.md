# graphql-schema-unused

Detect unused types in graphql schema.

# Install

Download binaries from [release pages](https://github.com/skaji/graphql-schema-unused/releases/latest).

# Usage

Let's say you have the following "schema.graphql":

```graphql
schema {
  query: Query
}

type Query {
  user(id: Int!): User
}

type User {
  name: String!
  age: Int!
}

type Animal {
  name: String!
}
```

Then:

```shell
❯ graphql-schema-unused schema.graphql
unused type Animal at schema.graphql line 14

❯ echo $?
1
```

# Author

Shoichi Kaji

# License

MIT
