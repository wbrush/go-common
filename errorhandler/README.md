# Errors System
errorhandler.Error extends basic go builtin error interface with some useful fields and functions 

You can create the new service errors within 3 ways:
1. NewError() -- creates an Error object with provided ErrorCode and optional value (which can be an user input)
2. NewErrorWithDesc() -- same as NewError() but adds description field before optional value (it is used when you want to say something more informative than only ErrorCode)
3. FromVanillaError() -- used to simple create a new filled Error from basic go error

ServiceError interface (which really extends go builtin error interface) has several functions to use:
1. ErrorCode() -- returns ErrorCode from the error itself
2. ToMap() -- used as util func, when you need to provide all Error fields as map for some reason