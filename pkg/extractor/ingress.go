package extractor

import (
	"strings"
)

var ingressYAMLPatterns = []string{
	"**/gateway*.yaml",
	"**/gateway*.yml",
	"**/httproute*.yaml",
	"**/httproute*.yml",
	"**/virtualservice*.yaml",
	"**/virtualservice*.yml",
	"**/destinationrule*.yaml",
	"**/destinationrule*.yml",
	"**/ingress*.yaml",
	"**/ingress*.yml",
	"**/route*.yaml",
	"**/route*.yml",
	"**/serviceentry*.yaml",
	"**/serviceentry*.yml",
	"config/**/ingress*.yaml",
	"charts/**/templates/*ingress*.yaml",
	"charts/**/templates/*route*.yaml",
}

// supportedIngressKinds lists Kubernetes resource kinds we extract as ingress/routing.
var supportedIngressKinds = map[string]bool{
	"Gateway":         true,
	"HTTPRoute":       true,
	"Ingress":         true,
	"VirtualService":  true,
	"DestinationRule": true,
	"ServiceEntry":    true,
	"Route":           true, // OpenShift Route
}

// extractIngress finds Gateway API, Istio, and Kubernetes ingress resources.
func extractIngress(repoPath string) []IngressResource {
	files := findYAMLFiles(repoPath, ingressYAMLPatterns)
	var resources []IngressResource

	for _, fpath := range files {
		for _, doc := range parseYAMLSafe(fpath) {
			kind, _ := doc["kind"].(string)
			if !supportedIngressKinds[kind] {
				continue
			}

			metadata, _ := doc["metadata"].(map[string]interface{})
			name := ""
			if metadata != nil {
				name, _ = metadata["name"].(string)
			}

			res := IngressResource{
				Kind:   kind,
				Name:   name,
				Source: relativePath(repoPath, fpath),
			}

			spec, _ := doc["spec"].(map[string]interface{})
			if spec == nil {
				resources = append(resources, res)
				continue
			}

			switch kind {
			case "Ingress":
				res.Hosts, res.Paths, res.Backend, res.TLS = parseIngress(spec)
			case "HTTPRoute":
				res.Hosts, res.Paths, res.Backend = parseHTTPRoute(spec)
			case "Gateway":
				res.Hosts, res.TLS = parseGateway(spec)
			case "VirtualService":
				res.Hosts, res.Paths = parseVirtualService(spec)
			case "Route":
				res.Hosts, res.Paths, res.Backend, res.TLS = parseOCPRoute(spec)
			case "DestinationRule":
				res.Backend = parseDestinationRule(spec)
			case "ServiceEntry":
				res.Hosts = parseServiceEntry(spec)
			}

			if res.Hosts == nil {
				res.Hosts = []string{}
			}
			if res.Paths == nil {
				res.Paths = []string{}
			}

			resources = append(resources, res)
		}
	}

	if resources == nil {
		resources = []IngressResource{}
	}
	return resources
}

func parseIngress(spec map[string]interface{}) (hosts, paths []string, backend string, tls bool) {
	if _, ok := spec["tls"]; ok {
		tls = true
	}
	rules, _ := spec["rules"].([]interface{})
	for _, r := range rules {
		rule, ok := r.(map[string]interface{})
		if !ok {
			continue
		}
		if h, ok := rule["host"].(string); ok {
			hosts = append(hosts, h)
		}
		httpBlock, _ := rule["http"].(map[string]interface{})
		if httpBlock == nil {
			continue
		}
		pathRules, _ := httpBlock["paths"].([]interface{})
		for _, p := range pathRules {
			pm, ok := p.(map[string]interface{})
			if !ok {
				continue
			}
			if path, ok := pm["path"].(string); ok {
				paths = append(paths, path)
			}
			if be, ok := pm["backend"].(map[string]interface{}); ok {
				if svc, ok := be["service"].(map[string]interface{}); ok {
					svcName, _ := svc["name"].(string)
					backend = svcName
				} else if svcName, ok := be["serviceName"].(string); ok {
					backend = svcName
				}
			}
		}
	}
	return
}

func parseHTTPRoute(spec map[string]interface{}) (hosts, paths []string, backend string) {
	hostnames, _ := spec["hostnames"].([]interface{})
	for _, h := range hostnames {
		if s, ok := h.(string); ok {
			hosts = append(hosts, s)
		}
	}
	rules, _ := spec["rules"].([]interface{})
	for _, r := range rules {
		rule, ok := r.(map[string]interface{})
		if !ok {
			continue
		}
		matches, _ := rule["matches"].([]interface{})
		for _, m := range matches {
			mm, ok := m.(map[string]interface{})
			if !ok {
				continue
			}
			if p, ok := mm["path"].(map[string]interface{}); ok {
				if v, ok := p["value"].(string); ok {
					paths = append(paths, v)
				}
			}
		}
		backendRefs, _ := rule["backendRefs"].([]interface{})
		for _, br := range backendRefs {
			bm, ok := br.(map[string]interface{})
			if !ok {
				continue
			}
			if n, ok := bm["name"].(string); ok && backend == "" {
				backend = n
			}
		}
	}
	return
}

func parseGateway(spec map[string]interface{}) (hosts []string, tls bool) {
	listeners, _ := spec["listeners"].([]interface{})
	for _, l := range listeners {
		lm, ok := l.(map[string]interface{})
		if !ok {
			continue
		}
		if h, ok := lm["hostname"].(string); ok {
			hosts = append(hosts, h)
		}
		if p, ok := lm["protocol"].(string); ok {
			if strings.EqualFold(p, "HTTPS") || strings.EqualFold(p, "TLS") {
				tls = true
			}
		}
		if _, ok := lm["tls"]; ok {
			tls = true
		}
	}
	return
}

func parseVirtualService(spec map[string]interface{}) (hosts, paths []string) {
	hostsField, _ := spec["hosts"].([]interface{})
	for _, h := range hostsField {
		if s, ok := h.(string); ok {
			hosts = append(hosts, s)
		}
	}
	httpRules, _ := spec["http"].([]interface{})
	for _, r := range httpRules {
		rule, ok := r.(map[string]interface{})
		if !ok {
			continue
		}
		matches, _ := rule["match"].([]interface{})
		for _, m := range matches {
			mm, ok := m.(map[string]interface{})
			if !ok {
				continue
			}
			if uri, ok := mm["uri"].(map[string]interface{}); ok {
				for _, v := range uri {
					if s, ok := v.(string); ok {
						paths = append(paths, s)
					}
				}
			}
		}
	}
	return
}

func parseOCPRoute(spec map[string]interface{}) (hosts, paths []string, backend string, tls bool) {
	if h, ok := spec["host"].(string); ok {
		hosts = append(hosts, h)
	}
	if p, ok := spec["path"].(string); ok {
		paths = append(paths, p)
	}
	if to, ok := spec["to"].(map[string]interface{}); ok {
		if n, ok := to["name"].(string); ok {
			backend = n
		}
	}
	if _, ok := spec["tls"]; ok {
		tls = true
	}
	return
}

func parseDestinationRule(spec map[string]interface{}) string {
	if h, ok := spec["host"].(string); ok {
		return h
	}
	return ""
}

func parseServiceEntry(spec map[string]interface{}) []string {
	var hosts []string
	hostsField, _ := spec["hosts"].([]interface{})
	for _, h := range hostsField {
		if s, ok := h.(string); ok {
			hosts = append(hosts, s)
		}
	}
	return hosts
}
