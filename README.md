# xerror

A simple library to allow structured contextual data to be included in error messages, including optional stack traces.
Because the xerror type implements Error() it can be used synonymously with the native error type and std lib helpers.

**For Example**
```
err := fmt.Errorf("wrapped err: %w", xerror.Errorf("not enough %d/30 in bucket", 19).AsKind("CUSTOM_ERROR"))
println(err.Error())
```
will output:

```
wrapped err: Kind: CUSTOM_ERROR | not enough 19/30 in bucket
```
