package enum

import "github.com/guregu/null/v6"

// Gender is type for gender enum.
type Gender int16

// Gender enum.
const (
	MALE   Gender = 1
	FEMALE Gender = 2
)

// ToNullInt16 converts Gender to nullable int16 type.
func (g *Gender) ToNullInt16() null.Int16 {
	if g == nil {
		return null.NewInt16(0, false)
	}
	return null.Int16From(int16(*g))
}

// GenderFromNullInt16 converts null.Int16 type to pointer of Gender.
func GenderFromNullInt16(value null.Int16) *Gender {
	genderNumeric := value.Ptr()
	var gender *Gender
	if genderNumeric != nil {
		genderValue := Gender(*genderNumeric)
		gender = &genderValue
	}

	return gender
}
