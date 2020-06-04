package bitcoinwalletwatcher

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

func TestNewInfoStorage(t *testing.T) {
	expected := &InfoFile{
		CurrentBlock: DefaultCurrentBlock,
	}

	s, err := NewInfoStorage("")
	if err != nil {
		t.Errorf("did not expect an error but got %v", err)
	}

	if !reflect.DeepEqual(s, expected) {
		t.Errorf("expected %v but got %v", expected, s)
	}
}

func TestNewInfoStorageWithFile(t *testing.T) {
	filename := "/tmp/.info"
	f, err := os.Create(filename)
	if err != nil {
		t.Errorf("error creating file %s", filename)
	}

	expected := &InfoFile{
		CurrentBlock: DefaultCurrentBlock,
	}

	b, err := json.Marshal(expected)
	if err != nil {
		t.Error(err)
	}

	if _, err := f.Write(b); err != nil {
		t.Error(err)
	}

	s, err := NewInfoStorage(filename)
	if err != nil {
		t.Errorf("did not expect an error but got %v", err)
	}

	if !reflect.DeepEqual(s, expected) {
		t.Errorf("expected %v but got %v", expected, s)
	}

	if err := os.Remove(filename); err != nil {
		t.Error(err)
	}
}

func TestNewInfoStorageWithJsonError(t *testing.T) {
	filename := "/tmp/.info"
	f, err := os.Create(filename)
	if err != nil {
		t.Errorf("error creating file %s", filename)
	}

	b := []byte("test")
	if _, err := f.Write(b); err != nil {
		t.Error(err)
	}

	_, err = NewInfoStorage(filename)
	if err == nil {
		t.Errorf("did expect an error but got nil")
	}

	if err := os.Remove(filename); err != nil {
		t.Error(err)
	}
}
