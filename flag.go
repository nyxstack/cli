package cli

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// FlagSet manages command flags
type FlagSet struct {
	flags []*Flag // Array storage for flags (pointers to preserve modifications)
}

// Flag represents a command flag
type Flag struct {
	names    []string      // All names: ["port", "p"] - first is primary, rest are aliases
	flagType reflect.Type  // Go type of the flag
	defValue interface{}   // Default value
	usage    string        // Help text
	value    reflect.Value // Pointer to actual variable
	required bool          // Whether flag is required (future)
	hidden   bool          // Whether to hide from help (future)
	set      bool          // Whether flag was actually set by user
}

// Getter methods
func (f *Flag) GetNames() []string {
	return f.names
}

func (f *Flag) GetType() string {
	if f.flagType == nil {
		return "string"
	}
	switch f.flagType.Kind() {
	case reflect.Bool:
		return "bool"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if f.flagType == reflect.TypeOf(time.Duration(0)) {
			return "duration"
		}
		return "int"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "uint"
	case reflect.Float32, reflect.Float64:
		return "float"
	case reflect.String:
		return "string"
	case reflect.Slice:
		return "array"
	default:
		return "string"
	}
}

func (f *Flag) GetDefault() interface{} {
	return f.defValue
}

func (f *Flag) GetUsage() string {
	return f.usage
}

func (f *Flag) GetValue() interface{} {
	if f.value.IsValid() && f.value.CanInterface() {
		// Dereference pointer to get actual value
		if f.value.Kind() == reflect.Ptr && !f.value.IsNil() {
			return f.value.Elem().Interface()
		}
		return f.value.Interface()
	}
	return nil
}

func (f *Flag) IsRequired() bool {
	return f.required
}

func (f *Flag) IsHidden() bool {
	return f.hidden
}

func (f *Flag) IsSet() bool {
	return f.set
}

// Helper methods
func (f *Flag) PrimaryName() string {
	if len(f.names) > 0 {
		return f.names[0]
	}
	return ""
}

func (f *Flag) ShortName() string {
	if len(f.names) > 1 {
		return f.names[1]
	}
	return ""
}

func (f *Flag) HasName(name string) bool {
	for _, n := range f.names {
		if n == name {
			return true
		}
	}
	return false
}

// Setter for value (internal use)
func (f *Flag) setValue(val reflect.Value) {
	f.value = val
}

// NewFlagSet creates a new flag set
func NewFlagSet() *FlagSet {
	return &FlagSet{
		flags: []*Flag{},
	}
}

// Add adds a flag to the flag set
func (fs *FlagSet) Add(ptr interface{}, name, shorthand string, defaultValue interface{}, usage string) {
	flagType, err := inferType(ptr)
	if err != nil {
		panic(fmt.Sprintf("failed to infer flag type for %s: %v", name, err))
	}

	flagValue := reflect.ValueOf(ptr).Elem()

	// Set default value if provided
	if defaultValue != nil {
		defaultVal := reflect.ValueOf(defaultValue)
		if defaultVal.Type().ConvertibleTo(flagType) {
			flagValue.Set(defaultVal.Convert(flagType))
		}
	}

	// Build names array: [primary, short] if shorthand provided
	names := []string{name}
	if shorthand != "" {
		names = append(names, shorthand)
	}

	flag := Flag{
		names:    names,
		flagType: flagType,
		defValue: defaultValue,
		usage:    usage,
		value:    flagValue,
		required: false,
		hidden:   false,
	}

	fs.flags = append(fs.flags, &flag)
}

// Get returns a flag by name
func (fs *FlagSet) GetFlag(name string) *Flag {
	for _, flag := range fs.flags {
		if flag.HasName(name) {
			return flag
		}
	}
	return nil
}

// GetFlags returns all flags
func (fs *FlagSet) GetFlags() []*Flag {
	return fs.flags
}

// GetAll is an alias for GetFlags for backward compatibility
func (fs *FlagSet) GetAll() []Flag {
	result := make([]Flag, len(fs.flags))
	for i, flag := range fs.flags {
		result[i] = *flag
	}
	return result
}

// Parse parses command line arguments and sets flag values
func (fs *FlagSet) Parse(args []string) ([]string, error) {
	remaining := make([]string, 0)

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if !strings.HasPrefix(arg, "-") {
			remaining = append(remaining, arg)
			continue
		}

		var flagName string
		var flagValue string
		var hasValue bool

		if strings.HasPrefix(arg, "--") {
			// Long flag: --flag=value or --flag (for booleans)
			name := arg[2:]
			if idx := strings.Index(name, "="); idx >= 0 {
				flagName = name[:idx]
				flagValue = name[idx+1:]
				hasValue = true
			} else {
				flagName = name
			}
		} else {
			// Short flag: -f=value or -f (for booleans)
			name := arg[1:]
			if idx := strings.Index(name, "="); idx >= 0 {
				flagName = name[:idx]
				flagValue = name[idx+1:]
				hasValue = true
			} else {
				flagName = name
			}
		}

		flag := fs.GetFlag(flagName)
		if flag == nil {
			return nil, fmt.Errorf("unknown flag: %s", flagName)
		}

		// Handle boolean flags
		if flag.flagType.Kind() == reflect.Bool {
			if hasValue {
				// Parse boolean value: --flag=true, --flag=1, etc.
				if err := fs.setValue(flag, flagValue); err != nil {
					return nil, fmt.Errorf("invalid value %q for flag %s: %v", flagValue, flagName, err)
				}
			} else {
				// Standalone boolean flag means true
				flag.value.SetBool(true)
			}
			flag.set = true
			continue
		}

		// Non-boolean flags MUST have a value with =
		if !hasValue {
			return nil, fmt.Errorf("flag %s requires a value (use --flag=value)", flagName)
		}

		// Parse and set the value
		if err := fs.setValue(flag, flagValue); err != nil {
			return nil, fmt.Errorf("invalid value %q for flag %s: %v", flagValue, flagName, err)
		}
		flag.set = true
	}

	return remaining, nil
}

