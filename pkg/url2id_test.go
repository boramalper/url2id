package pkg

import (
	"net/url"
	"reflect"
	"testing"
)

var cases = []struct{
	url string
	exRes Result
}{
	{
		"https://www.ncbi.nlm.nih.gov/pmc/articles/PMC5075450/",
		Result {
			Doi: "10.1016/j.jneumeth.2016.04.012",
			Pmid: "27102043",
			Pmcid: "PMC5075450",
			Camsid: "cams6038",
		},
	},
	{
		"https://www.ncbi.nlm.nih.gov/pubmed/21765645",
		Result {
			Doi: "10.1016/j.ijms.2010.08.003",
			Pmid: "21765645",
			Pmcid: "PMC3134971",
			Nihmsid: "NIHMS236863",
		},
	},
	{
		"https://www.ncbi.nlm.nih.gov/pmc/articles/mid/EMS48932/",
		Result {
			Doi: "10.1016/j.cub.2012.03.051",
			Pmid: "22521786",
			Pmcid: "PMC3780767",
			Emsid: "EMS48932",
		},
	},
	{
		"https://www.sciencedirect.com/science/article/pii/S0030666510000575",
		Result {
			Doi: "10.1016/j.otc.2010.04.002",
		},
	},
	{
		"https://www.sciencedirect.com/science/article/pii/S0196677403000762",
		Result {
			Doi: "10.1016/S0196-6774(03)00076-2",
		},
	},
	{
		"http://onlinelibrary.wiley.com/doi/10.1038/sj.bjp.0702844/full", // will redirect
		Result {
			Doi: "10.1038/sj.bjp.0702844",
		},
	},
	{
		"https://www.ingentaconnect.com/content/ben/mrmc/2007/00000007/00000006/art00004",
		Result {
			Doi: "10.2174/138955707780859431",
		},
	},
	{
		"https://link.springer.com/article/10.1007/BF01197436",
		Result {
			Doi: "10.1007/BF01197436",
		},
	},
	{
		"https://journals.sagepub.com/doi/abs/10.1177/0954407017737901",
		Result {
			Doi: "10.1177/0954407017737901",
		},
	},
	{
		"https://aip.scitation.org/doi/10.1063/1.871778",
		Result {
			Doi: "10.1063/1.871778",
		},
	},
	{
		"https://ajp.psychiatryonline.org/doi/abs/10.1176/appi.ajp.159.3.394",
		Result {
			Doi: "10.1176/appi.ajp.159.3.394",
		},
	},
}

func TestCases(t *testing.T) {
	for i, c := range cases {
		no := i + 1

		url_, err := url.Parse(c.url)
		if err != nil {
			t.Fatalf("couldn't parse URL #%d", no)
		}

		res, err2 := URL2ID(url_)
		if err2 != nil {
			t.Fatalf("#%d: %s\n%+v", no, err2.Error(), err2.PrevErr)
		}

		if !reflect.DeepEqual(c.exRes, *res) {
			t.Errorf("#%d yielded unexpected result:\nGot     : %+v\nExpected: %+v", no, *res, c.exRes)
		}
	}
}
