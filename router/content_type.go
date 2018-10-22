package router

type ContentType uint16

const (
	ContentTypeText ContentType = 1 << iota
	ContentTypeJson
	ContentTypeTar
	ContentTypeTarGzip
	ContentTypePkixCert

	ContentTypeAll ContentType = ContentTypeText | ContentTypeJson | ContentTypeTar | ContentTypeTarGzip | ContentTypePkixCert
)

func ContentTypeFromString(types ...string) ContentType {
	var ct ContentType
	for _, t := range types {
		switch t {
		case "text/plain":
			ct |= ContentTypeText
		case "application/json":
			ct |= ContentTypeJson
		case "application/tar":
			ct |= ContentTypeTar
		case "application/tar+gzip":
			ct |= ContentTypeTarGzip
		case "application/pkix-cert":
			ct |= ContentTypePkixCert
		}
	}
	return ct
}

func (c ContentType) String() string {
	buf := ""
	for i := 1; i < int(ContentTypeAll); i <<= 1 {
		if i == (i & int(c)) {
			switch ContentType(i) {
			case ContentTypeText:
				buf += "text/plain"
			case ContentTypeJson:
				buf += ", application/json"
			case ContentTypeTar:
				buf += ", application/tar"
			case ContentTypeTarGzip:
				buf += ", application/tar+gzip"
			case ContentTypePkixCert:
				buf += ", application/pkix-cert"
			}
		}
	}
	if len(buf) > 2 && buf[0] == ',' && buf[1] == ' ' {
		buf = buf[2:]
	}
	return buf
}
