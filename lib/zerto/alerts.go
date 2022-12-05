//
// Zerto API Interface Wrapper
//
// Author: Martin Weber <martin.weber@de.clara.net>
//

package zerto

import (
	"encoding/json"
	"net/url"
)

type ZertoAlert struct {
//	AffectedVpgs
//	AffectedZorgs
	Description	string
	Entity		string
	HelpIdentifier	string
	IsDismissed	bool
	Level		string
//	Link
//	Site
	TurnedOn	string
}

func (z *Zerto) ListAlerts() []ZertoAlert {
	resp, _ := z.makeRequest("GET", "/v1/alerts", RequestParams{})
	data := json.NewDecoder(resp.Body)

	var d []ZertoAlert
	data.Decode(&d)

	return d
}

func (z *Zerto) ListErrors() []ZertoAlert {
	v := url.Values{}
	v.Add("level", "error")

	resp, _ := z.makeRequest("GET", "/v1/alerts", RequestParams{params: v})
	data := json.NewDecoder(resp.Body)

	var d []ZertoAlert
	data.Decode(&d)

	return d
}

func (z *Zerto) ListWarnings() []ZertoAlert {
	v := url.Values{}
	v.Add("level", "warning")

	resp, _ := z.makeRequest("GET", "/v1/alerts", RequestParams{params: v})
	data := json.NewDecoder(resp.Body)

	var d []ZertoAlert
	data.Decode(&d)

	return d
}
