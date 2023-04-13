---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "forwardnetworks_checks Data Source - forwardnetworks"
subcategory: ""
description: |-
  
---

# forwardnetworks_checks (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `snapshot_id` (String)

### Optional

- `check_id` (String)
- `priority` (String)
- `status` (String)
- `type` (String)

### Read-Only

- `checks` (List of Object) (see [below for nested schema](#nestedatt--checks))
- `id` (String) The ID of this resource.

<a id="nestedatt--checks"></a>
### Nested Schema for `checks`

Read-Only:

- `check_type` (String)
- `creation_date_millis` (Number)
- `creator_id` (String)
- `definition_date_millis` (Number)
- `description` (String)
- `enabled` (Boolean)
- `execution_date_millis` (Number)
- `execution_duration_millis` (Number)
- `id` (String)
- `name` (String)
- `predefined_check_type` (String)
- `priority` (String)
- `status` (String)

