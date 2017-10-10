package controller

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"regexp"
	"text/tabwriter"

	"github.com/pbergman/caserver/ca"
	"github.com/pbergman/caserver/storage"
	"github.com/pbergman/logger"
)

type ApiListController struct {
	pattern *regexp.Regexp
	manager *ca.Manager
}

func (a ApiListController) Name() string {
	return "controller.api.list"
}

func (a ApiListController) Match(request *http.Request) bool {
	return a.pattern.MatchString(request.RequestURI)
}

func NewApiList(manager *ca.Manager) *ApiListController {
	return &ApiListController{
		pattern: regexp.MustCompile(`^(?i)/api/v1/list(?:/(ca|cert|csr))?(?:\?host=[^$]+)?$`),
		manager: manager,
	}
}

func (a ApiListController) getPath(req *http.Request) string {
	match := a.pattern.FindStringSubmatch(req.RequestURI)
	return match[1]
}

func (a ApiListController) Handle(resp http.ResponseWriter, req *http.Request, logger logger.LoggerInterface) {
	certs, err := a.getCerts(req)
	if err != nil {
		write_error(resp, err.Error(), http.StatusInternalServerError, logger)
		return
	}
	switch readAndSetContentType(resp, req, CONTENT_TYPE_TEXT|CONTENT_TYPE_JSON) {
	case CONTENT_TYPE_TEXT:
		writer := tabwriter.NewWriter(resp, 0, 0, 3, ' ', 0)
		for k, v := range certs {
			for c, i := len(v), 0; i < c; i++ {
				switch t := v[i].(type) {
				case *x509.CertificateRequest:
					writer.Write([]byte("[CERTIFICATE REQUEST]\t\n"))
					writer.Write([]byte(" id\t" + k + "\n"))
					a.writeMergeList(writer, " hosts", a.mergeHosts(t.DNSNames, t.IPAddresses))
					a.writeTextName(writer, t.Subject, "SUBJECT")
					writer.Write([]byte("\t\n"))
				case *x509.Certificate:
					writer.Write([]byte("[CERTIFICATE]\t\n"))
					writer.Write([]byte(" id\t" + k + "\n"))
					a.writeMergeList(writer, " hosts", a.mergeHosts(t.DNSNames, t.IPAddresses))
					a.writeTextName(writer, t.Subject, "SUBJECT")
					if !t.IsCA {
						a.writeTextName(writer, t.Issuer, "ISSUER")
					}
					writer.Write([]byte("\t\n"))
				}
			}
		}
		writer.Flush()
	case CONTENT_TYPE_JSON:
		data := make(map[string]map[string]interface{}, 0)
		for k, v := range certs {
			for c, i := len(v), 0; i < c; i++ {
				item := make(map[string]interface{}, 0)
				switch t := v[i].(type) {
				case *x509.CertificateRequest:
					item["hosts"] = a.mergeHosts(t.DNSNames, t.IPAddresses)
					item["subject"] = a.nameToMap(t.Subject)
					if data[k] == nil {
						data[k] = make(map[string]interface{})
					}
					data[k]["certificate_request"] = item
				case *x509.Certificate:
					item["hosts"] = a.mergeHosts(t.DNSNames, t.IPAddresses)
					item["subject"] = a.nameToMap(t.Subject)
					if !t.IsCA {
						item["issuer"] = a.nameToMap(t.Issuer)
					}
					if data[k] == nil {
						data[k] = make(map[string]interface{})
					}
					data[k]["certificate"] = item
				}
			}
		}
		encoder := json.NewEncoder(resp)
		encoder.SetIndent("", " ")
		if err := encoder.Encode(data); err != nil {
			write_error(resp, err.Error(), http.StatusInternalServerError, logger)
		}
	default:
		resp.WriteHeader(http.StatusNotAcceptable)
	}
}

