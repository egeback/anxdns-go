package anxdns

type Data struct {
	Domain  string `json:"domain"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	TTL     int    `json:"ttl"`
	Address string `json:"address,omitempty"`
	TxtData string `json:"txtData,omitempty"`
	Line    int    `json:"line,omitempty"`
}
