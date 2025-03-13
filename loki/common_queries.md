# Common Queries

Common LogQL queries.

> [!TIP]
> LogQL doesn't care about spaces, tabs, carriage returns, etc. These are all the same:
>
> ```
> {label="value"}|="str"
> 
> { label = "value" } |= "str"
> 
> { label = "value" } 
> |= "str"
>
> { label = "value" } 
>   |= "str"
> ```

## Basic log query

Search for logs that contain a specific string.

```txt
{ label="value" } |= "string match"
```

## Basic metric query

Count the number of log messages over time that match the search criteria.

```txt
count_over_time( 
  { label="value" } 
  |= "string match" [$__auto] 
)
```

## Parse JSON log message

```txt
{ label="value" } 
|= "string match" 
| json keyName = "key.name"
```

```txt
{ label="value" } |= "string match" 
| json name = "object.name" 
| name = "example"
```

## Parse Logfmt log message

```txt
{ label="value" } |= "string match" 
| logfmt objectName
```

```txt
{ label="value" } |= "string match" 
| logfmt objectName 
| objectName = "example"
```

## Pattern parser

```txt
{ label="value" } |= "string match" 
| pattern "<_>objectName=<objectName> <_>"
```

```txt
{ label="value" } |= "string match" 
| pattern "<_>objectName=<objectName> <_>" 
| objectName = "example"
```

## Regex parser

> [!NOTE]
> Special characters must be escaped twice. `\s` becomes `\\s`, `\?` becomes `\\?`, etc.

```txt
{ label="value" } |= "string match" 
| regexp "^.*objectName=(?P<objectName>.*)\\s.*$"
```

```txt
{ label="value" } |= "string match" 
| regexp "^.*objectName=(?P<objectName>.*)\\s.*$" 
| objectName = "example"
```

## Basic unwrap query

```txt
sum_over_time(
  { label="value" } |= "string match"
  | json number = "object.number"
  | unwrap number
  [$__auto]
)
```
