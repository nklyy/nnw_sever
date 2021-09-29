package Bitcoin

import "testing"

func TestRunTransactionV3(t *testing.T) {
	err := RunTransactionV3()
	if err != nil {
		t.Error(err)
	}
}
