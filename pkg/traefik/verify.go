package traefik

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/MizuchiLabs/mantrae/pkg/util"
)

func (r *Router) Verify() error {
	if r.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if r.RouterType == "" {
		return fmt.Errorf("router type cannot be empty")
	}
	if r.Rule == "" && r.RouterType != "udp" {
		return fmt.Errorf("rule cannot be empty")
	}
	if r.Provider == "" {
		return fmt.Errorf("provider cannot be empty")
	}

	r.Rule = strings.TrimSpace(r.Rule)
	r.Name = validateName(r.Name, r.Provider)

	// Ignore for now
	domain, _ := util.ExtractDomainFromRule(r.Rule)
	if domain != "" {
		if err := util.ValidSSLCert(domain); err != nil {
			r.SSLError = err.Error()
		}
	}
	if r.Service == "" {
		r.Service = r.Name
	}
	return nil
}

func (s *Service) Verify() error {
	if s.Name == "" {
		return fmt.Errorf("service name cannot be empty")
	}
	if s.ServiceType == "" {
		return fmt.Errorf("service type cannot be empty")
	}
	if s.Provider == "" {
		return fmt.Errorf("provider cannot be empty")
	}
	if s.Provider == "http" {
		if s.LoadBalancer != nil {
			validServers := make([]Server, 0)
			for _, server := range s.LoadBalancer.Servers {
				if server.Address != "" || server.URL != "" {
					validServers = append(validServers, server)
				}
			}

			if len(validServers) == 0 {
				return fmt.Errorf("no valid servers found in load balancer")
			}

			s.LoadBalancer.Servers = validServers
		} else {
			return fmt.Errorf("load balancer cannot be nil")
		}
	}

	s.Name = validateName(s.Name, s.Provider)
	return nil
}

func (m *Middleware) Verify() error {
	if m.Name == "" {
		return fmt.Errorf("middleware name cannot be empty")
	}
	if m.Type == "" {
		return fmt.Errorf("type cannot be empty")
	}
	if m.MiddlewareType == "" {
		return fmt.Errorf("middleware type cannot be empty")
	}
	if m.Provider == "" {
		return fmt.Errorf("provider cannot be empty")
	}

	m.Name = validateName(m.Name, m.Provider)

	// Hashes the password strings in the middleware
	if m.BasicAuth != nil {
		for i, u := range m.BasicAuth.Users {
			hash, err := util.HashBasicAuth(u)
			if err != nil {
				return fmt.Errorf("error hashing password: %s", err.Error())
			}
			m.BasicAuth.Users[i] = hash
		}
	}
	if m.DigestAuth != nil {
		for i, u := range m.DigestAuth.Users {
			hash, err := util.HashBasicAuth(u)
			if err != nil {
				return fmt.Errorf("error hashing password: %s", err.Error())
			}
			m.DigestAuth.Users[i] = hash
		}
	}

	// Make sure the correct fields are set in the new middleware
	setMiddlewareByType(m)
	cleanStruct(m)
	return nil
}

func validateName(s, p string) string {
	name := strings.ToLower(s)
	parts := strings.Split(name, "@")

	if len(parts) > 1 {
		name = parts[0] + "@" + parts[1]
	} else {
		name = parts[0] + "@" + p
	}

	return name
}

// Dynamically set the correct middleware field based on the type
func setMiddlewareByType(m *Middleware) {
	newMiddleware := &Middleware{
		Name:           m.Name,
		Provider:       m.Provider,
		Type:           m.Type,
		Status:         m.Status,
		MiddlewareType: m.MiddlewareType,
	}

	// Reflect on the original middleware struct
	v := reflect.ValueOf(m)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	newV := reflect.ValueOf(newMiddleware).Elem()

	for i := 0; i < v.NumField(); i++ {
		fieldName := v.Type().Field(i).Name

		// Check if the field name matches m.Type using case-insensitive comparison
		if strings.EqualFold(fieldName, m.Type) {
			originalField := v.Field(i)
			newField := newV.FieldByName(fieldName)

			if newField.CanSet() && originalField.IsValid() {
				newField.Set(originalField) // Set the field's value
			}
			break
		}
	}

	*m = *newMiddleware
}

func cleanStruct(v interface{}) {
	val := reflect.ValueOf(v).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		if !field.CanSet() {
			continue
		}

		switch field.Kind() {
		case reflect.String:
			if field.String() == "" {
				field.Set(reflect.Zero(field.Type()))
			}
		case reflect.Map:
			cleanedMap := reflect.MakeMap(field.Type())
			for _, key := range field.MapKeys() {
				value := field.MapIndex(key)
				if value.Kind() == reflect.String && value.String() == "" {
					// Skip empty string values
					continue
				}
				// Set the non-empty entry in the new map
				cleanedMap.SetMapIndex(key, value)
			}
			field.Set(cleanedMap)
		case reflect.Slice:
			if field.Len() == 0 {
				field.Set(reflect.Zero(field.Type()))
			} else {
				cleanSlice(field) // Clean individual elements inside the slice
			}
		case reflect.Ptr:
			if !field.IsNil() {
				cleanStruct(field.Interface())
				if reflect.DeepEqual(
					field.Elem().Interface(),
					reflect.Zero(field.Elem().Type()).Interface(),
				) {
					field.Set(reflect.Zero(field.Type()))
				}
			}
		case reflect.Struct:
			cleanStruct(field.Addr().Interface())
		}
	}
}

func cleanSlice(slice reflect.Value) {
	if slice.Kind() != reflect.Slice {
		return
	}

	// Create a new slice and filter out empty or zero-value elements
	newSlice := reflect.MakeSlice(slice.Type(), 0, slice.Len())

	for i := 0; i < slice.Len(); i++ {
		elem := slice.Index(i)
		switch elem.Kind() {
		case reflect.String:
			if elem.String() != "" { // Remove empty strings
				newSlice = reflect.Append(newSlice, elem)
			}
		case reflect.Ptr, reflect.Struct:
			// Recursively clean elements inside the slice
			if elem.Kind() == reflect.Ptr && !elem.IsNil() {
				cleanStruct(elem.Interface())
			} else if elem.Kind() == reflect.Struct {
				cleanStruct(elem.Addr().Interface())
			}

			// Check if the struct/pointer is non-zero after cleaning
			if !reflect.DeepEqual(elem.Interface(), reflect.Zero(elem.Type()).Interface()) {
				newSlice = reflect.Append(newSlice, elem)
			}
		default:
			newSlice = reflect.Append(newSlice, elem)
		}
	}

	// Set the cleaned slice back
	slice.Set(newSlice)
}
