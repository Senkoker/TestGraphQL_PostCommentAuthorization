package runtime

import (
	"context"
	"errors"
	"github.com/99designs/gqlgen/graphql"
	"reflect"
)

func InputUnionDirective(ctx context.Context, obj any, next graphql.Resolver) (res any, err error) {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Struct {
		panic("@inputUnion directive should only be used with objects")
	}
	valueFound := false
	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).IsNil() {
			if valueFound {
				return obj, errors.New("only one field of the input union should have a value")
			}
			valueFound = true
		}
	}
	if !valueFound {
		return obj, errors.New("one of the input union fields must have a value")
	}
	return obj, nil
}
