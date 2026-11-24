package datatypes

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strconv"
	"strings"
)

// DecimalMaxDigits is the most digits a Decimal holds, before and after the point together. Eighteen
// digits fit a 64-bit coefficient, mongo's Decimal128 and SQL NUMERIC(18,s) alike, which is what lets
// every storage provider keep the value exactly.
const DecimalMaxDigits = 18

// DecimalMaxScale is the most digits a Decimal holds after the point.
const DecimalMaxScale = 18

// DecimalTag is the struct tag key that carries a decimal field's declared scale, written
// `decimal:"scale=2"`. The build emits it on every generated decimal field, and the storage
// providers that convert records by reflection read it with DecimalScaleOf.
const DecimalTag = "decimal"

// decimalMaxCoefficient is the largest coefficient magnitude: eighteen nines.
const decimalMaxCoefficient int64 = 999_999_999_999_999_999

var (
	// ErrDecimalSyntax reports text, or a value, that is not a decimal number.
	ErrDecimalSyntax = errors.New("not a decimal number")
	// ErrDecimalRange reports a value with more than DecimalMaxDigits digits.
	ErrDecimalRange = errors.New("a decimal holds at most 18 digits")
	// ErrDecimalScale reports a scale outside 0 to DecimalMaxScale.
	ErrDecimalScale = errors.New("a decimal's scale is 0 to 18")
	// ErrDecimalPrecision reports a value with more digits after the point than the scale it is
	// asked to fit, where fitting it would drop a digit.
	ErrDecimalPrecision = errors.New("the decimal has more digits after the point than the scale allows")
	// ErrDecimalDivisionByZero reports a division by a zero decimal.
	ErrDecimalDivisionByZero = errors.New("decimal division by zero")
	// ErrDecimalRoundingMode reports a rounding mode that is not one of the declared modes.
	ErrDecimalRoundingMode = errors.New("unknown decimal rounding mode")
)

// RoundingMode says how an operation that must drop digits rounds the result.
type RoundingMode int

const (
	// RoundHalfEven rounds to the nearest value, and a tie to the even neighbour (banker's
	// rounding): 2.345 to two places is 2.34, 2.355 is 2.36.
	RoundHalfEven RoundingMode = iota
	// RoundHalfUp rounds to the nearest value, and a tie away from zero: 2.345 to two places is
	// 2.35, -2.345 is -2.35.
	RoundHalfUp
	// RoundDown drops the extra digits, rounding toward zero: 2.349 to two places is 2.34.
	RoundDown
)

// Decimal is an exact decimal number, such as an amount of money, a quantity or a rate. It holds at
// most DecimalMaxDigits digits and at most DecimalMaxScale of them after the point, and it never
// passes through binary floating point.
//
// The zero value is 0. A Decimal is a value: copy it freely.
//
// # Equality
//
// Values that differ only in trailing zeros are equal numbers but different representations: 1.5
// and 1.50 Cmp as 0, while `==` on the two structs is false. Compare numbers with Cmp. The build
// pads every stored value to its field's scale, so two values of one field compare equal with `==`
// as well.
//
// # Text form
//
// Text is an optional sign, digits, and optionally a point followed by digits: "1234.50", "-0.05",
// "7". No exponent, no grouping, no spaces. String writes the value at its own scale, so the number
// of digits after the point survives a round trip.
type Decimal struct {
	coef  int64
	scale uint8
}

