package nubarium_test

import (
	"testing"
	"time"

	"github.com/Idmission-LLC/nubarium-go"
	"github.com/stretchr/testify/assert"
)

func TestParseFechaF(t *testing.T) {
	tests := []struct {
		input string
		want  time.Time
		err   error
	}{
		{input: "01/01/2020", want: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local)},
		{input: "08/06/25", want: time.Date(2025, 6, 8, 0, 0, 0, 0, time.Local)},
		{input: "02/03/25", want: time.Date(2025, 3, 2, 0, 0, 0, 0, time.Local)},
		{input: "13/04/2025", want: time.Date(2025, 4, 13, 0, 0, 0, 0, time.Local)},
		{input: "19/06/25", want: time.Date(2025, 6, 19, 0, 0, 0, 0, time.Local)},
		{input: "24/05/2025", want: time.Date(2025, 5, 24, 0, 0, 0, 0, time.Local)},
		{input: "02/06/25", want: time.Date(2025, 6, 2, 0, 0, 0, 0, time.Local)},
		{input: "08/06/2025", want: time.Date(2025, 6, 8, 0, 0, 0, 0, time.Local)},
		{input: "29/05/25", want: time.Date(2025, 5, 29, 0, 0, 0, 0, time.Local)},
		{input: "22/06/25", want: time.Date(2025, 6, 22, 0, 0, 0, 0, time.Local)},
		{input: "20/04/25", want: time.Date(2025, 4, 20, 0, 0, 0, 0, time.Local)},
		{input: "06/04/25", want: time.Date(2025, 4, 6, 0, 0, 0, 0, time.Local)},
		{input: "22/05/25", want: time.Date(2025, 5, 22, 0, 0, 0, 0, time.Local)},
		{input: "07/06/25", want: time.Date(2025, 6, 7, 0, 0, 0, 0, time.Local)},
		{input: "20/03/25", want: time.Date(2025, 3, 20, 0, 0, 0, 0, time.Local)},
		{input: "10/Jun/o.20", want: time.Date(2020, 6, 10, 0, 0, 0, 0, time.Local)},
		{input: "22/05/2025", want: time.Date(2025, 5, 22, 0, 0, 0, 0, time.Local)},
		{input: "31/05/2025", want: time.Date(2025, 5, 31, 0, 0, 0, 0, time.Local)},
		{input: "15/06/25", want: time.Date(2025, 6, 15, 0, 0, 0, 0, time.Local)},
		{input: "04/jun/o.20", want: time.Date(2020, 6, 4, 0, 0, 0, 0, time.Local)},
		{input: "19/01/25", want: time.Date(2025, 1, 19, 0, 0, 0, 0, time.Local)},
		{input: "01/04/2025", want: time.Date(2025, 4, 1, 0, 0, 0, 0, time.Local)},
		{input: "20/06/2025", want: time.Date(2025, 6, 20, 0, 0, 0, 0, time.Local)},
		{input: "04/05/25", want: time.Date(2025, 5, 4, 0, 0, 0, 0, time.Local)},
		{input: "03/05/25", want: time.Date(2025, 5, 3, 0, 0, 0, 0, time.Local)},
		{input: "13/05/25", want: time.Date(2025, 5, 13, 0, 0, 0, 0, time.Local)},
		{input: "23/03/25AC", want: time.Date(2025, 3, 23, 0, 0, 0, 0, time.Local)},
		{input: "20/08/22", want: time.Date(2022, 8, 20, 0, 0, 0, 0, time.Local)},
		{input: "31/05/2025", want: time.Date(2025, 5, 31, 0, 0, 0, 0, time.Local)},
		{input: "19/05/25", want: time.Date(2025, 5, 19, 0, 0, 0, 0, time.Local)},
		{input: "03/05/2025", want: time.Date(2025, 5, 3, 0, 0, 0, 0, time.Local)},
		{input: "11/jun/25", want: time.Date(2025, 6, 11, 0, 0, 0, 0, time.Local)},
		{input: "05/06/25", want: time.Date(2025, 6, 5, 0, 0, 0, 0, time.Local)},
		{input: "07/05/25", want: time.Date(2025, 5, 7, 0, 0, 0, 0, time.Local)},
		{input: "06/jun/o.20", want: time.Date(2020, 6, 6, 0, 0, 0, 0, time.Local)},
		{input: "10/04/25", want: time.Date(2025, 4, 10, 0, 0, 0, 0, time.Local)},
		{input: "22/05/2025", want: time.Date(2025, 5, 22, 0, 0, 0, 0, time.Local)},
		{input: "02/06/25", want: time.Date(2025, 6, 2, 0, 0, 0, 0, time.Local)},
		{input: "01/05/25", want: time.Date(2025, 5, 1, 0, 0, 0, 0, time.Local)},
		{input: "21/04/25", want: time.Date(2025, 4, 21, 0, 0, 0, 0, time.Local)},
		{input: "17/04/25", want: time.Date(2025, 4, 17, 0, 0, 0, 0, time.Local)},
		{input: "21/02/2025", want: time.Date(2025, 2, 21, 0, 0, 0, 0, time.Local)},
		{input: "11/05/2025", want: time.Date(2025, 5, 11, 0, 0, 0, 0, time.Local)},
		{input: "22/03/25", want: time.Date(2025, 3, 22, 0, 0, 0, 0, time.Local)},
		{input: "01/05/2025", want: time.Date(2025, 5, 1, 0, 0, 0, 0, time.Local)},
		{input: "11/05/2025", want: time.Date(2025, 5, 11, 0, 0, 0, 0, time.Local)},
		{input: "08/06/2025", want: time.Date(2025, 6, 8, 0, 0, 0, 0, time.Local)},
		{input: "20/05/25", want: time.Date(2025, 5, 20, 0, 0, 0, 0, time.Local)},
		{input: "21/04/25", want: time.Date(2025, 4, 21, 0, 0, 0, 0, time.Local)},
		{input: "24/ene/o.20", want: time.Date(2020, 1, 24, 0, 0, 0, 0, time.Local)},
		{input: "01/05/25", want: time.Date(2025, 5, 1, 0, 0, 0, 0, time.Local)},
		{input: "20/06/25", want: time.Date(2025, 6, 20, 0, 0, 0, 0, time.Local)},
		{input: "20/06/25", want: time.Date(2025, 6, 20, 0, 0, 0, 0, time.Local)},
		{input: "24/05/25", want: time.Date(2025, 5, 24, 0, 0, 0, 0, time.Local)},
		{input: "01/06/25", want: time.Date(2025, 6, 1, 0, 0, 0, 0, time.Local)},
		{input: "25/05/25", want: time.Date(2025, 5, 25, 0, 0, 0, 0, time.Local)},
		{input: "27/04/25", want: time.Date(2025, 4, 27, 0, 0, 0, 0, time.Local)},
		{input: "01/06/25", want: time.Date(2025, 6, 1, 0, 0, 0, 0, time.Local)},
		{input: "25/05/25", want: time.Date(2025, 5, 25, 0, 0, 0, 0, time.Local)},
		{input: "26/05/2025", want: time.Date(2025, 5, 26, 0, 0, 0, 0, time.Local)},
		{input: "03/05/25", want: time.Date(2025, 5, 3, 0, 0, 0, 0, time.Local)},
		{input: "11/05/2025", want: time.Date(2025, 5, 11, 0, 0, 0, 0, time.Local)},
		{input: "08/05/25", want: time.Date(2025, 5, 8, 0, 0, 0, 0, time.Local)},
		{input: "24/05/2025", want: time.Date(2025, 5, 24, 0, 0, 0, 0, time.Local)},
		{input: "05/05/2025", want: time.Date(2025, 5, 5, 0, 0, 0, 0, time.Local)},
		{input: "06/06/25", want: time.Date(2025, 6, 6, 0, 0, 0, 0, time.Local)},
		{input: "02/05/20", want: time.Date(2020, 5, 2, 0, 0, 0, 0, time.Local)},
		{input: "04/abr/l.20", want: time.Date(2020, 4, 4, 0, 0, 0, 0, time.Local)},
		{input: "24/may/.202", want: time.Date(202, 5, 24, 0, 0, 0, 0, time.Local)},
		{input: "", want: time.Time{}, err: nubarium.ErrDateEmpty},
		{input: "//", want: time.Time{}, err: nubarium.ErrDateEmpty},
	}

	parser := nubarium.NewDateParser(nubarium.WithExpiryReferenceDate(time.Date(2025, 6, 12, 0, 0, 0, 0, time.Local)))

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			got, err := parser.Parse(test.input)
			assert.Equal(t, test.want, got)
			assert.Equal(t, test.err, err)
		})
	}
}
