package docopt_unmarshall

import (
	"errors"
	"fmt"
	"reflect"
)

func DocoptUnmarshall(arguments map[string]interface{}, options interface{}) error {
       var seen []string
       seen, err := docopt_unmarshall(arguments, options, seen)
       if err != nil { return err }
       for _, a := range seen {
               delete(arguments, a)
       }
       for leftover, _ := range arguments {
               return errors.New(fmt.Sprintf("%s is missing from options struct", leftover))
       }
       return nil
}
func docopt_unmarshall(arguments map[string]interface{}, options interface{}, seen []string) ([]string, error) {
       val := reflect.ValueOf(options).Elem()
       typ := val.Type()
       for i:=0; i<val.NumField(); i++ {
               f_val := val.Field(i)
               f_typ := typ.Field(i)
               flag := f_typ.Tag.Get("docopt")
               if flag != "" {
                       a, exist := arguments[flag]
                       if !exist {
                               return seen, errors.New(fmt.Sprintf("Struct member %s has no corresponding option %s in docopt", f_typ.Name, flag))
                       } else if a != nil {
                               a_typ := reflect.TypeOf(a)
                               if a_typ.Kind() == reflect.String && f_typ.Type.Kind() == reflect.Bool {
                                       f_val.SetBool(a != nil)
                               } else {
                                       f_val.Set(reflect.ValueOf(a))
                               }
                       }
                       seen = append(seen, flag)
               }
               if f_val.Type().Kind() == reflect.Struct {
                       if seen, err := docopt_unmarshall(arguments, f_val.Addr().Interface(), seen); err != nil {
                               return seen, err
                       }
               }
       }
       return seen, nil
}
