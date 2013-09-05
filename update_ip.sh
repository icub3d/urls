# Copyright 2013 Joshua Marsh. All rights reserved.
# Use of this source code is governed by the MIT
# license that count be found in the LICENSE file.
#
#!/bin/bash

# Get the file and unzip it.
wget -q http://geolite.maxmind.com/download/geoip/database/GeoIPCountryCSV.zip 
unzip -q GeoIPCountryCSV.zip 

# Echo out the top.
echo "package urls

type ipRange struct {
  start int
  end int
  country string
}

var ipLookup = []ipRange{"  > ips.go

# Echo out each line.
awk -F, '{print $3,$4,$5;}' GeoIPCountryWhois.csv | sed -e 's/"//g' | \
		while read START END COUNTRY; do
		
		echo "ipRange{start:$START,end:$END,country:\"$COUNTRY\"}," >> ips.go
done

echo "}" >> ips.go

go fmt ips.go >/dev/null 2>&1 

# Clean up.
rm GeoIPCountryCSV.zip GeoIPCountryWhois.csv