// ParseDecimal reads a decimal from its text form. It refuses text that is not a decimal number,
// and a value with more than DecimalMaxDigits digits; the value keeps as many digits after the point
// as the text has.
func ParseDecimal(s string) (Decimal, error) {
	t := s
	neg := false
	if t != "" && (t[0] == '+' || t[0] == '-') {
		neg = t[0] == '-'
		t = t[1:]
	}
	// split at the point; the integer part is required, the fraction only when a point is written
	intPart, frac, hasPoint := strings.Cut(t, ".")
	if intPart == "" || (hasPoint && frac == "") || !decimalDigitsOnly(intPart) || !decimalDigitsOnly(frac) {
		return Decimal{}, fmt.Errorf("%w: %q", ErrDecimalSyntax, s)
	}
	// leading zeros of the integer part are not digits of the value
	significant := strings.TrimLeft(intPart, "0")
	if len(frac) > DecimalMaxScale || len(significant)+len(frac) > DecimalMaxDigits {
		return Decimal{}, fmt.Errorf("%w: %q", ErrDecimalRange, s)
	}
	var coef int64
	if digits := significant + frac; digits != "" {
		// at most eighteen digits, so the value fits an int64 and ParseInt cannot overflow
		parsed, err := strconv.ParseInt(digits, 10, 64)
		if err != nil {
			return Decimal{}, fmt.Errorf("%w: %q", ErrDecimalSyntax, s)
		}
		coef = parsed
	}
	if neg {
		coef = -coef
	}
	return Decimal{coef: coef, scale: uint8(len(frac))}, nil
}

// MustParseDecimal is ParseDecimal for text known to be valid, such as a default the build has
// already checked. It panics on text ParseDecimal refuses.
func MustParseDecimal(s string) Decimal {
	d, err := ParseDecimal(s)
	if err != nil {
		panic(err)
	}
	return d
}

// DecimalFromInt answers a whole number as a Decimal with scale 0, refusing one with more than
// DecimalMaxDigits digits.
func DecimalFromInt(n int64) (Decimal, error) {
	if n > decimalMaxCoefficient || n < -decimalMaxCoefficient {
		return Decimal{}, fmt.Errorf("%w: %d", ErrDecimalRange, n)
	}
	return Decimal{coef: n}, nil
}

// DecimalFromFloat64 answers a float as the Decimal its shortest decimal representation states:
// 0.1 is 0.1, not the binary value nearest it. That is exact for any value written with up to 15
// significant digits, which is every float a person typed. NaN and the infinities are refused, and
// so is a float whose representation needs more than DecimalMaxDigits digits.
func DecimalFromFloat64(f float64) (Decimal, error) {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return Decimal{}, fmt.Errorf("%w: %v", ErrDecimalSyntax, f)
	}
	// 'f' with precision -1 is the fewest digits that read back as the same float, with no exponent
	return ParseDecimal(strconv.FormatFloat(f, 'f', -1, 64))
}

// DecimalFromUnits answers a count of the smallest unit at a scale as a Decimal: 123450 units at
// scale 2 is 1234.50. It is how a store holding the scaled whole number reads it back.
func DecimalFromUnits(units int64, scale int) (Decimal, error) {
	if scale < 0 || scale > DecimalMaxScale {
		return Decimal{}, fmt.Errorf("%w: %d", ErrDecimalScale, scale)
	}
	if units > decimalMaxCoefficient || units < -decimalMaxCoefficient {
		return Decimal{}, fmt.Errorf("%w: %d units", ErrDecimalRange, units)
	}
	return Decimal{coef: units, scale: uint8(scale)}, nil
}

