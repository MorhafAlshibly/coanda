package graphqlEnums

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
)

type TournamentInterval string

const (
	TournamentIntervalDaily     TournamentInterval = "DAILY"
	TournamentIntervalWeekly    TournamentInterval = "WEEKLY"
	TournamentIntervalMonthly   TournamentInterval = "MONTHLY"
	TournamentIntervalUnlimited TournamentInterval = "UNLIMITED"
)

var AllTournamentInterval = []TournamentInterval{
	TournamentIntervalDaily,
	TournamentIntervalWeekly,
	TournamentIntervalMonthly,
	TournamentIntervalUnlimited,
}

func (e TournamentInterval) IsValid() bool {
	switch e {
	case TournamentIntervalDaily, TournamentIntervalWeekly, TournamentIntervalMonthly, TournamentIntervalUnlimited:
		return true
	}
	return false
}

func (e TournamentInterval) String() string {
	return string(e)
}

func (e *TournamentInterval) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = TournamentInterval(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid TournamentInterval", str)
	}
	return nil
}

func (e TournamentInterval) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

func (e *TournamentInterval) UnmarshalJSON(b []byte) error {
	s, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}
	return e.UnmarshalGQL(s)
}

func (e TournamentInterval) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	e.MarshalGQL(&buf)
	return buf.Bytes(), nil
}
