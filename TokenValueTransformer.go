package lexer

/*
A TokenValueTransformer is a function definition used to provide
custom token value transformation from string to a typed-interface.
*/
type TokenValueTransformer func(tokenValue string) interface{}