// DecimalOf answers a value of any number type, or text holding a decimal, as a Decimal: a Decimal
// or a pointer to one as it is, text through ParseDecimal, integers exactly, and floats through
// their shortest representation (DecimalFromFloat64). It is the one conversion every caller holding
// an untyped value uses -- a filter operand, a value in an update patch, an argument to an
// expression function.
func DecimalOf(v interface{}) (Decimal, error) {
	switch n := v.(type) {
	case Decimal:
		return n, nil
	case *Decimal:
		if n == nil {
			return Decimal{}, fmt.Errorf("%w: a nil decimal", ErrDecimalSyntax)
		}
		return *n, nil
	case string:
		return ParseDecimal(n)
	case []byte:
		return ParseDecimal(string(n))
	case int:
		return DecimalFromInt(int64(n))
	case int8:
		return DecimalFromInt(int64(n))
	case int16:
		return DecimalFromInt(int64(n))
	case int32:
		return DecimalFromInt(int64(n))
	case int64:
		return DecimalFromInt(n)
	case uint:
		return decimalFromUint(uint64(n))
	case uint8:
		return decimalFromUint(uint64(n))
	case uint16:
		return decimalFromUint(uint64(n))
	case uint32:
		return decimalFromUint(uint64(n))
	case uint64:
		return decimalFromUint(n)
	case float32:
		// a float32's own shortest text, so 0.1 typed as a float32 is 0.1
		return ParseDecimal(strconv.FormatFloat(float64(n), 'f', -1, 32))
	case float64:
		return DecimalFromFloat64(n)
	}
	// a named type over a number or text, such as encoding/json's Number, by its kind
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.String:
		return ParseDecimal(rv.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return DecimalFromInt(rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return decimalFromUint(rv.Uint())
	case reflect.Float32:
		return ParseDecimal(strconv.FormatFloat(rv.Float(), 'f', -1, 32))
	case reflect.Float64:
		return DecimalFromFloat64(rv.Float())
	}
	return Decimal{}, fmt.Errorf("%w: %v (%T)", ErrDecimalSyntax, v, v)
}

// decimalFromUint answers an unsigned whole number as a Decimal, refusing one beyond the digit limit.
func decimalFromUint(n uint64) (Decimal, error) {
	if n > uint64(decimalMaxCoefficient) {
		return Decimal{}, fmt.Errorf("%w: %d", ErrDecimalRange, n)
	}
	return Decimal{coef: int64(n)}, nil
}

// Scale answers how many digits the value holds after the point.
func (d Decimal) Scale() int {
	return int(d.scale)
}

// Sign answers -1, 0 or +1 as the value is negative, zero or positive.
func (d Decimal) Sign() int {
	switch {
	case d.coef < 0:
		return -1
	case d.coef > 0:
		return 1
	}
	return 0
}

// IsZero answers whether the value is zero, at any scale.
func (d Decimal) IsZero() bool {
	return d.coef == 0
}

// Neg answers the value with its sign reversed. It cannot overflow: the digit limit is symmetric.
func (d Decimal) Neg() Decimal {
	return Decimal{coef: -d.coef, scale: d.scale}
}

// Abs answers the value without its sign.
func (d Decimal) Abs() Decimal {
	if d.coef < 0 {
		return d.Neg()
	}
	return d
}

// Cmp compares the numbers: -1 when d is less than o, 0 when they are equal, +1 when d is greater.
// Trailing zeros do not matter, so 1.5 and 1.50 compare equal.
func (d Decimal) Cmp(o Decimal) int {
	if d.scale == o.scale {
		switch {
		case d.coef < o.coef:
			return -1
		case d.coef > o.coef:
			return 1
		}
		return 0
	}
	s := max(int(d.scale), int(o.scale))
	// aligning the scales can exceed an int64, so compare the aligned coefficients as big integers
	return decimalAligned(d, s).Cmp(decimalAligned(o, s))
}

// String writes the value in its text form at its own scale: 1234.50 stays "1234.50".
func (d Decimal) String() string {
	return d.text(int(d.scale))
}

// StringPadded writes the value with at least scale digits after the point, padding with zeros and
// never dropping a digit: 1234.5 padded to 2 is "1234.50", and 1.234 padded to 2 stays "1.234". It
// is how a field's value is written at the field's declared scale.
func (d Decimal) StringPadded(scale int) string {
	return d.text(max(scale, int(d.scale)))
}

// text writes the value with exactly digits places after the point, which is at least its scale.
func (d Decimal) text(places int) string {
	neg := d.coef < 0
	u := d.coef
	if neg {
		u = -u
	}
	// padding to more places multiplies the coefficient by a power of ten, written as appended zeros
	digits := strconv.FormatInt(u, 10) + strings.Repeat("0", places-int(d.scale))
	if places > 0 {
		if len(digits) <= places {
			digits = strings.Repeat("0", places-len(digits)+1) + digits
		}
		digits = digits[:len(digits)-places] + "." + digits[len(digits)-places:]
	}
	if neg {
		return "-" + digits
	}
	return digits
}

// Float64 answers the float nearest the value. It is for handing a decimal to code that only knows
// floats; the answer is no longer exact.
func (d Decimal) Float64() float64 {
	// ParseFloat rounds correctly, so this is the nearest float to the exact value
	f, _ := strconv.ParseFloat(d.String(), 64)
	return f
}

// UnitsAt answers the value as a count of the smallest unit at a scale: 1234.5 at scale 2 is 123450.
// It refuses a value with more digits after the point than the scale (ErrDecimalPrecision) and a
// count beyond the digit limit (ErrDecimalRange). It is how a store holding the scaled whole number
// writes it.
func (d Decimal) UnitsAt(scale int) (int64, error) {
	r, err := d.Rescale(scale)
	if err != nil {
		return 0, err
	}
	return r.coef, nil
}

// Rescale answers the same number at another scale, padding with zeros or dropping only zeros. It
// refuses to drop a digit that is not zero (ErrDecimalPrecision) -- use Round for that -- and a
// padding that exceeds the digit limit (ErrDecimalRange).
func (d Decimal) Rescale(scale int) (Decimal, error) {
	if scale < 0 || scale > DecimalMaxScale {
		return Decimal{}, fmt.Errorf("%w: %d", ErrDecimalScale, scale)
	}
	if scale < int(d.scale) {
		div := decimalPow10(int(d.scale) - scale)
		// a remainder means a non-zero digit would be lost
		q, r := new(big.Int).QuoRem(big.NewInt(d.coef), div, new(big.Int))
		if r.Sign() != 0 {
			return Decimal{}, fmt.Errorf("%w: %s at scale %d", ErrDecimalPrecision, d.String(), scale)
		}
		return decimalFromBig(q, scale)
	}
	return decimalFromBig(decimalAligned(d, scale), scale)
}

// Round answers the value at a scale, rounding with mode when digits must be dropped and padding
// with zeros when the scale is larger. It refuses an unknown mode and a result beyond the digit
// limit.
func (d Decimal) Round(scale int, mode RoundingMode) (Decimal, error) {
	if scale < 0 || scale > DecimalMaxScale {
		return Decimal{}, fmt.Errorf("%w: %d", ErrDecimalScale, scale)
	}
	if scale >= int(d.scale) {
		return decimalFromBig(decimalAligned(d, scale), scale)
	}
	q, err := decimalDivRound(big.NewInt(d.coef), decimalPow10(int(d.scale)-scale), mode)
	if err != nil {
		return Decimal{}, err
	}
	return decimalFromBig(q, scale)
}

// Add answers d + o exactly, at the larger of their scales. It refuses a sum beyond the digit limit.
func (d Decimal) Add(o Decimal) (Decimal, error) {
	s := max(int(d.scale), int(o.scale))
	return decimalFromBig(new(big.Int).Add(decimalAligned(d, s), decimalAligned(o, s)), s)
}

// Sub answers d - o exactly, at the larger of their scales. It refuses a difference beyond the digit
// limit.
func (d Decimal) Sub(o Decimal) (Decimal, error) {
	s := max(int(d.scale), int(o.scale))
	return decimalFromBig(new(big.Int).Sub(decimalAligned(d, s), decimalAligned(o, s)), s)
}

// Mul answers d × o rounded to scale with mode. The product is computed exactly first, so rounding
// happens once. It refuses an unknown mode and a result beyond the digit limit.
func (d Decimal) Mul(o Decimal, scale int, mode RoundingMode) (Decimal, error) {
	if scale < 0 || scale > DecimalMaxScale {
		return Decimal{}, fmt.Errorf("%w: %d", ErrDecimalScale, scale)
	}
	product := new(big.Int).Mul(big.NewInt(d.coef), big.NewInt(o.coef))
	exact := int(d.scale) + int(o.scale)
	if scale >= exact {
		return decimalFromBig(product.Mul(product, decimalPow10(scale-exact)), scale)
	}
	q, err := decimalDivRound(product, decimalPow10(exact-scale), mode)
	if err != nil {
		return Decimal{}, err
	}
	return decimalFromBig(q, scale)
}

// Div answers d ÷ o rounded to scale with mode. It refuses a zero divisor, an unknown mode and a
// result beyond the digit limit.
func (d Decimal) Div(o Decimal, scale int, mode RoundingMode) (Decimal, error) {
	if scale < 0 || scale > DecimalMaxScale {
		return Decimal{}, fmt.Errorf("%w: %d", ErrDecimalScale, scale)
	}
	if o.coef == 0 {
		return Decimal{}, ErrDecimalDivisionByZero
	}
	// (d.coef / 10^ds) / (o.coef / 10^os) at scale s is d.coef·10^(os+s) / (o.coef·10^ds)
	num := new(big.Int).Mul(big.NewInt(d.coef), decimalPow10(int(o.scale)+scale))
	den := new(big.Int).Mul(big.NewInt(o.coef), decimalPow10(int(d.scale)))
	q, err := decimalDivRound(num, den, mode)
	if err != nil {
		return Decimal{}, err
	}
	return decimalFromBig(q, scale)
}

// MarshalJSON writes the value as a JSON string, "1234.50", so no JSON reader parses it as a
// binary float.
func (d Decimal) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(d.String())), nil
}

