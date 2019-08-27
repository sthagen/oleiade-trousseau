package trousseau

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

type Meta struct {
	CreatedAt        string   `json:"created_at"`
	LastModifiedAt   string   `json:"last_modified_at"`
	Recipients       []string `json:"recipients"`
	TrousseauVersion string   `json:"version"`
}

func (m *Meta) ListRecipients() []string {
	return m.Recipients
}

func (m *Meta) AddRecipient(recipient string) error {
	in, _ := m.containsRecipient(recipient)
	if in {
		errMsg := fmt.Sprintf("Recipient %s already mapped to store metadata", recipient)
		return errors.New(errMsg)
	} else {
		m.Recipients = append(m.Recipients, recipient)
	}

	return nil
}

func (m *Meta) RemoveRecipient(recipient string) error {
	in, idx := m.containsRecipient(recipient)
	if !in {
		errMsg := fmt.Sprintf("Unknown recipient: %s", recipient)
		return errors.New(errMsg)
	} else if m.isLastRecipient(recipient) {
		return errors.New("Forbidden: removing last data store recipient")
	} else {
		newRecipients := make([]string, len(m.Recipients)-1)
		copy(newRecipients[0:idx], m.Recipients[0:idx])
		copy(newRecipients[:idx], m.Recipients[:idx+1])
		m.Recipients = newRecipients
	}

	return nil

}

func (m *Meta) updateLastModificationMarker() {
	m.LastModifiedAt = time.Now().String()
}

func (m *Meta) containsRecipient(recipient string) (status bool, index int) {
	for index, r := range m.Recipients {
		if r == recipient {
			return true, index
		}
	}

	return false, -1
}

func (m *Meta) isLastRecipient(recipient string) bool {
	if len(m.Recipients) == 1 {
		return true
	}

	return false
}

func (m *Meta) String() string {
	var repr string
	metaType := reflect.TypeOf(*m)
	metaValue := reflect.ValueOf(*m)

	for i := 0; i < metaType.NumField(); i++ {
		key := metaType.Field(i).Name
		value := metaValue.FieldByName(key).Interface()
		repr += fmt.Sprintf("%s : %v\n", key, value)
	}

	return repr
}
