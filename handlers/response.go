package handlers

type response struct {
	Ok    bool
	Error string `json:",omitempty"`
}
