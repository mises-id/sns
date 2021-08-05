package models

import (
	"context"
	"testing"
)

func TestCreateStatus(t *testing.T) {
	status, err := CreateStatus(context.TODO(), &CreateStatusParams{
		UID:     uint64(1),
		Content: "test status",
	})
	if err != nil {
		t.Error(err)
	}
	if status.Content != "test status" {
		t.Errorf("status content = %s; expected %s", status.Content, "test status")
	}
}