// UnmarshalJSON reads a JSON string holding a decimal, or a JSON number: number text is read
// exactly, and one written with an exponent through its float. JSON null leaves the value
// unchanged.
func (d *Decimal) UnmarshalJSON(b []byte) error {
	s := strings.TrimSpace(string(b))
	if s == "null" {
		return nil
	}
	if strings.HasPrefix(s, `"`) {
		// a JSON string; Unquote handles the escapes JSON allows
		text, err := strconv.Unquote(s)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrDecimalSyntax, s)
		}
		v, err := ParseDecimal(text)
		if err != nil {
			return err
		}
		*d = v
		return nil
	}
	v, err := ParseDecimal(s)
	if errors.Is(err, ErrDecimalSyntax) {
		// a JSON number with an exponent; ParseFloat reads it and its shortest form is the value
		f, ferr := strconv.ParseFloat(s, 64)
		if ferr != nil {
			return err
		}
		v, err = DecimalFromFloat64(f)
	}
	if err != nil {
		return err
	}
	*d = v
	return nil
}

// MarshalText writes the value's text form.
func (d Decimal) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// UnmarshalText reads the value's text form.
func (d *Decimal) UnmarshalText(b []byte) error {
	v, err := ParseDecimal(string(b))
	if err != nil {
		return err
	}
	*d = v
	return nil
}

