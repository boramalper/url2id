# url2id
A client-side JavaScript library to extract *standard identifiers* from a given URL.

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
**url2id** will try matching the supplied URL to one of the patterns in its database to extract an identifier, and then
it will fetch the web-page referred by `url` and scan it for more identifiers.

## Documentation
**url2id** provides a single `url2id()` function with the following signature:

`url2id(url, [doFetch], function callback(err, result))`

- `url` must be either a [String](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/String),
  or a [URL](https://developer.mozilla.org/en-US/docs/Web/API/URL), or a
  [Location](https://developer.mozilla.org/en-US/docs/Web/API/Location).
  
- `doFetch` must be a [Boolean](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Boolean)
  which if true, will cause the web-page referred by `url` to be fetched and scanned for standard identifiers.
  
  True by default.
  
- `callback` must be a [Function](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Function)
  of type `function (err, result)`, which will be called upon a successful result or an error. The return value is
  ignored.
  
  - `err` is a [String](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/String) that
    describes the error (including exceptions). `null` if operation was successful.
    
  - `result` is an [Object](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Object) of
    the following type:
    
    ```javascript
    {
        doi: "10.1016/j.cub.2012.03.051",
        pmid: "22521786",
        pmcid: "PMC3780767",
        nihmsid: "NIHMS236863",
        emsid: "EMS48932",
        camsid: "cams6038"    
    }
    ```
    
    The keys shall be self-descriptive. Beware that `pmid` does *not* have any alphabetic prefix, whereas `pmcid`,
    `nihmsid`, `emsid`, and `camsid` do; this is by convention.
    
    `null` if operation was unsuccessful.

### Example 
```javascript
url2id("https://www.ncbi.nlm.nih.gov/pmc/articles/mid/EMS48932/", function(error, result) {
    if (error) {
        console.error("url2id error", error);
        return;
    }
    
    console.info("DOI:", result.doi);
    console.info("PMID:", result.pmid);
    console.info("PMCID:", result.pmcid);
    console.info("NIHMSID:", result.nihmsid);
    console.info("EMSID:", result.emsid);
    console.info("CAMSID:", result.camsid);
});
```

## License
ISC License, see [LICENSE](./LICENSE).

