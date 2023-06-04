package entity_test

import (
	"fmt"
	"testing"

	"github.com/zhangzh-pku/software-engineering/pkg/api/entity"
)

func TestCreate(t *testing.T) {
	var p1 entity.Pin
	p1.Type = "testing"
	p1.Metadata = make(map[string]interface{})
	p1.Metadata["prop1"] = "none"
	p1_id, err := p1.Create()
	if err != nil {
		t.Fatalf(`%v`, err)
	}
	fmt.Printf("Pin1 successfully created: id=%s\n", p1_id)

	var p2 entity.Pin
	p2.Type = "testing"
	p2.Metadata = make(map[string]interface{})
	p2.Metadata["prop2"] = "noneNone"
	p2_id, err := p2.Create()
	if err != nil {
		t.Fatalf(`%v`, err)
	}
	fmt.Printf("Pin2 successfully created: id=%s\n", p2_id)

	var l entity.Link
	l.From = p1_id
	l.To = p2_id
	l.Type = "link"
	l.Metadata = make(map[string]interface{})
	l.Metadata["prop3"] = "NONE"
	l_id, err := l.Create()
	if err != nil {
		t.Fatalf(`%v`, err)
	}
	fmt.Printf("Link successfully created: id=%s\n", l_id)
}