// MarshalBinary writes the value's text form as bytes. Binary encoders, such as the msgpack codec
// the KV stores use, store the value through it.
func (d Decimal) MarshalBinary() ([]byte, error) {
	return d.MarshalText()
}

// UnmarshalBinary reads what MarshalBinary wrote.
func (d *Decimal) UnmarshalBinary(b []byte) error {
	return d.UnmarshalText(b)
}

// Value hands the value to a SQL driver as its text form, which a NUMERIC column reads exactly.
func (d Decimal) Value() (driver.Value, error) {
	return d.String(), nil
}

// Scan reads a value a SQL driver hands back: text or bytes (a NUMERIC column), a whole number, or a
// float through its shortest form. A NULL reads as zero.
func (d *Decimal) Scan(src interface{}) error {
	if src == nil {
		*d = Decimal{}
		return nil
	}
	v, err := DecimalOf(src)
	if err != nil {
		return err
	}
	*d = v
	return nil
}

// Compare compares the value with another of any type DecimalOf accepts: -1, 0 or +1 as the value is
// less than, equal to or greater than it. It is the comparer the KV stores' query engines (bolthold,
// badgerhold) call when they filter and sort a field, so a decimal stored as text still orders as a
// number.
func (d Decimal) Compare(other interface{}) (int, error) {
	o, err := DecimalOf(other)
	if err != nil {
		return 0, err
	}
	return d.Cmp(o), nil
}

