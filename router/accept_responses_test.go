package router

import (
	"testing"
)

func TestAcceptResponses_MatchFor(t *testing.T) {
	accept := NewAcceptResponses("text/html,application/xml,application/*;q=0.9,*/*;q=0.8,application/json")

	if act := accept.MatchFor(ContentTypeAll); act != ContentTypeJson {
		t.Fatalf("expected %s got %s", ContentTypeJson, act)
	}

	if act := accept.MatchFor(ContentTypeAll ^ (ContentTypeJson | ContentTypeTar | ContentTypeTarGzip | ContentTypePkixCert)); act != ContentTypeText {
		t.Fatalf("expected %s got %s", ContentTypeText, act)
	}

	if act := accept.MatchFor(ContentTypeAll ^ (ContentTypeJson | ContentTypeText)); act != ContentTypeTar {
		t.Fatalf("expected %s got %s", ContentTypeTar, act)
	}

	if act := accept.MatchFor(ContentTypeAll ^ (ContentTypeJson | ContentTypeTar)); act != ContentTypeTarGzip {
		t.Fatalf("expected %s got %s", ContentTypeTarGzip, act)
	}

	if act := accept.MatchFor(ContentTypeAll ^ (ContentTypeJson | ContentTypeTar | ContentTypeTarGzip)); act != ContentTypePkixCert {
		t.Fatalf("expected %s got %s", ContentTypePkixCert, act)
	}
}

func TestAcceptResponses_MatchFor_no_match(t *testing.T) {
	accept := NewAcceptResponses("text/html")

	if act := accept.MatchFor(ContentTypeAll); act != ContentType(0) {
		t.Fatalf("expected %s got %s", ContentType(0), act)
	}

}

func TestNewAcceptResponses(t *testing.T) {

	accept := NewAcceptResponses("text/html,application/xml;q=0.9;charset='utf-8',application/*;q=0.9,*/*;q=0.8,application/json;q=1;")

	if l := accept.Len(); l != 5 {
		t.Fatalf("expected 4 items got %d", l)
	}

	accept.sort()

	if s := (*accept)[0].GetType(); s != "text/html" {
		t.Fatalf("expected text/html got %s", s)
	}

	if s := (*accept)[1].GetType(); s != "application/json" {
		t.Fatalf("expected application/json got %s", s)
	}

	if s := (*accept)[2].GetType(); s != "application/xml" {
		t.Fatalf("expected application/xml got %s", s)
	}

	if s := (*accept)[3].GetType(); s != "application/*" {
		t.Fatalf("expected application/* got %s", s)
	}

	if s := (*accept)[4].GetType(); s != "*/*" {
		t.Fatalf("expected */* got %s", s)
	}
}