// setValue parses a string value and sets it on the flag
func (fs *FlagSet) setValue(flag *Flag, value string) error {
	switch flag.flagType.Kind() {
	case reflect.String:
		flag.value.SetString(value)
	case reflect.Bool:
		if val, err := strconv.ParseBool(value); err != nil {
			return err
		} else {
			flag.value.SetBool(val)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if flag.flagType == reflect.TypeOf(time.Duration(0)) {
			if dur, err := time.ParseDuration(value); err != nil {
				return err
			} else {
				flag.value.SetInt(int64(dur))
			}
		} else {
			if val, err := strconv.ParseInt(value, 10, flag.flagType.Bits()); err != nil {
				return err
			} else {
				flag.value.SetInt(val)
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if val, err := strconv.ParseUint(value, 10, flag.flagType.Bits()); err != nil {
			return err
		} else {
			flag.value.SetUint(val)
		}
	case reflect.Float32, reflect.Float64:
		if val, err := strconv.ParseFloat(value, flag.flagType.Bits()); err != nil {
			return err
		} else {
			flag.value.SetFloat(val)
		}
	case reflect.Slice:
		if flag.flagType.Elem().Kind() == reflect.String {
			// Handle string slices
			currentSlice := flag.value
			newSlice := reflect.Append(currentSlice, reflect.ValueOf(value))
			flag.value.Set(newSlice)
		} else {
			return fmt.Errorf("unsupported slice type: %v", flag.flagType.Elem().Kind())
		}
	default:
		// Try to handle custom types that implement flag.Value interface
		if flag.value.Addr().Type().Implements(reflect.TypeOf((*interface{ Set(string) error })(nil)).Elem()) {
			method := flag.value.Addr().MethodByName("Set")
			if method.IsValid() {
				result := method.Call([]reflect.Value{reflect.ValueOf(value)})
				if len(result) > 0 && !result[0].IsNil() {
					return result[0].Interface().(error)
				}
			}
		} else {
			return fmt.Errorf("unsupported flag type: %v", flag.flagType.Kind())
		}
	}
	return nil
}

// BindStruct binds struct fields as flags using struct tags
func (fs *FlagSet) BindStruct(structPtr interface{}) {
	v := reflect.ValueOf(structPtr)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		panic("BindStruct requires a pointer to a struct")
	}

	structValue := v.Elem()
	structType := structValue.Type()

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := structValue.Field(i)

		// Skip unexported fields
		if !fieldValue.CanSet() {
			continue
		}

		// Get the field's address for the flag binding
		fieldPtr := fieldValue.Addr().Interface()

		// Parse struct tag
		tag := field.Tag.Get("cli")
		if tag == "" {
			continue
		}

		// Parse tag format: "name,shorthand"
		parts := strings.Split(tag, ",")
		name := strings.TrimSpace(parts[0])
		shorthand := ""
		if len(parts) > 1 {
			shorthand = strings.TrimSpace(parts[1])
		}

		// Get usage and default from tags
		usage := field.Tag.Get("usage")
		defaultTag := field.Tag.Get("default")

		var defaultValue interface{}
		if defaultTag != "" {
			defaultValue = parseDefaultValue(defaultTag, field.Type)
		}

		fs.Add(fieldPtr, name, shorthand, defaultValue, usage)
	}
}

// parseDefaultValue parses a default value string to the appropriate type
func parseDefaultValue(value string, targetType reflect.Type) interface{} {
	switch targetType.Kind() {
	case reflect.String:
		return value
	case reflect.Bool:
		if val, err := strconv.ParseBool(value); err == nil {
			return val
		}
		return false
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if targetType == reflect.TypeOf(time.Duration(0)) {
			if dur, err := time.ParseDuration(value); err == nil {
				return dur
			}
			return time.Duration(0)
		} else {
			if val, err := strconv.ParseInt(value, 10, 64); err == nil {
				return reflect.ValueOf(val).Convert(targetType).Interface()
			}
			return reflect.Zero(targetType).Interface()
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if val, err := strconv.ParseUint(value, 10, 64); err == nil {
			return reflect.ValueOf(val).Convert(targetType).Interface()
		}
		return reflect.Zero(targetType).Interface()
	case reflect.Float32, reflect.Float64:
		if val, err := strconv.ParseFloat(value, 64); err == nil {
			return reflect.ValueOf(val).Convert(targetType).Interface()
		}
		return reflect.Zero(targetType).Interface()
	default:
		return nil
	}
}

// inferType infers the type from a pointer using reflection
func inferType(ptr interface{}) (reflect.Type, error) {
	v := reflect.ValueOf(ptr)
	if v.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("expected pointer, got %v", v.Kind())
	}
	return v.Elem().Type(), nil
}
