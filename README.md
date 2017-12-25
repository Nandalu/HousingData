## Data
http://plvr.land.moi.gov.tw/DownloadOpenData

## Publishing to Jinma
### Prepare the geocoding cache (Optional)
In the case where the geocoding service is too slow,
run cmd/geocode to generate a list of addresses,
which can then be used to prefetch the location of these addresses.
The result of the prefetching should be presented in the following form
for use in the next steps:

```
Addr      string
Lat       float64
Lng       float64
Precision float64
```

### Parse the raw data
Run cmd/parse

### Upload to Jinma
Run cmd/pub.
*IMPORTANT* Note that in order to prevent duplicate jinma.Msgs resulting
from the rerun of the job in face of failures, we have to set
the infileOffset flag to the line number which we should continue from.
The line number can be read from the logs of the previous aborted run.
