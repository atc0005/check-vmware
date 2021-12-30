<!-- omit in toc -->
# Type conversion notes

<!-- omit in toc -->
## Table of contents

- [About](#about)
- [Notes](#notes)
  - [ManagedEntity, CustomValue, BaseCustomFieldValue](#managedentity-customvalue-basecustomfieldvalue)
  - [BaseCustomFieldValue interface, CustomFieldValue](#basecustomfieldvalue-interface-customfieldvalue)
  - [CustomFieldValue](#customfieldvalue)
  - [CustomFieldValue, CustomFieldStringValue](#customfieldvalue-customfieldstringvalue)

## About

This is my scratch notes from trying to understand how I about about
retrieving the necessary values associated with Custom Attributes. I initially
could not find a way to get at the final key/value pair, but ended up
realizing that I needed to perform a type conversion. This was both
non-obvious to me and was a bit of a stretch for my current understanding of
the topic.

## Notes

### ManagedEntity, CustomValue, BaseCustomFieldValue

`ManagedEntity` has a `CustomValue` field.
The `CustomValue` field is of type `[]types.BaseCustomFieldValue`.
<https://pkg.go.dev/github.com/vmware/govmomi@v0.24.0/vim25/mo#ManagedEntity.CustomValue>

### BaseCustomFieldValue interface, CustomFieldValue

The `BaseCustomFieldValue` is an interface with a single
`GetCustomFieldValue()` method. This method returns a `*CustomFieldValue`
type.
<https://pkg.go.dev/github.com/vmware/govmomi@v0.24.0/vim25/types#BaseCustomFieldValue>

### CustomFieldValue

The `CustomFieldValue` type is a struct. This struct has a
`GetCustomFieldValue()` method that returns a `*CustomFieldValue`.

NOTE: Returning itself seems an odd choice (ignorance speaking), unless this
is being done to satisfy the interface for conversion purposes?

<https://pkg.go.dev/github.com/vmware/govmomi@v0.24.0/vim25/types#CustomFieldValue>

### CustomFieldValue, CustomFieldStringValue

By using `%+v` and `%#v` against a `types.BaseCustomFieldValue` and the
 returned `*types.CustomFieldValue` from `func
 (types.BaseCustomFieldValue).GetCustomFieldValue()`, I can tell that we want
 to end up with a `&types.CustomFieldStringValue` type in order to expose the
 `Value` field.

```golang
type CustomFieldStringValue struct {
  CustomFieldValue

  Value string `xml:"value"`
}
```

```golang
type CustomFieldValue struct {
  DynamicData

  Key int32 `xml:"key"`
}
```

Since the `CustomFieldStringValue` type contains a `CustomFieldValue`, I
believe I remember that this satisfies an "IsA" relationship? If so, we should
be able to convert "up"?

This seems to work:

```golang
  // ...
  var hs []mo.HostSystem
  err = v.Retrieve(ctx, []string{"HostSystem"}, nil,&hs)
  // ...
  for _, host := range hs {
    fmt.Println("Host:", host.Name)
    for _, cf := range host.CustomValue {
      fmt.Printf("Type assert to CustomFieldStringValue: %+v\n", cf.(*types.CustomFieldStringValue))
      fmt.Printf("CustomFieldStringValue (Key): %+v\n", cf.(*types.CustomFieldStringValue).Key)
      fmt.Printf("CustomFieldStringValue (Value): %+v\n", cf.(*types.CustomFieldStringValue).Value)
    }

    for _, af := range host.AvailableField {
      fmt.Printf("AvailableField.Key %v\n", af.Key)
      fmt.Printf("AvailableField.Name: %v\n", af.Name)
    }
```

Looking at the actual govmomi package examples, and core source code, this
type assertion is used quite a bit. Some examples:

```golang
// simulator\custom_fields_manager.go

// Iterates through all entities of passed field type;
// Removes found field from their custom field properties.
func entitiesFieldRemove(field types.CustomFieldDef) {
  entities := Map.All(field.ManagedObjectType)
  for _, e := range entities {
    entity := e.Entity()
    Map.WithLock(entity, func() {
      aFields := entity.AvailableField
      for i, aField := range aFields {
        if aField.Key == field.Key {
          entity.AvailableField = append(aFields[:i], aFields[i+1:]...)
          break
        }
      }

      values := e.Entity().Value
      for i, value := range values {
        if value.(*types.CustomFieldStringValue).Key == field.Key {
          entity.Value = append(values[:i], values[i+1:]...)
          break
        }
      }

      cValues := e.Entity().CustomValue
      for i, cValue := range cValues {
        if cValue.(*types.CustomFieldStringValue).Key == field.Key {
          entity.CustomValue = append(cValues[:i], cValues[i+1:]...)
          break
        }
      }
    })
  }
}
```

```golang
// object\example_test.go

    // filter used to find objects with "backup=true"
    filter := property.Filter{"customValue": &types.CustomFieldStringValue{
      CustomFieldValue: types.CustomFieldValue{Key: field.Key},
      Value:            "true",
    }}

    var objs []mo.ManagedEntity
    err = v.RetrieveWithFilter(ctx, any, []string{"name", "customValue"}, &objs, filter)
    if err != nil {
      return err
    }
```

This `object\example_test.go` test shows the struct hierarchy used to set the
`Value` field.

```golang
// simulator\object_test.go

  if vm.CustomValue[0].(*types.CustomFieldStringValue).Key != field.Key {
    t.Fatalf("vm.CustomValue[0].Key expected %d, got %d",
      field.Key, vm.CustomValue[0].(*types.CustomFieldStringValue).Key)
  }
  if vm.CustomValue[0].(*types.CustomFieldStringValue).Value != fieldValue {
    t.Fatalf("vm.CustomValue[0].Value expected %s, got %s",
      fieldValue, vm.CustomValue[0].(*types.CustomFieldStringValue).Value)
  }

  if vm.Value[0].(*types.CustomFieldStringValue).Key != field.Key {
    t.Fatalf("vm.Value[0].Key expected %d, got %d",
      field.Key, vm.Value[0].(*types.CustomFieldStringValue).Key)
  }
  if vm.Value[0].(*types.CustomFieldStringValue).Value != fieldValue {
    t.Fatalf("vm.Value[0].Value expected %s, got %s",
      fieldValue, vm.Value[0].(*types.CustomFieldStringValue).Value)
  }
```

I've yet to find here or elsewhere any "guards" against potential interface
conversion failure, so obviously there is a very high degree of confidence
that converting a `*types.BaseCustomFieldValue` to a
`*types.CustomFieldStringValue` will succeed.

One last example, saying pretty much the same thing:

```golang
// simulator\custom_fields_manager_test.go

  values := vm.Entity().CustomValue
  if len(values) != 1 {
    t.Fatalf("expect CustomValue has 1 item; got %d", len(values))
  }
  fkey := values[0].GetCustomFieldValue().Key
  if fkey != field.Key {
    t.Fatalf("expect value.Key == field.Key; got %d != %d", fkey, field.Key)
  }
  value := values[0].(*types.CustomFieldStringValue).Value
  if value != "value" {
    t.Fatalf("expect value.Value to be %q; got %q", "value", value)
  }
```

From The Go Programming Language book:

> Second, if instead the asserted type T is an interface type, then the type
> assertion checks whether x 's dynamic type satisfies T . If this check
> succeeds, the dynamic value is not extracted; the result is still an
> interface value with the same type and value components, but the result has
> the interface type T . In other words, a type assertion to an interface type
> changes the type of the expression, making a different (and usually larger)
> set of methods accessible, but it preserves the dynamic type and value
> components inside the interface value.
