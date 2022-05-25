# data_model_ci (CI Resource) 

```hcl

resource "data_model_ci" "test" {
  provider     = ucmdb.cms
  type = "nt"
  properties {
    name = "test.com"
    description = "Created by terraform ucmdb provider"
    ...
  }
}

```


## Argument Reference

- `type` (String) Valid UCMDB CI type.
- `name` (String) Name of the CI item.
- `description` (String) Long description for the CI.




