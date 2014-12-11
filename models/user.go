package models

import(
    "github.com/wcl48/valval"
    "time"
    "regexp"
)

type User struct {
    Id int64
    Name string `sql:"size:255"`
    CreatedAt time.Time
    UpdatedAT time.Time
    DeletedAt time.Time
}

func UserValidate(user User)(error){
    Validator := valval.Object(valval.M{
        "Name": valval.String(
            valval.MaxLength(20),
            valval.Regexp(regexp.MustCompile(`^[a-z ]+$`)),
        ),
    })
    
    return Validator.Validate(user)
}
