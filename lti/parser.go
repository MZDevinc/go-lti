package lti

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/dgrijalva/jwt-go"
)

// ParseLaunchMessage marshals unstructured json claims into a LaunchMessage struct, and checks for required fields
func ParseLaunchMessage(claims jwt.MapClaims) (LaunchMessage, error) {
	msg := LaunchMessage{}
	serializedClaims, err := json.Marshal(claims)
	if err != nil {
		return msg, err
	}
	err = json.Unmarshal(serializedClaims, &msg)
	if err != nil {
		return msg, err
	}

	err = validateLaunchMessage(msg)
	if err != nil {
		return msg, err
	}

	return msg, nil
}

func validateLaunchMessage(msg LaunchMessage) error {
	v := reflect.ValueOf(msg)
	return validateRecursive(v, true, "")
}

func validateRecursive(v reflect.Value, required bool, fieldName string) error {
	exists := !(!v.IsValid() || v.IsZero())
	t := v.Type()

	if required && !exists {
		// A field must be present if it's required
		return fmt.Errorf("Field %s was missing", fieldName)
	}

	// Deference a pointer to a struct
	if t.Kind() == reflect.Ptr && exists && t.Elem().Kind() == reflect.Struct {
		v = v.Elem()
		t = t.Elem()
	}

	if t.Kind() == reflect.Struct && exists {
		for i := 0; i < t.NumField(); i++ {
			sv := v.Field(i)
			st := t.Field(i)

			required := st.Tag.Get("required") == "true"

			subFieldName := st.Name
			if fieldName != "" {
				subFieldName = fieldName + " / " + st.Name
			}

			err := validateRecursive(sv, required, subFieldName)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ParseLaunchMessage materializes json claims into a ResourceLinkMessage struct
// func ParseLaunchMessage(claims jwt.MapClaims) (LaunchMessage, error) {
// 	linkMessage := LaunchMessage{}

// 	t := reflect.TypeOf(linkMessage)
// 	ps := reflect.ValueOf(&linkMessage)
// 	s := ps.Elem()

// 	for i := 0; i < t.NumField(); i++ {
// 		f := t.Field(i)
// 		v := s.Field(i)

// 		if f.Name == "Custom" || f.Name == "Extensions" {
// 			continue
// 		}

// 		tag := f.Tag.Get("json")
// 		if tag == "" {
// 			return linkMessage, fmt.Errorf("No 'json' tag for field %s", f.Name)
// 		}

// 		switch f.Type.Kind() {
// 		case reflect.Map:
// 			subclaims := claims[tag].(map[string]interface{})
// 			m := s.Field(i)
// 			for key, value := range subclaims {
// 				m.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(value))
// 			}
// 		case reflect.Struct:
// 			if subclaims, ok := claims[tag].(map[string]interface{}); ok {
// 				err := parseStruct(v, subclaims)
// 				if err != nil {
// 					return linkMessage, err
// 				}
// 			} else {
// 				fmt.Println("COULD NOT PARSE SUBCLAIMS", f.Name, f.Type.Name(), tag, claims[tag])
// 			}
// 		case reflect.Ptr:
// 			switch f.Type.Elem().Kind() {
// 			case reflect.Struct:
// 				if subclaims, ok := claims[tag]; ok {
// 					if subclaimsMap, ok := subclaims.(map[string]interface{}); ok {
// 						fmt.Println("PARSING STRUCT", f.Name)
// 						structPointer := reflect.New(f.Type.Elem())
// 						//structDef := reflect.Zero(f.Type.Elem())
// 						err := parseStruct(structPointer.Elem(), subclaimsMap)
// 						if err != nil {
// 							return linkMessage, err
// 						}
// 						v.Set(structPointer)
// 					} else {
// 						fmt.Println("COULD NOT CONVERT TO MAP", subclaims)
// 					}
// 				} else {
// 					fmt.Println("COULD NOT FIND", tag)
// 				}
// 				//If claim wasn't present, that's fine, because this was an optional field
// 			case reflect.String:
// 				if str, ok := claims[tag]; ok {
// 					def := str.(string)
// 					v.Set(reflect.ValueOf(&def))
// 				}
// 				//If claim wasn't present, that's fine, because this was an optional field
// 			case reflect.Slice:
// 				if in, ok := claims[tag]; ok {
// 					slice := castStringSlice(in)
// 					v.Set(reflect.ValueOf(&slice))
// 				}
// 				//If claim wasn't present, that's fine, because this was an optional field
// 			default:
// 				return linkMessage, fmt.Errorf("Field %s was unexpected type %s", f.Name, f.Type.Name())
// 			}
// 		case reflect.String:
// 			if str, ok := claims[tag]; ok {
// 				v.SetString(str.(string))
// 			} else {
// 				return linkMessage, fmt.Errorf("Required field %s was not present in claims", f.Name)
// 			}
// 		case reflect.Slice:
// 			if in, ok := claims[tag]; ok {
// 				slice := castStringSlice(in)
// 				v.Set(reflect.ValueOf(slice))
// 			} else {
// 				return linkMessage, fmt.Errorf("Required field %s was not present in claims", f.Name)
// 			}
// 		}
// 	}

// 	return linkMessage, nil
// }

// func parseStruct(v reflect.Value, subclaims map[string]interface{}) error {
// 	st := v.Type()
// 	for j := 0; j < st.NumField(); j++ {
// 		sf := st.Field(j)
// 		sv := v.Field(j)
// 		fmt.Println("PARSING STRUCT FIELD", sf.Name)

// 		tag := sf.Tag.Get("json")
// 		if tag == "" {
// 			return fmt.Errorf("No 'json' tag for field %s", sf.Name)
// 		}

// 		switch sf.Type.Kind() {
// 		case reflect.Ptr:
// 			if sf.Type.Elem().Kind() == reflect.String {
// 				if str, ok := subclaims[tag]; ok {
// 					def := fmt.Sprintf("%v", str)
// 					sv.Set(reflect.ValueOf(&def))
// 					fmt.Println("set string pointer value", &def, def)
// 				}
// 			} else if sf.Type.Elem().Kind() == reflect.Bool {
// 				if val, ok := subclaims[tag]; ok {
// 					def, ok := val.(bool)
// 					if ok {
// 						sv.Set(reflect.ValueOf(&def))
// 						fmt.Println("set bool pointer value", &def, def)
// 					} else {
// 						return fmt.Errorf("Field %s expected bool", tag)
// 					}
// 				}
// 			} else {
// 				return fmt.Errorf("Field %s was not string or bool pointer", sf.Name)
// 			}
// 			//If claim wasn't present, that's fine, because this was an optional field
// 		case reflect.String:
// 			if str, ok := subclaims[tag]; ok {
// 				def := fmt.Sprintf("%v", str)
// 				sv.Set(reflect.ValueOf(def))
// 				fmt.Println("set string value", def)
// 			} else {
// 				return fmt.Errorf("Required field %s was not present in claims", sf.Name)
// 			}
// 		case reflect.Slice:
// 			if in, ok := subclaims[tag]; ok {
// 				slice := castStringSlice(in)
// 				sv.Set(reflect.ValueOf(slice))
// 			} else {
// 				return fmt.Errorf("Required field %s was not present in claims", sf.Name)
// 			}
// 		default:
// 			return fmt.Errorf("Field %s was not string or pointer (%s)", sf.Name, sf.Type.Kind().String())
// 		}
// 	}

// 	return nil
// }

// func castStringSlice(in interface{}) []string {
// 	slice := in.([]interface{})
// 	out := []string{}
// 	for _, v := range slice {
// 		str := v.(string)
// 		out = append(out, str)
// 	}

// 	return out
// }
