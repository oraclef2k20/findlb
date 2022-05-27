package util

import "testing"

func TestGetDomain(t *testing.T) {

	input := "api.test.com"
	wanteddomain := "test.com"

	wantedhost := "api.test.com"
	gotdomain, gothost := GetDomain(input)

	if wanteddomain != gotdomain {
		t.Errorf("\ninput = %v\ngot = %v\nwanted =%v\n", input, gotdomain, wanteddomain)
	}
	if wantedhost != gothost {
		t.Errorf("\ninput = %v\ngot = %v\nwanted =%v\n", input, gothost, wantedhost)
	}

	input = "http://api.test.com"
	wanteddomain = "test.com"
	wantedhost = "api.test.com"
	gotdomain, gothost = GetDomain(input)

	if wanteddomain != gotdomain {
		t.Errorf("\ninput = %v\ngot = %v\nwanted =%v\n", input, gotdomain, wanteddomain)
	}
	if wantedhost != gothost {
		t.Errorf("\ninput = %v\ngot = %v\nwanted =%v\n", input, gothost, wantedhost)
	}
}
