package pkg

import (
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"time"
)

type Result struct {
	Doi string `json:"doi,omitempty"`
	Pmid string `json:"pmid,omitempty"`
	Pmcid string `json:"pmcid,omitempty"`
	Nihmsid string `json:"nihmsid,omitempty"`
	Emsid string `json:"emsid,omitempty"`
	Camsid string `json:"camsid,omitempty"`
}

type URL2IDError struct {
	Short string
	PrevErr error
}

var doiRE = regexp.MustCompile(`doi\.org/([^"]+)"`)
var pmidRE1 = regexp.MustCompile(`/pubmed/(\d+)"`)
var pmidRE2 = regexp.MustCompile(`/pmid/(\d+)/?`)
var pmcidRE = regexp.MustCompile(`(?i)PMC(\d+)`)
var nihmsidRE = regexp.MustCompile(`(?i)NIHMS(\d+)`)
var emsidRE = regexp.MustCompile(`(?i)EMS(\d+)`)
var camsidRE = regexp.MustCompile(`(?i)cams(\d+)`)

func URL2ID(url_ *url.URL) (*Result, *URL2IDError) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, &URL2IDError{
			Short: "error with initialising HTTP client",
			PrevErr: err,
		}
	}

	client := &http.Client{
		Jar: jar,
		Timeout: 10 * time.Second,
	}

	req, _ := http.NewRequest("GET", url_.String(), nil)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Host", url_.Host)
	req.Header.Set("User-Agent", "what")
	req.Close = true

	response, err := client.Do(req)
	if err != nil {
		return nil, &URL2IDError{
			Short: "could not fetch response",
			PrevErr: err,
		}
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, &URL2IDError{
			Short: "could not read the body",
			PrevErr: err,
		}
	}

	var result Result
	if matches := doiRE.FindSubmatch(body); matches != nil {
		result.Doi, err = url.QueryUnescape(string(matches[1]))
		if err != nil {
			return nil, &URL2IDError{
				Short: "could not decode DOI",
				PrevErr: nil,
			}
		}
	}

	if matches := pmidRE1.FindSubmatch(body); matches != nil {
		result.Pmid = string(matches[1])
	} else if matches := pmidRE2.FindSubmatch(body); matches != nil {
		result.Pmid = string(matches[1])
	}

	if matches := pmcidRE.FindSubmatch(body); matches != nil {
		result.Pmcid = "PMC" + string(matches[1])
	}

	// A paper can have only ONE of the following manuscript identifiers.
	//
	// NIH website has a quirk where the assets of a paper are prefixed with "nihms".
	// For instance, the paper with CAMSID 6038 <https://www.ncbi.nlm.nih.gov/pmc/articles/mid/CAMS6038>
	// has its PDF at https://www.ncbi.nlm.nih.gov/pmc/articles/mid/CAMS6038/pdf/nihms6038.pdf
	// which causes us to return wrong results since there is another (different) article whose NIHMSID is 6038.
	//
	// Therefore, a paper can have only ONE of the following manuscript identifiers, we assume. There might be a few
	// exceptions to this assumption of ours, but they are probably insignificantly small in number to worth the
	// trade-off.
	if matches := camsidRE.FindSubmatch(body); matches != nil {
		result.Camsid = "cams" + string(matches[1])
	} else if matches := emsidRE.FindSubmatch(body); matches != nil {
		result.Emsid = "EMS" + string(matches[1])
	} else if matches := nihmsidRE.FindSubmatch(body); matches != nil {
		result.Nihmsid = "NIHMS" + string(matches[1])
	}

	return &result, nil
}

func (ue URL2IDError) Error() string {
	return ue.Short
}
