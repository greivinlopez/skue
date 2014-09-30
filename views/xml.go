package views

import (
	"bytes"
	"encoding/xml"
	"github.com/greivinlopez/skue"
	"io"
	"io/ioutil"
	"net/http"
)

type XmlProducer struct {
}

func (producer XmlProducer) MimeType() string {
	return skue.MIME_XML
}

func (producer XmlProducer) Out(w http.ResponseWriter, statusCode int, value interface{}) {
	output, err := xml.MarshalIndent(value, " ", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
	} else {
		w.Header().Set(skue.HEADER_ContentType, skue.MIME_JSON)
		w.WriteHeader(statusCode)
		w.Write(output)
	}
}

type XmlConsumer struct {
}

func (consumer XmlConsumer) MimeType() string {
	return skue.MIME_JSON
}

func (consumer XmlConsumer) In(r *http.Request, value interface{}) error {
	buffer, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(buffer)
	return xml.NewDecoder(reader).Decode(value)
}

type XmlView skue.ViewLayer

func NewXmlView() *skue.ViewLayer {
	return skue.NewViewLayer(XmlProducer{}, XmlConsumer{})
}
