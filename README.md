# url2id
A Go library + web server to extract *standard identifiers* from a given URL.

Intended for *all* scientific literature *only*.

## Standard Identifiers
- DOI ([Digital Object Identifier)](https://www.doi.org/)
- PMID ([PubMeb](https://www.ncbi.nlm.nih.gov/pubmed/) Identifier)
- PMCID ([PubMeb Central](http://www.ncbi.nlm.nih.gov/pmc/) Identifier)
- *Manuscript IDs*
  - NIHMSID ([U.S. National Institutes of Health Manuscripts Identifier](https://www.nihms.nih.gov/db/sub.cgi?page=overview))
  - EMSID ([Europe PubMed Central Manuscripts Identifier](https://europepmc.org))
  - CAMSID ([PubMed Central Canada Manuscripts Identifier](http://pubmedcentralcanada.ca/pmcc/static/aboutUs/))

## High-Level Description
**url2id** will fetch the web-page referred by `url` and scan it for identifiers.

## Documentation
For the library, see the source code, it's really short.

The web server will listen on `:8080` by default, but you can change the port by `PORT` environment variable.

Any GET request to `/` with parameter `url` will be responded by a JSON of the following type on success:
    
```json
{
    "doi": "10.1016/j.cub.2012.03.051",
    "pmid": "22521786",
    "pmcid": "PMC3780767",
    "nihmsid": "NIHMS236863",
    "emsid": "EMS48932",
    "camsid": "cams6038"    
}
```

On failure, appropriate HTTP status code will be set (often either 400 or 500), with a high-level error message (details
will be logged by the server).

Empty keys will be omitted.

## Notice
**url2id** does *not* rely on the structure of the websites but uses few simple regexes to extract information. This
makes it resilient to changes (on the websites) in future, fast and also simple, **but it also means that errors are probable**.

DOI extraction seems to be the most reliable of all other standard identifiers, and it's widely supported by variety of
platforms. Thus it might be wise to fetch the title of the paper using the DOI you have got from **url2id** and
cross-check (if you are intending to use this in a bibliography software for instance).

## License
ISC License, see [LICENSE](./LICENSE).

