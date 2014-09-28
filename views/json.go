package views

import (
	"bytes"
	"encoding/json"
	"github.com/greivinlopez/skue"
	"io"
	"io/ioutil"
	"net/http"
)

type JSONProducer struct {
}

func (producer JSONProducer) MimeType() string {
	return skue.MIME_JSON
}

func (producer JSONProducer) Out(w http.ResponseWriter, statusCode int, value interface{}) {
	output, err := json.MarshalIndent(value, " ", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
	} else {
		w.Header().Set(skue.HEADER_ContentType, skue.MIME_JSON)
		w.WriteHeader(statusCode)
		w.Write(output)
	}
}

type JSONConsumer struct {
}

func (consumer JSONConsumer) MimeType() string {
	return skue.MIME_JSON
}

func (consumer JSONConsumer) In(r *http.Request, value interface{}) error {
	buffer, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(buffer)
	decoder := json.NewDecoder(reader)
	decoder.UseNumber()
	return decoder.Decode(value)
}

type JSONView skue.ViewLayer
