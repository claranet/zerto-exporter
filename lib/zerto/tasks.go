//
// Zerto API Interface Wrapper
//
// Author: Martin Weber <martin.weber@de.clara.net>
//

package zerto

import "encoding/json"

type ZertoTask struct {
	CompleteReason	string
	Completed				string
	InitiatedBy			string
	IsCancellable		bool
	Started					string
	Status					*ZertoTaskStatus
	TaskIdentifier	string
	Type						string
}

type ZertoTaskStatus struct {
	Progress	int
	State		int
}

func (z *Zerto) ListTasks() []ZertoTask {
	resp, _ := z.makeRequest("GET", "/v1/tasks", RequestParams{})
	data := json.NewDecoder(resp.Body)

	var d []ZertoTask
	data.Decode(&d)

	return d
}
