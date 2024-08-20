package events

type Comment struct {
	ID      int            `json:"id"`
	Content CommentContent `json:"content"`
	User    struct {
		DisplayName string `json:"display_name"`
	} `json:"user"`
}

type CommentContent struct {
	Raw string `json:"raw"`
}

type PullRequestSource struct {
	Branch Branch `json:"branch"`
}

type Branch struct {
	Name string `json:"name"`
}

type PullRequestLinks struct {
	HTML HTML `json:"html"`
}

type HTML struct {
	Href string `json:"href"`
}

type Destination struct {
	Branch Branch `json:"branch"`
}

type Reviewer struct {
	Approved    bool   `json:"approved"`
	DisplayName string `json:"display_name"`
}

type PullRequest struct {
	ID          int               `json:"id"`
	Title       string            `json:"title"`
	State       string            `json:"state"`
	Links       PullRequestLinks  `json:"links"`
	Source      PullRequestSource `json:"source"`
	Destination Destination       `json:"destination"`
	Reviewers   []Reviewer        `json:"reviewers"`
}

type Actor struct {
	DisplayName string `json:"display_name"`
}

type Repository struct {
	Name string `json:"name"`
}

type PullRequestEvent struct {
	Repository  Repository  `json:"repository"`
	Actor       Actor       `json:"actor"`
	PullRequest PullRequest `json:"pullRequest"`
	Comment     Comment     `json:"comment"`
}