// DecimalScaleOf reads a struct field's declared scale from its DecimalTag, `decimal:"scale=2"`. It
// answers false with no error when the field carries no such tag, and an error when the tag is
// malformed or the scale is outside 0 to DecimalMaxScale.
func DecimalScaleOf(field reflect.StructField) (int, bool, error) {
	// Lookup tells an absent tag from an empty one
	tag, ok := field.Tag.Lookup(DecimalTag)
	if !ok {
		return 0, false, nil
	}
	for _, part := range strings.Split(tag, ",") {
		key, val, found := strings.Cut(strings.TrimSpace(part), "=")
		if !found || key != "scale" {
			continue
		}
		scale, err := strconv.Atoi(val)
		if err != nil || scale < 0 || scale > DecimalMaxScale {
			return 0, true, fmt.Errorf("%w: field %s tag %q", ErrDecimalScale, field.Name, tag)
		}
		return scale, true, nil
	}
	return 0, true, fmt.Errorf("%w: field %s tag %q names no scale", ErrDecimalScale, field.Name, tag)
}

// decimalDigitsOnly answers whether s holds only the digits 0-9. The empty string answers true.
func decimalDigitsOnly(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

// decimalPow10 answers 10^n as a big integer.
func decimalPow10(n int) *big.Int {
	return new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(n)), nil)
}

// decimalAligned answers d's coefficient at a scale at least its own, as a big integer.
func decimalAligned(d Decimal, scale int) *big.Int {
	c := big.NewInt(d.coef)
	if scale > int(d.scale) {
		c.Mul(c, decimalPow10(scale-int(d.scale)))
	}
	return c
}

// decimalFromBig answers a coefficient and scale as a Decimal, refusing a coefficient beyond the
// digit limit.
func decimalFromBig(c *big.Int, scale int) (Decimal, error) {
	if !c.IsInt64() || c.Int64() > decimalMaxCoefficient || c.Int64() < -decimalMaxCoefficient {
		return Decimal{}, fmt.Errorf("%w: %s at scale %d", ErrDecimalRange, c.String(), scale)
	}
	return Decimal{coef: c.Int64(), scale: uint8(scale)}, nil
}

// decimalDivRound answers num ÷ den as a whole number rounded with mode. den is never zero here.
func decimalDivRound(num, den *big.Int, mode RoundingMode) (*big.Int, error) {
	if mode != RoundHalfEven && mode != RoundHalfUp && mode != RoundDown {
		return nil, fmt.Errorf("%w: %d", ErrDecimalRoundingMode, mode)
	}
	neg := (num.Sign() < 0) != (den.Sign() < 0)
	n := new(big.Int).Abs(num)
	dd := new(big.Int).Abs(den)
	// QuoRem truncates toward zero; on magnitudes that is the rounded-down quotient
	q, r := new(big.Int).QuoRem(n, dd, new(big.Int))
	if r.Sign() != 0 && mode != RoundDown {
		// twice the remainder against the divisor says whether the dropped part is below, at or
		// above one half
		half := new(big.Int).Lsh(r, 1).Cmp(dd)
		if half > 0 || (half == 0 && (mode == RoundHalfUp || q.Bit(0) == 1)) {
			q.Add(q, big.NewInt(1))
		}
	}
	if neg {
		q.Neg(q)
	}
	return q, nil
}
