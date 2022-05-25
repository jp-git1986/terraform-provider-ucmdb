# ucmdb_list (Data Source)

### Get the list of CI irems based on filter.

## Example Usage

```hcl
data "ucmdb_list" "my_nt_servers" {
    provider     = ucmdb.cms
    filter {
        type = "nt"
        names = ["mimi.ecb.de","dada.ecb.de"]
    }
}

}
```


## Argument Reference


- `filter` (Block) Filter block must be provided with a pair of "type" and "names" attributes.
- `type` (String) Valid UCMDB type.
- `name` (String) List of CI item names to retrieve.