func (a ApiListController) getCerts(req *http.Request) (map[string][]interface{}, error) {
	var path string = a.getPath(req)
	var certs map[string][]interface{} = make(map[string][]interface{})
	err := a.manager.Each(func(r storage.Record) bool {
		if path == "ca" && !r.IsCa() {
			return true
		}
		if (path == "cert" || path == "csr") && r.IsCa() {
			return true
		}
		items := make([]interface{}, 0)
		if path == "cert" || path == "ca" || path == "" {
			if cert := r.GetCertificate(); cert != nil {
				if host := req.URL.Query().Get("host"); path == "cert" && host != "" {
					if e := cert.VerifyHostname(host); e == nil {
						items = append(items, cert)
					}
				} else {
					items = append(items, cert)
				}
			}
		}
		if path == "csr" {
			if cert := r.GetCertificateRequest(); cert != nil {
				items = append(items, cert)
			}
		}
		if len(items) > 0 {
			certs[r.GetId().String()] = items
		}
		return true
	})
	return certs, err
}

func (a ApiListController) mergeHosts(dns []string, ip []net.IP) []string {
	hosts := []string{}
	for f, k := 0, len(dns); f < k; f++ {
		hosts = append(hosts, dns[f])
	}
	for f, k := 0, len(ip); f < k; f++ {
		hosts = append(hosts, ip[f].String())
	}
	return hosts
}

func (a ApiListController) nameToMap(name pkix.Name) map[string]interface{} {
	ret := make(map[string]interface{})
	ret["common_name"] = name.CommonName
	ret["serial_number"] = name.SerialNumber
	if s := len(name.Country); s > 0 {
		ret["country"] = make([]string, s)
		for i := 0; i < s; i++ {
			ret["country"].([]string)[i] = name.Country[i]
		}
	}
	if s := len(name.Organization); s > 0 {
		ret["organization"] = make([]string, s)
		for i := 0; i < s; i++ {
			ret["organization"].([]string)[i] = name.Organization[i]
		}
	}
	if s := len(name.OrganizationalUnit); s > 0 {
		ret["organizational_unit"] = make([]string, s)
		for i := 0; i < s; i++ {
			ret["organizational_unit"].([]string)[i] = name.OrganizationalUnit[i]
		}
	}
	if s := len(name.Locality); s > 0 {
		ret["locality"] = make([]string, s)
		for i := 0; i < s; i++ {
			ret["locality"].([]string)[i] = name.Locality[i]
		}
	}
	if s := len(name.Province); s > 0 {
		ret["province"] = make([]string, s)
		for i := 0; i < s; i++ {
			ret["province"].([]string)[i] = name.Province[i]
		}
	}
	if s := len(name.StreetAddress); s > 0 {
		ret["street_address"] = make([]string, s)
		for i := 0; i < s; i++ {
			ret["street_address"].([]string)[i] = name.StreetAddress[i]
		}
	}
	if s := len(name.PostalCode); s > 0 {
		ret["postalcode"] = make([]string, s)
		for i := 0; i < s; i++ {
			ret["postalcode"].([]string)[i] = name.PostalCode[i]
		}
	}
	return ret
}

func (a ApiListController) writeTextName(writer io.Writer, name pkix.Name, header string) {
	writer.Write([]byte("[" + header + "]\t\n"))
	writer.Write([]byte(" common name\t" + name.CommonName + "\n"))
	writer.Write([]byte(" serial number\t" + name.SerialNumber + "\n"))
	a.writeMergeList(writer, " country", name.Country)
	a.writeMergeList(writer, " organization", name.Organization)
	a.writeMergeList(writer, " organizational unit", name.OrganizationalUnit)
	a.writeMergeList(writer, " locality", name.Locality)
	a.writeMergeList(writer, " province", name.Province)
	a.writeMergeList(writer, " street address", name.StreetAddress)
	a.writeMergeList(writer, " postalcode", name.PostalCode)
}

func (a ApiListController) writeMergeList(writer io.Writer, name string, list []string) {
	if l := len(list); l > 0 {
		writer.Write([]byte(name + "\t"))
		for i := 0; i < l; i++ {
			if i == 0 {
				writer.Write([]byte(list[i]))
			} else {
				writer.Write([]byte(", " + list[i]))
			}
		}
		writer.Write([]byte("\n"))
	}
}
